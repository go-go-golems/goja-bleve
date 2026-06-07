---
Title: Research logbook
Ticket: RAGEVAL-GOJA-BLEVE-XGOJA-MIGRATION-LOG
Status: active
Topics:
    - goja
    - xgoja
    - research
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: go-go-goja/cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md
      Note: Primary migration guide evaluated for provider API rename coverage and gaps
    - Path: go-go-goja/pkg/xgoja/app/runtime_spec.go
      Note: Current RuntimeSpec source used to diagnose generated app.Spec breakage
    - Path: go-go-goja/pkg/xgoja/providerapi/module.go
      Note: Current provider Module and ModuleSetupContext source used to verify signatures
    - Path: go-go-goja/pkg/xgoja/providerapi/provider_registry.go
      Note: Current ProviderRegistry source used to verify registry constructor and ResolveModule behavior
    - Path: goja-bleve/README.md
      Note: Public snippet updated from NewRegistry to NewProviderRegistry
    - Path: goja-bleve/cmd/goja-bleve/main.go
      Note: Generated command source that exposed NewProviderRegistry and RuntimeSpec migration needs
    - Path: goja-bleve/pkg/provider.go
      Note: Root goja-bleve provider registration migrated from old xgoja API
    - Path: goja-bleve/pkg/xgoja/providers/bleve/bleve.go
      Note: Provider package migrated from old xgoja API
ExternalSources: []
Summary: Resource-by-resource log for the goja-bleve migration to the renamed xgoja provider and runtime APIs.
LastUpdated: 2026-06-05T10:45:00-04:00
WhatFor: Track which migration resources were useful, stale, misleading, or need documentation updates after moving goja-bleve to the upgraded go-go-goja API.
WhenToUse: Use before continuing xgoja provider migrations, refreshing generated xgoja binaries, or updating go-go-goja migration documentation.
---


# Research logbook

## Goal

Keep an evidence trail for the `goja-bleve` migration to the upgraded `go-go-goja` / `xgoja` provider API. This log records which resources were consulted, why each one was chosen, what it clarified, what was stale, and what should be updated for the next migration.

## Context

The immediate compile errors were:

```text
# github.com/go-go-golems/goja-bleve/pkg
pkg/provider.go:14:45: undefined: providerapi.Registry
pkg/provider.go:19:3: unknown field New in struct literal of type providerapi.Module
pkg/provider.go:19:25: undefined: providerapi.ModuleContext
```

The package was migrated from the old provider API names to the new ones:

| Old API | New API used in goja-bleve |
| --- | --- |
| `providerapi.Registry` | `providerapi.ProviderRegistry` |
| `providerapi.NewRegistry()` | `providerapi.NewProviderRegistry()` |
| `providerapi.Module.New` | `providerapi.Module.NewModuleFactory` |
| `providerapi.ModuleContext` | `providerapi.ModuleSetupContext` |
| `app.Spec` | `app.RuntimeSpec` |

Validation commands that passed after the migration:

```bash
cd goja-bleve && go test ./...
cd goja-bleve && GOWORK=off go test ./pkg/...
cd goja-bleve/cmd/goja-bleve && GOWORK=off go test ./...
```

## Resource log

### 1. `go-go-goja/cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md`

- **What I was researching:** The intentional breaking changes in the xgoja provider and engine API.
- **What I was looking for in this document:** Direct name replacements for `providerapi.Registry`, `providerapi.ModuleContext`, and `providerapi.Module.New`, because those names appeared verbatim in the compile errors.
- **Why I chose it:** The user explicitly pointed to this migration document as likely helpful.
- **How I found the resource itself:** User-provided path in the prompt.
- **What I found useful:** The direct replacement table exactly covered the provider compile errors: `NewRegistry` → `NewProviderRegistry`, `Registry` → `ProviderRegistry`, `ModuleContext` → `ModuleSetupContext`, and `Module.New` → `Module.NewModuleFactory`.
- **What I didn't find useful:** It did not cover the generated xgoja binary's `decodeSpec()` helper using the removed `app.Spec` type.
- **What is out of date / what was wrong:** The document is accurate for provider API renames, but incomplete for generated-binary migration failures.
- **What would need updating:** Add a troubleshooting row for `undefined: app.Spec`, with the fix `app.RuntimeSpec`. Also add a short note that generated `cmd/<binary>/go.mod` files may need their `github.com/go-go-golems/go-go-goja` version bumped alongside the root module.

