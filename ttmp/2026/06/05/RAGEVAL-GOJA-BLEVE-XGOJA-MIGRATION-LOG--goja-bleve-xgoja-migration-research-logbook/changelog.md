# Changelog

## 2026-06-05

- Initial workspace created


## 2026-06-05

Created research logbook for goja-bleve xgoja provider API migration, including useful/stale resource notes and proposed xgoja documentation updates.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md — Migration document assessed for follow-up updates
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/05/RAGEVAL-GOJA-BLEVE-XGOJA-MIGRATION-LOG--goja-bleve-xgoja-migration-research-logbook/reference/01-research-logbook.md — Research logbook deliverable


## 2026-06-05

Uploaded research logbook PDF to reMarkable at /ai/2026/06/05/RAGEVAL-GOJA-BLEVE-XGOJA-MIGRATION-LOG.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/05/RAGEVAL-GOJA-BLEVE-XGOJA-MIGRATION-LOG--goja-bleve-xgoja-migration-research-logbook/reference/01-research-logbook.md — Uploaded logbook source


## 2026-06-05

Updated xgoja migration guide with goja-bleve follow-ups: app.Spec to app.RuntimeSpec, generated-binary checklist, standalone go.mod warning, provider example, and troubleshooting rows.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md — Migration guide updated from goja-bleve migration findings
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/ttmp/2026/06/05/RAGEVAL-GOJA-BLEVE-XGOJA-MIGRATION-LOG--goja-bleve-xgoja-migration-research-logbook/reference/01-research-logbook.md — Source logbook that identified guide updates


## 2026-06-05

Completed the goja-bleve/xgoja migration implementation: added Geppetto embeddings bindings, xgoja `go.env` build support, vector-tagged xgoja specs with FAISS linker env, a Geppetto+Bleve `rag index-query` jsverb, and an xgoja runtime-section schema preservation fix. Validated the full generated vector binary against local Ollama `all-minilm` embeddings.

### Validation

- `go test ./cmd/xgoja/internal/buildexec ./cmd/xgoja/internal/buildspec ./cmd/xgoja ./pkg/xgoja/app -count=1`
- `go test ./pkg/js/modules/geppetto ./pkg/js/modules/geppetto/provider -count=1`
- `GOWORK=off go test ./... -count=1`
- `GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1`
- `./dist/goja-bleve-vectors rag index-query --profilePath ../../../geppetto/examples/js/geppetto/profiles/40-embeddings.yaml --embeddingProfile ollama-all-minilm-embedding --output json privacy`

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/geppetto/pkg/js/modules/geppetto/api_embeddings.go — New JavaScript embeddings API used by the RAG jsverb
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/cmd/xgoja/internal/buildexec/buildexec.go — xgoja build execution path now accepts build-time environment variables
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/xgoja/app/module_sections.go — Runtime section attachment now preserves existing command schemas
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/rag.js — Geppetto+Bleve RAG indexing/querying jsverb
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.yaml — Vector-enabled xgoja spec with build tags, rpath, and FAISS CGO linker environment

## 2026-06-05

Expanded go-go-goja jsverbs naming behavior so top-level JavaScript parameter fields use idiomatic kebab-case CLI flags while preserving JavaScript parameter names at invocation time. Verified `profilePath` -> `--profile-path` and `foo_bar` -> `--foo-bar`, regenerated goja-bleve binaries, and reran the RAG smoke test with kebab-case flags.

### Validation

- `go test ./pkg/jsverbs -run 'TestTopLevelFieldNamesUseKebabCaseCLI|TestCommandDescriptionForVerb|TestCommandForVerbWithInvokerUsesCustomInvoker|TestCommandsWithInvokerNilFallsBackToDefaultExecution' -count=1`
- `go test ./cmd/xgoja/internal/buildexec ./cmd/xgoja/internal/buildspec ./cmd/xgoja ./pkg/xgoja/app -count=1`
- `./dist/goja-bleve-vectors rag index-query --profile-path ../../../geppetto/examples/js/geppetto/profiles/40-embeddings.yaml --embedding-profile ollama-all-minilm-embedding --output json privacy`

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/jsverbs/command.go — Normalizes top-level jsverb field names for CLI exposure
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/jsverbs/runtime.go — Looks up normalized CLI field names while preserving positional JavaScript argument delivery
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/jsverbs/jsverbs_test.go — Regression for `profilePath` and `foo_bar` field names
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/doc/11-jsverbs-example-reference.md — Documents top-level field naming behavior and current section-field caveat
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/rag.js — Plan output updated to kebab-case flags
