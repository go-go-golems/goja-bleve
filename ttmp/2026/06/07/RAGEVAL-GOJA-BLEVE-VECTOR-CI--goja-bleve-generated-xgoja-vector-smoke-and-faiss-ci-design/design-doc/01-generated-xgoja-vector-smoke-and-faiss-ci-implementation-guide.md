---
Title: Generated xgoja vector smoke and FAISS CI implementation guide
Ticket: RAGEVAL-GOJA-BLEVE-VECTOR-CI
Status: active
Topics:
    - goja
    - xgoja
    - bleve
    - ci
    - vector-search
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: .github/workflows/dependency-scanning.yml
      Note: Existing security workflow used as CI context for the optional vector workflow
    - Path: .github/workflows/push.yml
      Note: Existing default Go test pipeline that should remain separate from optional FAISS CI
    - Path: Makefile
      Note: Defines current test-vectors target and will receive generated xgoja smoke targets
    - Path: cmd/goja-bleve/jsverbs/vector.js
      Note: Deterministic generated-runtime KNN and hybrid smoke commands
    - Path: cmd/goja-bleve/xgoja-vectors.yaml
      Note: Vector generated-host spec with build tags
    - Path: docs/faiss-xgoja-playbook.md
      Note: Operational FAISS and xgoja vector-linking playbook referenced by CI design
ExternalSources: []
Summary: Design and intern-oriented implementation guide for adding generated xgoja vector smoke targets and optional FAISS-backed CI to goja-bleve.
LastUpdated: 2026-06-07T13:40:00-04:00
WhatFor: Use this before implementing Makefile targets and GitHub Actions jobs that validate goja-bleve's generated xgoja vector binary and FAISS-backed vector tests.
WhenToUse: 'Read when hardening goja-bleve after PR #1, especially when working on generated xgoja smoke coverage, FAISS installation in CI, or vector build troubleshooting.'
---


# Generated xgoja vector smoke and FAISS CI implementation guide

## Executive summary

This document describes the next hardening slice for `goja-bleve`: add a generated xgoja smoke target and an optional FAISS-backed CI job. The target proves that the checked-in xgoja spec can build a vector-enabled generated binary and run the embedded JavaScript smoke verbs. The CI job proves the same vector path on a clean runner by installing or caching FAISS, running `make test-vectors`, and optionally running the generated xgoja vector smoke target.

The work is deliberately focused. It does not add new public search APIs. It hardens the path that already exists:

```text
xgoja-vectors.yaml
  -> xgoja build
  -> generated goja-bleve-vectors binary
  -> embedded jsverbs
  -> Bleve vector mapping + FAISS KNN
  -> hybrid RRF smoke result
```

The intern implementing this work should understand five parts of the system:

1. `goja-bleve` is the native module that exposes Bleve as `require("bleve")`.
2. Bleve vector support is behind the `vectors` Go build tag and requires FAISS native libraries.
3. xgoja reads `xgoja.yaml` files, generates a Go host, and runs `go build` with declared tags, ldflags, and build-time environment variables.
4. jsverbs are JavaScript functions embedded into the generated binary and mounted as CLI commands.
5. GitHub Actions runners do not have the local FAISS install, sibling checkout layout, or `/usr/local/lib` state from the development machine unless the workflow creates them.

The recommended implementation is two-stage:

- First add local Makefile targets that build and run the generated vector host using the current local workspace assumptions.
- Then add a non-required GitHub Actions workflow that builds FAISS and runs vector tests on `workflow_dispatch` and a schedule. Once it is stable, decide whether to run it on pull requests or keep it opt-in.

## Problem statement

The default CI already proves that normal non-vector builds work. It runs package tests, linting, logcopter checks, dependency scanning, Go vulnerability checks, GoSec, and CodeQL. That is enough to protect the default module path where `bleve.vectorSupport` is false and vector APIs return clear `-tags=vectors` errors.

It does not prove the vector path.

The vector path has additional dependencies and failure modes:

- `-tags=vectors` must be passed to Go.
- `libfaiss_c.so` must be available.
- `libfaiss.so` must also be linked, because this local FAISS C API build leaves C++ FAISS symbols to be resolved by the final executable.
- `libstdc++` and `libm` must be included through `CGO_LDFLAGS`.
- The runtime loader must be able to find FAISS shared libraries after build.
- xgoja must pass `go.env.CGO_LDFLAGS` from the YAML spec into `go build`.
- jsverbs must remain mounted with their fields and arguments after provider runtime sections are appended.
- The generated binary must include the same provider modules used by the RAG tool: `bleve`, `geppetto`, `path`, `yaml`, and `fs`.

A plain package test cannot catch all of those. A generated xgoja smoke target catches the generated-host layer. A FAISS CI job catches the clean-runner native-library layer.

## Current-state architecture

### The module layer

The native package under `pkg/` registers the `bleve` module. JavaScript sees builders such as `bleve.field()`, `bleve.mapping()`, `bleve.search()`, `bleve.create(path)`, and `bleve.open(path)`. The implementation stores Go references on wrapper objects rather than representing Bleve state as plain JavaScript objects.

The vector build uses build-tag split files:

| File | Build tag | Role |
|---|---:|---|
| `pkg/vector_api.go` | `!vectors` | Provides stubs that return clear unavailable errors. |
| `pkg/vector_api_vectors.go` | `vectors` | Calls Bleve vector constructors, attaches KNN clauses, and validates KNN fields. |
| `pkg/vector_support.go` | `!vectors` | Sets `vectorSupportEnabled = false`. |
| `pkg/vector_support_vectors.go` | `vectors` | Sets `vectorSupportEnabled = true`. |

The current vector package test target is:

```makefile
test-vectors:
	GOWORK=off CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1
```

This target validates the package-level vector tests, including vector field construction, KNN ranking, hybrid score fusion, reopened-index mapping behavior, and vector input validation.

### The xgoja generated-host layer

The generated vector host is described by `cmd/goja-bleve/xgoja-vectors.yaml`:

```yaml
name: goja-bleve-vectors
appName: goja-bleve
go:
  tags:
    - vectors
  ldflags:
    - -r
    - /usr/local/lib
  env:
    CGO_LDFLAGS: "-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
target:
  kind: xgoja
  output: dist/goja-bleve-vectors
```

The important xgoja API is `go.env`. In `go-go-goja`, the build spec type is:

```go
type GoSpec struct {
    Version string            `yaml:"version" json:"version"`
    Module  string            `yaml:"module" json:"module"`
    Tags    []string          `yaml:"tags" json:"tags,omitempty"`
    LDFlags []string          `yaml:"ldflags" json:"ldflags,omitempty"`
    Env     map[string]string `yaml:"env" json:"env,omitempty"`
    Imports []GoImportSpec    `yaml:"imports" json:"imports,omitempty"`
}
```

The xgoja build executor threads that map into `go build`:

```go
func GoBuild(ctx context.Context, dir string, output string, tags []string, ldflags []string, env map[string]string) (Result, error) {
    args := []string{"build", "-o", output}
    if len(tags) > 0 {
        args = append(args, "-tags", joinSpace(tags))
    }
    if len(ldflags) > 0 {
        args = append(args, "-ldflags", joinSpace(ldflags))
    }
    args = append(args, ".")
    return run(ctx, dir, env, "go", args...)
}
```

That is the reason the vector xgoja spec can be self-contained. The spec carries `vectors`, `-r /usr/local/lib`, and `CGO_LDFLAGS` together.

### The jsverb layer

The generated binary embeds JavaScript verbs from `cmd/goja-bleve/jsverbs`. The vector smoke verbs live in `vector.js`; the RAG integration verb lives in `rag.js`.

The vector smoke commands are expected to run as:

```bash
./dist/goja-bleve-vectors vector knn --output json
./dist/goja-bleve-vectors vector hybrid --output json
```

The RAG command is heavier because it needs a real embedding profile and provider:

```bash
./dist/goja-bleve-vectors rag index-query \
  --profile-path ./profiles.yaml \
  --embedding-profile assistant \
  privacy
```

The generated smoke target should run `vector knn` and `vector hybrid` by default. It should not require an external embedding provider. The RAG command can remain a manual integration check because it depends on local model/profile availability.

## Proposed solution

Add two Makefile targets and one optional CI workflow.

### Target 1: build the generated vector binary

Add a target that builds `cmd/goja-bleve/xgoja-vectors.yaml` using the current xgoja version:

```makefile
XGOJA_VERSION ?= v0.8.3
XGOJA_VECTOR_WORK_DIR ?= /tmp/goja-bleve-vector-work

xgoja-build-vectors:
	cd cmd/goja-bleve && GOWORK=off go run github.com/go-go-golems/go-go-goja/cmd/xgoja@$(XGOJA_VERSION) build \
		-f xgoja-vectors.yaml \
		--work-dir $(XGOJA_VECTOR_WORK_DIR) \
		--keep-work \
		--xgoja-version $(XGOJA_VERSION)
```

This target validates the generated-host build path. It does not run the binary.

### Target 2: run generated vector smoke verbs

Add a target that depends on the build target and runs the two vector smoke verbs:

```makefile
xgoja-smoke-vectors: xgoja-build-vectors
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector knn --output json
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector hybrid --output json
```

The target should be separate from `test-vectors`. Package tests and generated-host smoke tests validate different layers.

```text
make test-vectors
  -> package-level Go tests under ./pkg with -tags=vectors

make xgoja-smoke-vectors
  -> xgoja build from YAML
  -> generated binary
  -> embedded JavaScript smoke verbs
```

### Optional CI workflow

Add a GitHub Actions workflow that is initially non-required. It should run on `workflow_dispatch` and a weekly schedule. Pull request execution can be added later if setup time is acceptable.

Recommended name:

```text
.github/workflows/vector-faiss.yml
```

Recommended first version:

```yaml
name: Vector FAISS Smoke

on:
  workflow_dispatch:
  schedule:
    - cron: '0 4 * * 0'

jobs:
  vector-faiss:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v6
        with:
          go-version-file: go.mod
          cache: true

      - name: Install FAISS build dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y cmake g++ make libopenblas-dev libgomp1

      - name: Build and install Bleve-compatible FAISS
        run: |
          git clone https://github.com/blevesearch/faiss.git /tmp/faiss
          cd /tmp/faiss
          git checkout fff814d
          cmake -B build \
            -DFAISS_ENABLE_GPU=OFF \
            -DFAISS_ENABLE_C_API=ON \
            -DBUILD_SHARED_LIBS=ON \
            -DFAISS_ENABLE_PYTHON=OFF \
            -DCMAKE_INSTALL_PREFIX=/usr/local \
            -DCMAKE_CXX_FLAGS="-I$PWD" \
            .
          make -C build -j2 faiss faiss_c
          sudo make -C build install
          sudo cp build/faiss/libfaiss.so /usr/local/lib/ || true
          sudo cp build/c_api/libfaiss_c.so /usr/local/lib/ || true
          sudo ldconfig

      - name: Verify FAISS installation
        run: |
          ls -lh /usr/local/lib/libfaiss.so /usr/local/lib/libfaiss_c.so
          ls -lh /usr/local/include/faiss/c_api/IndexBinary_c_ex.h
          ldconfig -p | grep -E 'libfaiss(_c)?\.so'

      - name: Run vector package tests
        run: make test-vectors

      - name: Run generated xgoja vector smoke
        run: make xgoja-smoke-vectors
```

The first version should use `-j2` rather than `-j$(nproc)` to reduce runner memory pressure. If the build is too slow, add caching after correctness is proven.

## Important CI portability issue: local replaces in xgoja specs

The current `xgoja-vectors.yaml` contains local replaces:

