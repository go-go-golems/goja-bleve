---
Title: Investigation diary
Ticket: RAGEVAL-GOJA-BLEVE-VECTOR-CI
Status: active
Topics:
    - goja
    - xgoja
    - bleve
    - ci
    - vector-search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: Makefile
      Note: Current vector test target and proposed generated smoke target location
    - Path: cmd/goja-bleve/xgoja-vectors.yaml
      Note: Evidence for local vector xgoja build settings and CI portability issue
    - Path: docs/faiss-xgoja-playbook.md
      Note: FAISS linker and runtime-loader evidence used by design
ExternalSources: []
Summary: Diary for the generated xgoja vector smoke and FAISS CI design ticket.
LastUpdated: 2026-06-07T13:45:00-04:00
WhatFor: Use this to understand why the ticket was created, what evidence shaped the design, and how to continue implementation.
WhenToUse: Read before implementing generated xgoja vector smoke targets or optional FAISS-backed CI for goja-bleve.
---


# Diary

## Goal

This diary records the setup and design work for adding generated xgoja vector smoke coverage and optional FAISS-backed CI to `goja-bleve`.

## Step 1: Create the vector smoke and FAISS CI design ticket

I created a dedicated docmgr ticket for the first two hardening items the user selected: generated xgoja smoke coverage and optional FAISS CI. The work here is design-only; it explains the system boundaries, implementation sequence, and failure modes for an intern who has not worked in this repository before.

The key design conclusion is that package vector tests and generated xgoja smoke tests validate different layers. `make test-vectors` proves the Go package can link FAISS and run vector-tagged tests. A new `make xgoja-smoke-vectors` target should prove that xgoja can build the vector host from YAML and run embedded JavaScript verbs through the generated CLI.

### Prompt Context

**User prompt (verbatim):** "Actually, let's create 6 and 7 first. Create a new docmgr ticket, and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket for hardening tasks 6 and 7, write a detailed intern-oriented design and implementation guide, store it in the ticket workspace, and upload it to reMarkable.

**Inferred user intent:** Preserve a clear implementation plan before writing code so a new contributor can safely add generated xgoja vector smoke coverage and optional FAISS CI without rediscovering the xgoja, FAISS, jsverb, and CI constraints.

**Commit (code):** N/A — documentation/design work only.

### What I did

- Created ticket `RAGEVAL-GOJA-BLEVE-VECTOR-CI` with topics `goja,xgoja,bleve,ci,vector-search`.
- Added the design document `design-doc/01-generated-xgoja-vector-smoke-and-faiss-ci-implementation-guide.md`.
- Added this diary document.
- Inspected the current `Makefile`, GitHub Actions workflows, `cmd/goja-bleve/xgoja-vectors.yaml`, and existing FAISS/xgoja docs.
- Wrote an intern-oriented guide covering:
  - current module and generated-host architecture,
  - package vector tests vs generated xgoja smoke tests,
  - FAISS build and link requirements,
  - CI portability problems caused by local `replace` paths,
  - proposed Makefile targets,
  - proposed optional GitHub Actions workflow,
  - pseudocode, diagrams, API references, failure modes, file references, and review checklist.

### Why

The previous hardening discussion identified generated xgoja smoke coverage and FAISS CI as useful next steps. They should be designed before implementation because they cross several subsystem boundaries: native Go modules, xgoja code generation, jsverb command mounting, CGO linker flags, FAISS native libraries, nested generated modules, and GitHub Actions runner setup.

### What worked

- The existing repository already has `make test-vectors`, so the design can build on a known-good local vector validation command.
- `cmd/goja-bleve/xgoja-vectors.yaml` already encodes the correct vector tag, rpath, and `CGO_LDFLAGS` through xgoja `go.env`.
- The generated vector smoke commands are deterministic: `vector knn` and `vector hybrid` do not need an external embedding provider.
- The FAISS playbook created earlier provides the exact native-library troubleshooting content needed by the CI design.

### What didn't work

- No implementation was attempted in this step.
- A clean GitHub Actions checkout will not automatically satisfy the local sibling `replace` paths currently present in `xgoja-vectors.yaml`. The design calls this out as a portability issue and recommends either a CI-specific spec or sibling repository checkouts.

### What I learned

- The generated xgoja smoke target should not be folded into `test-vectors`; each target validates a distinct layer and should fail independently.
- The first FAISS CI workflow should be non-required and manually/scheduled, because source-building FAISS can be slow or sensitive to runner changes.
- The RAG Geppetto command is not the right first CI smoke because it needs real embedding provider configuration; deterministic vector jsverbs are better for CI.

### What was tricky to build

