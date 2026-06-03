# goja-bleve

`goja-bleve` is a native Go module for the go-go-goja runtime. It exposes Bleve full-text and vector search through `require("bleve")`.

The implementation is currently in Phase 0/1:

- the module name and native-loader registration are in place
- JavaScript can `require("bleve")`
- the module exposes scaffolded builder factories
- JavaScript wrapper objects carry non-enumerable Go-backed references via `__bleve_ref`
- vector support is detected at build time through the `vectors` build tag

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

The mapping factories now expose terminal `.build()` methods. Later phases will add index lifecycle operations, query execution, KNN search, hybrid score fusion, provider integration, and TypeScript declarations.

## Mapping API scope in the current phase

Phase 2 exposes the core mapping surface needed for text-first indexes:

- index mappings: `bleve.mapping()` / `bleve.indexMapping()`
- document mappings: `bleve.docMapping()` / `bleve.documentMapping()`
- field mappings: `text`, `keyword`, `number`, `datetime`, `boolean`, `geoPoint`, `geoShape`, `ip`, and `disabled`
- common field options: `name`, `analyzer`, `store`, `index`, `docValues`, `includeTermVectors`, `includeInAll`, and `dateFormat`

The first implementation intentionally does not expose custom analyzers, custom token filters, custom tokenizers, custom date parsers, synonym sources, scoring-model configuration, or vector field options. Custom analysis is a larger Bleve registry concern and needs its own validation/error model. Vector mappings are deferred to the vector/KNN phase because they require `-tags=vectors` and FAISS setup.

## Batch lifecycle

`index.newBatch()` returns a batch bound to exactly one open index. A batch supports `.index(id, doc)`, `.delete(id)`, `.size()`, `.operationCount()`, `.reset()`, and `.execute()`.

Batches are single-use after execution. Once `.execute()` succeeds, later mutation or reset attempts throw `bleve: batch has already been executed`. This keeps lifecycle behavior explicit and avoids ambiguity about whether a batch should retain or clear its queued operations after being submitted to Bleve.

## Vector / KNN support

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
```

## Development validation

```bash
go test ./... -count=1
GOWORK=off go test ./... -count=1
cd cmd/goja-bleve && GOWORK=off go test ./... -count=1
```
