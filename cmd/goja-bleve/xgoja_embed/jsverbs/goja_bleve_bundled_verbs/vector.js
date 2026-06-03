const bleve = require("bleve");

function buildVectorIndex(suffix) {
  const text = bleve.field().text().store(true).build();
  const embedding = bleve.field()
    .vector(4)
    .similarity("cosine")
    .optimizedFor("recall")
    .build();
  const doc = bleve.docMapping()
    .dynamic(false)
    .field("text", text)
    .field("embedding", embedding)
    .build();
  const mapping = bleve.mapping().defaultMapping(doc).build();
  const idx = bleve.create(`/tmp/goja-bleve-vector-jsverb-index-${suffix}-${Date.now()}`).mapping(mapping).build();
  idx.index("chunk-1", { text: "alpha", embedding: [1, 0, 0, 0] });
  idx.index("chunk-2", { text: "beta", embedding: [0, 1, 0, 0] });
  idx.index("chunk-3", { text: "gamma", embedding: [0.9, 0.1, 0, 0] });
  return idx;
}

function knn() {
  const idx = buildVectorIndex("knn");
  const req = bleve.search()
    .query(bleve.matchNone())
    .knnOperator("or")
    .knn("embedding", [1, 0, 0, 0], 2, 1.0)
    .fields(["text"])
    .build();
  const result = idx.search(req);
  idx.close();
  return result.hits.map((hit, rank) => ({
    rank: rank + 1,
    id: hit.id,
    score: hit.score,
    text: hit.fields.text,
    total: result.total,
    vectorSupport: bleve.vectorSupport
  }));
}

function hybrid() {
  const idx = buildVectorIndex("hybrid");
  const req = bleve.search()
    .query(bleve.match("alpha").field("text"))
    .knn("embedding", [0.9, 0.1, 0, 0], 2, 1.0)
    .score("rrf")
    .scoreRankConstant(1)
    .scoreWindowSize(4)
    .size(3)
    .fields(["text"])
    .build();
  const result = idx.search(req);
  idx.close();
  return result.hits.map((hit, rank) => ({
    rank: rank + 1,
    id: hit.id,
    score: hit.score,
    text: hit.fields.text,
    total: result.total,
    scoreMode: "rrf",
    vectorSupport: bleve.vectorSupport
  }));
}

__verb__("knn", {
  short: "Run a KNN vector search smoke test (requires -tags=vectors and FAISS)"
});

__verb__("hybrid", {
  short: "Run a BM25+KNN RRF hybrid search smoke test (requires -tags=vectors and FAISS)"
});
