---
Title: Investigation diary
Ticket: RAGEVAL-GOJA-RAG-STRATEGIES
Status: active
Topics:
    - goja
    - bleve
    - search
    - embeddings
    - rag
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: 2026-05-27--rag-evaluation-system/cmd/experiments/bleve-knn/main.go
      Note: Standalone bleve vector/KNN experiment now runs end-to-end with local FAISS libraries
    - Path: faiss/tests/test_hamming.cpp
      Note: Patched hamming test ids to use faiss::int_maxheap_array_t::TI so full FAISS build succeeds
    - Path: goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/reference/01-investigation-diary.md
      Note: Recorded FAISS full build patch and validation
ExternalSources: []
Summary: Chronological investigation diary for goja-bleve design, FAISS build, and vector search experiment.
LastUpdated: 2026-06-02T22:30:00-04:00
WhatFor: Track design decisions, FAISS build steps, and experiment progress
WhenToUse: When resuming work on the goja-bleve module or FAISS/vector experiments
---



# Diary

## Goal

Capture the design investigation, FAISS build process, and vector search experiment for goja-bleve: a native Go module that exposes bleve full-text and vector search to the goja JavaScript runtime.

## Step 1: Study existing goja module patterns and bleve API

### Prompt Context

**User prompt (verbatim):** "We're working on RAGEVAL-GOJA-RAG-STRATEGIES. We're leaving the idea of a custom goja module aside for now, and instead will be using the existing modules: fs/yaml/etc... (standard functionality), db (for db access), geppetto (for llm and embeddings), goja-text (for chunking), express (for serving http page custom webui). However we still need a set of bleve bindings for accessing bleve. We want to have a similar fluent API with the builder pattern, where most objects are actually go side objects to have proper type safety at runtime instead of dealing with hashmaps and js objects. Study bleve and its APIs and design an elegant goja API. Create a new docmgr ticket for implement the goja-bleve, and Create a detailed analysis / design / implementation guide..."

**Assistant interpretation:** Design a goja-bleve module with fluent builder API backed by Go structs, create a docmgr ticket, write a comprehensive design document.

**Inferred user intent:** Enable RAG pipeline composition in JavaScript by exposing bleve search primitives (BM25, vector/KNN, hybrid) through the goja module system.

### What I did

- Studied goja module ecosystem: NativeModule, Loader, SetExport, TypeScriptDeclarer, providerapi
- Studied geppetto module: moduleRuntime, attachRef/getRef, engine/schema/inference settings builders
- Studied goja-text extract module: simpler ExtractOptionsBuilder pattern
- Studied bleve Go API: Index, SearchRequest, KNNRequest, mapping hierarchy, query factories, rescorer
- Studied rag-eval search service: bm25.go, vector.go, hybrid.go
- Created docmgr ticket RAGEVAL-GOJA-RAG-STRATEGIES
- Fetched bleve docs via defuddle into sources/ folder
- Wrote 56KB design document

### Why

Thorough understanding of all three layers needed before designing the JS-facing API.

### What worked

- geppetto provides excellent reference for attachRef/getRef and builder conventions
- rag-eval bm25.go shows exactly how bleve is used in practice
- Bleve's KNN + rescorer maps well to the builder pattern

### What didn't work

- Several blevesearch.com doc URLs returned 404
- Had to rely on bleve source code directly

### What I learned

- rag-eval does NOT use bleve's KNN — vector search is brute-force cosine similarity in Go over SQLite-stored embeddings
- Hybrid fusion is also manual RRF, not bleve's built-in rescoring
- Bleve vector support requires `vectors` build tag + FAISS shared library
- Vector fields auto-normalize for cosine similarity during indexing

### What was tricky to build

- Mapping bleve's query breadth to a coherent JS API. Solution: flat namespace for query factories, builders for complex types.

### What warrants a second pair of eyes

- KNN search flow: is `.score("rrf")` + `.knn(...)` sufficient to trigger bleve's two-phase search?

### What should be done in the future

