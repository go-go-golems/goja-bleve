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

**Commit (code):** 50e7eacf392f8ede73054522fafba72354338799 — "Document xgoja vector smoke target"

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

**Commit (code):** 257b91e6083a6c22acfb7b49401b25056d4ef556 — "Add optional FAISS vector workflow"

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


## Step 5: Add a CI-compatible xgoja vector spec and enable generated smoke in the FAISS workflow

I removed the main limitation from the first workflow version by adding a CI-specific xgoja vector spec. The new spec keeps the `goja-bleve` package replaced to the current checkout, but it uses released `geppetto` and `go-go-goja` versions instead of local sibling workspace paths. That makes generated xgoja vector smoke viable on a clean GitHub Actions checkout.

After adding the spec, I changed the FAISS workflow so generated xgoja smoke runs by default after `make test-vectors`. Locally, the CI spec built successfully and both deterministic vector jsverbs passed.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Continue hardening step by step by making the generated xgoja smoke target usable from CI, not just from the local workspace.

**Inferred user intent:** Close the gap between local generated-host validation and CI validation so the new workflow covers both requested hardening items.

**Commit (code):** 0b1661b4345b64ee4575688d3e412e2228c82f25 — "Enable generated vector smoke in FAISS workflow"

### What I did

- Added `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.ci.yaml`.
- Copied the vector build settings from `xgoja-vectors.yaml`:
  - `go.tags: [vectors]`,
  - `go.ldflags: [-r, /usr/local/lib]`,
  - `go.env.CGO_LDFLAGS` with `-lfaiss_c -lfaiss -lstdc++ -lm`.
- Kept the current repository replacement:
  - `goja-bleve` version `v0.0.0`, `replace: ../..`.
- Replaced local sibling paths with released versions:
  - `geppetto v0.11.7`,
  - `go-go-goja-core v0.8.3`,
  - `go-go-goja-host v0.8.3`.
- Updated `.github/workflows/vector-faiss.yml` to run:
  - `make test-vectors`,
  - `make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml`.
- Updated `docs/faiss-xgoja-playbook.md` to explain the local spec vs CI spec split.
- Validated locally with:
  - `ruby -e 'require "yaml"; p YAML.load_file(".github/workflows/vector-faiss.yml").keys'`,
  - `make test-vectors`,
  - `make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml XGOJA_VECTOR_WORK_DIR=/tmp/goja-bleve-vector-ci-spec-work`.

### Why

The local `xgoja-vectors.yaml` is correct for active multi-repository development, but its sibling `replace` paths are not a clean CI contract. A separate CI spec makes the intended environment explicit: use the pull request checkout for `goja-bleve`, and use released dependency versions for external providers.

### What worked

The CI spec validated with xgoja:

```text
validated 18 check(s) for xgoja-vectors.ci.yaml
generated build workspace: /tmp/goja-bleve-vector-ci-spec-work
xgoja build ok: .../cmd/goja-bleve/dist/goja-bleve-vectors
```

Both smoke commands passed and returned `vectorSupport: true`:

```text
./dist/goja-bleve-vectors vector knn --output json
./dist/goja-bleve-vectors vector hybrid --output json
```

`make test-vectors` also passed.

### What didn't work

No command failure occurred in this step. The only caveat is that local validation still used the machine's existing FAISS install; the full FAISS build/install path must be proven by running the GitHub Actions workflow.

### What I learned

The released `geppetto v0.11.7` and `go-go-goja v0.8.3` versions are sufficient for the generated vector smoke host. The CI spec therefore does not need sibling checkouts for this hardening slice.

### What was tricky to build

The subtle part was choosing which replacement to keep. The CI spec must keep `goja-bleve` replaced to `../..`; otherwise the generated nested module would fetch the last published module instead of testing the current checkout. But it should not replace `geppetto` or `go-go-goja` to sibling paths because those paths do not exist on a single-repo runner.

### What warrants a second pair of eyes

- Whether maintaining two mostly identical xgoja specs is acceptable, or whether xgoja should grow a parameterized replace strategy in the future.
- Whether the CI workflow should run generated smoke by default immediately, or whether it should stay package-only until the first workflow run proves FAISS build time is acceptable.
- Whether the output binary path `dist/goja-bleve-vectors` could collide with local-spec builds in developer machines; both specs intentionally produce the same executable name.

### What should be done in the future

- Run the new workflow on GitHub and record the result.
- If runtime is high, add caching for the FAISS build or replace the source build with a maintained prebuilt artifact.
- If released dependency versions drift, update `xgoja-vectors.ci.yaml` alongside `cmd/goja-bleve/go.mod` and the local xgoja spec.

