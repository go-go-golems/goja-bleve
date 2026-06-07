# Tasks

## DONE

- [x] Migrate goja-bleve provider code and generated command module to the new xgoja provider API.
- [x] Add Geppetto JavaScript embeddings API with tests and TypeScript declaration updates.
- [x] Extend xgoja build specs/executor to support build-time `go.env`, including `CGO_LDFLAGS` for Bleve vector builds.
- [x] Update goja-bleve xgoja specs to include Geppetto, core, host, root-mounted jsverbs, and vector build tags/linking.
- [x] Add `rag plan` and `rag index-query` jsverbs that combine Geppetto embeddings with Bleve vector/hybrid search.
- [x] Fix xgoja runtime-section attachment so command schemas survive section augmentation.
- [x] Diagnose and avoid the Geppetto `profile` runtime flag collision by using `embeddingProfile` in the RAG verb.
- [x] Validate non-vector tests, vector tests, generated binary builds, vector jsverbs, and the full local Ollama RAG smoke test.
- [x] Normalize top-level jsverb parameter field CLI flags to kebab-case and update `rag.js` plan output/docs.

## TODO

- [ ] Decide whether to publish a follow-up xgoja release containing `go.env`, runtime-section schema preservation, and top-level jsverb field-name normalization.
- [ ] Decide whether section field flags should also gain kebab-case aliases while preserving JavaScript object keys.