- Implement Phase 1 (core index + BM25) to validate architecture
- Profile vector conversion overhead (JS Float32Array → Go []float32)

### Code review instructions

- Read design doc at `design-doc/01-goja-bleve-api-design-and-implementation-guide.md`
- Cross-reference "Key File Reference Map" against actual source files

### Technical details

- Ticket: `goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/`
- Uploaded to reMarkable at `/ai/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES`

---

## Step 2: Analyze current vector search implementation

### Prompt Context

**User prompt (verbatim):** "how do we currently do vector search in 2026-05-27--rag-evaluation-system/ ?"

**Assistant interpretation:** Explain the current vector search architecture.

**Inferred user intent:** Understand the gap between current vector search and what bleve KNN would provide.

### What I did

- Read `internal/services/search/vector.go`, `embedding/similarity.go`, `embedding/provider.go`, `db/search_queries.go`
- Documented the full flow

### What I learned

- **Brute-force cosine similarity** — no vector index at all
- Embeddings stored as binary blobs in SQLite (`chunk_embeddings.embedding` → `[]byte` via little-endian `[]float32`)
- `QueryVector()` flow: embed query via geppetto → load up to 500 candidate rows from SQLite → decode each → compute cosine similarity one-by-one in Go → sort → truncate
- Hybrid search is manual RRF in Go: run BM25 + vector independently → merge via `1/(k+rank)` per retriever → sort
- **bleve KNN would be a genuine upgrade**: FAISS-backed ANN indexes (flat, IVF, BIVF) for scalable search + built-in rescoring with score breakdowns

### What was tricky to build

- Nothing tricky — just reading code

### What warrants a second pair of eyes

- The manual RRF in hybrid.go vs bleve's built-in rescoring — they may produce different results because bleve's rescoring handles the two-phase search properly (pre-search KNN + FTS merge + score normalization)

### What should be done in the future

- Replace brute-force vector search with bleve KNN once the module is ready

---

## Step 3: Set up FAISS + bleve vectors build for vector search experiment

### Prompt Context

**User prompt (verbatim):** "look if you can get a test command in ./2026-05-27--rag-evaluation-system/cmd/experiments/... to run vector search using bleve, which means we'll need to properly set up and compile and set the right tags to have faiss support in bleve (look at doc online if necessary, remember to save useful resources in the sources/ folder of the ticket)"

**Assistant interpretation:** Create an experiment command that uses bleve's KNN vector search, which requires building FAISS and compiling with the `vectors` build tag.

**Inferred user intent:** Validate that bleve KNN works end-to-end with our data before investing in the full goja-bleve module.

### What I did

- Checked that FAISS shared library (`libfaiss_c.so`) was already installed at `/usr/local/lib`
- Discovered the installed FAISS is missing `IndexBinary_c_ex.h` header required by `go-faiss v1.1.3`
- Fetched bleve vectors doc and go-faiss README via defuddle → saved to sources/
- Cloned `blevesearch/faiss` at commit `fff814d` (required by bleve v2.6.0) into the workspace
- Built FAISS from source: `cmake -B build -DFAISS_ENABLE_GPU=OFF -DFAISS_ENABLE_C_API=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_CXX_FLAGS="-I$PWD" . && make -C build faiss faiss_c`
- Attempted CGO compilation with local include/library paths — failed due to C API header relative include paths (`"impl/AuxIndexStructures_c.h"`, etc.)
- Tried creating a merged include directory tree — still missing transitive includes
- Rebuilt with `CMAKE_INSTALL_PREFIX=/usr/local` — both `libfaiss.so` and `libfaiss_c.so` built successfully
- Install step requires sudo — handed off to user

### Why

Need FAISS properly installed at system level so that `go test -tags=vectors` works, which is a prerequisite for any bleve KNN experiment.

### What worked

- FAISS builds cleanly from blevesearch/faiss@fff814d with the extra `-DCMAKE_CXX_FLAGS="-I$PWD"` flag
- Both `libfaiss.so` and `libfaiss_c.so` produced in `build/`
- Bleve v2.6.0 and go-faiss v1.1.3 are already in the rag-eval go.mod

