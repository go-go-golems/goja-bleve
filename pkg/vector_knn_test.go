//go:build vectors

package pkg

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestJSCanCreateVectorFieldAndRunKNN(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	value, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).build();
		const embedding = bleve.field().vector(4).similarity("cosine").optimizedFor("recall").build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("text", text)
			.field("embedding", embedding)
			.build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.create(%q).mapping(mapping).build();
		idx.index("chunk-1", { text: "alpha", embedding: [1, 0, 0, 0] });
		idx.index("chunk-2", { text: "beta", embedding: [0, 1, 0, 0] });
		idx.index("chunk-3", { text: "gamma", embedding: [0.9, 0.1, 0, 0] });
		const req = bleve.search()
			.query(bleve.matchNone())
			.knnOperator("or")
			.knn("embedding", [1, 0, 0, 0], 2, 1.0)
			.fields(["text"])
			.build();
		const result = idx.search(req);
		idx.close();
		({ total: result.total, firstID: result.hits[0] && result.hits[0].id, secondID: result.hits[1] && result.hits[1].id });
	`, indexPath))
	if err != nil {
		t.Fatalf("vector KNN script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["total"] != int64(2) && got["total"] != uint64(2) && got["total"] != float64(2) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["firstID"] != "chunk-1" {
		t.Fatalf("firstID = %#v", got["firstID"])
	}
}

func TestVectorSearchBuilderRejectsInvalidKNNInputs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		script   string
		contains []string
	}{
		{
			name:     "invalid k",
			script:   `const bleve = require("bleve"); bleve.search().knn("embedding", [1, 0, 0, 0], 0);`,
			contains: []string{"k must be positive"},
		},
		{
			name:     "non finite vector",
			script:   `const bleve = require("bleve"); bleve.search().knn("embedding", [1, NaN, 0, 0], 2);`,
			contains: []string{"finite number"},
		},
		{
			name:     "invalid operator",
			script:   `const bleve = require("bleve"); bleve.search().knnOperator("xor");`,
			contains: []string{"operator", "or", "and"},
		},
	} {
		vm := newBleveTestVM()
		_, err := vm.RunString(tc.script)
		if err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
		if got := err.Error(); !containsAll(got, tc.contains) {
			t.Fatalf("%s: error = %q", tc.name, got)
		}
	}
}

func TestVectorMappingRejectsUnsupportedSimilarity(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		const embedding = bleve.field().vector(4).similarity("not-a-similarity").build();
		const doc = bleve.docMapping().dynamic(false).field("embedding", embedding).build();
		bleve.mapping().defaultMapping(doc).build();
	`)
	if err == nil {
		t.Fatalf("expected unsupported vector similarity to fail")
	}
	if got := err.Error(); !containsAll(got, []string{"invalid similarity"}) {
		t.Fatalf("unsupported similarity error = %q", got)
	}
}

func TestMissingVectorFieldFailsBeforeSearch(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	_, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).build();
		const doc = bleve.docMapping().dynamic(false).field("text", text).build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.create(%q).mapping(mapping).build();
		idx.index("chunk-1", { text: "alpha" });
		const req = bleve.search().query(bleve.matchNone()).knn("embedding", [1, 0, 0, 0], 2).build();
		idx.search(req);
	`, indexPath))
	if err == nil {
		t.Fatalf("expected missing vector field to fail")
	}
	if got := err.Error(); !containsAll(got, []string{"kNN field", "not mapped"}) {
		t.Fatalf("missing vector field error = %q", got)
	}
}

func TestVectorDimensionMismatchFails(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	_, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const embedding = bleve.field().vector(4).similarity("cosine").build();
		const doc = bleve.docMapping().dynamic(false).field("embedding", embedding).build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.create(%q).mapping(mapping).build();
		idx.index("chunk-1", { embedding: [1, 0, 0, 0] });
		const req = bleve.search().query(bleve.matchNone()).knn("embedding", [1, 0, 0], 2).build();
		idx.search(req);
	`, indexPath))
	if err == nil {
		t.Fatalf("expected vector dimension mismatch to fail")
	}
	if got := err.Error(); !containsAll(got, []string{"kNN", "vector"}) {
		t.Fatalf("dimension mismatch error = %q", got)
	}
}