### Code review instructions

- Compare `cmd/goja-bleve/xgoja-vectors.yaml` and `cmd/goja-bleve/xgoja-vectors.ci.yaml`.
- Ensure only local sibling `replace` paths were removed from the CI spec.
- Review `.github/workflows/vector-faiss.yml` and confirm it uses `XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml`.
- Validate locally with:
  - `make test-vectors`,
  - `make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml XGOJA_VECTOR_WORK_DIR=/tmp/goja-bleve-vector-ci-spec-work`.

### Technical details

The CI spec dependency policy is:

```text
goja-bleve provider
  -> replace ../.. so CI tests the current checkout

geppetto provider
  -> released module v0.11.7

go-go-goja core/host providers
  -> released module v0.8.3
```


## Step 6: Final local validation and handoff notes

I ran the final local validation set after the Makefile target, documentation updates, workflow, and CI-specific xgoja spec were committed. The repository now has a local generated xgoja vector smoke path and an optional GitHub Actions workflow that should validate both package vector tests and generated vector smoke after building FAISS.

This step did not change runtime behavior. It records the commands that passed and the remaining caveat: the new GitHub Actions workflow still needs to be executed on GitHub to prove the Ubuntu runner's FAISS source-build path.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Finish this implementation slice with validation, diary updates, and a clear continuation point.

**Inferred user intent:** Leave the branch in a reviewable state with enough evidence for the next person to understand what was changed and what remains unproven.

**Commit (code):** N/A — validation and diary handoff only.

### What I did

- Ran the vector package test target:
  - `make test-vectors`
- Ran generated xgoja vector smoke using the CI spec:
  - `make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml XGOJA_VECTOR_WORK_DIR=/tmp/goja-bleve-vector-ci-spec-work`
- Ran the default package test suite:
  - `GOWORK=off go test ./...`
- Checked diary commit references and updated Step 5 with its commit hash.

### Why

The final validation should cover the three important paths touched by this work: vector-tag package tests, generated xgoja vector smoke, and the default non-vector package test suite.

### What worked

All local validation commands passed.

`make test-vectors` returned:

```text
ok  github.com/go-go-golems/goja-bleve/pkg  0.143s
```

The CI-spec generated xgoja smoke built successfully and both commands returned JSON results with `vectorSupport: true`:

```text
validated 18 check(s) for xgoja-vectors.ci.yaml
xgoja build ok: .../cmd/goja-bleve/dist/goja-bleve-vectors
```

`GOWORK=off go test ./...` passed for the module.

### What didn't work

The GitHub Actions workflow was not executed in this local session. The workflow's FAISS build commands are based on the documented local playbook, but they still need a real runner execution.

### What I learned

The CI spec works locally without sibling replacements. That means the workflow can run generated xgoja vector smoke by default instead of keeping it as a manual experiment.

### What was tricky to build

The remaining tricky point is outside local validation: FAISS source builds can fail or run slowly on hosted CI even when the same commands are correct locally. The workflow uses `-j2` to reduce memory pressure, but actual runtime and stability are unknown until the first scheduled or manual run.

### What warrants a second pair of eyes

- The exact FAISS build/install commands in `.github/workflows/vector-faiss.yml`.
- Whether `go 1.26.4` in `go.mod` is accepted by `actions/setup-go@v6` on current GitHub runners.
- Whether the historical `cmd/XXX` package reported by `GOWORK=off go test ./...` should be cleaned up in a separate release-plumbing change.

### What should be done in the future

- Push the branch and run the `Vector FAISS Smoke` workflow manually.
- Record workflow runtime and failure/success in this ticket if further work is needed.
- If stable, consider adding pull-request path filters later.

### Code review instructions

Start with these files:

- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.ci.yaml`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml`
- `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md`

Validate with:

```text
make test-vectors
make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml XGOJA_VECTOR_WORK_DIR=/tmp/goja-bleve-vector-ci-spec-work
GOWORK=off go test ./...
```

### Technical details

Final local validation proves local command correctness, not hosted-runner correctness. The workflow still needs GitHub execution because it is the only place that will exercise:

```text
apt install native dependencies
clone blevesearch/faiss
cmake FAISS with C API and shared libs
make -j2 faiss faiss_c
sudo make install
sudo ldconfig
make test-vectors
make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml
```

## Step 7: Push branch and attempt remote workflow dispatch