### What didn't work

- Old system FAISS (`/usr/local/lib/libfaiss_c.so` from May 2025) lacks `IndexBinary_c_ex.h` header → CGO compilation fails
- CGO_CFLAGS/CGO_LDFLAGS approach with local paths failed because the C API headers use relative includes (`"impl/X.h"`) that need `make install` to resolve properly
- `sudo` not available from this session — install step deferred to user

### What I learned

- Bleve v2.6.0 requires blevesearch/faiss at commit `fff814d` (modified v1.13.2)
- The `vectors` build tag activates all vector/KNN code in bleve
- FAISS is a C++ library with a C API wrapper; `go-faiss` uses CGO to call it
- The include path problem is solved by `make install` which copies headers into the correct `/usr/local/include/faiss/` tree
- Sanity check command: `go test -ldflags "-r /usr/local/lib" ./... -tags=vectors`

### What was tricky to build

- The FAISS header layout: the blevesearch fork restructures headers differently from upstream. The C API headers include `"impl/X.h"` which only resolves if installed properly. Manual CGO flag fiddling couldn't replicate the full install layout.

### What warrants a second pair of eyes

- After install, verify the sanity check passes: `cd bleve && go test -ldflags "-r /usr/local/lib" -tags=vectors -run TestKNN .`

### What should be done in the future

- After FAISS install: write the experiment command
- Create `cmd/experiments/bleve-knn/main.go` that builds a bleve index with vector fields, indexes documents, runs KNN search

### Code review instructions

- Check FAISS build at `./faiss/build/` — both `.so` files should exist
- After install, verify: `ls /usr/local/include/faiss/c_api/IndexBinary_c_ex.h` and `ls /usr/local/lib/libfaiss_c.so`

### Technical details

- FAISS source: `./faiss/` (cloned from blevesearch/faiss@fff814d)
- Build command: `cmake -B build -DFAISS_ENABLE_GPU=OFF -DFAISS_ENABLE_C_API=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_CXX_FLAGS="-I$PWD" . && make -C build faiss faiss_c`
- Install command (requires sudo): `sudo make -C build install && sudo ldconfig`
- Bleve sanity check: `cd bleve && GOWORK=off go test -ldflags "-r /usr/local/lib" -tags=vectors -run TestKNN .`
- Sources saved: `bleve-vectors-doc.md`, `go-faiss-readme.md`, `go-faiss-pkg-doc.md` in ticket sources/

---

## Step 4: Patch FAISS hamming test type mismatch and verify full build

The full FAISS build was failing after the library targets had already built because the test target compiled `tests/test_hamming.cpp`. The failure was a platform-specific C++ type mismatch: `faiss::int_maxheap_array_t` expects its `ids` buffer to use the heap's own `TI` pointer type, while the test hard-coded `std::vector<long long>`.

I patched the test to derive the id type directly from `faiss::int_maxheap_array_t::TI`, then reran `make -C build -j$(nproc)`. The initial rerun surfaced a stale root-owned dependency file from an earlier sudo attempt; after removing that file, the full build completed successfully, including `faiss_test`.

### Prompt Context

**User prompt (verbatim):** "try make -C build

----

Fix the compile error, see : Yes — this is the same class of FAISS build failure as the old `test_hamming.cpp` heap-id type mismatch. There is even an upstream issue for `test_hamming.cpp` with the same `ids_gen.data(), dist_gen.data()` initializer, though in that report the mismatch was reversed: `long int*` passed where `long long int*` was expected. ([GitHub][1])

What is happening:

`faiss::int_maxheap_array_t` is defined as a heap array over `CMax<int, int64_t>`, and `HeapArray`’s `ids` pointer type is `C::TI`. ([Faiss][2]) On your platform/build, that resolves to `long int*`, but the test has `ids_gen` as `long long int`. Those are both 64-bit on many platforms, but they are distinct C++ types, so `std::vector<long long>::data()` cannot initialize a `long*`.

