const bleve = require("bleve");

function factories() {
  return {
    version: bleve.version,
    vectorSupport: bleve.vectorSupport,
    exports: [
      "mapping",
      "indexMapping",
      "docMapping",
      "documentMapping",
      "field",
      "search",
      "searchRequest",
      "create",
      "open",
      "memory",
      "matchAll",
      "matchNone"
    ].map((name) => ({ name, kind: typeof bleve[name] }))
  };
}

__verb__("factories", {
  short: "Inspect goja-bleve factory exports"
});

function buildBasic() {
  const text = bleve.field()
    .text()
    .store(true)
    .includeTermVectors(true)
    .build();

  const keyword = bleve.field()
    .keyword()
    .store(true)
    .build();

  const number = bleve.field()
    .number()
    .store(true)
    .build();

  const doc = bleve.docMapping()
    .dynamic(false)
    .field("text", text)
    .field("source_id", keyword)
    .field("chunk_index", number)
    .build();

  const mapping = bleve.mapping()
    .defaultMapping(doc)
    .defaultAnalyzer("standard")
    .defaultField("text")
    .storeDynamic(false)
    .indexDynamic(false)
    .build();

  return {
    mappingType: mapping.type,
    docType: doc.type,
    textFieldType: text.type,
    hiddenRefEnumerable: Object.keys(mapping).includes("__bleve_ref"),
    json: JSON.stringify(mapping)
  };
}

__verb__("buildBasic", {
  short: "Build a basic text/keyword/number Bleve mapping"
});

function wrongWrapperError() {
  try {
    const fieldBuilder = bleve.field().text();
    bleve.docMapping().field("text", fieldBuilder);
    return { ok: true, error: "" };
  } catch (err) {
    return { ok: false, error: String(err) };
  }
}

__verb__("wrongWrapperError", {
  short: "Verify mapping builders reject unbuilt wrapper objects"
});
