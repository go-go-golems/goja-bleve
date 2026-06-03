---
Title: Source: go-faiss-pkg-doc.md
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

## README

### go-faiss

[![Go 
Reference](https://pkg.go.dev/badge/github.com/DataIntelligenceCrew/go-faiss.svg)](https://pkg.go.de
v/github.com/DataIntelligenceCrew/go-faiss)

Go bindings for [Faiss](https://github.com/facebookresearch/faiss), a library for vector similarity 
search.

#### Install

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
[https://github.com/facebookresearch/faiss/issues/2111](https://github.com/facebookresearch/faiss/is
sues/2111):

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

#### Usage

API documentation is available at 
[https://pkg.go.dev/github.com/DataIntelligenceCrew/go-faiss](https://pkg.go.dev/github.com/DataInte
lligenceCrew/go-faiss). See the [Faiss wiki](https://github.com/facebookresearch/faiss/wiki) for 
more information.

Examples can be found in the 
[\_example](https://github.com/blevesearch/go-faiss/blob/v1.1.3/_example) directory.

## Documentation

### Overview

Package faiss provides bindings to Faiss, a library for vector similarity search. More detailed 
documentation can be found at the Faiss wiki: 
[https://github.com/facebookresearch/faiss/wiki](https://github.com/facebookresearch/faiss/wiki).

### Index

### Constants

```
const (
    MetricInnerProduct  = C.METRIC_INNER_PRODUCT
    MetricL2            = C.METRIC_L2
    MetricL1            = C.METRIC_L1
    MetricLinf          = C.METRIC_Linf
    MetricLp            = C.METRIC_Lp
    MetricCanberra      = C.METRIC_Canberra
    MetricBrayCurtis    = C.METRIC_BrayCurtis
    MetricJensenShannon = C.METRIC_JensenShannon
)
```

Metric type

```
const (
    IOFlagMmap         = C.FAISS_IO_FLAG_MMAP
    IOFlagReadOnly     = C.FAISS_IO_FLAG_READ_ONLY
    IOFlagReadMmap     = C.FAISS_IO_FLAG_READ_MMAP | C.FAISS_IO_FLAG_ONDISK_IVF
    IOFlagSkipPrefetch = C.FAISS_IO_FLAG_SKIP_PREFETCH
)
```

### Variables

```
var (
    ErrCreateIndexFailed    = errors.New("create index failed")
    ErrCreateSelectorFailed = errors.New("create selector failed")

    ErrCreateParamsFailed = errors.New("create search params failed")
    ErrSetParamsFailed    = errors.New("set index params failed")

    ErrAddFailed          = errors.New("add vectors failed")
    ErrTrainFailed        = errors.New("train index failed")
    ErrSearchFailed       = errors.New("search index failed")
    ErrReconstructFailed  = errors.New("reconstruct vector failed")
    ErrResetIndexFailed   = errors.New("reset index failed")
    ErrSetQuantizerFailed = errors.New("set quantizer failed")
    ErrMergeFromFailed    = errors.New("merge from index failed")
    ErrRemoveIDsFailed    = errors.New("remove IDs failed")

    ErrInspectIndexFailed = errors.New("inspect index failed")

    ErrWriteIndexFailed = errors.New("write index failed")
    ErrReadIndexFailed  = errors.New("read index failed")

    ErrNoUsableGPUDevices = errors.New("no GPU usable devices available")
    ErrGPUCloneFailed     = errors.New("GPU clone failed")
    ErrGPUSetupFailed     = errors.New("GPU setup failed")
    ErrGPUContextFailed   = errors.New("GPU context init failed")
    ErrGPUOutOfMemory     = errors.New("GPU out of memory")

    ErrIndexNil      = errors.New("index is nil")
    ErrSelectorNil   = errors.New("selector is nil")
    ErrNotIDMapIndex = errors.New("index is not an IDMap index")
    ErrNotIVFIndex   = errors.New("index is not an IVF index")
    ErrNotBIVFIndex  = errors.New("index is not a binary IVF index")

    ErrMergeFromNotSupported    = errors.New("merge from is not supported for this index type")
    ErrSetQuantizerNotSupported = errors.New("set quantizer not supported for this index type")
)
```

FAISS error types for categorizing errors returned by the C API.

### Functions

#### added in v1.0.22

```
func NormalizeVector(vector []float32) []float32
```

In-place normalization of provided vector (single)

#### added in v1.0.17

```
func SetOMPThreads(n uint)
```

#### added in v1.0.28

```
func WriteBinaryIndexIntoBuffer(idx BinaryIndex) ([]byte, error)
```

#### func WriteIndex ¶

```
func WriteIndex(idx Index, filename string) error
```

WriteIndex writes an index to a file.

#### added in v1.0.0

```
func WriteIndexIntoBuffer(idx Index) ([]byte, error)
```

### Types

#### added in v1.0.28

```
type BinaryIndex interface {
    // D returns the dimension of the indexed vectors.
    D() int

    // MetricType returns the metric type of the index.
    MetricType() int

    // Ntotal returns the total number of vectors currently stored in the index.
    Ntotal() int64

    // set the direct map type for IVF indexes.
    // 0 for No Map
    // 1 for Array
    // 2 for Hash
    SetDirectMap(maptype int) error

    // set the number of probes for IVF indexes
    SetNProbe(nprobe int32)

    // returns true if the underlying index is an IVF index
    IsIVFIndex() bool

    // IVFParams returns the nlist and nprobe parameters for IVF indexes
    IVFParams() (nprobe int, nlist int)

    // trains the index on a representative set of vectors
    Train(xb []uint8) error

    // adds vectors to the index
    Add(xb []uint8) error

    // sets the qunatizers from the source index, supposed to be used only for
    // BIVF indexes and returns error otherwise
    SetQuantizers(srcIndex BinaryIndex) error

    // merges another binary index into this one, currently applicable only for
    // IVF indexes returns an error
    MergeFrom(other BinaryIndex, add_id int64) error

    // queries the index with the vectors in xb
    // returns the IDs of the k nearest neighbors for each query vector and
    // their corresponding distances
    Search(xb []uint8, k int64) (distances []int32, labels []int64, err error)

    // SearchWithOptions performs a search with additional optional constraints.
    // - Selector can be used to restrict the search to a subset of the indexed vectors based on 
their IDs.
    // - params is a JSON object that can contain additional search parameters specific to the 
index type, such as IVF search parameters.
    SearchWithOptions(xb []uint8, k int64, sel Selector, params json.RawMessage) (distances 
[]int32, labels []int64, err error)

    // returns a slice where each index corresponds to a cluster in an IVF
    // index, and the value at each index is the count of vectors in that
    // cluster, considering only the vectors specified in the include selector.
    ObtainClusterVectorCountsFromIVFIndex(include Selector, nlist int) (
        []int64, error)

    // returns the IDs and distances of the closest numCentroids centroids to
    // the query vector xb, considering only the centroids specified in the
    // includedCentroids selector.
    ObtainClustersWithDistancesFromIVFIndex(xb []uint8, includedCentroids Selector,
        numCentroids int64) ([]int64, []int32, error)

    // Applicable only to IVF indexes: Returns the top k centroid cardinalities and
    // their vectors in chosen order (descending or ascending)
    ObtainKCentroidCardinalitiesFromIVFIndex(limit int, descending bool) ([]uint64, [][]uint8, 
error)

    // searches the specified clusters in an IVF index for the k nearest neighbors
    // of the query vector xb, considering only the vectors specified in the include selector
    // and additional search parameters passed as a JSON object.
    SearchClustersFromIVFIndex(eligibleCentroidIDs []int64, centroidDis []int32,
        centroidsToProbe int, xb []uint8, k int64, include Selector,
        params json.RawMessage) ([]int32, []int64, error)

    // Size estimates the memory footprint of the index in bytes
    // if the underlying faiss index is memory-mapped and not fully loaded into memory.
    Size() uint64

    // frees the memory associated with the index
    Close()

    // CodeSize returns the size of the produced codes in bytes.
    CodeSize() (uint64, error)
    // contains filtered or unexported methods
}
```

#### added in v1.0.28

```
type BinaryIndexImpl struct {
    BinaryIndex
}
```

#### added in v1.0.28

```
func BinaryIndexFactory(dims int, description string) (*BinaryIndexImpl, error)
```

#### added in v1.0.28

```
func ReadBinaryIndexFromBuffer(buf []byte, ioflags int) (*BinaryIndexImpl, error)
```

#### added in v1.0.31

```
type GPUIndexImpl struct{}
```

GPUIndexImpl is an opaque type when not built with GPU support.

#### added in v1.0.31

```
func CloneToGPU(_ *IndexImpl) (*GPUIndexImpl, error)
```

CloneToGPU is not available without the gpu build tag.

#### added in v1.0.31

```
func (g *GPUIndexImpl) Add(x []float32) error
```

#### added in v1.0.31

```
func (g *GPUIndexImpl) Close()
```

#### added in v1.0.31

```
func (g *GPUIndexImpl) Search(x []float32, k int64) ([]float32, []int64, error)
```

#### added in v1.1.2

```
func (g *GPUIndexImpl) Size() uint64
```

#### added in v1.0.31

```
func (g *GPUIndexImpl) Train(x []float32) error
```

#### type IDSelector ¶

```
type IDSelector struct {
    // contains filtered or unexported fields
}
```

IDSelector represents a set of IDs to remove.

#### func (\*IDSelector) Delete ¶

```
func (s *IDSelector) Delete()
```

Delete frees the memory associated with s.

#### added in v1.0.33

```
func (s *IDSelector) ExcludeFilter() bool
```

#### added in v1.0.23

```
func (s *IDSelector) Get() *C.FaissIDSelector
```

#### type Index ¶

```
type Index interface {
    // D returns the dimension of the indexed vectors.
    D() int

    // IsTrained returns true if the index has been trained or does not require
    // training.
    IsTrained() bool

    // Ntotal returns the number of indexed vectors.
    Ntotal() int64

    // set the direct map type for IVF indexes.
    // 0 for No Map
    // 1 for Array
    // 2 for Hash
    SetDirectMap(maptype int) error

    // set the number of probes for IVF indexes
    SetNProbe(nprobe int32)

    // MetricType returns the metric type of the index.
    MetricType() int

    // Train trains the index on a representative set of vectors.
    Train(x []float32) error

    // Add adds vectors to the index.
    Add(x []float32) error

    // AddWithIDs is like Add, but stores xids instead of sequential IDs.
    AddWithIDs(x []float32, xids []int64) error

    // Returns true if the index is an IVF index.
    IsIVFIndex() bool

    // Returns true if the index is a scalar quantization (SQ) index.
    IsSQIndex() bool

    // Returns true if the index has RaBitQ
    HasRaBitQ() bool

    // Returns the IVF parameters nprobe and nlist for IVF indexes.
    IVFParams() (nprobe, nlist int)

    // Applicable only to IVF indexes: Returns a slice where each index represents
    // a cluster (list) ID and the value is the count of selected vectors belonging
    // to that cluster. Only vectors specified by the given Selector are considered.
    ObtainClusterVectorCountsFromIVFIndex(include Selector, nlist int) ([]int64, error)

    // Applicable only to IVF indexes: Returns the centroid IDs in the selector in
    // decreasing order of proximity to query 'x' and their distance from 'x'
    ObtainClustersWithDistancesFromIVFIndex(x []float32, centroids Selector, numCentroids int64) (
        []int64, []float32, error)

    // Applicable only to IVF indexes: Returns the top k centroid cardinalities and
    // their vectors in chosen order (descending or ascending)
    ObtainKCentroidCardinalitiesFromIVFIndex(limit int, descending bool) ([]uint64, [][]float32, 
error)

    // fetch centroid count
    Nlist() int

    // Search queries the index with the vectors in x.
    // Returns the IDs of the k nearest neighbors for each query vector and the
    // corresponding distances.
    Search(x []float32, k int64) (distances []float32, labels []int64, err error)

    // SearchWithOptions performs a search with additional optional constraints.
    // - Selector can be used to restrict the search to a subset of the indexed vectors based on 
their IDs.
    // - params is a JSON object that can contain additional search parameters specific to the 
index type, such as IVF search parameters.
    SearchWithOptions(x []float32, k int64, sel Selector, params json.RawMessage) (distances 
[]float32, labels []int64, err error)

    // Applicable only to IVF indexes: Search clusters whose IDs are in eligibleCentroidIDs
    SearchClustersFromIVFIndex(eligibleCentroidIDs []int64, centroidDis []float32, centroidsToProbe 
int,
        x []float32, k int64, include Selector, params json.RawMessage) ([]float32, []int64, error)

    Reconstruct(key int64) ([]float32, error)

    ReconstructBatch(keys []int64, recons []float32) ([]float32, error)

    MergeFrom(other Index, add_id int64) error

    // RangeSearch queries the index with the vectors in x.
    // Returns all vectors with distance < radius.
    RangeSearch(x []float32, radius float32) (*RangeSearchResult, error)

    // DistCompute computes the distance between the query vector and the vectors specified by ids.
    DistCompute(x []float32, labels []int64) ([]float32, error)

    // Reset removes all vectors from the index.
    Reset() error

    // RemoveIDs removes the vectors specified by sel from the index.
    // Returns the number of elements removed and error.
    RemoveIDs(sel *IDSelector) (int, error)

    // Close frees the memory used by the index.
    Close()

    // Size estimates the memory footprint of the index in bytes,
    // if the underlying faiss index is memory-mapped and not fully loaded into memory.
    Size() uint64

    // set the quantizers from a source index into this index, applicable only
    // for IVF indexes
    SetQuantizers(source Index) error

    // CodeSize returns the size of the produced codes in bytes.
    CodeSize() (uint64, error)
    // contains filtered or unexported methods
}
```

Index is a Faiss index.

Note that some index implementations do not support all methods. Check the Faiss wiki to see what 
operations an index supports.

#### type IndexFlat ¶

```
type IndexFlat struct {
    Index
}
```

IndexFlat is an index that stores the full vectors and performs exhaustive search.

#### func NewIndexFlat ¶

```
func NewIndexFlat(d int, metric int) (*IndexFlat, error)
```

NewIndexFlat creates a new flat index.

#### func NewIndexFlatIP ¶

```
func NewIndexFlatIP(d int) (*IndexFlat, error)
```

NewIndexFlatIP creates a new flat index with the inner product metric type.

#### func NewIndexFlatL2 ¶

```
func NewIndexFlatL2(d int) (*IndexFlat, error)
```

NewIndexFlatL2 creates a new flat index with the L2 metric type.

#### func (\*IndexFlat) Xb ¶

```
func (idx *IndexFlat) Xb() []float32
```

Xb returns the index's vectors. The returned slice becomes invalid after any add or remove 
operation.

#### type IndexImpl ¶

```
type IndexImpl struct {
    Index
}
```

IndexImpl is an abstract structure for an index.

#### added in v1.0.31

```
func CloneToCPU(_ *GPUIndexImpl) (*IndexImpl, error)
```

CloneToCPU is not available without the gpu build tag.

#### func IndexFactory ¶

```
func IndexFactory(d int, description string, metric int) (*IndexImpl, error)
```

IndexFactory builds a composite index. description is a comma-separated list of components.

#### func ReadIndex ¶

```
func ReadIndex(filename string, ioflags int) (*IndexImpl, error)
```

ReadIndex reads an index from a file.

#### added in v1.0.0

```
func ReadIndexFromBuffer(buf []byte, ioflags int) (*IndexImpl, error)
```

#### type ParameterSpace ¶

```
type ParameterSpace struct {
    // contains filtered or unexported fields
}
```

#### func NewParameterSpace ¶

```
func NewParameterSpace() (*ParameterSpace, error)
```

NewParameterSpace creates a new ParameterSpace.

#### func (\*ParameterSpace) Delete ¶

```
func (p *ParameterSpace) Delete()
```

Delete frees the memory associated with p.

#### func (\*ParameterSpace) SetIndexParameter ¶

```
func (p *ParameterSpace) SetIndexParameter(idx Index, name string, val float64) error
```

SetIndexParameter sets one of the parameters.

#### type RangeSearchResult ¶

```
type RangeSearchResult struct {
    // contains filtered or unexported fields
}
```

RangeSearchResult is the result of a range search.

#### func (\*RangeSearchResult) Delete ¶

```
func (r *RangeSearchResult) Delete()
```

Delete frees the memory associated with r.

#### func (\*RangeSearchResult) Labels ¶

```
func (r *RangeSearchResult) Labels() (labels []int64, distances []float32)
```

Labels returns the unsorted IDs and respective distances for each query. The result for query i is 
labels\[lims\[i\]:lims\[i+1\]\].

#### func (\*RangeSearchResult) Lims ¶

```
func (r *RangeSearchResult) Lims() []int
```

Lims returns a slice containing start and end indices for queries in the distances and labels 
slices returned by Labels.

#### func (\*RangeSearchResult) Nq ¶

```
func (r *RangeSearchResult) Nq() int
```

Nq returns the number of queries.

#### added in v1.0.20

```
type SearchParams struct {
    // contains filtered or unexported fields
}
```

#### added in v1.0.28

```
func NewBinarySearchParams(idx BinaryIndex, params json.RawMessage, selector Selector,
    defaultParams *defaultSearchParamsIVF) (*SearchParams, error)
```

#### added in v1.0.20

```
func NewSearchParams(idx Index, params json.RawMessage, selector Selector,
    defaultParams *defaultSearchParamsIVF) (*SearchParams, error)
```

Returns a valid SearchParams object, configured according to the provided parameters and selector. 
The returned SearchParams object is allocated, thus caller must clean up the object by invoking 
Delete() method.

#### added in v1.0.27

```
func NewStandardSearchParams(selector Selector) (*SearchParams, error)
```

Returns a standard SearchParams object without any special settings with the provided selector. The 
returned SearchParams object is allocated, thus caller must clean up the object by invoking 
Delete() method.

#### added in v1.0.20

```
func (s *SearchParams) Delete()
```

Delete frees the memory associated with s.

#### added in v1.0.23

```
type Selector interface {
    ExcludeFilter() bool
    Get() *C.FaissIDSelector
    Delete()
}
```

Note: currently we have only one implementation, but we keep the interface for future extensibility

#### func NewIDSelectorBatch ¶

```
func NewIDSelectorBatch(indices []int64) (Selector, error)
```

NewIDSelectorBatch creates a new batch selector.

#### added in v1.0.27

```
func NewIDSelectorBatchNot(exclude []int64) (Selector, error)
```

NewIDSelectorBatchNot creates a new Not selector, wrapped around a batch selector, with the IDs in 
'exclude'.

#### added in v1.0.27

```
func NewIDSelectorBitmap(bitmap []byte) (Selector, error)
```

NewIDSelectorBitmap creates a selector using a bitset, where each bit indicates whether the 
corresponding ID is to be selected. NOTE: This function assumes that len(bitmap)\*8 covers the full 
range of IDs in the index, and only works when we have vector IDs ranging from 0 to N-1, where N is 
the number of vectors in the index. The length of the bitmap should be at least ceil(N/8).

#### added in v1.0.27

```
func NewIDSelectorBitmapNot(bitmap []byte) (Selector, error)
```

NewIDSelectorBitmapNot creates a NOT selector using a bitset, where each bit indicates whether the 
corresponding ID is NOT to be selected. NOTE: This function assumes that len(bitmap)\*8 covers the 
full range of IDs in the index, and only works when we have vector IDs ranging from 0 to N-1, where 
N is the number of vectors in the index. The length of the bitmap should be at least ceil(N/8).

#### func NewIDSelectorRange ¶

```
func NewIDSelectorRange(imin, imax int64) (Selector, error)
```

NewIDSelectorRange creates a selector that removes IDs on \[imin, imax).

## Directories

| Path | Synopsis |
| --- | --- |
| \_example |  |
| [flat](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/flat) command  Usage 
example for IndexFlat. | Usage example for IndexFlat. |
| [hnsw](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/hnsw) command |  |
| [io](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/io) command |  |
| [ivfflat](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/ivfflat) command  
Usage example for IndexIVFFlat. | Usage example for IndexIVFFlat. |
| [ivfpq](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/ivfpq) command  Usage 
example for IndexIVFPQ. | Usage example for IndexIVFPQ. |
| [misc](https://pkg.go.dev/github.com/blevesearch/go-faiss@v1.1.3/_example/misc) command  Usage 
example for IndexIVFFlat. | Usage example for IndexIVFFlat. |