Patch it by using the heap’s own id type instead of spelling `long`, `long long`, or even assuming what `int64_t` aliases to:

```cpp
using Heap = faiss::int_maxheap_array_t;
using HeapId = Heap::TI;

std::vector<HeapId> ids_gen(na * k);
std::vector<int> dist_gen(na * k);

Heap res{na, k, ids_gen.data(), dist_gen.data()};
```

Do the same for the other id buffers in that test, probably `ids_ham_knn` as well:

```cpp
std::vector<HeapId> ids_ham_knn(na * k);
```

Avoid fixing this with a cast. The safe fix is to make the vector element type exactly match `faiss::int_maxheap_array_t::TI`. This keeps the test portable across platforms where `int64_t` is `long` versus `long long`.

[1]: https://github.com/facebookresearch/faiss/issues/4267?utm_source=chatgpt.com "test_hamming.cpp:277:36: error: cannot convert 'long int*' to 'faiss"
[2]: https://faiss.ai/cpp_api/file/Heap_8h.html?utm_source=chatgpt.com "File Heap.h""

**Assistant interpretation:** Apply the portable heap id type fix to FAISS's hamming test and verify the full `make -C build` target succeeds.

**Inferred user intent:** Avoid relying only on partial `faiss`/`faiss_c` targets; make the checked-out FAISS tree build cleanly end-to-end so install and future rebuilds are less surprising.

### What I did

- Edited `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss/tests/test_hamming.cpp`
- Added aliases:
  - `using HammingHeap = faiss::int_maxheap_array_t;`
  - `using HammingHeapId = HammingHeap::TI;`
- Replaced hard-coded `std::vector<long long>` hamming id buffers with `std::vector<HammingHeapId>`
- Replaced `std::set<long>` expected id storage with `std::set<HammingHeapId>`
- Ran `make -C build -j$(nproc)`
- Removed stale root-owned dependency file `build/tests/CMakeFiles/faiss_test.dir/test_hamming.cpp.o.d`
- Reran `make -C build -j$(nproc)`, which completed successfully

### Why

The heap result struct requires exact pointer type compatibility. `long` and `long long` are both 64-bit on many Linux builds, but C++ treats them as distinct types, so using the heap's own `TI` type is the portable fix.

### What worked

- The type mismatch in `test_hamming.cpp` was fixed by deriving the id vector type from `faiss::int_maxheap_array_t::TI`
- Full `make -C build -j$(nproc)` completed and built `faiss_test`

### What didn't work

- First rerun after patch failed with:
  - `fatal error: opening dependency file CMakeFiles/faiss_test.dir/test_hamming.cpp.o.d: Permission denied`
- Cause: that `.o.d` file was owned by `root:root`, probably created by a previous sudo invocation
- Fix: remove the stale dependency file from the user-owned build directory and rerun make

### What I learned

- The FAISS library targets were already sufficient for bleve, but the test target needed a source patch for this platform's `int64_t` alias
- Root-owned intermediate files in a user build tree can mask source-level fixes with misleading compiler errors

### What was tricky to build

- The compile failure moved from a type mismatch to a permission error after the source patch. The second error was unrelated to C++; it came from a root-owned generated dependency file under `build/tests/CMakeFiles/faiss_test.dir/`.

### What warrants a second pair of eyes

- The local patch is source-level and should not be forgotten if the FAISS clone is reset; if this becomes permanent, keep it as a small patch file or fork branch.

### What should be done in the future

- Run `sudo make -C build install && sudo ldconfig`
- Re-test `go-faiss`/bleve vector compilation after install

### Code review instructions

- Review `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss/tests/test_hamming.cpp` around the `HammingHeapId` alias and `test_hamming_knn`
- Validate with `cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss && make -C build -j$(nproc)`

### Technical details

- Successful command: `cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss && make -C build -j$(nproc)`
- Permission fix command: `rm -f build/tests/CMakeFiles/faiss_test.dir/test_hamming.cpp.o.d`
- Relevant diff: `git diff -- tests/test_hamming.cpp`

---