The tricky part was designing a CI plan that does not assume the developer workstation layout. The local xgoja vector spec uses sibling `replace` paths for `geppetto` and `go-go-goja`, which is useful during active workspace development but not portable to a clean single-repository GitHub Actions checkout. The design therefore separates the local smoke target from the CI portability strategy and recommends a CI-specific spec with released versions when possible.

### What warrants a second pair of eyes

- Whether the project should add a CI-specific `xgoja-vectors.ci.yaml`, or whether the workflow should checkout sibling repositories while the related modules are still evolving.
- Whether FAISS source build time is acceptable for a scheduled/manual workflow, or whether a prebuilt cache/container should be introduced after the first implementation.
- Whether `xgoja-smoke-vectors` should eventually be required on pull requests touching `pkg/**` or `cmd/goja-bleve/**`.

### What should be done in the future

- Implement the Makefile targets from the design guide.
- Add the optional FAISS workflow and observe its runtime/stability before making it a required PR check.
- Update README, quickstart, and FAISS playbook once the targets exist.

### Code review instructions

- Start with `design-doc/01-generated-xgoja-vector-smoke-and-faiss-ci-implementation-guide.md`.
- Check the proposed Makefile snippets against the current `Makefile`.
- Check the proposed workflow against `.github/workflows/push.yml` and `.github/workflows/dependency-scanning.yml` so the FAISS job stays separate from default CI.
- Validate future implementation with:
  - `make test-vectors`
  - `make xgoja-smoke-vectors`
  - `GOWORK=off go test ./...`

### Technical details

Primary local references:

- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/push.yml`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/dependency-scanning.yml`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.yaml`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/vector.js`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/rag.js`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md`

## Step 2: Add and validate local generated xgoja vector smoke targets

I implemented the first code slice from the design guide: local Makefile targets that build the vector-enabled xgoja host and run the deterministic vector JavaScript smoke verbs. This gives contributors one command that validates the generated binary layer, not only the package-level vector tests.

The target passed locally with the existing FAISS installation. The generated binary built from `cmd/goja-bleve/xgoja-vectors.yaml`, and both `vector knn` and `vector hybrid` returned JSON results with `vectorSupport: true` and `chunk-1` ranked first.

### Prompt Context

**User prompt (verbatim):** "go ahead, step by step, commit at appropriat eintervals, keep a frequent detailed diary as you go"

**Assistant interpretation:** Start implementing the ticket incrementally, commit focused changes at useful checkpoints, and keep the diary updated after each meaningful step.

**Inferred user intent:** Make steady, reviewable progress while preserving a clear continuation trail and avoiding one large opaque change.

**Commit (code):** 10e75adee33d149391e5966310228fbee8160453 — "Add generated xgoja vector smoke target"

### What I did

- Edited `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile`.
- Added `xgoja-build-vectors` and `xgoja-smoke-vectors` to `.PHONY`.
- Added configurable variables:
  - `XGOJA_VERSION ?= v0.8.3`
  - `XGOJA_VECTOR_SPEC ?= xgoja-vectors.yaml`
  - `XGOJA_VECTOR_WORK_DIR ?= /tmp/goja-bleve-vector-work`
- Added `xgoja-build-vectors` to run xgoja against `cmd/goja-bleve/xgoja-vectors.yaml`.
- Added `xgoja-smoke-vectors` to run:
  - `./dist/goja-bleve-vectors vector knn --output json`
  - `./dist/goja-bleve-vectors vector hybrid --output json`
- Ran `make xgoja-smoke-vectors`.

### Why

The package target `make test-vectors` proves the Go package can compile and test under `-tags=vectors`. It does not prove the xgoja YAML spec, generated host, embedded jsverbs, provider mounting, or CLI command path. The new target covers that generated-runtime layer with deterministic smoke commands.

### What worked

The command passed:

```text
make xgoja-smoke-vectors
```

Key output:

```text
validated 21 check(s) for xgoja-vectors.yaml
generated build workspace: /tmp/goja-bleve-vector-work
xgoja build ok: /home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/dist/goja-bleve-vectors
```

The KNN smoke returned `chunk-1` then `chunk-3`, both with `vectorSupport: true`. The hybrid smoke also returned `chunk-1` then `chunk-3` with `scoreMode: "rrf"` and `vectorSupport: true`.

### What didn't work

No command failure occurred in this step.

### What I learned

The current local xgoja vector spec is already sufficient for the developer workstation path. The new Makefile target did not need additional shell `CGO_LDFLAGS` because the spec's `go.env` entry is correctly passed into `go build` by xgoja.

### What was tricky to build

The main trick was preserving unrelated working-tree changes. `Makefile` already had uncommitted release/install placeholder fixes before this step. The new target was added on top of that file, but the commit should stage only the xgoja vector smoke hunks and leave unrelated Makefile hunks unstaged.

### What warrants a second pair of eyes

