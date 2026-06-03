const bleve = require("bleve");

function buildIndex() {
  const text = bleve.field().text().store(true).build();
  const keyword = bleve.field().keyword().store(true).build();
  const doc = bleve.docMapping()
    .dynamic(false)
    .field("text", text)
    .field("source_id", keyword)
    .build();
  const mapping = bleve.mapping().defaultMapping(doc).build();
  return bleve.memory().mapping(mapping).build();
}

function indexAndSearch(query) {
  const idx = buildIndex();
  const batch = idx.newBatch()
    .index("chunk-1", { text: "privacy screen trees", source_id: "tree-center" })
    .index("chunk-2", { text: "flowering shrubs", source_id: "tree-center" })
    .index("chunk-3", { text: "evergreen privacy hedge", source_id: "tree-center" });
  const sizeBefore = batch.size();
  const operations = batch.operationCount();
  batch.execute();

  const req = bleve.search()
    .query(bleve.match(query || "privacy").field("text"))
    .fields(["text", "source_id"])
    .build();
  const result = idx.search(req);
  const docCount = idx.docCount();
  idx.close();

  return result.hits.map((hit, rank) => ({
    rank: rank + 1,
    id: hit.id,
    score: hit.score,
    text: hit.fields.text,
    total: result.total,
    docCount,
    batchSizeBeforeExecute: sizeBefore,
    batchOperations: operations
  }));
}

__verb__("indexAndSearch", {
  short: "Batch-index documents and run a text search",
  fields: {
    query: { argument: true, default: "privacy", help: "Search query text" }
  }
});

function reuseError() {
  const idx = bleve.memory().build();
  const batch = idx.newBatch().index("chunk-1", { text: "x" });
  batch.execute();
  try {
    batch.index("chunk-2", { text: "y" });
    return { ok: true, error: "" };
  } catch (err) {
    idx.close();
    return { ok: false, error: String(err) };
  }
}

__verb__("reuseError", {
  short: "Verify executed batches cannot be reused"
});