## Step 5: Run the bleve KNN vector experiment

After FAISS headers were installed, the experiment progressed past CGO header compilation and surfaced two ordinary Go issues in the newly written experiment: an unused import and an incorrect assumption that `index.Document` exposes a `Fields` slice. I fixed the experiment by keeping the generated vectors in an in-memory `map[string][]float32` and using `chunk-042`'s original vector as the query vector.

The experiment then linked and ran successfully when pointed at the freshly built local FAISS libraries under `faiss/build/`. The pure KNN search returned `chunk-042` as the top hit with score `1.0`, confirming that bleve vector indexing and KNN query execution are working with the `vectors` build tag.

### Prompt Context

**User prompt (verbatim):** "ok, installed. run experiment"

**Assistant interpretation:** Verify the installed FAISS setup by running the new bleve KNN experiment command end-to-end.

**Inferred user intent:** Confirm that the FAISS + bleve vector stack works in this repository before moving on to goja-bleve implementation.

### What I did

- Verified `/usr/local/include/faiss/c_api/IndexBinary_c_ex.h` exists
- Ran `GOWORK=off go run -tags=vectors -ldflags "-r /usr/local/lib" ./cmd/experiments/bleve-knn/`
- Fixed experiment compile errors:
  - removed unused `mapping` import
  - stopped trying to access nonexistent `index.Document.Fields`
  - retained generated vectors in `vectorsByChunkID` for query construction
- Reran with explicit local FAISS library flags:
  - `CGO_LDFLAGS="-L$FAISS_SRC/build/c_api -L$FAISS_SRC/build/faiss -lfaiss_c -lfaiss -lstdc++ -lm"`
  - `-ldflags "-r $FAISS_SRC/build/c_api:$FAISS_SRC/build/faiss:/usr/local/lib"`

### Why

The installed headers were enough to compile CGO, but the linker needed a matching `libfaiss.so` paired with the new `libfaiss_c.so`. The system install showed `/usr/local/lib/libfaiss_c.so` updated, but `/usr/local/lib/libfaiss.so` still looked old, so the successful run used the freshly built local library paths.

### What worked

- Pure KNN search returned 5 hits; top hit was `chunk-042` with score `1.000000`
- Hybrid BM25 + KNN with `searchRequest.Score = "rrf"` returned fused results
- BM25-only comparison returned text-only results
- End-to-end bleve vector indexing and search worked with `-tags=vectors`

### What didn't work

- First run failed with Go compile errors:
  - unused import: `github.com/blevesearch/bleve/v2/mapping`
  - `queryDoc.Fields undefined (type index.Document has no field or method Fields)`
- Running against only `/usr/local/lib` failed at link time because `libfaiss_c.so` had unresolved references to FAISS C++ symbols; explicit `-lfaiss` and the fresh local `libfaiss.so` fixed it

### What I learned

- Bleve's `index.Document` API uses visitors, not a direct `Fields` slice
- For experiments it is cleaner to keep original vectors in application memory rather than read them back from bleve's internal document API
- The installed `/usr/local/lib/libfaiss.so` may still be the older May 2025 library while `/usr/local/lib/libfaiss_c.so` is new

### What was tricky to build

- The link failure looked like a go-faiss/bleve problem, but it was actually a library pairing issue: new C API library plus old or implicitly missing C++ FAISS library. Supplying both local build directories and explicit `-lfaiss_c -lfaiss -lstdc++ -lm` resolved it.

### What warrants a second pair of eyes

- Confirm whether `/usr/local/lib/libfaiss.so` should be overwritten with the freshly built `faiss/build/faiss/libfaiss.so`; the current experiment works with local rpaths but system-only builds may still fail.

### What should be done in the future

- Either install the fresh `libfaiss.so` to `/usr/local/lib` or document the required `CGO_LDFLAGS` + rpath for local FAISS builds
- Convert the standalone experiment into a reusable integration smoke test if bleve KNN becomes a production dependency

### Code review instructions

