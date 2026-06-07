package pkg

import "testing"

func TestJSCanRunCompoundQueriesSortHighlightAndExplain(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).includeTermVectors(true).build();
		const keyword = bleve.field().keyword().store(true).build();
		const doc = bleve.docMapping()
			.dynamic(false)
			.field("text", text)
			.field("source_id", keyword)
			.build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.memory().mapping(mapping).build();
		idx.index("chunk-1", { text: "privacy screen trees", source_id: "b" });
		idx.index("chunk-2", { text: "flowering ornamental shrubs", source_id: "a" });
		idx.index("chunk-3", { text: "evergreen privacy hedge", source_id: "c" });
		const q = bleve.bool()
			.addMust(bleve.match("privacy").field("text"))
			.addMustNot(bleve.term("a").field("source_id"));
		const req = bleve.search()
			.query(q)
			.fields(["text", "source_id"])
			.highlight(["text"])
			.explain(true)
			.sort(["source_id"])
			.size(10)
			.build();
		const result = idx.search(req);
		idx.close();
		({
			total: result.total,
			ids: result.hits.map((h) => h.id).join(","),
			hasFragments: !!result.hits[0].fragments,
			hasExplanation: !!result.hits[0].explanation,
		});
	`)
	if err != nil {
		t.Fatalf("compound query script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["total"] != int64(2) && got["total"] != uint64(2) && got["total"] != float64(2) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["ids"] != "chunk-1,chunk-3" {
		t.Fatalf("ids = %#v", got["ids"])
	}
	if got["hasFragments"] != true {
		t.Fatalf("hasFragments = %#v", got["hasFragments"])
	}
	if got["hasExplanation"] != true {
		t.Fatalf("hasExplanation = %#v", got["hasExplanation"])
	}
}

func TestSearchSizeZeroReturnsCountOnly(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const idx = bleve.memory().build();
		idx.index("chunk-1", { text: "privacy screen trees" });
		idx.index("chunk-2", { text: "privacy hedge" });
		const req = bleve.search()
			.query(bleve.match("privacy"))
			.size(0)
			.build();
		const result = idx.search(req);
		idx.close();
		({ total: result.total, hitCount: result.hits.length });
	`)
	if err != nil {
		t.Fatalf("size zero search script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["total"] != int64(2) && got["total"] != uint64(2) && got["total"] != float64(2) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["hitCount"] != int64(0) && got["hitCount"] != float64(0) {
		t.Fatalf("hitCount = %#v", got["hitCount"])
	}
}

func TestAdditionalQueryFactoriesAreExportedAndUsable(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const names = ["matchPhrase", "prefix", "fuzzy", "regexp", "wildcard", "bool", "conj", "disj"];
		({ ok: names.every((name) => typeof bleve[name] === "function") });
	`)
	if err != nil {
		t.Fatalf("query factory export script: %v", err)
	}
	if got := value.Export().(map[string]any); got["ok"] != true {
		t.Fatalf("factory exports = %#v", got)
	}
}