```yaml
packages:
  - id: geppetto
    import: github.com/go-go-golems/geppetto/pkg/js/modules/geppetto/provider
    version: v0.0.0
    replace: ../../../geppetto
  - id: go-go-goja-core
    import: github.com/go-go-golems/go-go-goja/pkg/xgoja/providers/core
    replace: ../../../go-go-goja
```

Those paths exist in the development workspace. They will not exist on a clean GitHub Actions checkout of only `goja-bleve`.

There are three possible solutions.

### Option A: create a CI-specific vector spec

Add `cmd/goja-bleve/xgoja-vectors.ci.yaml` that removes sibling `replace` entries and uses released versions for `geppetto` and `go-go-goja`:

```yaml
packages:
  - id: goja-bleve
    import: github.com/go-go-golems/goja-bleve/pkg/xgoja/providers/bleve
    version: v0.0.0
    replace: ../..
  - id: geppetto
    import: github.com/go-go-golems/geppetto/pkg/js/modules/geppetto/provider
    version: v0.11.7
  - id: go-go-goja-core
    import: github.com/go-go-golems/go-go-goja/pkg/xgoja/providers/core
    version: v0.8.3
  - id: go-go-goja-host
    import: github.com/go-go-golems/go-go-goja/pkg/xgoja/providers/host
    version: v0.8.3
```

Then make the target configurable:

```makefile
XGOJA_VECTOR_SPEC ?= xgoja-vectors.yaml

xgoja-build-vectors:
	cd cmd/goja-bleve && GOWORK=off go run github.com/go-go-golems/go-go-goja/cmd/xgoja@$(XGOJA_VERSION) build \
		-f $(XGOJA_VECTOR_SPEC) \
		--work-dir $(XGOJA_VECTOR_WORK_DIR) \
		--keep-work \
		--xgoja-version $(XGOJA_VERSION)
```

CI can run:

```bash
make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml
```

This is the recommended option once the needed provider versions are tagged.

### Option B: checkout sibling repositories in CI

The workflow can reproduce the local workspace layout:

```yaml
- uses: actions/checkout@v6
  with:
    repository: go-go-golems/go-go-goja
    path: ../go-go-goja
- uses: actions/checkout@v6
  with:
    repository: go-go-golems/geppetto
    path: ../geppetto
```

This keeps one spec, but it makes the workflow depend on multiple repository branches and path assumptions. It is useful while waiting for releases, but less stable as a long-term CI contract.

### Option C: only run package vector tests in CI

The workflow can stop after `make test-vectors` and leave `make xgoja-smoke-vectors` local-only. This is easiest but does not fulfill the generated xgoja smoke goal in CI.

Recommended order:

1. Add local `xgoja-smoke-vectors` first.
2. Add FAISS CI running `make test-vectors` only.
3. Add a CI-compatible xgoja spec and enable generated smoke in the workflow.

## Design decisions

### Decision 1: keep package vector tests and generated smoke separate

**Decision:** `make test-vectors` should remain package-level, and `make xgoja-smoke-vectors` should validate the generated binary.

**Rationale:** These tests fail for different reasons. Package tests catch binding and Bleve behavior. Generated smoke catches xgoja build-spec handling, provider registration, jsverb mounting, embedded source loading, and CLI command execution. Combining them into one target makes failures harder for interns to triage.

**Consequence:** CI can choose either or both targets depending on native dependency setup.

### Decision 2: make FAISS CI non-required first

**Decision:** Start with `workflow_dispatch` and a weekly schedule. Do not make the job required on every PR until setup time and stability are known.

**Rationale:** FAISS is a C++ dependency. Source builds can be slow and may be sensitive to runner image changes. A flaky required job would slow down unrelated non-vector changes.

**Consequence:** The job still catches scheduled regressions and can be run manually before releases or vector-heavy merges.

### Decision 3: prefer a CI-specific xgoja spec over sibling checkouts

**Decision:** Use released module versions in a `xgoja-vectors.ci.yaml` spec when possible.

**Rationale:** The checked-in local spec has sibling `replace` directives because it supports active workspace development. CI should test what a standalone consumer can build. Released module versions better represent that environment.