- Review `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/experiments/bleve-knn/main.go`
- Validate with:
  - `FAISS_SRC=/home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss`
  - `GOWORK=off CGO_LDFLAGS="-L$FAISS_SRC/build/c_api -L$FAISS_SRC/build/faiss -lfaiss_c -lfaiss -lstdc++ -lm" go run -tags=vectors -ldflags "-r $FAISS_SRC/build/c_api:$FAISS_SRC/build/faiss:/usr/local/lib" ./cmd/experiments/bleve-knn/`

### Technical details

Successful command:

```bash
FAISS_SRC=/home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss
cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system
GOWORK=off \
CGO_LDFLAGS="-L$FAISS_SRC/build/c_api -L$FAISS_SRC/build/faiss -lfaiss_c -lfaiss -lstdc++ -lm" \
go run -tags=vectors -ldflags "-r $FAISS_SRC/build/c_api:$FAISS_SRC/build/faiss:/usr/local/lib" ./cmd/experiments/bleve-knn/
```

Key output:

```text
KNN search returned 5 hits (total=5, maxScore=1.0000)
1. id=chunk-042 score=1.000000
Hybrid search returned 5 hits
BM25 search returned 5 hits
Summary: bleve KNN vector search works with the 'vectors' build tag.
```

### Follow-up: System FAISS install validation

After `/usr/local/lib/libfaiss.so` was updated to the fresh 14MB build, the experiment succeeded using only system library paths, but still required explicit `CGO_LDFLAGS` to include `-lfaiss`:

```bash
GOWORK=off \
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
go run -tags=vectors -ldflags "-r /usr/local/lib" ./cmd/experiments/bleve-knn/
```

Plain `go run -tags=vectors -ldflags "-r /usr/local/lib" ...` still fails because the Go link command only includes `-lfaiss_c`, and this FAISS C API shared library does not advertise `libfaiss.so` as a dynamic dependency. The practical fix for builds is to export `CGO_LDFLAGS` with both `-lfaiss_c` and `-lfaiss`.

---

## Step 6: Write FAISS compilation how-to

I wrote a durable repository document that captures the exact FAISS build, patch, install, and Bleve vector experiment workflow. The goal was to preserve the lessons from the setup session in a normal `docs/` location, not only in the ticket diary.

The how-to includes the `blevesearch/faiss@fff814d` checkout, CMake flags, the `test_hamming.cpp` heap-id type patch, the stale root-owned dependency-file trap, install verification, and the final known-good `CGO_LDFLAGS` command needed for Bleve vector builds.

### Prompt Context

**User prompt (verbatim):** "write a "howto compile faiss" document in the docs folder so we can retrace our steps later"

**Assistant interpretation:** Create a standalone documentation file in the rag-evaluation-system docs folder explaining how to reproduce the FAISS setup.

**Inferred user intent:** Make the FAISS build process repeatable for future agents/developers without reading the whole session transcript.

### What I did

- Created `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md`
- Included exact CMake/build/install commands
- Included troubleshooting for missing headers, `test_hamming.cpp`, stale root-owned `.o.d` files, and Go linker undefined references
- Related the doc to the ticket index and updated the docmgr changelog

### Why

The successful setup required several non-obvious details: using the Bleve FAISS fork, adding `-DCMAKE_CXX_FLAGS="-I$PWD"`, patching a test-only C++ type mismatch, and exporting explicit `CGO_LDFLAGS` with both `-lfaiss_c` and `-lfaiss`.

### What worked

- The document was written under the repository `docs/` folder and is ready for future reference
- Ticket index/changelog were updated through docmgr

### What didn't work

- `docmgr doc relate --doc <diary-path>` did not resolve the diary document path in this context, so I related the how-to to the ticket index with `--ticket RAGEVAL-GOJA-RAG-STRATEGIES` instead

### What I learned

- The reproducible command summary should include explicit `CGO_LDFLAGS`; relying on go-faiss's default `-lfaiss_c` alone is insufficient on this install

### What was tricky to build

- The how-to needed to distinguish between what is required for Bleve (`faiss` and `faiss_c` targets only) and what was required to make the full FAISS tree build (`test_hamming.cpp` patch). Mixing those would make the future workflow look more complicated than necessary.

