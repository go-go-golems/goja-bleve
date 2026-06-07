const bleve = require("bleve");

const text = bleve.field().text().store(true).includeTermVectors(true).build();
const keyword = bleve.field().keyword().store(true).build();
const doc = bleve.docMapping()
  .dynamic(false)
  .field("text", text)
  .field("source", keyword)
  .build();
const mapping = bleve.mapping().defaultMapping(doc).build();

const idx = bleve.memory().mapping(mapping).build();
idx.index("chunk-1", { text: "privacy preserving retrieval", source: "paper-a" });
idx.index("chunk-2", { text: "vector search over embeddings", source: "paper-b" });
idx.index("chunk-3", { text: "privacy and ranking evaluation", source: "paper-c" });

const req = bleve.search()
  .query(bleve.match("privacy").field("text"))
  .fields(["text", "source"])
  .highlight(["text"])
  .size(10)
  .build();

const result = idx.search(req);
idx.close();
result;
