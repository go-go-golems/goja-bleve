# Tasks

## TODO

### Phase 0: Repository baseline and implementation scaffolding

- [x] Create the package layout under `goja-bleve/pkg/` with files aligned to the design document: `module.go`, `api_types.go`, `api_mapping.go`, `api_index.go`, `api_query.go`, `api_search.go`, `api_batch.go`, `api_knn.go`, `codec.go`, `provider.go`, and test files.
- [x] Add or verify module metadata in `go.mod`: Bleve v2 dependency, `go-go-goja` provider API dependency, and any test-only dependencies needed for goja runtime integration tests.
- [x] Define the top-level module name as `require("bleve")` and ensure the Go package exposes a `modules.NativeModule` implementation consistent with `go-go-goja/modules/common.go` and `goja-text/pkg/extract/module.go`.
- [x] Add a minimal smoke test that starts a `goja.Runtime`, loads the native module, and verifies that `require("bleve")` exposes the expected top-level factories.
- [x] Add README or package documentation describing that vector/KNN support requires `-tags=vectors` and the FAISS linker setup documented in the RAG evaluation system.

**Done when:** `go test ./...` passes for the non-vector build and the module can be required from JavaScript without exposing any vector-only APIs incorrectly.

---

### Phase 1: Runtime object model and Go-backed references

- [x] Implement `moduleRuntime` with references to the `goja.Runtime`, exported module object, open index registry, vector support flag, and cleanup hooks.
- [x] Implement the hidden-reference helper pattern based on geppetto: `attachRef(value, ref)` and `getRef[T](value, expectedType)` using a non-enumerable `__bleve_ref` property.
- [x] Define Go-backed reference structs for each JS wrapper: `indexRef`, `mappingRef`, `docMappingRef`, `fieldMappingRef`, `queryRef`, `searchRequestRef`, `batchRef`, `knnRef`, and result wrappers where needed.
- [x] Define consistent error behavior for wrong wrapper types, missing refs, closed indexes, double-close, and attempts to reuse batches after execution.
- [x] Add tests that pass the wrong JS objects into builder methods and verify that errors are clear and type-specific rather than panics.
- [x] Add tests that prove wrapper objects do not expose implementation refs during normal JS enumeration or JSON serialization.

**Done when:** every public builder method accepts only valid wrapper objects, rejects invalid values with useful errors, and preserves type safety through Go refs rather than dynamic JS maps.

---

### Phase 2: Mapping builder API

- [x] Implement `bleve.mapping()` / `bleve.indexMapping()` builder entrypoint with `.defaultMapping()`, `.typeField()`, `.defaultAnalyzer()`, `.storeDynamic()`, `.indexDynamic()`, and `.build()` methods.
- [x] Implement `bleve.documentMapping()` builder with `.field(name, fieldBuilderOrFieldRef)`, `.dynamic(boolean)`, `.enabled(boolean)`, `.properties(object)`, and `.build()`.
- [x] Implement field builders for text, keyword, numeric, datetime, boolean, geospatial, composite, and disabled fields where supported by Bleve.
- [x] Implement field options: `.store()`, `.index()`, `.includeTermVectors()`, `.includeInAll()`, `.analyzer()`, `.dateFormat()`, `.docValues()`, and equivalent Bleve field-mapping settings.
- [x] Ensure builder `.build()` returns a Go-backed mapping wrapper and not a plain JS object.
- [x] Add tests that create explicit mappings, index documents, and verify expected behavior for text versus keyword fields.
- [x] Add tests that verify dynamic mapping behavior: unmapped fields included/excluded according to mapping settings.
- [x] Document which Bleve mapping options are intentionally not exposed in v1 and why.

**Done when:** a JS script can define an explicit Bleve index mapping with multiple field types and use that mapping to create a working index.

---

### Phase 3: Index lifecycle, document indexing, and BM25 text search

- [x] Implement index creation APIs: `bleve.create(path, mapping)`, `bleve.open(path)`, and `bleve.memory(mapping)` or the final names chosen in the design document.
- [x] Implement `index.close()`, `index.docCount()`, `index.index(id, doc)`, `index.delete(id)`, `index.search(request)`, and `index.batch()`.
- [x] Track open indexes in `moduleRuntime` so module shutdown can close resources and so closed index refs reject later operations.
- [x] Implement document conversion from JS values to Go maps/struct-compatible values while preserving `[]float32` support for later vector fields.
- [x] Implement text query factories: `match`, `matchPhrase`, `term`, `queryString`, `prefix`, `fuzzy`, `regexp`, `wildcard`, `matchAll`, `matchNone`, `conjunction`, `disjunction`, and boolean query composition.
- [x] Implement search request builder methods: `.query(queryRef)`, `.size(n)`, `.from(n)`, `.fields([...])`, `.sort([...])`, `.highlight(...)`, `.explain(boolean)`, and `.build()`.
- [x] Implement result conversion: total hits, max score, took, hit id, score, fields, fragments, locations/explanation where available.
- [x] Add an end-to-end JS integration test that creates an in-memory text index, indexes documents, runs a match query, and asserts ranked hits.
- [x] Add error-path tests for missing index path, invalid mapping object, invalid document id, and invalid query object.

**Done when:** the module can create/open indexes, index ordinary documents, run BM25 text queries, and return stable JS result objects.

---

### Phase 4: Batch API and operational ergonomics