I pushed the implementation branch after local validation. The first push attempt failed because the pre-push hook runs a GoReleaser snapshot against the current working tree, and there were unrelated unstaged release-plumbing edits in `.goreleaser.yaml` and `Makefile`. I stashed only those unrelated file changes, pushed successfully with a clean working tree, then restored the stash.

After the branch was pushed, I attempted to trigger the new `Vector FAISS Smoke` workflow with `gh workflow run`. GitHub returned `HTTP 404: Not Found` for the workflow file. This is expected for a newly added workflow that is not present on the repository default branch yet; GitHub generally only exposes `workflow_dispatch` for workflows known on the default branch.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Push the committed incremental work when ready and try to exercise the remote CI path if possible.

**Inferred user intent:** Move beyond local validation and get as close as possible to proving the GitHub Actions workflow.

**Commit (code):** N/A — push/remote-dispatch diary only.

### What I did

- Attempted:
  - `git push origin task/rag-evaluation-system`
- The first push failed in the local pre-push hook during the `release` step.
- Temporarily stashed unrelated dirty files:
  - `.goreleaser.yaml`
  - `Makefile`
- Re-ran the push successfully.
- Restored the unrelated stash.
- Attempted remote workflow dispatch:
  - `gh workflow run vector-faiss.yml --ref task/rag-evaluation-system`

### Why

The branch needed to be pushed so the new workflow file exists on GitHub. A manual workflow run would be the best validation of the FAISS source-build path, because local validation cannot prove hosted-runner behavior.

### What worked

The branch pushed successfully after temporarily stashing unrelated dirty release-plumbing edits. The pre-push hook passed on the clean tree:

```text
✔️ test
✔️ lint
✔️ release
To ssh://github.com/go-go-golems/goja-bleve
 * [new branch]      task/rag-evaluation-system -> task/rag-evaluation-system
```

The unrelated working-tree edits were restored afterward.

### What didn't work

The first push attempt failed because GoReleaser read the unrelated dirty `.goreleaser.yaml` change and tried to build `./cmd/goja-bleve`, which does not currently exist:

```text
build failed: exit status 1: main module (github.com/go-go-golems/goja-bleve) does not contain package github.com/go-go-golems/goja-bleve/cmd/goja-bleve
make: *** [Makefile:61: goreleaser] Error 1
error: failed to push some refs to 'ssh://github.com/go-go-golems/goja-bleve'
```

Manual workflow dispatch failed:

```text
gh workflow run vector-faiss.yml --ref task/rag-evaluation-system
HTTP 404: Not Found (https://api.github.com/repos/go-go-golems/goja-bleve/actions/workflows/vector-faiss.yml)
```

Likely cause: the workflow file is new on this branch and not yet present on the default branch, so GitHub does not expose it as a dispatchable workflow by filename.

### What I learned

The unrelated release-plumbing edits affect local pre-push behavior even when they are not staged or committed. For this branch, stashing them was safer than committing unrelated changes or bypassing hooks.

### What was tricky to build

The tricky part was separating the branch's committed implementation from the dirty local environment. The clean-tree push showed that the committed vector CI work passes the repository pre-push hook. The dirty-tree failure should not be attributed to the vector CI commits.

### What warrants a second pair of eyes

- Whether the unrelated `.goreleaser.yaml`/`Makefile` release-plumbing changes should be kept, fixed, or discarded separately.
- Whether the new workflow should be merged first and then manually dispatched from default branch, or tested through a PR by temporarily adding another trigger.

### What should be done in the future

- After the workflow file exists on the default branch, run `gh workflow run vector-faiss.yml` again.
- Alternatively, open a PR and let maintainers decide whether to add a temporary branch trigger for first-run testing.

### Code review instructions

- Do not review the unrelated dirty `.goreleaser.yaml` and `Makefile` release-plumbing edits as part of this vector CI slice.
- Confirm the pushed commits end at `fa5587a` on `origin/task/rag-evaluation-system`.

### Technical details

The push workaround was:

```text
git stash push -m "temporary unrelated release plumbing before vector CI push" -- .goreleaser.yaml Makefile
git push origin task/rag-evaluation-system
git stash pop
```

## Step 8: Upload updated implementation diary bundle to reMarkable

I uploaded an updated reMarkable bundle after the implementation and validation diary entries were committed and pushed. This gives the reMarkable copy the implementation trail, not only the initial design document.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Keep the ticket documentation and external review copy current as the implementation progresses.

**Inferred user intent:** Preserve an up-to-date review artifact for reading outside the repository.

**Commit (code):** N/A — documentation delivery only.

### What I did

- Uploaded a bundle containing:
  - the design guide,
  - the updated investigation diary.