### 2. `go-go-goja/pkg/xgoja/providerapi/module.go`

- **What I was researching:** The current shape of provider module registration.
- **What I was looking for in this document:** The authoritative `providerapi.Module` struct fields and the context type expected by module factories.
- **Why I chose it:** After reading the migration table, I needed the concrete type definitions to avoid guessing function signatures.
- **How I found the resource itself:** Repository search for `type Module` and `ModuleSetupContext` under `go-go-goja/pkg/xgoja/providerapi`.
- **What I found useful:** It confirmed `ModuleSetupContext` fields and that `Module` now has `NewModuleFactory func(ModuleSetupContext) (require.ModuleLoader, error)`.
- **What I didn't find useful:** The type comments explain runtime setup but do not include an example provider snippet.
- **What is out of date / what was wrong:** Nothing observed in the source itself.
- **What would need updating:** Optional: add a brief example in the package docs or keep relying on the migration document and tutorial docs.

### 3. `go-go-goja/pkg/xgoja/providerapi/provider_registry.go`

- **What I was researching:** The replacement registry constructor and registry type.
- **What I was looking for in this document:** `ProviderRegistry`, `NewProviderRegistry()`, `Package(...)`, and `ResolveModule(...)` behavior used by goja-bleve tests.
- **Why I chose it:** Tests and runtime setup needed to instantiate and query the provider registry.
- **How I found the resource itself:** Repository search for `type ProviderRegistry` and `func NewProviderRegistry`.
- **What I found useful:** It confirmed that `ResolveModule(...)` still returns `providerapi.Module`, so existing test assertions could stay the same apart from the renamed factory field.
- **What I didn't find useful:** It does not describe migration behavior; it is implementation source, not migration guidance.
- **What is out of date / what was wrong:** Nothing observed.
- **What would need updating:** No source update needed. Documentation could link to this registry source from provider tutorials for maintainers.

### 4. `go-go-goja/pkg/xgoja/testprovider/provider.go`

- **What I was researching:** A known-good provider implementation already using the new API.
- **What I was looking for in this document:** Real examples of `Register(registry *providerapi.ProviderRegistry)`, `NewModuleFactory`, `ModuleSetupContext`, `SectionRequest`, and `RuntimeInitializerHandle.EngineRuntime()`.
- **Why I chose it:** Tests are usually the most reliable migration examples because they compile against the current API.
- **How I found the resource itself:** Search for `NewModuleFactory` and `ProviderRegistry` in `go-go-goja`.
- **What I found useful:** It validated the exact provider registration style used in `goja-bleve`: a `Register` function receiving `*providerapi.ProviderRegistry` and returning `registry.Package(...)` entries.
- **What I didn't find useful:** It contains broader fixture command/provider examples beyond the simple native-module-only case.
- **What is out of date / what was wrong:** Nothing observed.
- **What would need updating:** No code update needed. A reduced version of the simple module example could be copied into `cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md` or the provider tutorial.

### 5. `goja-bleve/pkg/provider.go`

- **What I was researching:** The exact failing provider registration in the root `pkg` package.
- **What I was looking for in this document:** Old API usages matching the compile errors.
- **Why I chose it:** The compiler identified this file as the immediate failure.
- **How I found the resource itself:** Compiler output and direct file read.
- **What I found useful:** The provider was small and isolated; all required changes were localized to the register signature and module factory field.
- **What I didn't find useful:** It did not reveal whether there were additional old API usages elsewhere.
- **What is out of date / what was wrong:** It still used `providerapi.Registry`, `providerapi.Module.New`, and `providerapi.ModuleContext`.
- **What would need updating:** Updated to `*providerapi.ProviderRegistry`, `NewModuleFactory`, and `ModuleSetupContext`.

### 6. `goja-bleve/pkg/xgoja/providers/bleve/bleve.go`

