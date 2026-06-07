package pkg

import (
	"testing"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/dop251/goja"
)

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

func TestBatchExecuteFailureDoesNotMarkBatchExecuted(t *testing.T) {
	vm := goja.New()
	rt := newRuntime(vm)
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		t.Fatalf("new memory index: %v", err)
	}
	batch := idx.NewBatch()
	if err := batch.Index("chunk-1", map[string]any{"text": "privacy"}); err != nil {
		t.Fatalf("index into batch: %v", err)
	}
	if err := idx.Close(); err != nil {
		t.Fatalf("close index before batch execute: %v", err)
	}

	ref := &batchRef{
		refBase: refBase{api: rt, kind: refKindBatch},
		index: &indexRef{
			refBase: refBase{api: rt, kind: refKindIndex},
			name:    "closed-underlying-index",
			index:   idx,
			mapping: mapping,
		},
		batch:     batch,
		operation: 1,
	}
	obj := rt.batchObject(ref)
	execute, ok := goja.AssertFunction(obj.Get("execute"))
	if !ok {
		t.Fatalf("execute is not a function")
	}
	if _, err := execute(obj); err == nil {
		t.Fatalf("expected underlying batch execution failure")
	}
	if ref.executed {
		t.Fatalf("batch was marked executed after failed execute")
	}
	reset, ok := goja.AssertFunction(obj.Get("reset"))
	if !ok {
		t.Fatalf("reset is not a function")
	}
	if _, err := reset(obj); err != nil {
		t.Fatalf("reset after failed execute: %v", err)
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