- Remote path:
  - `/ai/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI`
- Uploaded filename:
  - `goja bleve vector ci implementation update.pdf`

### Why

The first upload happened before implementation. The updated bundle includes the step-by-step implementation notes, validation commands, push attempt, and workflow-dispatch caveat.

### What worked

The upload command returned:

```text
OK: uploaded goja bleve vector ci implementation update.pdf -> /ai/2026/06/07/RAGEVAL-GOJA-BLEVE-VECTOR-CI
```

### What didn't work

N/A

### What I learned

N/A

### What was tricky to build

N/A

### What warrants a second pair of eyes

N/A

### What should be done in the future

If the GitHub Actions workflow is run and produces useful logs or failures, append another diary step and upload a refreshed bundle.

### Code review instructions

N/A

### Technical details

Upload command used `remarquee upload bundle` with `--toc-depth 2` and `--non-interactive`.

## Step 9: Fix GoReleaser and install targets for nested xgoja module

I investigated the release failure reported after the branch started using the real `goja-bleve` binary name. The failing command was not caused by the package name inside the generated xgoja module. The root cause was that `cmd/goja-bleve` is its own nested Go module, so a root-module build target like `main: ./cmd/goja-bleve` asks the root module to build a package it does not contain.

The fix is to make GoReleaser build from the nested module directory with `dir: cmd/goja-bleve` and `main: .`. I also updated the Makefile `install` target to build from inside the nested module instead of using `go build ./cmd/goja-bleve` from the root module.

### Prompt Context

**User prompt (verbatim):** "• building                                       binary=dist/linux_amd64/goja-bleve-linux_linux_amd64_v1/goja-bleve
  ⨯ release failed after 0s                         
    error=
    │ build failed: exit status 1: main module (github.com/go-go-golems/goja-bleve) does not contain package github.com/go-go-golems/goja-bleve/cmd/goja-bleve
    target=linux_amd64_v1
make: *** [Makefile:61: goreleaser] Error 1

We need to update the packakge name of the go.mod generated with xgoja maybe?"

**Assistant interpretation:** Diagnose and fix the GoReleaser failure caused by trying to release the generated xgoja binary from the wrong module boundary.

**Inferred user intent:** Make the release/pre-push GoReleaser snapshot work with the generated `cmd/goja-bleve` binary and real project naming.

**Commit (code):** e9a545820c5e4e3ceb60a351555fb68ae2448455 — "Fix release build for generated xgoja binary"

### What I did

- Confirmed the failing root-module build reproduces:
  - `GOWORK=off go build ./cmd/goja-bleve`
- Confirmed the nested-module build works:
  - `cd cmd/goja-bleve && GOWORK=off go build .`
- Updated `.goreleaser.yaml`:
  - Linux build now uses `dir: cmd/goja-bleve`, `main: .`.
  - Darwin build now uses `dir: cmd/goja-bleve`, `main: .`.
- Updated `Makefile` `install` target:
  - builds with `cd cmd/goja-bleve && GOWORK=off go build -o ../../dist/goja-bleve .`.
- Ran:
  - `make goreleaser`

### Why

The generated xgoja host lives in `cmd/goja-bleve`, which has its own `go.mod`. The root module cannot import or build packages under a nested module by path. GoReleaser must enter the nested module before invoking `go build`.

### What worked

`make goreleaser` passed after the change:

```text
building binary=dist/linux_amd64/goja-bleve-linux_linux_amd64_v1/goja-bleve
archiving name=dist/linux_amd64/goja-bleve_0.0.2-next_linux_amd64.tar.gz
creating package=goja-bleve format=rpm
creating package=goja-bleve format=deb
release succeeded after 13s
```

### What didn't work

The root-module build is still invalid, as expected:

```text
GOWORK=off go build ./cmd/goja-bleve
main module (github.com/go-go-golems/goja-bleve) does not contain package github.com/go-go-golems/goja-bleve/cmd/goja-bleve
```

This is not a generated-module-name problem; it is a nested-module build-context problem.

### What I learned

The `module example.com/generated/goja-bleve-vectors` line in `cmd/goja-bleve/go.mod` is not what makes GoReleaser fail. GoReleaser failed before entering that module. Once GoReleaser is told to use `dir: cmd/goja-bleve`, `main: .` builds successfully.

### What was tricky to build

The confusing symptom is that the generated module path looks suspicious, but changing it is not the necessary fix. The important invariant is the working directory for `go build`: root-module package paths stop at nested `go.mod` boundaries, while `cd cmd/goja-bleve && go build .` uses the generated module correctly.

