// Requires a host binary built with -tags=vectors and linked against FAISS.
const bleve = require("bleve");

const text = bleve.field().text().store(true).build();
const embedding = bleve.field().vector(4).similarity("cosine").optimizedFor("recall").build();
const doc = bleve.docMapping()
  .dynamic(false)
  .field("text", text)
  .field("embedding", embedding)
  .build();
const mapping = bleve.mapping().defaultMapping(doc).build();
const idx = bleve.create(`/tmp/goja-bleve-vector-example-${Date.now()}`).mapping(mapping).build();

idx.index("chunk-1", { text: "alpha", embedding: [1, 0, 0, 0] });
idx.index("chunk-2", { text: "beta", embedding: [0, 1, 0, 0] });
idx.index("chunk-3", { text: "gamma", embedding: [0.9, 0.1, 0, 0] });

const result = idx.search(
  bleve.search()
    .query(bleve.matchNone())
    .knn("embedding", [1, 0, 0, 0], 2, 1.0)
    .fields(["text"])
    .build()
);
idx.close();
result;
