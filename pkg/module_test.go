package pkg

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
)

func TestRequireBleveExportsPhaseZeroFactories(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	modules.DefaultRegistry.Enable(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const bleve = require("bleve");
		({
			version: bleve.version,
			vectorSupport: bleve.vectorSupport,
			mapping: typeof bleve.mapping,
			indexMapping: typeof bleve.indexMapping,
			docMapping: typeof bleve.docMapping,
			documentMapping: typeof bleve.documentMapping,
			field: typeof bleve.field,
			search: typeof bleve.search,
			searchRequest: typeof bleve.searchRequest,
			create: typeof bleve.create,
			open: typeof bleve.open,
			memory: typeof bleve.memory,
			matchAll: typeof bleve.matchAll,
			matchNone: typeof bleve.matchNone,
		});
	`)
	if err != nil {
		t.Fatalf("require bleve: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["version"] != "0.1.0" {
		t.Fatalf("version = %#v", got["version"])
	}
	if _, ok := got["vectorSupport"].(bool); !ok {
		t.Fatalf("vectorSupport = %#v, want bool", got["vectorSupport"])
	}
	for _, name := range []string{"mapping", "indexMapping", "docMapping", "documentMapping", "field", "search", "searchRequest", "create", "open", "memory", "matchAll", "matchNone"} {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function; all exports: %#v", name, got[name], got)
		}
	}
}

func TestHiddenRefsAreNonEnumerableAndNonJSON(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	modules.DefaultRegistry.Enable(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const bleve = require("bleve");
		const field = bleve.field();
		({
			keys: Object.keys(field),
			names: Object.getOwnPropertyNames(field),
			json: JSON.stringify(field),
			type: field.type,
		});
	`)
	if err != nil {
		t.Fatalf("inspect field wrapper: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["type"] != string(refKindFieldBuilder) {
		t.Fatalf("type = %#v", got["type"])
	}
	for _, name := range got["keys"].([]any) {
		if name == hiddenRefKey {
			t.Fatalf("hidden ref key should not be enumerable: %#v", got)
		}
	}
	if json := got["json"].(string); json != `{"type":"fieldBuilder"}` {
		t.Fatalf("JSON.stringify(field) = %s", json)
	}
}

func TestTypedRefErrorsAreSpecific(t *testing.T) {
	vm := goja.New()
	rt := newRuntime(vm)
	plain := vm.NewObject()
	if _, err := getTypedRef[fieldBuilderRef](rt, plain, "field builder"); err == nil || err.Error() != "bleve: expected field builder wrapper, got value without Go reference" {
		t.Fatalf("plain object error = %v", err)
	}
	mapping := rt.mappingBuilder()
	if _, err := getTypedRef[fieldBuilderRef](rt, mapping, "field builder"); err == nil {
		t.Fatalf("expected wrong-type error")
	}
}
