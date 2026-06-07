# Changelog

## 2026-06-07

- Initial workspace created


## 2026-06-07

Created intern-oriented design guide and diary for generated xgoja vector smoke targets plus optional FAISS CI.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/design-doc/01-generated-xgoja-vector-smoke-and-faiss-ci-implementation-guide.md — Primary implementation guide
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/reference/01-investigation-diary.md — Chronological design diary


## 2026-06-07

Uploaded generated xgoja vector smoke and FAISS CI design bundle to reMarkable at /ai/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/design-doc/01-generated-xgoja-vector-smoke-and-faiss-ci-implementation-guide.md — Uploaded design guide


## 2026-06-07

Step 2: Added local generated xgoja vector smoke Makefile targets and validated make xgoja-smoke-vectors.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile — New xgoja-build-vectors and xgoja-smoke-vectors targets
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/reference/01-investigation-diary.md — Step 2 implementation notes


## 2026-06-07

Step 3: Documented make xgoja-smoke-vectors in README, quickstart, docs index, and FAISS/xgoja playbook; validated make test-vectors and make xgoja-smoke-vectors.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/README.md — Documents generated xgoja vector smoke validation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md — Documents target usage
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/reference/01-investigation-diary.md — Step 3 documentation notes


## 2026-06-07

Step 4: Added optional Vector FAISS Smoke GitHub Actions workflow and documented its manual/scheduled behavior.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml — Optional FAISS-backed vector test workflow
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md — Documents optional workflow and xgoja smoke caveat
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI--goja-bleve-generated-xgoja-vector-smoke-and-faiss-ci-design/reference/01-investigation-diary.md — Step 4 workflow notes


## 2026-06-07

Step 5: Added xgoja-vectors.ci.yaml and enabled generated xgoja vector smoke in the optional FAISS workflow.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml — Runs make test-vectors and generated xgoja smoke with CI spec
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.ci.yaml — Clean-checkout xgoja vector spec for CI
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md — Documents local vs CI xgoja vector specs