**Consequence:** When go-go-goja or Geppetto changes are required, those repos must be tagged before CI can use the new behavior without sibling checkouts.

### Decision 4: treat `rag index-query` as manual integration, not CI smoke

**Decision:** CI smoke should run `vector knn` and `vector hybrid`, not the Geppetto-backed `rag index-query` command.

**Rationale:** `rag index-query` needs an embedding provider profile and possibly network/model availability. CI should avoid external model dependencies for this hardening slice.

**Consequence:** The generated binary still includes the RAG verb, but CI validates the vector search path with deterministic embedded smoke data.

## Implementation plan for the intern

### Step 1: add Makefile variables and local vector smoke targets

Edit `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile`.

Add phony targets:

```makefile
.PHONY: gifs logcopter-generate logcopter-check test-vectors xgoja-build-vectors xgoja-smoke-vectors
```

Add variables near the other build variables:

```makefile
XGOJA_VERSION ?= v0.8.3
XGOJA_VECTOR_SPEC ?= xgoja-vectors.yaml
XGOJA_VECTOR_WORK_DIR ?= /tmp/goja-bleve-vector-work
```

Add targets:

```makefile
xgoja-build-vectors:
	cd cmd/goja-bleve && GOWORK=off go run github.com/go-go-golems/go-go-goja/cmd/xgoja@$(XGOJA_VERSION) build \
		-f $(XGOJA_VECTOR_SPEC) \
		--work-dir $(XGOJA_VECTOR_WORK_DIR) \
		--keep-work \
		--xgoja-version $(XGOJA_VERSION)

xgoja-smoke-vectors: xgoja-build-vectors
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector knn --output json
	cd cmd/goja-bleve && ./dist/goja-bleve-vectors vector hybrid --output json
```

Validation:

```bash
make xgoja-smoke-vectors
```

Expected output includes JSON arrays whose first hit is `chunk-1`.

### Step 2: add a CI-compatible xgoja vector spec if needed

If a clean checkout cannot use the local sibling replaces in `xgoja-vectors.yaml`, add:

```text
/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.ci.yaml
```

Start by copying `xgoja-vectors.yaml`, then remove local sibling replaces for `geppetto`, `go-go-goja-core`, and `go-go-goja-host`. Keep the local `goja-bleve` replace because the nested generated module must use the current PR checkout:

```yaml
  - id: goja-bleve
    import: github.com/go-go-golems/goja-bleve/pkg/xgoja/providers/bleve
    version: v0.0.0
    replace: ../..
```

Validate locally with network access:

```bash
make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml
```

If this fails because released versions do not contain needed xgoja features, use sibling checkouts in CI until releases are available.

### Step 3: add optional FAISS workflow

Create:

```text
/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml
```

Initial workflow should run manually and weekly. Do not include `pull_request` until it is stable.

Use the workflow skeleton in the proposed solution. If the FAISS build takes too long, add caching in a second pass rather than over-optimizing the first version.

### Step 4: validate workflow commands locally where possible

Local validation commands:

```bash
make test-vectors
make xgoja-smoke-vectors
GOWORK=off go test ./...
GOWORK=off golangci-lint run ./...
```

If the CI spec exists:

```bash
make xgoja-smoke-vectors XGOJA_VECTOR_SPEC=xgoja-vectors.ci.yaml
```

### Step 5: document the targets

Update:

- `README.md`
- `docs/quickstart.md`
- `docs/faiss-xgoja-playbook.md`

The docs should distinguish these commands:

```bash
make test-vectors          # package vector tests
make xgoja-smoke-vectors   # generated vector binary and JS smoke verbs
```

### Step 6: decide required-vs-optional status after observation

After the workflow runs successfully several times, decide whether to add:

```yaml
on:
  pull_request:
    paths:
      - 'pkg/**'
      - 'cmd/goja-bleve/**'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/vector-faiss.yml'
```

Do not make the job required immediately. First measure runtime and stability.