- Whether `--keep-work` should remain in the default target or become configurable; it is useful for debugging but leaves `/tmp/goja-bleve-vector-work` behind.
- Whether the target should delete `cmd/goja-bleve/dist/goja-bleve-vectors` before rebuilding to make stale-binary failures more obvious.
- Whether the target should assert the JSON output shape or keep command success as the smoke assertion.

### What should be done in the future

- Add documentation for the new Makefile target.
- Add a CI workflow for package vector tests.
- Decide whether to add a CI-compatible xgoja spec or sibling checkouts before enabling generated xgoja smoke in GitHub Actions.

### Code review instructions

- Review `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile` near the `test-vectors` target.
- Validate with:
  - `make xgoja-smoke-vectors`
  - `make test-vectors`

### Technical details

The new smoke path is:

```text
make xgoja-smoke-vectors
  -> make xgoja-build-vectors
  -> cd cmd/goja-bleve
  -> go run github.com/go-go-golems/go-go-goja/cmd/xgoja@v0.8.3 build -f xgoja-vectors.yaml
  -> ./dist/goja-bleve-vectors vector knn --output json
  -> ./dist/goja-bleve-vectors vector hybrid --output json
```


## Step 3: Document the generated xgoja vector smoke target

I updated the public documentation so the new Makefile target is discoverable from the README, quickstart, documentation index, and FAISS/xgoja playbook. The docs now distinguish package vector tests from generated xgoja smoke tests instead of asking readers to copy the long xgoja build command by hand.

The validation commands still pass after the documentation update. This confirms the documented target is not merely aspirational: it builds the vector host from the checked-in YAML spec and runs the KNN and hybrid jsverbs against the local FAISS setup.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Continue implementing incrementally, update docs as soon as a working target exists, and commit this as a separate reviewable slice.

**Inferred user intent:** Keep the repository usable for future contributors by making every new workflow discoverable and backed by validation evidence.

**Commit (code):** TBD — pending documentation commit.

### What I did

- Updated `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/README.md`:
  - replaced the long manual xgoja vector build recipe with `make xgoja-smoke-vectors`,
  - kept the equivalent direct binary commands for readers who want to inspect the generated host,
  - added `make xgoja-smoke-vectors` to development validation.
- Updated `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/quickstart.md`:
  - explained that `make test-vectors` validates package-level vector tests,
  - added `make xgoja-smoke-vectors` for generated host and jsverb smoke coverage.
- Updated `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md`:
  - renamed the xgoja section to build-and-smoke-test,
  - documented the Makefile target and its configurable variables,
  - updated the future CI note and final checklist.
- Updated `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/README.md` to mention both vector validation targets.
- Ran:
  - `make test-vectors`
  - `make xgoja-smoke-vectors`

### Why

Once a workflow target exists, docs should point to it instead of duplicating a long command sequence. This reduces drift and makes it easier for an intern or CI workflow author to run the same validation that maintainers run locally.

### What worked

Both validation commands passed:

```text
make test-vectors
make xgoja-smoke-vectors
```

`make test-vectors` passed with `ok github.com/go-go-golems/goja-bleve/pkg`. The generated smoke target built `dist/goja-bleve-vectors`, then returned successful JSON for both vector commands.

### What didn't work

No validation failure occurred in this step.

### What I learned

The docs previously mixed generated-host validation with raw xgoja invocation. Replacing the raw path with the Makefile target makes the user-facing path shorter while still preserving enough direct command detail for troubleshooting.

### What was tricky to build

The tricky part was avoiding over-documenting the target as a CI guarantee. The current local xgoja spec still contains sibling `replace` paths, so the docs describe `make xgoja-smoke-vectors` as the local generated-host smoke path and keep CI language cautious until a clean-checkout spec or sibling checkout strategy exists.

### What warrants a second pair of eyes

- Whether README should keep any full raw xgoja build command, or whether the FAISS playbook is the only right place for that detail.
- Whether the quickstart should mention that FAISS must already be installed before either vector target can pass.

### What should be done in the future

- Add the optional FAISS workflow.
- Decide how to make generated xgoja smoke portable in CI.
- Re-upload or refresh the reMarkable bundle if the ticket docs need to include implementation updates, not just the initial design.

### Code review instructions

- Review `README.md` and `docs/faiss-xgoja-playbook.md` first to ensure the target names are clear.
- Validate the documented commands:
  - `make test-vectors`
  - `make xgoja-smoke-vectors`

### Technical details

The documentation now treats the two vector targets as distinct:

```text
make test-vectors
  -> Go package tests under ./pkg with -tags=vectors

make xgoja-smoke-vectors
  -> xgoja build from cmd/goja-bleve/xgoja-vectors.yaml
  -> generated dist/goja-bleve-vectors binary
  -> vector knn and vector hybrid jsverbs
```


