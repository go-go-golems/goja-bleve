package pkg

import (
	"strings"
	"testing"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
)

func TestMappingBuildersCreateUsableBleveMapping(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).includeTermVectors(true).build();
		const keyword = bleve.field().keyword().store(true).build();
		const number = bleve.field().number().store(true).build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("text", text)
			.field("source_id", keyword)
			.field("chunk_index", number)
			.build();
		bleve.mapping()
			.defaultMapping(doc)
			.defaultAnalyzer("standard")
			.defaultField("text")
			.storeDynamic(false)
			.indexDynamic(false)
			.build();
	`)
	if err != nil {
		t.Fatalf("build JS mapping: %v", err)
	}
	rt := newRuntime(vm)
	mappingRef, err := getTypedRef[mappingRef](rt, value, "index mapping")
	if err != nil {
		t.Fatalf("extract mapping ref: %v", err)
	}

	idx, err := bleve.NewMemOnly(mappingRef.mapping)
	if err != nil {
		t.Fatalf("create mem index: %v", err)
	}
	defer func() { _ = idx.Close() }()

	if err := idx.Index("chunk-1", map[string]any{
		"text":        "privacy screen trees",
		"source_id":   "tree-center",
		"chunk_index": 1,
		"unmapped":    "this field should not be indexed",
	}); err != nil {
		t.Fatalf("index doc: %v", err)
	}

	textResult, err := idx.Search(bleve.NewSearchRequest(bleve.NewMatchQuery("privacy")))
	if err != nil {
		t.Fatalf("search mapped text: %v", err)
	}
	if textResult.Total != 1 {
		t.Fatalf("mapped text total = %d, want 1", textResult.Total)
	}

	unmappedResult, err := idx.Search(bleve.NewSearchRequest(bleve.NewMatchQuery("indexed")))
	if err != nil {
		t.Fatalf("search unmapped text: %v", err)
	}
	if unmappedResult.Total != 0 {
		t.Fatalf("unmapped text total = %d, want 0", unmappedResult.Total)
	}
}

func TestMappingBuildersRejectWrongWrapperTypes(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		const fieldBuilder = bleve.field().text();
		bleve.docMapping().field("text", fieldBuilder);
	`)
	if err == nil {
		t.Fatalf("expected docMapping.field to reject unbuilt field builder")
	}
	if got := err.Error(); got == "" || !containsAll(got, []string{"expected field mapping wrapper", "fieldBuilder"}) {
		t.Fatalf("wrong wrapper error = %q", got)
	}
}

func TestFieldBuilderRequiresTypeBeforeBuild(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		bleve.field().build();
	`)
	if err == nil {
		t.Fatalf("expected field build without type to fail")
	}
	if got := err.Error(); got == "" || !containsAll(got, []string{"field type is required"}) {
		t.Fatalf("field build error = %q", got)
	}
}

func newBleveTestVM() *goja.Runtime {
	vm := goja.New()
	reg := require.NewRegistry()
	modules.DefaultRegistry.Enable(reg)
	reg.Enable(vm)
	return vm
}

func containsAll(s string, parts []string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