## Pseudocode: target and workflow responsibilities

### Local target flow

```text
function xgoja_smoke_vectors():
    run xgoja build with:
        spec = cmd/goja-bleve/xgoja-vectors.yaml
        xgoja version = v0.8.3
        work dir = /tmp/goja-bleve-vector-work
    assert dist/goja-bleve-vectors exists
    run dist/goja-bleve-vectors vector knn --output json
    assert command exits 0
    run dist/goja-bleve-vectors vector hybrid --output json
    assert command exits 0
```

### CI workflow flow

```text
function vector_faiss_ci():
    checkout goja-bleve
    setup Go from go.mod
    install native build dependencies
    clone blevesearch/faiss at fff814d
    configure FAISS with C API and shared libs
    build faiss and faiss_c
    install headers and shared libs into /usr/local
    run ldconfig
    verify headers and libs exist
    run make test-vectors
    if generated smoke is enabled:
        run make xgoja-smoke-vectors with CI-compatible spec
```

## Expected failure modes and diagnosis

### `undefined reference to faiss::...`

Cause: the build linked `libfaiss_c.so` but did not link `libfaiss.so`.

Fix:

```bash
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
```

### `fatal error: faiss/c_api/IndexBinary_c_ex.h: No such file or directory`

Cause: FAISS headers are not installed or the wrong FAISS fork was used.

Fix: build and install `blevesearch/faiss@fff814d`, then verify:

```bash
ls /usr/local/include/faiss/c_api/IndexBinary_c_ex.h
```

### `replacement directory ../../../geppetto does not exist`

Cause: `xgoja-vectors.yaml` uses local workspace replaces that do not exist in clean CI.

Fix: use `xgoja-vectors.ci.yaml` with released module versions, or checkout sibling repositories into the expected paths.

### `bleve.vectorSupport` is false inside generated binary

Cause: xgoja build did not pass the `vectors` tag.

Fix: ensure the spec contains:

```yaml
go:
  tags:
    - vectors
```

and rebuild the binary.

### `rag` group appears but child verbs are missing

Cause: jsverb command fields may have been dropped when provider runtime sections were appended, or the verb source failed scanning.

Fixes to inspect:

- Check `pkg/xgoja/app/module_sections.go` in go-go-goja preserves existing command sections.
- Run xgoja validation output and inspect jsverb scanner errors.
- Ensure `__verb__` is attached to a real JavaScript function identifier such as `indexQuery`, not a hyphenated string that cannot be a function name.

## Testing strategy

### Unit tests affected by this work

The Makefile target itself does not need Go unit tests, but the xgoja changes it depends on already have tests in go-go-goja:

- `cmd/xgoja/internal/buildexec/buildexec_test.go` for deterministic env handling.
- `pkg/jsverbs/jsverbs_test.go` for field-name remapping and bound section behavior.
- `pkg/xgoja/app/module_sections_test.go` for preserving command sections while appending runtime sections.

For goja-bleve, use command-level tests:

```bash
make test-vectors
make xgoja-smoke-vectors
```

### CI acceptance criteria

A successful FAISS CI run must prove all of the following:

- FAISS headers are installed.
- `libfaiss_c.so` and `libfaiss.so` are visible to the linker and runtime loader.
- `make test-vectors` passes on the runner.
- If generated smoke is enabled, the generated vector binary builds from YAML and both vector jsverbs pass.

### What not to assert

Do not assert exact floating-point scores in shell workflow tests. Let Go tests assert detailed behavior. The workflow should assert command success. Exact JSON score values can vary with Bleve/FAISS internals.

## File reference map

