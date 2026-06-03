# goja-bleve

`goja-bleve` is a native Go module for the go-go-goja runtime. It exposes Bleve full-text and vector search through `require("bleve")`.

The implementation currently includes the Phase 7 core surface:

- JavaScript can `require("bleve")`
- JavaScript wrapper objects carry non-enumerable Go-backed references via `__bleve_ref`
- mapping, in-memory/persistent index, query, search request, batch, vector/KNN, and hybrid scoring builders are implemented
- vector support is detected at build time through the `vectors` build tag and reports clear errors in non-vector builds
- native-module and xgoja provider registration are available for host applications

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

The mapping factories expose terminal `.build()` methods, and search requests can combine ordinary Bleve queries with KNN clauses when the host binary is built with vector support. Hybrid score fusion is available through Bleve's RRF/RSF scoring modes. A minimal TypeScript descriptor is available for provider discovery; Phase 8 will expand it into full API documentation and golden declaration tests.

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

const hybrid = bleve.search()
  .query(bleve.match("privacy").field("text"))
  .knn("embedding", [1, 0, 0, 0], 10, 1.0)
  .score("rrf")
  .scoreRankConstant(60)
  .scoreWindowSize(50)
  .build()
```

Supported mapping helpers normalize common aliases for cosine, dot product, and L2/euclidean similarity, then let Bleve validate the final mapping. `idx.search()` validates KNN vector length against the index mapping before executing the search so JavaScript callers get a clear dimension mismatch error.

Hybrid score fusion uses Bleve's request-level scoring options:

- `.score("rrf")` for Reciprocal Rank Fusion
- `.score("rsf")` for Relative Score Fusion
- `.score("none")` or `.score("default")` for non-fusion modes
- `.scoreRankConstant(n)` for RRF's rank constant
- `.scoreWindowSize(n)` for the fusion window, which must be at least the request size

This differs from the current RAG evaluation service in `2026-05-27--rag-evaluation-system/internal/services/search/hybrid.go`. That service runs BM25 and vector retrieval as two separate service calls, merges candidates by `ChunkID`, and computes manual RRF scores with `1/(rrfK + rank)`. `goja-bleve` instead builds one Bleve `SearchRequest` containing both the text query and one or more KNN clauses, then lets Bleve perform RRF/RSF rescoring inside the search engine. The native Bleve path keeps text, vector, pagination/windowing, KNN boosts, and score-fusion parameters in one request object; the manual rag-eval path still exposes component ranks/scores explicitly through `RetrievalResult.Components`.

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

## Provider and host integration

Host applications can register the module directly with a `goja_nodejs/require.Registry`:

```go
reg := require.NewRegistry()
bleve.Register(reg)
reg.Enable(vm)
```

xgoja hosts should use the provider package:

```go
registry := providerapi.NewRegistry()
err := bleveprovider.Register(registry)
```

The provider package id is `goja-bleve`, and the JavaScript module name is `bleve`, so xgoja specs should mount it as:

```yaml
packages:
  - id: goja-bleve
    import: github.com/go-go-golems/goja-bleve/pkg/xgoja/providers/bleve
runtimes:
  main:
    modules:
      - package: goja-bleve
        name: bleve
        as: bleve
```

There is currently no provider-level configuration schema. Path policy is deliberately a host concern: scripts can call `bleve.create(path)` and `bleve.open(path)`, so applications embedding the module should decide whether paths are unrestricted, sandboxed to a root, or mediated through a future host wrapper. Index lifecycle is explicit: scripts should call `index.close()` when done. The module runtime tracks open indexes internally for cleanup support, but the public host policy remains explicit-close-first.

RAG evaluation scripts can load `bleve` alongside other runtime modules such as `fs`, database modules, `geppetto`, and `goja-text`. Use `fs`/database modules to read source chunks, construct explicit `bleve.mapping()` definitions, batch-index chunks through `idx.newBatch()`, and run BM25/KNN/hybrid requests through one `idx.search(req)` call.

## Examples and TypeScript declarations

- Quickstart: `docs/quickstart.md`
- Text search: `examples/text-search.js`
- Batch indexing: `examples/batch-indexing.js`
- Pure vector KNN: `examples/vector-knn.js`
- Hybrid RRF: `examples/hybrid-rrf.js`

The module implements `modules.TypeScriptDeclarer`. A declaration snapshot is tested in `pkg/testdata/bleve.d.ts.golden`; update it whenever the public JS API changes.

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
./dist/goja-bleve-vectors vector hybrid --output json
```

## Development validation

```bash
go test ./... -count=1
GOWORK=off go test ./... -count=1
GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
  go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1
cd cmd/goja-bleve && GOWORK=off go test ./... -count=1
```