- [x] Implement `batch.index(id, doc)`, `batch.delete(id)`, `batch.size()`, `batch.reset()`, and `batch.execute()` or the final API names chosen in the design doc.
- [x] Ensure batch objects are bound to a single index and cannot be executed against another index.
- [x] Decide whether a batch is reusable after execution; implement and document the chosen behavior.
- [x] Add batch-size and operation-count metadata where Bleve exposes it or where the wrapper can track it safely.
- [x] Add tests for mixed index/delete batches, large-ish batches, duplicate ids, and batch execution errors.
- [x] Add examples showing batch indexing for chunk documents with text and metadata fields.

**Done when:** batch indexing works from JavaScript and all lifecycle constraints are explicit in tests and docs.

---

### Phase 5: Vector field mappings and KNN search

- [x] Add build-tag-aware vector support detection so the module knows whether it was compiled with Bleve vector support.
- [x] Implement vector field builder APIs: `.vector(dims)`, `.vectorBase64(dims)`, `.similarity("cosine" | "dot_product" | "l2_norm")`, and `.optimizedFor("recall" | "memory" | "latency")` as supported by Bleve.
- [x] Make vector builder methods fail with a clear error when the module was not compiled with `-tags=vectors`.
- [x] Implement JS vector conversion in `codec.go`: regular JS arrays, typed arrays where practical, and validation for finite numeric values and exact dimension count.
- [x] Implement KNN request builder or search-request `.knn(field, vector, k, boost)` method using Bleve's `SearchRequest.AddKNN` path.
- [x] Implement KNN operator support if exposed by Bleve: OR/AND semantics for combining multiple KNN clauses.
- [x] Add vector integration tests gated by `//go:build vectors` that create a vector field index, index embeddings, query with a known vector, and assert nearest-neighbor ranking.
- [x] Reuse the validated FAISS command pattern from the RAG evaluation system: `CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"` plus `-tags=vectors`.
- [x] Add tests for dimension mismatch, unsupported similarity value, missing vector field, invalid KNN `k`, and non-finite vector values.

**Done when:** JavaScript can create a vector-enabled Bleve index, index embeddings, and run pure KNN search under `-tags=vectors`.

---

### Phase 6: Hybrid BM25 + vector search and score fusion

- [x] Add search request scoring methods: `.score("rrf" | "rsf" | "none" | default)`, `.scoreRankConstant(n)`, and `.scoreWindowSize(n)` based on Bleve's request parameters.
- [x] Ensure hybrid requests can combine a normal text query with one or more KNN clauses.
- [x] Expose KNN boosts and text-query boosts in a way that matches Bleve's semantics and does not pretend scores are directly comparable before fusion.
- [x] Return score breakdowns when Bleve provides them, and document when `scoreBreakdown` may be empty.
- [x] Add integration tests comparing BM25-only, KNN-only, and hybrid RRF results on the same small corpus.
- [x] Add an example modeled after `cmd/experiments/bleve-knn/main.go`, but written as a goja script using the native module API.
- [x] Document how this differs from the current rag-eval manual RRF implementation over separate BM25 and brute-force vector result sets.

**Done when:** JS scripts can run hybrid text+vector search through one Bleve search request and select RRF/RSF scoring.

---

### Phase 7: Provider registration and host integration

- [x] Implement `pkg/provider.go` with `providerapi.Module` registration so host applications can load the module through the standard go-go-goja provider registry.
- [x] Add an `init()` registration path matching existing modules and ensure there are no import cycles.
- [x] Define configuration options for host services if needed: default index root, allowed paths, vector support policy, and index cleanup behavior.
- [x] Wire the module into an integration host or test provider that starts a goja engine and uses `require("bleve")`.
- [x] Add tests that verify module registration, module name, exports, and TypeScript declaration availability through provider APIs.
- [x] Document how RAG evaluation scripts should import and use the module alongside `fs`, `db`, `geppetto`, and `goja-text`.

**Done when:** a host application can register goja-bleve and JavaScript can `require("bleve")` without manual module wiring.

---

### Phase 8: TypeScript declarations, examples, and API documentation

- [x] Implement `TypeScriptDeclarer` for the full public API: mappings, fields, queries, search requests, indexes, batches, KNN, results, and error-relevant options.
- [x] Keep TypeScript declarations aligned with builder terminal `.build()` methods and Go-backed wrapper object types.
- [x] Add examples for text-only search, explicit mapping, batch indexing, pure KNN search, and hybrid RRF search.
- [x] Add a quickstart document that starts with the smallest working text index and then adds vector search under `-tags=vectors`.
- [x] Add a vector setup note linking to the FAISS how-to in `2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md`.
- [x] Add API snapshots or golden tests for TypeScript declaration output.

**Done when:** an implementer can discover the API from declarations and examples without reading the Go wrapper source first.

---

### Phase 9: Hardening, performance, and production-readiness

- [ ] Add concurrency and lifecycle tests for closing indexes while searches are in flight, repeated module initialization, and cleanup on runtime shutdown.
- [ ] Add benchmarks for JS array to `[]float32` conversion, indexing throughput, batch throughput, and KNN query latency on representative dimensions.
- [ ] Add memory-safety checks for large vectors and large batches so the module fails clearly instead of exhausting memory unexpectedly.
- [ ] Add path-safety policy if the module can open indexes from arbitrary filesystem paths in host applications.
- [ ] Add compatibility tests for non-vector builds to ensure vector APIs report clear unavailable errors rather than failing during module import.
- [ ] Decide whether persistent indexes should be auto-closed by module shutdown or left to explicit user control; document and test the final policy.
- [ ] Add CI jobs for non-vector tests and a separately gated vector/FAISS test job if the runner can install FAISS.

**Done when:** the module has explicit lifecycle, error, performance, and CI coverage sufficient for use in RAG evaluation scripts beyond a prototype.
