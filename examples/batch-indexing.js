const bleve = require("bleve");

const text = bleve.field().text().store(true).build();
const doc = bleve.docMapping().dynamic(false).field("text", text).build();
const mapping = bleve.mapping().defaultMapping(doc).build();
const idx = bleve.memory().mapping(mapping).build();

const batch = idx.newBatch();
for (const chunk of [
  ["chunk-1", "batch indexing improves ingestion throughput"],
  ["chunk-2", "bleve supports full text search"],
  ["chunk-3", "goja scripts can build indexes"]
]) {
  batch.index(chunk[0], { text: chunk[1] });
}
batch.execute();

const result = idx.search(
  bleve.search().query(bleve.match("indexing").field("text")).fields(["text"]).build()
);
idx.close();
result;
