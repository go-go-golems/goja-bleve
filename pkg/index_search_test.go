package pkg

import "testing"

func TestJSCanIndexAndSearchInMemoryBleveIndex(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).build();
		const keyword = bleve.field().keyword().store(true).build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("text", text)
			.field("source_id", keyword)
			.build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.memory().mapping(mapping).build();
		idx.index("chunk-1", { text: "privacy screen trees", source_id: "tree-center" });
		idx.index("chunk-2", { text: "flowering ornamental shrubs", source_id: "tree-center" });
		const req = bleve.search()
			.query(bleve.match("privacy").field("text"))
			.fields(["text", "source_id"])
			.size(5)
			.build();
		const result = idx.search(req);
		const count = idx.docCount();
		idx.close();
		({
			count,
			total: result.total,
			firstID: result.hits[0].id,
			firstText: result.hits[0].fields.text,
		});
	`)
	if err != nil {
		t.Fatalf("index/search script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["count"] != int64(2) && got["count"] != uint64(2) && got["count"] != float64(2) {
		t.Fatalf("doc count = %#v", got["count"])
	}
	if got["total"] != int64(1) && got["total"] != uint64(1) && got["total"] != float64(1) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["firstID"] != "chunk-1" {
		t.Fatalf("firstID = %#v", got["firstID"])
	}
}

func TestIndexRejectsSearchBuilderInsteadOfBuiltSearchRequest(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		const idx = bleve.memory().build();
		idx.search(bleve.search().query(bleve.matchAll()));
	`)
	if err == nil {
		t.Fatalf("expected search builder wrapper to be rejected")
	}
	if got := err.Error(); !containsAll(got, []string{"search request is not built"}) {
		t.Fatalf("search builder error = %q", got)
	}
}