- **What I was researching:** The generated/provider-facing package that xgoja imports for the `goja-bleve` provider.
- **What I was looking for in this document:** Duplicate old API usages not shown in the initial `pkg/provider.go` compiler error.
- **Why I chose it:** Repository search found old `providerapi.Registry` and `ModuleContext` usages here too.
- **How I found the resource itself:** `rg` search for `providerapi.Registry`, `ModuleContext`, and `New:` in `goja-bleve`.
- **What I found useful:** It showed the second provider registration path that needed the same migration.
- **What I didn't find useful:** The code only wraps a native module from `modules.DefaultRegistry`; it did not cover runtime config or cleanup migration cases.
- **What is out of date / what was wrong:** It used the old registry and module factory names.
- **What would need updating:** Updated to the new provider API. If this file is generated by xgoja in the future, regenerate it with the upgraded generator rather than editing manually.

### 7. `goja-bleve/pkg/provider_test.go` and `goja-bleve/pkg/xgoja/providers/bleve/bleve_test.go`

- **What I was researching:** How the provider tests exercised registry construction and module loader creation.
- **What I was looking for in these documents:** Test breakages caused by old constructor/factory names.
- **Why I chose them:** After changing provider source, tests needed to compile and still verify module registration.
- **How I found the resources themselves:** Search results for `providerapi.NewRegistry`, `entry.New`, and `providerapi.ModuleContext`.
- **What I found useful:** The tests provided a quick local validation path for the migrated API and confirmed `require("bleve")` still works.
- **What I didn't find useful:** They only test module loading; they do not exercise xgoja command generation or generated binary wiring.
- **What is out of date / what was wrong:** Test setup used old API names.
- **What would need updating:** Updated test setup to `NewProviderRegistry()` and `entry.NewModuleFactory(providerapi.ModuleSetupContext{...})`.

### 8. `goja-bleve/cmd/goja-bleve/main.go`

- **What I was researching:** Whether the generated xgoja command binary still compiled after the provider source migration.
- **What I was looking for in this document:** Old provider registry construction and any additional compile errors introduced by upgrading `go-go-goja`.
- **Why I chose it:** `goja-bleve/cmd/goja-bleve` is a nested module and may compile against module versions independent of the root package.
- **How I found the resource itself:** Search for `providerapi.NewRegistry` and then `GOWORK=off go test ./...` inside the nested command module.
- **What I found useful:** It exposed a second migration issue: generated helper `decodeSpec()` referenced `app.Spec`, which no longer exists in the upgraded app package.
- **What I didn't find useful:** The file is generated and has a `DO NOT EDIT` header, so manual edits are a tactical fix rather than the ideal long-term maintenance path.
- **What is out of date / what was wrong:** It used `providerapi.NewRegistry()` and `app.Spec`.
- **What would need updating:** Updated to `providerapi.NewProviderRegistry()` and `app.RuntimeSpec`. Longer term, regenerate this file with the current xgoja generator so the generated source is canonical.

### 9. `go-go-goja/pkg/xgoja/app/runtime_spec.go`

- **What I was researching:** The replacement for removed `app.Spec`.
- **What I was looking for in this document:** The current runtime spec DTO type used by generated xgoja binaries.
- **Why I chose it:** `cmd/goja-bleve/main.go` failed with `undefined: app.Spec` after upgrading the nested module to `go-go-goja v0.8.1`.
- **How I found the resource itself:** Search for `type RuntimeSpec` and `type Spec` under `go-go-goja/pkg/xgoja/app`.
- **What I found useful:** The comments explicitly state `RuntimeSpec` is the normalized embedded runtime spec decoded by generated xgoja binaries.
- **What I didn't find useful:** The migration document did not point here, so this had to be discovered from source search.
- **What is out of date / what was wrong:** Nothing observed in the source. The generated goja-bleve command source was stale relative to this current type.
- **What would need updating:** Add this rename to xgoja migration documentation and ensure generated templates no longer emit `app.Spec`.

### 10. `goja-bleve/go.mod`, `goja-bleve/go.sum`, `goja-bleve/cmd/goja-bleve/go.mod`, and `goja-bleve/cmd/goja-bleve/go.sum`

