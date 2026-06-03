package pkg

import "testing"

func TestJSBatchIndexDeleteAndExecute(t *testing.T) {
	vm := newBleveTestVM()
	value, err := vm.RunString(`
		const bleve = require("bleve");
		const text = bleve.field().text().store(true).build();
		const doc = bleve.docMapping().dynamic(false).field("text", text).build();
		const mapping = bleve.mapping().defaultMapping(doc).build();
		const idx = bleve.memory().mapping(mapping).build();
		const batch = idx.newBatch()
			.index("chunk-1", { text: "privacy screen trees" })
			.index("chunk-2", { text: "flowering shrubs" })
			.delete("missing-doc");
		const sizeBefore = batch.size();
		const opsBefore = batch.operationCount();
		batch.execute();
		const req = bleve.search().query(bleve.match("privacy").field("text")).fields(["text"]).build();
		const result = idx.search(req);
		const count = idx.docCount();
		idx.close();
		({ sizeBefore, opsBefore, count, total: result.total, firstID: result.hits[0].id });
	`)
	if err != nil {
		t.Fatalf("batch script: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["opsBefore"] != int64(3) && got["opsBefore"] != float64(3) {
		t.Fatalf("opsBefore = %#v", got["opsBefore"])
	}
	if got["count"] != int64(2) && got["count"] != uint64(2) && got["count"] != float64(2) {
		t.Fatalf("count = %#v", got["count"])
	}
	if got["total"] != int64(1) && got["total"] != uint64(1) && got["total"] != float64(1) {
		t.Fatalf("total = %#v", got["total"])
	}
	if got["firstID"] != "chunk-1" {
		t.Fatalf("firstID = %#v", got["firstID"])
	}
}

func TestBatchCannotBeReusedAfterExecute(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		const idx = bleve.memory().build();
		const batch = idx.newBatch().index("chunk-1", { text: "x" });
		batch.execute();
		batch.index("chunk-2", { text: "y" });
	`)
	if err == nil {
		t.Fatalf("expected batch reuse to fail")
	}
	if got := err.Error(); !containsAll(got, []string{"batch has already been executed"}) {
		t.Fatalf("batch reuse error = %q", got)
	}
}

func TestBatchFailsAfterIndexClose(t *testing.T) {
	vm := newBleveTestVM()
	_, err := vm.RunString(`
		const bleve = require("bleve");
		const idx = bleve.memory().build();
		const batch = idx.newBatch();
		idx.close();
		batch.index("chunk-1", { text: "x" });
	`)
	if err == nil {
		t.Fatalf("expected batch operation after index close to fail")
	}
	if got := err.Error(); !containsAll(got, []string{"batch index is closed"}) {
		t.Fatalf("closed-index batch error = %q", got)
	}
}
