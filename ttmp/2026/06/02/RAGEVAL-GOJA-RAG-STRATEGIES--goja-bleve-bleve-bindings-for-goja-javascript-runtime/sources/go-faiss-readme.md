---
Title: Source: go-faiss-readme.md
Ticket: RAGEVAL-GOJA-RAG-STRATEGIES
Status: active
Topics: [bleve]
DocType: reference
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Source material fetched for goja-bleve FAISS/vector investigation"
LastUpdated: 2026-06-02T22:00:00.000000-04:00
WhatFor: "Reference source material for FAISS vector build requirements"
WhenToUse: "When setting up bleve vector support and FAISS compilation"
---

## go-faiss

[![Go 
Reference](https://camo.githubusercontent.com/c5c625b6569f0c7b15ce9fd274402c9f779a507c7a3fb908573945
fdecad2899/68747470733a2f2f706b672e676f2e6465762f62616467652f6769746875622e636f6d2f44617461496e74656
c6c6967656e6365437265772f676f2d66616973732e737667)](https://pkg.go.dev/github.com/DataIntelligenceCr
ew/go-faiss)

Go bindings for [Faiss](https://github.com/facebookresearch/faiss), a library for vector similarity 
search.

## Install

First you will need to build and install Faiss:

```
git clone https://github.com/blevesearch/faiss.git
cd faiss
# Install minimal build dependencies
# Ubuntu/Debian based systems: apt-get install libblas-dev liblapack-dev swig build-essential

export Python_INCLUDE_DIRS=<PYTHON_INSTALLED_DIR>/include
export Python_LIBRARIES=<PYTHON_INSTALLED_DIR>/lib/libpython3.so
cmake -B build -DFAISS_ENABLE_GPU=OFF -DFAISS_ENABLE_C_API=ON -DBUILD_SHARED_LIBS=ON .
make -C build
sudo make -C build install
```

On osX ARM64, the instructions needed to be slightly adjusted based on 
[facebookresearch/faiss#2111](https://github.com/facebookresearch/faiss/issues/2111):

```
LDFLAGS="-L/opt/homebrew/opt/llvm/lib" CPPFLAGS="-I/opt/homebrew/opt/llvm/include" 
CXX=/opt/homebrew/opt/llvm/bin/clang++ CC=/opt/homebrew/opt/llvm/bin/clang cmake -B build 
-DFAISS_ENABLE_GPU=OFF -DFAISS_ENABLE_C_API=ON -DBUILD_SHARED_LIBS=ON .
// set FAISS_ENABLE_PYTHON to OFF in CMakeLists.txt to ignore libpython dylib
make -C build
sudo make -C build install
```

Building will produce the dynamic library `faiss_c`. You will need to install it in a place where 
your system will find it (e.g. `/usr/local/lib` on mac or `/usr/lib` on Linux). You can do this 
with:

```
sudo cp build/c_api/libfaiss_c.so /usr/local/lib
```

Now you can install the Go module:

```
go get github.com/blevesearch/go-faiss
```

## Usage

API documentation is available at 
[https://pkg.go.dev/github.com/DataIntelligenceCrew/go-faiss](https://pkg.go.dev/github.com/DataInte
lligenceCrew/go-faiss). See the [Faiss wiki](https://github.com/facebookresearch/faiss/wiki) for 
more information.

Examples can be found in the 
[\_example](https://github.com/blevesearch/go-faiss/blob/master/_example) directory.
