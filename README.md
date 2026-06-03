# goja-bleve

`goja-bleve` is a native Go module for the go-go-goja runtime. It exposes Bleve full-text and vector search through `require("bleve")`.

The implementation currently includes the Phase 5 core surface:

- JavaScript can `require("bleve")`
- JavaScript wrapper objects carry non-enumerable Go-backed references via `__bleve_ref`
- mapping, in-memory/persistent index, query, search request, batch, and vector/KNN builders are implemented
- vector support is detected at build time through the `vectors` build tag and reports clear errors in non-vector builds

## Minimal JavaScript shape

```javascript
const bleve = require("bleve")

const mapping = bleve.mapping()
const docs = bleve.docMapping()
const field = bleve.field()
const indexBuilder = bleve.memory()
const query = bleve.matchAll()
const request = bleve.search()
```

The mapping factories expose terminal `.build()` methods, and search requests can now combine ordinary Bleve queries with KNN clauses when the host binary is built with vector support. Later phases will add hybrid score fusion convenience APIs and TypeScript declarations.

## Mapping API scope in the current phase

Phase 2 exposes the core mapping surface needed for text-first indexes:

- index mappings: `bleve.mapping()` / `bleve.indexMapping()`
- document mappings: `bleve.docMapping()` / `bleve.documentMapping()`
- field mappings: `text`, `keyword`, `number`, `datetime`, `boolean`, `geoPoint`, `geoShape`, `ip`, and `disabled`
- common field options: `name`, `analyzer`, `store`, `index`, `docValues`, `includeTermVectors`, `includeInAll`, and `dateFormat`

The first implementation intentionally does not expose custom analyzers, custom token filters, custom tokenizers, custom date parsers, synonym sources, or scoring-model configuration. Custom analysis is a larger Bleve registry concern and needs its own validation/error model. Vector fields are available through `field().vector(dims)` and `field().vectorBase64(dims)` when the binary is compiled with `-tags=vectors`.

## Batch lifecycle

`index.newBatch()` returns a batch bound to exactly one open index. A batch supports `.index(id, doc)`, `.delete(id)`, `.size()`, `.operationCount()`, `.reset()`, and `.execute()`.

Batches are single-use after execution. Once `.execute()` succeeds, later mutation or reset attempts throw `bleve: batch has already been executed`. This keeps lifecycle behavior explicit and avoids ambiguity about whether a batch should retain or clear its queued operations after being submitted to Bleve.

## Vector / KNN support

Phase 5 exposes vector field mappings and KNN search:

```javascript
const embedding = bleve.field()
  .vector(4)
  .similarity("cosine")
  .optimizedFor("recall")
  .build()

const request = bleve.search()
  .query(bleve.matchNone())
  .knn("embedding", [1, 0, 0, 0], 2, 1.0)
  .build()
```

Supported mapping helpers normalize common aliases for cosine, dot product, and L2/euclidean similarity, then let Bleve validate the final mapping. `idx.search()` validates KNN vector length against the index mapping before executing the search so JavaScript callers get a clear dimension mismatch error.

Bleve vector search requires the host Go binary to be compiled with:

```bash
-tags=vectors
```

It also requires FAISS to be installed and linked for CGO. The validated linker flags from the RAG evaluation system are:

```bash
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
```

The detailed FAISS build instructions live in:

```text
/home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md
```

## xgoja/jsverb validation

The generated `cmd/goja-bleve` binary embeds small JavaScript verb smoke tests. Use these while building each API phase so the module is validated through the same generated xgoja runtime shape that downstream users will exercise.

```bash
cd cmd/goja-bleve
GOWORK=off go generate ./...
./dist/goja-bleve mapping factories --output json
./dist/goja-bleve mapping build-basic --output json
./dist/goja-bleve mapping wrong-wrapper-error --output json
./dist/goja-bleve search bm25 privacy --output json
./dist/goja-bleve batch index-and-search privacy --output json
```

For vector jsverbs, build the vector-specific xgoja spec:

```bash
cd cmd/goja-bleve
GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
  go run github.com/go-go-golems/go-go-goja/cmd/xgoja@v0.7.4 build \
  -f xgoja-vectors.yaml \
  --work-dir /tmp/goja-bleve-vector-work \
  --keep-work \
  --xgoja-version v0.7.4
./dist/goja-bleve-vectors vector knn --output json
```

## Development validation

```bash
go test ./... -count=1
GOWORK=off go test ./... -count=1
GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
  go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1
cd cmd/goja-bleve && GOWORK=off go test ./... -count=1
```