- **What I was researching:** Whether the package and nested generated command module were actually using the upgraded `go-go-goja` version.
- **What I was looking for in these documents:** `github.com/go-go-golems/go-go-goja` version pins and dependency changes needed for standalone builds.
- **Why I chose them:** Workspace tests can compile against the local `go-go-goja` checkout, while standalone consumers use `go.mod` versions.
- **How I found the resources themselves:** `go test` with and without `GOWORK=off`, plus direct reads of `go.mod` files.
- **What I found useful:** `GOWORK=off` exposed that `v0.7.4` did not contain the new API. Bumping to `v0.8.1` made standalone package tests compile.
- **What I didn't find useful:** The nested command module has a large indirect dependency set, so `go mod tidy` produces noisy `go.sum` changes.
- **What is out of date / what was wrong:** Both root and nested modules pinned `go-go-goja v0.7.4`, which is incompatible with the migrated source.
- **What would need updating:** Keep both module files on the same upgraded `go-go-goja` version. If the generated command module remains checked in, it must be included in future API migration sweeps.

### 11. `goja-bleve/README.md`

- **What I was researching:** User-facing examples that still referenced old provider API names.
- **What I was looking for in this document:** Code snippets containing `providerapi.NewRegistry()`.
- **Why I chose it:** Public documentation should not teach newly removed API names after source migration.
- **How I found the resource itself:** Repository search for `providerapi.NewRegistry`.
- **What I found useful:** It contained a short xgoja host snippet that could be updated directly.
- **What I didn't find useful:** It does not explain the provider API migration or generated binary caveats.
- **What is out of date / what was wrong:** The snippet used `providerapi.NewRegistry()`.
- **What would need updating:** Updated to `providerapi.NewProviderRegistry()`. Consider adding a short compatibility note requiring `go-go-goja >= v0.8.1`.

### 12. Repository search output (`rg` over `goja-bleve` and `go-go-goja`)

- **What I was researching:** Exhaustive occurrences of removed provider API names.
- **What I was looking for in this resource:** Remaining `providerapi.Registry`, `providerapi.NewRegistry`, `providerapi.ModuleContext`, `entry.New(...)`, and `app.Spec` references outside ticket notes.
- **Why I chose it:** API migrations are easy to miss if only compiler-reported files are changed.
- **How I found the resource itself:** Direct shell searches with `rg`.
- **What I found useful:** It identified the nested command module, README snippet, and provider-specific package that the first compiler error did not enumerate.
- **What I didn't find useful:** Ticket history under `ttmp/` intentionally contains old names as historical documentation; those should not necessarily be rewritten.
- **What is out of date / what was wrong:** Active source/docs outside `ttmp/` had stale provider API names before the migration.
- **What would need updating:** Keep the search pattern as a migration checklist for other packages. Exclude `ttmp/` when looking for active code/doc references.

### 13. `geppetto/pkg/js/modules/geppetto/api_embeddings.go`

- **What I was researching:** How to expose Geppetto embeddings to xgoja scripts so a JavaScript verb can embed documents and queries before indexing vectors in Bleve.
- **What I was looking for in this document:** A JS-facing API parallel to existing `geppetto` engine/profile helpers, with a small object returned by `geppetto.embeddings(settings)`.
- **Why I chose it:** The requested RAG tool needed `goja-bleve` plus Geppetto bindings in one generated xgoja binary, and the existing Geppetto JS module did not expose embeddings directly.
- **How I found the resource itself:** By reading the existing `geppetto/pkg/js/modules/geppetto/api_*` files and matching their module-runtime pattern.
- **What I found useful:** The existing Geppetto module style made it straightforward to add `embed(text)`, `embedBatch(texts)`, and `model()` while preserving the module's native Go-backed behavior.
- **What I didn't find useful:** The existing APIs focused on inference and profile resolution, so embeddings needed new tests and TypeScript declarations rather than a small wrapper-only change.
- **What is out of date / what was wrong:** Geppetto's hardcut docs/tests did not yet demonstrate embeddings through a resolved registry profile.
- **What would need updating:** Keep the embedding API in the TypeScript declaration template and generated declaration in sync whenever embedding settings change.

### 14. `go-go-goja/cmd/xgoja/internal/buildexec` and `go-go-goja/cmd/xgoja/internal/buildspec`

