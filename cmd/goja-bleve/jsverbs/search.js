const bleve = require("bleve");

function buildTextIndex() {
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
  idx.index("chunk-3", { text: "evergreen privacy hedge", source_id: "tree-center" });
  return idx;
}

function bm25(query) {
  const idx = buildTextIndex();
  const req = bleve.search()
    .query(bleve.match(query || "privacy").field("text"))
    .fields(["text", "source_id"])
    .size(10)
    .build();
  const result = idx.search(req);
  const count = idx.docCount();
  idx.close();
  return result.hits.map((hit, rank) => ({
    rank: rank + 1,
    id: hit.id,
    score: hit.score,
    text: hit.fields.text,
    sourceID: hit.fields.source_id,
    total: result.total,
    docCount: count
  }));
}

__verb__("bm25", {
  short: "Run a BM25 text-search smoke test over an in-memory Bleve index",
  fields: {
    query: { argument: true, default: "privacy", help: "Search query text" }
  }
});
