# Changelog

## 2026-06-02

- Initial workspace created


## 2026-06-02

Step 1: Studied goja module patterns, bleve API, and rag-eval search service. Created ticket and comprehensive design document for goja-bleve module with fluent builder API, vector/KNN support, and hybrid search.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/design-doc/01-goja-bleve-api-design-and-implementation-guide.md — Design document


## 2026-06-02

Step 2-3: Analyzed current vector search (brute-force SQLite). Cloned blevesearch/faiss@fff814d, built libfaiss.so + libfaiss_c.so. Wrote bleve-knn experiment at cmd/experiments/bleve-knn/main.go. Install of FAISS to /usr/local pending sudo.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/experiments/bleve-knn/main.go — Bleve KNN experiment command
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss/ — Cloned blevesearch/faiss at commit fff814d for bleve v2.6.0


## 2026-06-02

Step 4: Patched FAISS test_hamming.cpp heap-id type mismatch using int_maxheap_array_t::TI and verified full make -C build -j8 succeeds.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss/tests/test_hamming.cpp — Portable hamming test heap id type fix


## 2026-06-02

Step 5: Ran bleve-knn experiment successfully with vectors tag using freshly built local FAISS libraries; pure KNN returned chunk-042 as top hit with score 1.0.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/experiments/bleve-knn/main.go — Bleve KNN experiment validation


## 2026-06-02

Validated bleve-knn experiment against system-installed FAISS; run succeeds with CGO_LDFLAGS='-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm'.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/experiments/bleve-knn/main.go — System FAISS vector experiment validation


## 2026-06-02

Added docs/howto-compile-faiss-for-bleve-vectors.md to preserve the FAISS build, patch, install, and Bleve vector run workflow.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md — FAISS build how-to


## 2026-06-02

Added Makefile bleve-knn-experiment target wrapping the validated FAISS/CGO vector-search command and updated the FAISS how-to.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/Makefile — Added bleve-knn-experiment target with FAISS_LIB_DIR override
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md — Updated final command summary to use Makefile target


## 2026-06-02

Expanded goja-bleve implementation roadmap into 10 phases with detailed task checklists, done criteria, and validation expectations.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/tasks.md — Detailed phase/task implementation roadmap for goja-bleve


## 2026-06-02

Recorded Step 7 in the investigation diary: committed Bleve KNN experiment and expanded goja-bleve phase task roadmap.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/reference/01-investigation-diary.md — Step 7 planning and commit diary


## 2026-06-02

Implemented Phase 0 and the core Phase 1 scaffold: module path/dependencies, require('bleve') native loader, runtime state, hidden Go refs, vector support detection, provider wrapper, README, and smoke/reference tests.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/README.md — Phase 0/1 usage and vector setup documentation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/api_types.go — Go-backed wrapper refs and hidden __bleve_ref helpers
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/module.go — Native module loader
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/module_test.go — require('bleve') and hidden-reference tests


## 2026-06-02

Recorded Step 8 diary entry for Phase 0/1 scaffold implementation and validation.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/reference/01-investigation-diary.md — Step 8 implementation diary


## 2026-06-02

Implemented Phase 2 mapping builder API: index/document/field builders, field options, wrong-wrapper errors, and JS-to-real-Bleve mapping integration tests.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/api_mapping.go — Chainable mapping/document/field builders with build() terminal refs
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/api_types.go — Typed mapping refs now carry concrete Bleve mapping types
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/mapping_test.go — Mapping builder integration and wrong-wrapper tests


## 2026-06-02

Implemented Phase 2 mapping builders and added xgoja/jsverb validation harness modeled after goja-text.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/mapping.js — Bundled mapping jsverb smoke tests
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/api_mapping.go — Phase 2 mapping/document/field builder implementation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/mapping_test.go — JS-built mapping integration tests against real Bleve index
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/pkg/xgoja/providers/bleve/bleve.go — xgoja provider wrapper

