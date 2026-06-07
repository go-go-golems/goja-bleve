# FAISS + goja-bleve + xgoja Playbook

This playbook explains how to configure, build, install, diagnose, and use FAISS so that `goja-bleve` links successfully when Bleve's vector support is enabled. It also covers the xgoja case, where the generated host binary must receive the same Go build tags, CGO linker flags, and runtime library path.

## Executive summary

Bleve vector/KNN support is behind the Go build tag:

```bash
-tags=vectors
```

When that tag is enabled, Bleve imports `github.com/blevesearch/go-faiss`, which calls FAISS through CGO. A working local setup needs all of the following:

1. The Bleve-compatible FAISS fork installed, not an arbitrary upstream FAISS build.
2. `libfaiss_c.so` and `libfaiss.so` available to the linker.
3. FAISS C API headers installed under an include path CGO can find.
4. Explicit Go linker inputs for both the C API and the C++ library:

   ```bash
   CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
   ```

5. A runtime library path or loader configuration so the compiled test/binary can find the shared libraries:

   ```bash
   -ldflags "-r /usr/local/lib"
   ```

The canonical local validation command is:

```bash
make test-vectors
```

which expands to:

```bash
GOWORK=off \
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1
```

## Why plain `go test -tags=vectors` may fail

A tempting command is:

```bash
GOWORK=off go test -tags=vectors ./pkg
```

That can fail at link time with many errors like:

```text
/usr/bin/ld: /usr/local/lib/libfaiss_c.so: undefined reference to `faiss::IndexPreTransform::IndexPreTransform(...)'
/usr/bin/ld: /usr/local/lib/libfaiss_c.so: undefined reference to `faiss::IndexFlat::search(...) const'
collect2: error: ld returned 1 exit status
```

The reason is that `go-faiss` contributes the FAISS C API link input, but this local shared-library build leaves many C++ FAISS symbols unresolved in `libfaiss_c.so` until the final executable also links `libfaiss.so`. Therefore the final Go link must include:

```bash
-lfaiss_c -lfaiss -lstdc++ -lm
```

The order matters: put `-lfaiss_c` before `-lfaiss` so references from the C API library can be resolved by the C++ library.

## Known-good local layout

This repository assumes the following default local install paths:

```text
/usr/local/include/faiss/...
/usr/local/lib/libfaiss_c.so
/usr/local/lib/libfaiss.so
/usr/local/lib/libfaiss_avx512.so   # optional, depending on the build
```

Verify the installed files:

```bash
ls -lh \
  /usr/local/include/faiss/c_api/IndexBinary_c_ex.h \
  /usr/local/lib/libfaiss_c.so \
  /usr/local/lib/libfaiss.so

ldconfig -p | grep -E 'libfaiss(_c)?\.so'
```

Inspect dynamic dependencies:

```bash
ldd /usr/local/lib/libfaiss_c.so
ldd /usr/local/lib/libfaiss.so
```

A common observation is:

- `libfaiss_c.so` depends on libc/libstdc++, but does not list `libfaiss.so`.
- `libfaiss.so` depends on BLAS/OpenMP/C++ runtime libraries.

That is why explicit `CGO_LDFLAGS` are needed.

## Build FAISS from source

Use the Bleve-maintained fork:

```bash
cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system
git clone https://github.com/blevesearch/faiss.git faiss
cd faiss
git checkout fff814d
```

The `fff814d` commit is the checkpoint recorded as compatible with Bleve v2.6.0 in the existing RAG evaluation system setup notes.

Configure a CPU-only shared-library build with the C API enabled:

```bash
cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system/faiss
rm -rf build

cmake -B build \
  -DFAISS_ENABLE_GPU=OFF \
  -DFAISS_ENABLE_C_API=ON \
  -DBUILD_SHARED_LIBS=ON \
  -DFAISS_ENABLE_PYTHON=OFF \
  -DCMAKE_INSTALL_PREFIX=/usr/local \
  -DCMAKE_CXX_FLAGS="-I$PWD" \
  .
```

Build only the targets needed by Bleve/go-faiss:

```bash
make -C build -j$(nproc) faiss faiss_c
```

Install and refresh the dynamic linker cache:

```bash
sudo make -C build install
sudo ldconfig
```

If `make install` does not leave a fresh `libfaiss.so` in `/usr/local/lib`, copy it explicitly:

```bash
sudo cp build/faiss/libfaiss.so /usr/local/lib/
sudo cp build/c_api/libfaiss_c.so /usr/local/lib/
sudo ldconfig
```

Then re-run the verification commands from the previous section.

## Build and test goja-bleve with vectors

From the repository root:

```bash
make test-vectors
```

Equivalent explicit command:

```bash
GOWORK=off \
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm" \
go test -tags=vectors -ldflags "-r /usr/local/lib" ./pkg -count=1
```

What each part does:

- `GOWORK=off` keeps this module isolated from any surrounding workspace.
- `CGO_LDFLAGS=...` tells cgo which FAISS libraries to pass to the external linker.
- `-tags=vectors` enables Bleve/goja-bleve vector code.
- `-ldflags "-r /usr/local/lib"` embeds an ELF runtime search path in the test binary.
- `./pkg` is where the vector build-tag tests live.

If this command passes, the local FAISS installation is good enough for `goja-bleve`'s vector tests.

## Build and smoke-test an xgoja vector host

The xgoja host must receive the same three categories of build configuration:

1. Go build tag: `vectors`
2. CGO linker flags: `-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm`
3. Runtime rpath: `-r /usr/local/lib`

The vector-specific spec in this repository already encodes those settings:

```yaml
# cmd/goja-bleve/xgoja-vectors.yaml
go:
  tags:
    - vectors
  ldflags:
    - -r
    - /usr/local/lib
  env:
    CGO_LDFLAGS: "-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
```

Use the repository target for the normal smoke path:

```bash
make xgoja-smoke-vectors
```

That target builds the vector-enabled xgoja binary from `cmd/goja-bleve/xgoja-vectors.yaml`, then runs the deterministic vector smoke verbs:

```bash
cd cmd/goja-bleve
./dist/goja-bleve-vectors vector knn --output json
./dist/goja-bleve-vectors vector hybrid --output json
```

The build step is configurable for local debugging or CI experiments:

```bash
make xgoja-smoke-vectors \
  XGOJA_VERSION=v0.8.3 \
  XGOJA_VECTOR_SPEC=xgoja-vectors.yaml \
  XGOJA_VECTOR_WORK_DIR=/tmp/goja-bleve-vector-work
```

If xgoja is invoked another way, make sure the generated Go package sees the same effective settings. A non-vector xgoja build will still load `require("bleve")`, but `bleve.vectorSupport` will be `false` and vector APIs will return explicit `-tags=vectors` errors.

## Optional GitHub Actions vector workflow

This repository includes an optional workflow at `.github/workflows/vector-faiss.yml`. It is intentionally not part of the default pull-request pipeline yet. It runs on a weekly schedule and can be started manually from the GitHub Actions UI.

The workflow builds the Bleve-compatible FAISS fork on an Ubuntu runner, installs `libfaiss_c.so` and `libfaiss.so` under `/usr/local/lib`, verifies headers and dynamic linker visibility, then runs:

```bash
make test-vectors
```

The manual trigger has a `run-xgoja-smoke` input. Leave it at `false` unless the xgoja vector spec resolves on the runner. The local `cmd/goja-bleve/xgoja-vectors.yaml` still contains sibling workspace `replace` paths for active development, so generated xgoja smoke is opt-in until a CI-compatible spec or sibling checkout strategy is added.

## Use FAISS installed somewhere else

If FAISS is not installed under `/usr/local`, change all three places consistently.

For example, with FAISS under `/opt/faiss`:

```bash
FAISS_PREFIX=/opt/faiss
FAISS_LIB_DIR="$FAISS_PREFIX/lib"

GOWORK=off \
CGO_CFLAGS="-I$FAISS_PREFIX/include" \
CGO_LDFLAGS="-L$FAISS_LIB_DIR -lfaiss_c -lfaiss -lstdc++ -lm" \
go test -tags=vectors -ldflags "-r $FAISS_LIB_DIR" ./pkg -count=1
```

For xgoja, mirror that in the spec:

```yaml
go:
  tags:
    - vectors
  ldflags:
    - -r
    - /opt/faiss/lib
  env:
    CGO_CFLAGS: "-I/opt/faiss/include"
    CGO_LDFLAGS: "-L/opt/faiss/lib -lfaiss_c -lfaiss -lstdc++ -lm"
```

## Troubleshooting

### Missing FAISS C API header

Symptom:

```text
fatal error: faiss/c_api/IndexBinary_c_ex.h: No such file or directory
```

Likely causes:

- FAISS was not installed.
- You built upstream `facebookresearch/faiss` instead of `blevesearch/faiss`.
- The include prefix is not visible to CGO.

Fix:

```bash
sudo make -C build install
ls /usr/local/include/faiss/c_api/IndexBinary_c_ex.h
```

If using a custom prefix, set `CGO_CFLAGS="-I/path/to/include"`.

### Missing `IndexIVFRaBitQ.h` while building FAISS

Symptom:

```text
fatal error: faiss/IndexIVFRaBitQ.h: No such file or directory
```

Fix: configure with the repository root on the C++ include path:

```bash
-DCMAKE_CXX_FLAGS="-I$PWD"
```

### Undefined `faiss::...` references during Go link

Symptom:

```text
/usr/bin/ld: /usr/local/lib/libfaiss_c.so: undefined reference to `faiss::...'
collect2: error: ld returned 1 exit status
```

Fix: use the full linker flags:

```bash
CGO_LDFLAGS="-L/usr/local/lib -lfaiss_c -lfaiss -lstdc++ -lm"
```

Then run through `make test-vectors` or the explicit command above.

### Runtime loader cannot find `libfaiss_c.so` or `libfaiss.so`

Symptom:

```text
error while loading shared libraries: libfaiss_c.so: cannot open shared object file
```

Fix options:

1. Run `sudo ldconfig` after installing to `/usr/local/lib`.
2. Add an rpath while building:

   ```bash
   -ldflags "-r /usr/local/lib"
   ```

3. As a last-resort development workaround, set:

   ```bash
   export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
   ```

Prefer `ldconfig` or rpath over relying on `LD_LIBRARY_PATH` for generated xgoja binaries.

### xgoja binary reports `bleve.vectorSupport === false`

Symptom:

```javascript
const bleve = require("bleve");
bleve.vectorSupport; // false
```

Likely cause: the xgoja binary was built without `go.tags: [vectors]`.

Fix: build from `cmd/goja-bleve/xgoja-vectors.yaml`, or add the equivalent `go.tags`, `go.env.CGO_LDFLAGS`, and `go.ldflags` entries to your own xgoja spec.

### xgoja build links, but jsverb execution fails at startup

Likely cause: the build found FAISS at link time but the final binary cannot find the shared libraries at runtime.

Fix: ensure the spec contains:

```yaml
go:
  ldflags:
    - -r
    - /usr/local/lib
```

or install FAISS into a directory already known to the dynamic loader and run `sudo ldconfig`.

## Documentation setup in this repository

The current documentation setup is intentionally small:

- `README.md` is the broad overview and API reference.
- `docs/quickstart.md` is the fastest path from zero to a working text/vector example.
- `docs/faiss-xgoja-playbook.md` is this operational playbook for FAISS and xgoja linking.
- `examples/*.js` are runnable JavaScript examples.
- `cmd/goja-bleve/xgoja*.yaml` are executable xgoja build specifications.

Suggested improvements:

1. Keep `README.md` short and link out to focused docs instead of embedding long operational instructions.
2. Treat `docs/quickstart.md` as user-facing first-run material.
3. Treat this playbook as maintainer/operator material for local machines and CI runners.
4. Add a `docs/README.md` index if more documents are added.
5. Keep `.github/workflows/vector-faiss.yml` optional until the FAISS build runtime is known to be stable enough for pull requests.

## Final checklist

Use this checklist when preparing a machine for vector-enabled `goja-bleve` development:

- [ ] `blevesearch/faiss` is checked out at a Bleve-compatible commit.
- [ ] FAISS was configured with `FAISS_ENABLE_C_API=ON` and `BUILD_SHARED_LIBS=ON`.
- [ ] `libfaiss_c.so` and `libfaiss.so` exist under the chosen library directory.
- [ ] FAISS C API headers exist under the chosen include directory.
- [ ] `sudo ldconfig` was run after installing to `/usr/local/lib`.
- [ ] `make test-vectors` passes.
- [ ] `make xgoja-smoke-vectors` passes.
- [ ] xgoja vector specs include `tags: [vectors]`.
- [ ] xgoja vector specs include `CGO_LDFLAGS` with both `-lfaiss_c` and `-lfaiss`.
- [ ] xgoja vector specs include an rpath or the runtime loader can otherwise find FAISS.