### What warrants a second pair of eyes

- Whether xgoja should generate a stable module path such as `github.com/go-go-golems/goja-bleve/cmd/goja-bleve` instead of `example.com/generated/...` for aesthetics and traceability.
- Whether the repository should remove the old `cmd/XXX` stub now that release packaging targets `cmd/goja-bleve`.
- Whether `Makefile install` should fail early with a clearer error when `which goja-bleve` returns empty.

### What should be done in the future

- Consider cleaning up `cmd/XXX` in a separate focused commit.
- Consider teaching xgoja specs to set the generated module path explicitly if that is supported or useful.

### Code review instructions

- Review `.goreleaser.yaml` build entries and ensure both Linux and Darwin builds use `dir: cmd/goja-bleve` and `main: .`.
- Review `Makefile` `install` and ensure it builds from inside `cmd/goja-bleve`.
- Validate with:
  - `make goreleaser`

### Technical details

Correct release build shape:

```yaml
builds:
  - id: goja-bleve-linux
    dir: cmd/goja-bleve
    main: .
    binary: goja-bleve
```

Correct manual build shape:

```bash
cd cmd/goja-bleve
GOWORK=off go build .
```

## Step 10: Fix invalid release workflow disabled reusable job

I fixed the release workflow validation error caused by a disabled reusable `publish-docs` job. GitHub validates reusable-workflow permission requirements even when the job has `if: ${{ false }}`, so the workflow was invalid because the called workflow requested `id-token: write` while the caller allowed `id-token: none`.

The fix was to remove the disabled `uses:` job from the active workflow and leave it as explanatory comments. This keeps the release workflow valid without granting OIDC permissions to all release runs just to satisfy a disabled template.

### Prompt Context

**User prompt (verbatim):** "1 error
Invalid workflow file: .github/workflows/release.yaml#L80
The workflow is not valid. .github/workflows/release.yaml (Line: 80, Col: 3): Error calling workflow 'go-go-golems/infra-tooling/.github/workflows/publish-docsctl.yml@main'. The workflow is requesting 'id-token: write', but is only allowed 'id-token: none'."

**Assistant interpretation:** Diagnose and fix the GitHub Actions workflow validation error in `release.yaml`.

**Inferred user intent:** Make the release workflow valid without accidentally enabling the disabled docs publishing path or granting unnecessary permissions.

**Commit (code):** f13dc6b46ead0e8731d06fb6c6cfb48334156657 — "Fix release workflow docs publishing template"

### What I did

- Edited `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/release.yaml`.
- Removed the active disabled `publish-docs` reusable-workflow job.
- Kept a commented template explaining how to re-enable docs publishing later.
- Added a warning that GitHub validates reusable workflow permissions even when a job is guarded with `if: ${{ false }}`.
- Parsed the workflow locally with Ruby YAML and confirmed active jobs are now only:
  - `goreleaser-linux`
  - `goreleaser-darwin`
  - `goreleaser-merge`

### Why

Adding `id-token: write` would also make the workflow valid, but it would grant OIDC permissions for a job that is intentionally disabled. Removing the disabled reusable job is safer and more explicit.

### What worked

Local YAML parsing succeeded, and the active job list no longer includes `publish-docs`:

```text
["goreleaser-linux", "goreleaser-darwin", "goreleaser-merge"]
```

### What didn't work

No remote GitHub validation was run in this step.

### What I learned

GitHub Actions validates called reusable workflow permissions at workflow-parse time. An `if: false` guard prevents execution but does not prevent reusable-workflow permission validation.

### What was tricky to build

The tricky part is that the workflow looked disabled, but GitHub still rejected it. The safe pattern is not to keep disabled `uses:` jobs around as templates unless the caller grants every permission the reusable workflow requests.

### What warrants a second pair of eyes

- Whether the commented template should be moved to a separate documentation file instead of living in `release.yaml`.
- Whether future docs publishing should be enabled in this repository at all.

### What should be done in the future

If docs publishing is enabled later, add a real `publish-docs` job and explicitly grant its required permissions in the top-level or job-level `permissions` block.

### Code review instructions

- Review `.github/workflows/release.yaml` around the former `publish-docs` block.
- Confirm there is no active reusable `publish-docs` job unless the required permissions are also granted.

### Technical details

The invalid pattern was:

```yaml
publish-docs:
  if: ${{ false && startsWith(github.ref, 'refs/tags/v') }}
  uses: go-go-golems/infra-tooling/.github/workflows/publish-docsctl.yml@main
```

The corrected workflow keeps only comments and no active reusable job.