- **What I was researching:** How to express FAISS/Bleve vector linker requirements in an xgoja YAML spec instead of relying on shell-local environment variables.
- **What I was looking for in these documents:** Build-spec support for environment variables passed into `go build`, especially `CGO_LDFLAGS` for `-lfaiss_c -lfaiss -lstdc++ -lm`.
- **Why I chose them:** The vector xgoja binary must be reproducible from `xgoja-vectors.yaml`, including `go.tags: [vectors]`, `go.ldflags`, and CGO linker environment.
- **How I found the resources themselves:** By inspecting the xgoja build command path from YAML parsing through `go build` execution.
- **What I found useful:** Adding `go.env` kept the vector build declarative and let `xgoja-vectors.yaml` capture the known-good FAISS linker incantation.
- **What I didn't find useful:** The old build-spec reference documented tags and ldflags but not build-time environment variables.
- **What is out of date / what was wrong:** Without `go.env`, vector builds required a hidden external `CGO_LDFLAGS` shell setup.
- **What would need updating:** Keep `cmd/xgoja/doc/06-buildspec-reference.md` as the source of truth for `go.env`, and validate vector specs with `xgoja build -f xgoja-vectors.yaml`.

### 15. `go-go-goja/pkg/xgoja/app/module_sections.go`

- **What I was researching:** Why root-mounted jsverb groups disappeared after adding the Geppetto provider to the generated goja-bleve binary.
- **What I was looking for in this document:** The path that attaches provider runtime config sections to jsverb command descriptions before Glazed/Cobra mounting.
- **Why I chose it:** The scanner correctly discovered `rag index-query`, but the generated app displayed `rag` as a help-only parent with no child commands.
- **How I found the resource itself:** By adding a temporary command-tree inspection test in `goja-bleve/cmd/goja-bleve` and comparing bare jsverb mounting with app-level mounting.
- **What I found useful:** It revealed that provider runtime sections were being added to jsverb command descriptions, and that section merging needed to preserve existing command schemas.
- **What I didn't find useful:** The first suspected bug, command schema replacement, was real enough to deserve a regression test, but it was not the only cause of the missing `rag` subcommand.
- **What is out of date / what was wrong:** `appendSectionsToCommandDescription` replaced existing command sections with runtime sections, which could drop jsverb arguments/fields after runtime config was attached.
- **What would need updating:** The helper now preserves original command sections before appending runtime sections. The new regression test verifies a command argument survives runtime section attachment.

### 16. `goja-bleve/cmd/goja-bleve/jsverbs/rag.js`

- **What I was researching:** A generated JavaScript verb tool that combines Geppetto embeddings with Bleve vector indexing and hybrid querying.
- **What I was looking for in this document:** A safe `plan` command plus an `index-query` command that resolves an embedding profile, embeds documents, indexes vectors, and queries with KNN or hybrid RRF.
- **Why I chose it:** It is the requested user-facing tool surface for exercising `goja-bleve` and Geppetto together.
- **How I found the resource itself:** Created during the migration as a new jsverb file under the existing `cmd/goja-bleve/jsverbs` bundle.
- **What I found useful:** The `plan` command validates wiring without external provider calls. The `index-query` command ran successfully with local Ollama `all-minilm` embeddings and returned three ranked hits.
- **What I didn't find useful:** Hyphenated function names are invalid for `__verb__`; the scanner requires the first argument to resolve to an actual JavaScript function identifier. Field defaults also need to match jsverbs field typing expectations.
- **What is out of date / what was wrong:** The first version used `__verb__("index-query")`, which failed scanning because no JavaScript function can have that identifier. It also used a local `profile` field that collided with Geppetto's runtime `profile` flag, causing the `rag` group to mount without children.
- **What would need updating:** Keep the exported function as `indexQuery`, let jsverbs expose it as `index-query`, and use `embeddingProfile` for the verb-specific embedding profile to avoid the Geppetto runtime section's `profile` flag.

### 17. `goja-bleve/cmd/goja-bleve/xgoja.yaml` and `xgoja-vectors.yaml`

