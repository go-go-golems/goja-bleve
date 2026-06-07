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

func TestReopenedVectorIndexUsesStoredMapping(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	value, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).build();
		const embedding = bleve.field().vector(4).similarity("cosine").build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("text", text)
			.field("embedding", embedding)
			.build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const created = bleve.create(%q).mapping(mapping).build();
		created.index("chunk-1", { text: "alpha", embedding: [1, 0, 0, 0] });
		created.index("chunk-2", { text: "beta", embedding: [0, 1, 0, 0] });
		created.close();

		const reopened = bleve.open(%q).build();
		const req = bleve.search()
			.query(bleve.matchNone())
			.knnOperator("or")
			.knn("embedding", [1, 0, 0, 0], 1)
			.fields(["text"])
			.build();
		const result = reopened.search(req);
		reopened.close();
		({ total: result.total, firstID: result.hits[0] && result.hits[0].id });
	`, indexPath, indexPath))
	if err != nil {
		t.Fatalf("reopened vector KNN script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["total"] != int64(1) && got["total"] != uint64(1) && got["total"] != float64(1) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["firstID"] != "chunk-1" {
		t.Fatalf("firstID = %#v", got["firstID"])
	}
}

func TestJSCanRunHybridRRFFusion(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	value, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const color = bleve.field().text().store(true).includeTermVectors(true).build();
		const colorvect = bleve.field().vector(3).similarity("l2_norm").optimizedFor("recall").build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("color", color)
			.field("colorvect", colorvect)
			.build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.create(%q).mapping(mapping).build();
		idx.index("dark blue", { color: "dark blue", colorvect: [0, 0, 139] });
		idx.index("dark slate blue", { color: "dark slate blue", colorvect: [72, 61, 139] });
		idx.index("navy", { color: "navy", colorvect: [0, 0, 128] });
		idx.index("blue", { color: "blue", colorvect: [0, 0, 255] });
		idx.index("medium blue", { color: "medium blue", colorvect: [0, 0, 205] });
		idx.index("royal blue", { color: "royal blue", colorvect: [65, 105, 225] });
		const req = bleve.search()
			.query(bleve.matchPhrase("dark").field("color"))
			.knn("colorvect", [0, 0, 129], 5, 1.0)
			.knn("colorvect", [0, 0, 250], 5, 1.0)
			.score("rrf")
			.scoreRankConstant(1)
			.scoreWindowSize(10)
			.size(5)
			.fields(["color"])
			.explain(true)
			.build();
		const result = idx.search(req);
		idx.close();
		({ total: result.total, ids: result.hits.map(h => h.id), scores: result.hits.map(h => h.score), hasExplanation: !!result.hits[0].explanation });
	`, indexPath))
	if err != nil {
		t.Fatalf("hybrid RRF script: %v", err)
	}
	got := value.Export().(map[string]any)
	ids := got["ids"].([]any)
	if len(ids) < 3 {
		t.Fatalf("expected at least 3 hybrid hits, got %#v", ids)
	}
	if ids[0] != "dark blue" {
		t.Fatalf("first hybrid hit = %#v, want dark blue (ids=%#v)", ids[0], ids)
	}
	scores := got["scores"].([]any)
	if scores[0].(float64) <= scores[1].(float64) {
		t.Fatalf("expected first fused score to be greater than second, scores=%#v ids=%#v", scores, ids)
	}
}

func TestJSCanRunHybridRSFFusion(t *testing.T) {
	vm := newBleveTestVM()
	indexPath := filepath.Join(t.TempDir(), "idx")
	value, err := vm.RunString(fmt.Sprintf(`
		const bleve = require("bleve");
		const color = bleve.field().text().store(true).build();
		const colorvect = bleve.field().vector(3).similarity("l2_norm").optimizedFor("recall").build();
		const doc = bleve.docMapping().dynamic(false).field("color", color).field("colorvect", colorvect).build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.create(%q).mapping(mapping).build();
		idx.index("dark blue", { color: "dark blue", colorvect: [0, 0, 139] });
		idx.index("dark slate blue", { color: "dark slate blue", colorvect: [72, 61, 139] });
		idx.index("navy", { color: "navy", colorvect: [0, 0, 128] });
		idx.index("blue", { color: "blue", colorvect: [0, 0, 255] });
		idx.index("medium blue", { color: "medium blue", colorvect: [0, 0, 205] });
		const req = bleve.search()
			.query(bleve.matchPhrase("dark").field("color"))
			.knn("colorvect", [0, 0, 129], 5, 1.0)
			.score("rsf")
			.scoreWindowSize(10)
			.size(5)
			.fields(["color"])
			.build();
		const result = idx.search(req);
		idx.close();
		({ total: result.total, ids: result.hits.map(h => h.id), scores: result.hits.map(h => h.score) });
	`, indexPath))
	if err != nil {
		t.Fatalf("hybrid RSF script: %v", err)
	}
	got := value.Export().(map[string]any)
	ids := got["ids"].([]any)
	if len(ids) < 3 {
		t.Fatalf("expected at least 3 hybrid hits, got %#v", ids)
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
