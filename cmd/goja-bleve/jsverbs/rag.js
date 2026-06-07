const bleve = require("bleve");
const geppetto = require("geppetto");

const DEFAULT_DOCS = [
  { id: "chunk-1", text: "privacy preserving retrieval for evaluation systems", source: "demo" },
  { id: "chunk-2", text: "flowering shrubs and ornamental trees", source: "demo" },
  { id: "chunk-3", text: "vector search and hybrid reciprocal rank fusion", source: "demo" }
];

function _parseDocs(docsJson) {
  if (!docsJson) return DEFAULT_DOCS;
  const docs = JSON.parse(docsJson);
  if (!Array.isArray(docs)) throw new Error("docsJson must encode an array");
  return docs.map((doc, i) => ({
    id: String(doc.id || `doc-${i + 1}`),
    text: String(doc.text || ""),
    source: String(doc.source || "input")
  })).filter(doc => doc.text.length > 0);
}

function _resolveEmbedder(profilePath, profile) {
  if (!profilePath) throw new Error("profilePath is required; pass a Geppetto profile registry YAML path");
  const settings = geppetto.inferenceProfiles.load(profilePath).resolve(profile || "assistant");
  const embedder = geppetto.embeddings(settings);
  const model = embedder.model();
  if (!model || !model.dimensions) throw new Error("resolved embedding model did not report dimensions");
  return { embedder, model };
}

function _buildIndex(indexPath, dims) {
  const text = bleve.field().text().store(true).includeTermVectors(true).build();
  const source = bleve.field().keyword().store(true).build();
  const embedding = bleve.field().vector(dims).similarity("cosine").optimizedFor("recall").build();
  const doc = bleve.docMapping()
    .dynamic(false)
    .field("text", text)
    .field("source", source)
    .field("embedding", embedding)
    .build();
  const mapping = bleve.mapping().defaultMapping(doc).build();
  return bleve.create(indexPath || `/tmp/goja-bleve-geppetto-rag-${Date.now()}`).mapping(mapping).build();
}

function plan() {
  return {
    ok: true,
    modules: {
      bleve: { vectorSupport: bleve.vectorSupport, search: typeof bleve.search, field: typeof bleve.field },
      geppetto: { embeddings: typeof geppetto.embeddings, inferenceProfiles: typeof geppetto.inferenceProfiles }
    },
    command: "goja-bleve-vectors rag index-query --profile-path ./profiles.yaml --embedding-profile assistant privacy",
    note: "index-query calls geppetto.embeddings(settings).embed(...), so it needs a real embedding-capable Geppetto profile."
  };
}

function indexQuery(profilePath, query, docsJson, embeddingProfile, indexPath, mode, limit) {
  if (!bleve.vectorSupport) throw new Error("rag index-query requires a vector-enabled binary built with -tags=vectors");
  const docs = _parseDocs(docsJson);
  if (docs.length === 0) throw new Error("at least one document with text is required");
  const { embedder, model } = _resolveEmbedder(profilePath, embeddingProfile || "assistant");
  const idx = _buildIndex(indexPath, model.dimensions);
  const batch = idx.newBatch();
  for (const doc of docs) {
    const vector = embedder.embed(doc.text);
    batch.index(doc.id, { text: doc.text, source: doc.source, embedding: vector });
  }
  batch.execute();

  const queryText = query || "privacy";
  const queryVector = embedder.embed(queryText);
  const requestBuilder = bleve.search()
    .query((mode || "hybrid") === "knn" ? bleve.matchNone() : bleve.match(queryText).field("text"))
    .knn("embedding", queryVector, Number(limit || 5), 1.0)
    .fields(["text", "source"])
    .size(Number(limit || 5));
  if ((mode || "hybrid") !== "knn") {
    requestBuilder.score("rrf").scoreRankConstant(60).scoreWindowSize(Math.max(Number(limit || 5), docs.length));
  }
  const result = idx.search(requestBuilder.build());
  const count = idx.docCount();
  idx.close();
  return {
    ok: true,
    query: queryText,
    mode: mode || "hybrid",
    model,
    docCount: count,
    hits: result.hits.map((hit, rank) => ({
      rank: rank + 1,
      id: hit.id,
      score: hit.score,
      text: hit.fields.text,
      source: hit.fields.source,
      total: result.total
    }))
  };
}

__verb__("plan", {
  short: "Show the geppetto+bleve RAG indexing/querying tool wiring without calling an embedding provider"
});

__verb__("indexQuery", {
  short: "Embed documents with Geppetto, index them with Bleve vectors, and query with KNN or hybrid RRF",
  fields: {
    profilePath: { help: "Geppetto profile registry YAML path with embeddings settings" },
    embeddingProfile: { default: "assistant", help: "Geppetto profile slug to resolve for embeddings" },
    query: { argument: true, default: "privacy", help: "Query text to embed and search" },
    docsJson: { help: "JSON array of {id,text,source} documents; defaults to a small demo corpus" },
    indexPath: { help: "Optional persistent Bleve index path; defaults to a /tmp path" },
    mode: { default: "hybrid", help: "hybrid or knn" },
    limit: { default: "5", help: "Number of nearest/fused hits" }
  }
});