- **What I was researching:** How to generate one non-vector xgoja binary and one vector-enabled binary that include both `goja-bleve` and `geppetto` providers.
- **What I was looking for in these documents:** Provider package selection, module aliases, root-mounted jsverbs, local replaces, and vector build settings.
- **Why I chose them:** These specs are the reproducible build contract for `dist/goja-bleve` and `dist/goja-bleve-vectors`.
- **How I found the resources themselves:** Existing xgoja command specs under `goja-bleve/cmd/goja-bleve`.
- **What I found useful:** Adding Geppetto plus core/host modules allowed `rag.js` to call `require("geppetto")`, `require("bleve")`, and host filesystem/path helpers in one runtime.
- **What I didn't find useful:** The generated binary only used the patched local xgoja app code after explicit `replace: ../../../go-go-goja` entries were added for the go-go-goja provider packages.
- **What is out of date / what was wrong:** The first vector build pattern still depended on shell-provided `CGO_LDFLAGS`; the updated vector spec now carries it under `go.env`.
- **What would need updating:** Keep both specs' provider replaces aligned while this workspace depends on local `go-go-goja`, `goja-bleve`, and `geppetto` changes.

### 18. Validation commands and final RAG smoke test

- **What I was researching:** Whether the migrated provider API, xgoja build env, vector tags/linking, and Geppetto+Bleve jsverb tool work together end-to-end.
- **What I was looking for in these commands:** Passing unit tests, successful non-vector and vector xgoja builds, root-mounted jsverb commands, and a real embedding-backed vector/hybrid query.
- **Why I chose them:** Unit tests validate package-level behavior, while the generated binary smoke test validates the requested user-facing tool path.
- **How I found the resources themselves:** Known-good validation commands from prior goja-bleve phases plus the new `rag index-query` command.
- **What I found useful:** The final smoke test used local Ollama `all-minilm` embeddings through Geppetto, indexed three demo documents into Bleve, and returned `chunk-1` as the top hit for `privacy`.
- **What I didn't find useful:** `--profile-path` and `--embedding-profile` were not accepted because current jsverbs exposes camelCase field keys as `--profilePath` and `--embeddingProfile`.
- **What is out of date / what was wrong:** The initial `plan` output showed kebab-case flags; it now shows the actual generated camelCase flags.
- **What would need updating:** If jsverbs later normalizes camelCase fields to kebab-case flags, update `rag.js` plan output and examples accordingly.

## Recommended xgoja documentation updates

Now that an existing package has been migrated, these xgoja docs would be worth updating:

1. **Update `cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md`**
   - Add `app.Spec` → `app.RuntimeSpec` to the replacement table.
   - Add troubleshooting row: `undefined: app.Spec` means generated xgoja source is stale or using old generated templates; use `app.RuntimeSpec` or regenerate with the current xgoja generator.
   - Add a note that nested generated command modules can pin an older `go-go-goja` version even when the workspace build succeeds.

2. **Add a generated-binary migration checklist**
   - Search active source and docs, excluding `ttmp/`, for:
     - `providerapi.NewRegistry`
     - `providerapi.Registry`
     - `providerapi.ModuleContext`
     - `.New(providerapi.ModuleContext`
     - `app.Spec`
   - Run both workspace and standalone validation:
     - `go test ./...`
     - `GOWORK=off go test ./...`
     - for nested generated modules, `cd cmd/<binary> && GOWORK=off go test ./...`

3. **Add a short provider-only example**
   - The migration doc already has one `providerapi.Module` snippet, but a complete minimal `Register(registry *providerapi.ProviderRegistry) error` function copied from the successful pattern would make migrations faster.

4. **Document checked-in generated files**
   - If generated xgoja command modules are committed, document whether maintainers should manually patch them during migrations or regenerate them with `xgoja`.

## Quick reference: migration checklist for similar packages

```bash
# Find active old provider/runtime API references, excluding ticket history.
rg -n "providerapi\.(Registry|ModuleContext|NewRegistry)|\.New\(providerapi\.ModuleContext|app\.Spec" \
  --glob '!ttmp/**' \
  .

# Validate in workspace mode.
go test ./...

# Validate standalone root module mode.
GOWORK=off go test ./...

# Validate nested generated command module mode, if present.
cd cmd/<generated-binary> && GOWORK=off go test ./...
```

## Related

- `go-go-goja/cmd/xgoja/doc/10-migrating-xgoja-provider-engine-api.md`
- `go-go-goja/pkg/xgoja/providerapi/module.go`
- `go-go-goja/pkg/xgoja/providerapi/provider_registry.go`
- `go-go-goja/pkg/xgoja/app/runtime_spec.go`
- `goja-bleve/pkg/provider.go`
- `goja-bleve/pkg/xgoja/providers/bleve/bleve.go`
- `goja-bleve/cmd/goja-bleve/main.go`
