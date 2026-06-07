# goja-bleve Quickstart

`goja-bleve` exposes Bleve as `require("bleve")` inside goja/xgoja runtimes. The API is builder-oriented: create mappings, call terminal `.build()`, create an index, index documents, build a search request, and call `idx.search(req)`.

## Smallest text index

```javascript
const bleve = require("bleve");

const text = bleve.field().text().store(true).build();
const doc = bleve.docMapping().dynamic(false).field("text", text).build();
const mapping = bleve.mapping().defaultMapping(doc).build();
const idx = bleve.memory().mapping(mapping).build();

idx.index("chunk-1", { text: "privacy preserving retrieval" });
idx.index("chunk-2", { text: "vector search over embeddings" });

const req = bleve.search()
  .query(bleve.match("privacy").field("text"))
  .fields(["text"])
  .build();

const result = idx.search(req);
idx.close();
result;
```

## Batch indexing

Use `idx.newBatch()` when indexing many chunks. Batches are single-use after `.execute()`.

```javascript
const batch = idx.newBatch();
batch.index("chunk-1", { text: "first chunk" });
batch.index("chunk-2", { text: "second chunk" });
batch.execute();
```

## Vector KNN

Vector search requires a host binary built with Bleve's vector tag and FAISS linked:

```bash
make test-vectors
```

That target wraps the required `-tags=vectors`, `CGO_LDFLAGS`, and runtime rpath settings for package-level tests. To validate the generated xgoja vector host and embedded JavaScript smoke verbs, run:

```bash
make xgoja-smoke-vectors
```

For machine setup, troubleshooting, and xgoja vector builds, see [FAISS + goja-bleve + xgoja Playbook](faiss-xgoja-playbook.md).

```javascript
const embedding = bleve.field().vector(4).similarity("cosine").optimizedFor("recall").build();
const doc = bleve.docMapping().dynamic(false).field("embedding", embedding).build();
const mapping = bleve.mapping().defaultMapping(doc).build();
const idx = bleve.create("/tmp/my-vector-index").mapping(mapping).build();

idx.index("a", { embedding: [1, 0, 0, 0] });
idx.index("b", { embedding: [0, 1, 0, 0] });

const result = idx.search(
  bleve.search()
    .query(bleve.matchNone())
    .knn("embedding", [1, 0, 0, 0], 1)
    .build()
);
idx.close();
```

## Hybrid RRF

Hybrid search combines a normal text query and one or more KNN clauses in one Bleve request.

```javascript
const result = idx.search(
  bleve.search()
    .query(bleve.match("privacy").field("text"))
    .knn("embedding", [1, 0, 0, 0], 10, 1.0)
    .score("rrf")
    .scoreRankConstant(60)
    .scoreWindowSize(50)
    .build()
);
```

Use `.score("rsf")` for Relative Score Fusion. `scoreWindowSize` must be at least the request size.

## Examples

See:

- `examples/text-search.js`
- `examples/batch-indexing.js`
- `examples/vector-knn.js`
- `examples/hybrid-rrf.js`