| File | Role in this work |
|---|---|
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/Makefile` | Add `xgoja-build-vectors` and `xgoja-smoke-vectors`; already contains `test-vectors`. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/xgoja-vectors.yaml` | Local vector generated-host spec with `vectors`, rpath, `CGO_LDFLAGS`, providers, modules, and jsverbs. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/vector.js` | Deterministic generated-runtime vector KNN and hybrid smoke verbs. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/cmd/goja-bleve/jsverbs/rag.js` | Heavier Geppetto + Bleve RAG integration verb; useful for manual smoke, not first CI target. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/docs/faiss-xgoja-playbook.md` | Operational FAISS build and xgoja link guide; update with new targets. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/push.yml` | Existing default Go pipeline; do not add FAISS work here initially. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/goja-bleve/.github/workflows/vector-faiss.yml` | Proposed optional workflow for FAISS-backed vector validation. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/cmd/xgoja/internal/buildexec/buildexec.go` | xgoja build executor that passes `go.env` into `go build`. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/cmd/xgoja/internal/buildspec/build_spec.go` | `GoSpec.Env` schema used by `xgoja-vectors.yaml`. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/jsverbs/runtime.go` | jsverb argument/field remapping path used by generated commands. |
| `/home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja/pkg/xgoja/app/module_sections.go` | Section merge helper that must preserve command fields when adding provider runtime controls. |

## Review checklist

Before merging an implementation of this ticket, review these questions:

- Does `make test-vectors` still pass locally?
- Does `make xgoja-smoke-vectors` pass locally?
- Does the Makefile let CI override `XGOJA_VECTOR_SPEC` and `XGOJA_VERSION`?
- Does any CI-specific spec avoid non-existent sibling `replace` paths?
- Does the FAISS workflow start as non-required (`workflow_dispatch` and schedule only)?
- Does the workflow verify installed FAISS files before running Go tests?
- Does the workflow use `make test-vectors` rather than duplicating the vector test command inline?
- If generated smoke runs in CI, does it avoid embedding-provider-dependent commands?
- Are README and playbook docs updated with the new target names?

## Alternatives considered

### Add FAISS work to the default `golang-pipeline`

Rejected for the first implementation. The default pipeline should remain fast and reliable. FAISS source builds are native dependency work and should prove themselves in a separate workflow before becoming required.

### Use Docker image with prebuilt FAISS immediately

Deferred. A prebuilt image can make CI faster, but it introduces image maintenance. The first workflow should show the native dependency steps explicitly so failures are understandable.

### Run only xgoja smoke and skip package vector tests

Rejected. Generated smoke is intentionally small. Package tests cover more edge cases: invalid inputs, reopened mappings, missing KNN fields, dimension mismatches, and hybrid score validation. CI should run `make test-vectors` before any generated smoke.

### Run the RAG Geppetto command in CI

Rejected for the first implementation. It requires a real embedding provider profile and possibly model availability. Deterministic vector jsverbs provide enough generated-host coverage without external model dependencies.

## Open questions

1. Should CI use a separate `xgoja-vectors.ci.yaml` or checkout sibling repositories? The recommended long-term answer is a CI-specific spec with released versions, but that depends on current go-go-goja and Geppetto releases containing the needed provider features.
2. Should the FAISS build be cached? Start without cache; add cache only after measuring workflow runtime.
3. Should the workflow run on pull requests with path filters? Start manual/scheduled; add PR triggers only after stability is proven.
4. Should generated xgoja smoke include non-vector `dist/goja-bleve` commands too? That can be a separate target, but it is not required for this vector/FAISS slice.
5. Should `make xgoja-smoke-vectors` delete or reuse `/tmp/goja-bleve-vector-work`? Reuse is faster locally; `--keep-work` helps debugging. CI can use a fresh runner.

## Definition of done

This ticket is complete when:

- `make xgoja-smoke-vectors` exists and passes on a machine with the documented FAISS setup.
- The target builds `cmd/goja-bleve/xgoja-vectors.yaml` with xgoja v0.8.3 or a configurable newer version.
- The target runs `vector knn` and `vector hybrid` successfully.
- A non-required FAISS workflow exists and successfully runs `make test-vectors` on GitHub Actions.
- If generated smoke is included in CI, the workflow uses a spec that works on a clean checkout.
- Documentation explains the difference between package vector tests and generated xgoja vector smoke.
