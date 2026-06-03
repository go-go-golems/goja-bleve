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

The factories currently return Go-backed wrapper objects. Later phases will add terminal `.build()` methods, index lifecycle operations, query execution, KNN search, hybrid score fusion, provider integration, and TypeScript declarations.

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

## Development validation

```bash
go test ./... -count=1
GOWORK=off go test ./... -count=1
```