### What warrants a second pair of eyes

- Confirm whether we want a Makefile target for the final Bleve vector experiment command so users do not need to remember `CGO_LDFLAGS` manually.

### What should be done in the future

- Optionally add a `make bleve-knn-experiment` or `make test-vectors` target that wraps the known-good command.

### Code review instructions

- Review `docs/howto-compile-faiss-for-bleve-vectors.md`
- Validate its final command by running the Bleve KNN experiment with system-installed FAISS

### Technical details

- New doc: `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/docs/howto-compile-faiss-for-bleve-vectors.md`

---

## Step 7: Commit Bleve KNN experiment and expand goja-bleve phase tasks

I committed the RAG evaluation system changes that make the Bleve vector experiment repeatable: the experiment command, the FAISS how-to, the Makefile wrapper target, and the module dependency updates. After that, I returned to the goja-bleve ticket and replaced the coarse six-item phase checklist with a detailed implementation roadmap.

The new task list breaks the work into ten phases, from repository scaffolding and Go-backed reference objects through mapping builders, index lifecycle, BM25 search, batch operations, vector/KNN search, hybrid scoring, provider registration, TypeScript declarations, and hardening. Each phase now has concrete implementation tasks plus a "Done when" criterion so an implementer can work through the project without reverse-engineering the design document.

### Prompt Context

**User prompt (verbatim):** "commit. then let's go back to building out goja-bleve .

Create phases and for each phase detailed tasks in the ticket."

**Assistant interpretation:** First commit the current Bleve KNN experiment/how-to work, then update the goja-bleve ticket with detailed phase-based tasks.

**Inferred user intent:** Preserve the validated vector-search experiment as a repository commit, then make the goja-bleve implementation plan actionable for continued development.

**Commit (code):** 16a50162277efaeb81b9853f40f1e93dd0cb7498 — "Add Bleve KNN vector experiment"

### What I did

- Committed `Makefile`, `cmd/experiments/bleve-knn/main.go`, `docs/howto-compile-faiss-for-bleve-vectors.md`, `go.mod`, and `go.sum`
- Rewrote `goja-bleve/.../tasks.md` with ten detailed implementation phases
- Added detailed task checkboxes and done criteria for each phase
- Updated the docmgr changelog for the task expansion
- Verified `docmgr task list --ticket RAGEVAL-GOJA-RAG-STRATEGIES` sees the expanded task list

### Why

The design document was comprehensive, but the ticket task list was too coarse for step-by-step implementation. The expanded checklist turns the design into actionable work packages.

### What worked

- The code commit succeeded with hash `16a50162277efaeb81b9853f40f1e93dd0cb7498`
- `docmgr task list` now enumerates the detailed phase tasks

### What didn't work

- N/A

### What I learned

- Docmgr reads checkbox tasks directly from `tasks.md`, so structured phase headings can coexist with a flat task list view.

### What was tricky to build

- The task list needed enough detail to be actionable without duplicating the entire design document. I kept each phase focused on concrete implementation and validation tasks, then used "Done when" criteria to define phase boundaries.

### What warrants a second pair of eyes

- Review whether Phase 3 and Phase 4 should be merged. Index lifecycle, BM25 search, and batch ergonomics are closely related, but separating batch work keeps the first text-search milestone smaller.

### What should be done in the future

- Start implementation at Phase 0 and Phase 1 in `goja-bleve/pkg/`.
- Consider committing the ticket docs once the goja-bleve repository bookkeeping policy is clear.

### Code review instructions

- Review `goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/tasks.md`
- Verify with `docmgr task list --ticket RAGEVAL-GOJA-RAG-STRATEGIES`

### Technical details

- Code commit: `16a50162277efaeb81b9853f40f1e93dd0cb7498`
- Task file: `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/02/RAGEVAL-GOJA-RAG-STRATEGIES--goja-bleve-bleve-bindings-for-goja-javascript-runtime/tasks.md`