## Step 4: Add optional FAISS-backed GitHub Actions workflow

I added the first CI implementation slice as a new optional workflow rather than extending the default pull-request pipeline. The workflow builds the Bleve-compatible FAISS fork on an Ubuntu runner, verifies the installed native headers and shared libraries, and then runs `make test-vectors`.

Generated xgoja smoke remains opt-in in the workflow because the current local xgoja vector spec contains sibling workspace `replace` paths. The workflow exposes a manual `run-xgoja-smoke` input for experiments, but the scheduled/default path only runs package vector tests after installing FAISS.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Continue the implementation sequence by adding the optional FAISS CI workflow after the local target and docs are in place.

**Inferred user intent:** Get real CI coverage for vector-tagged package tests while avoiding unstable required checks or clean-runner failures from local xgoja replacement paths.

**Commit (code):** TBD — pending workflow commit.

### What I did

- Added `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml`.
- Configured triggers:
  - `workflow_dispatch`, with an optional `run-xgoja-smoke` input,
  - weekly schedule at `0 4 * * 0`.
- Kept permissions minimal with `contents: read`.
- Added steps to:
  - checkout the repository,
  - set up Go from `go.mod`,
  - install FAISS build dependencies,
  - clone `blevesearch/faiss`,
  - checkout commit `fff814d`,
  - configure FAISS with C API and shared libraries,
  - build `faiss` and `faiss_c` with `-j2`,
  - install/copy the shared libraries and run `ldconfig`,
  - verify headers and shared libraries,
  - run `make test-vectors`.
- Added an opt-in generated xgoja smoke step guarded by:
  - `github.event_name == 'workflow_dispatch'`
  - `inputs.run-xgoja-smoke == 'true'`
- Updated `docs/faiss-xgoja-playbook.md` with the new optional workflow behavior.
- Parsed the workflow with Ruby/Psych as a syntax sanity check.

### Why

The vector package tests need a clean-runner native dependency check. Adding the job as optional gives maintainers scheduled and manual signal without making all pull requests depend on FAISS build time or runner-specific C++ toolchain behavior.

### What worked

The workflow file was created and the YAML parsed locally with:

```text
ruby -e 'require "yaml"; p YAML.load_file(".github/workflows/vector-faiss.yml").keys'
```

The parser returned keys for the workflow document. No local GitHub Actions execution was attempted.

### What didn't work

No local workflow runner execution was attempted, so the FAISS build commands still need to be proven by GitHub Actions. Ruby's YAML 1.1 parser reads the `on` key as boolean `true`, which is a known YAML-version mismatch for GitHub workflows rather than a GitHub Actions syntax problem.

### What I learned

A first CI version can safely validate `make test-vectors` without solving xgoja clean-checkout portability in the same step. The generated xgoja smoke path can remain manual/experimental until the replacement strategy is chosen.

### What was tricky to build

The tricky part was exposing generated xgoja smoke without accidentally making the scheduled workflow fail on clean runners. The solution was to include the step but guard it behind a manual input whose default is `false`. That preserves an experimentation hook while keeping the normal CI signal focused on package vector tests.

### What warrants a second pair of eyes

- Whether `libopenblas-dev` and `libgomp1` are sufficient for the FAISS build on the current `ubuntu-latest` image.
- Whether `make -C build -j2 faiss faiss_c` is the right concurrency/runtime tradeoff.
- Whether the workflow should cache `/tmp/faiss/build` or installed artifacts after the first successful runs.
- Whether the manual xgoja smoke input should be removed until a CI-compatible spec exists, to avoid user confusion.

### What should be done in the future

- Run the workflow manually on GitHub and record the result.
- If it passes consistently, consider adding a pull-request trigger with path filters.
- Add a CI-compatible `xgoja-vectors.ci.yaml` or sibling checkout strategy before enabling generated xgoja smoke by default.

### Code review instructions

- Review `.github/workflows/vector-faiss.yml` from top to bottom.
- Compare the FAISS build commands with `docs/faiss-xgoja-playbook.md`.
- Confirm the workflow is not required on pull requests.
- Validate local prerequisites with:
  - `make test-vectors`
  - `make xgoja-smoke-vectors`

### Technical details

The workflow's normal path is:

```text
workflow_dispatch or schedule
  -> checkout goja-bleve
  -> setup Go from go.mod
  -> apt install cmake/g++/make/libopenblas-dev/libgomp1
  -> clone blevesearch/faiss
  -> checkout fff814d
  -> cmake with FAISS_ENABLE_C_API=ON and BUILD_SHARED_LIBS=ON
  -> make -j2 faiss faiss_c
  -> sudo make install + copy shared libs + ldconfig
  -> verify headers/libs
  -> make test-vectors
```
