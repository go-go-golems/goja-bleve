---
Title: goja-bleve API Design and Implementation Guide
Ticket: RAGEVAL-GOJA-RAG-STRATEGIES
Status: active
Topics:
    - goja
    - bleve
    - search
    - embeddings
    - rag
    - api-design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: 2026-05-27--rag-evaluation-system/internal/services/search/bm25.go
      Note: Current bleve usage pattern
    - Path: 2026-05-27--rag-evaluation-system/internal/services/search/hybrid.go
      Note: Current manual RRF fusion
    - Path: 2026-05-27--rag-evaluation-system/internal/services/search/vector.go
    - Path: bleve/document/field_vector.go
    - Path: bleve/index.go
      Note: Bleve Index interface
    - Path: bleve/mapping.go
    - Path: bleve/mapping/field.go
    - Path: bleve/mapping/index.go
    - Path: bleve/mapping/mapping_vectors.go
    - Path: bleve/mapping_vector.go
      Note: NewVectorFieldMapping
    - Path: bleve/query.go
      Note: All query factory functions
    - Path: bleve/rescorer.go
      Note: RRF/RSF fusion rescoring
    - Path: bleve/search.go
    - Path: bleve/search/query/knn.go
    - Path: bleve/search_knn.go
      Note: KNNRequest
    - Path: geppetto/pkg/js/modules/geppetto/api_engine_builder.go
      Note: Builder pattern with cloneFor
    - Path: geppetto/pkg/js/modules/geppetto/api_inference_settings.go
    - Path: geppetto/pkg/js/modules/geppetto/api_schema_builders.go
    - Path: geppetto/pkg/js/modules/geppetto/module.go
      Note: moduleRuntime pattern
    - Path: go-go-goja/modules/common.go
      Note: NativeModule interface
    - Path: go-go-goja/modules/database/database.go
    - Path: go-go-goja/modules/exports.go
    - Path: go-go-goja/modules/typing.go
    - Path: go-go-goja/pkg/xgoja/providerapi/module.go
    - Path: go-go-goja/pkg/xgoja/providerapi/registry.go
    - Path: goja-text/pkg/extract/module.go
    - Path: goja-text/pkg/extract/types.go
ExternalSources:
    - sources/bleve-github-readme.md
    - sources/bleve-index-mapping.md
    - sources/bleve-pkg-go-dev.md
    - sources/bleve-query.md
    - sources/blevesearch-homepage.md
    - sources/rag-eval-bm25-search.go.md
    - sources/rag-eval-hybrid-search.go.md
Summary: Comprehensive design for goja-bleve, a Go-backed native module exposing bleve full-text and vector search through a fluent builder-pattern JavaScript API.
LastUpdated: 2026-06-02T21:30:00-04:00
WhatFor: Implement goja-bleve module with type-safe Go-backed objects for bleve index management, document indexing, BM25 search, vector/KNN search, and hybrid retrieval
WhenToUse: When implementing the goja-bleve native module, understanding bleve's Go API surface, or designing the JavaScript-facing builder API
---


# goja-bleve API Design and Implementation Guide

## Executive Summary

This document specifies **goja-bleve**, a native Go module that exposes the Bleve full-text and vector search library to JavaScript runtimes built on goja. The module follows the established builder-pattern conventions from the geppetto and goja-text modules, where Go-backed reference objects are attached to JavaScript wrapper objects for runtime type safety. The design prioritizes vector similarity search (KNN) and hybrid retrieval—core operations in RAG pipelines—while also exposing the full breadth of Bleve's query, mapping, and indexing capabilities.

The module will be consumed by RAG evaluation scripts that run inside the goja engine, replacing the current Go-native search service layer (`internal/services/search/bm25.go`, `vector.go`, `hybrid.go`) with equivalent JavaScript-callable primitives. This lets pipeline authors compose retrieval strategies directly in JS, iterate faster, and leverage the existing goja module ecosystem (geppetto for embeddings, goja-text for chunking, database for metadata, express for serving web UI).

## Problem Statement

The RAG evaluation system currently implements search operations—BM25 index building, vector similarity search, hybrid RRF fusion—as compiled Go services in `internal/services/search/`. Each new retrieval strategy requires writing, compiling, and deploying Go code. This creates a slow iteration loop for experiment-driven RAG evaluation.

Meanwhile, the goja JavaScript runtime already provides a rich module ecosystem:

- **geppetto** — LLM inference and embeddings
- **goja-text** — text extraction and chunking
- **database** — SQL access for metadata
- **express** — HTTP serving for custom web UI
- **fs / yaml / etc.** — standard file and config operations

What is missing is a **search primitive**: the ability to create Bleve indexes, index documents, run text queries, perform vector KNN searches, and combine results with hybrid fusion—all from JavaScript. The goja-bleve module fills this gap.

### Goals

- Expose Bleve's core operations as a fluent, builder-pattern JavaScript API
- Provide first-class support for vector/KNN search, since embeddings are produced by geppetto
- Support hybrid BM25 + KNN retrieval with RRF/RSF fusion (matching bleve's built-in rescoring)
- Ensure type safety at runtime via Go-backed reference objects (not raw JS hashmaps)
- Follow the established goja module conventions (NativeModule, TypeScriptDeclarer, provider registration)

### Non-goals

- Re-implementing bleve's internal search algorithms in JS
- Exposing every bleve analysis/registry customization (custom analyzers, char filters, etc.) in v1
- Supporting multi-index aliases or distributed search in v1
- Providing a query-string parser for end-user search input (can be added later)

---

## Background: What Is Bleve?

Bleve is a Go library for full-text and vector search. It provides:

- **Index management** — create, open, close persistent or memory-only indexes
- **Document mapping** — declarative schemas that define how Go structs or JSON objects map to indexed fields
- **Field types** — text, number, datetime, boolean, geopoint, geoshape, IP, and **vector** (with the `vectors` build tag)
- **Text queries** — match, term, phrase, prefix, fuzzy, regexp, wildcard, boolean combinations, query strings
- **Vector search** — KNN (k-nearest neighbors) search on vector fields, with configurable similarity metrics (cosine, dot product, euclidean) and index optimizations (flat, IVF, BIVF)
- **Hybrid search** — combining text and KNN queries with score fusion (RRF: reciprocal rank fusion, RSF: relative score fusion)
- **Batch indexing** — grouping many index/delete operations for throughput
- **Facets** — aggregation buckets for numeric ranges, date ranges, and term counts
- **Highlighting** — marking matching fragments in stored field content

Bleve uses the **scorch** index type by default (an upsidedown-based KV store) and supports in-memory indexes for testing. When compiled with the `vectors` build tag, it links against `go-faiss` for efficient vector similarity search.

### Key Bleve Concepts for the Implementer

**IndexMapping** is the top-level schema object. It contains a `DefaultMapping` (a `DocumentMapping`), type-specific document mappings, and configuration for analyzers, date parsers, and synonym sources. When you index a document, bleve uses the mapping to decide which fields to index, how to analyze them, and whether to store the original values.

**DocumentMapping** describes how a particular document type maps to fields. It can contain `FieldMappings` (for leaf fields) and nested `SubDocumentMappings` (for object properties). The `Dynamic` flag controls whether unmapped fields are indexed using default rules or ignored.

**FieldMapping** describes a single field: its type (text, number, datetime, boolean, vector, vector_base64), the analyzer to use, whether to store/index/include term vectors, and—for vector fields—the dimensionality (`Dims`), similarity metric (`Similarity`), and vector index optimization (`VectorIndexOptimizedFor`).

**SearchRequest** bundles a query, pagination (Size, From), field loading, highlighting, facets, sort order, and—critically for us—KNN requests and a KNN operator (AND/OR). When KNN requests are present alongside a text query, bleve runs a two-phase search: first collect KNN hits, then run the text query, then merge/rescore.

**KNNRequest** specifies a vector field name, a query vector (float32 array or base64 string), k (number of nearest neighbors), an optional boost, optional search parameters (e.g., IVF nprobe), and an optional filter query for pre-filtering.

**SearchResult** contains hits (each with ID, Score, Fields, Fragments, ScoreBreakdown), total hit count, max score, duration, facet results, and status information.

**Batch** groups multiple Index and Delete operations for efficient writes.

---

## Architecture Overview

### Module Position in the goja Ecosystem

```
┌─────────────────────────────────────────────────┐
│              JavaScript Runtime (goja)            │
│                                                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │ geppetto │ │goja-text │ │   goja-bleve     │  │
│  │          │ │          │ │                  │  │
│  │ - LLM    │ │ - chunk  │ │ - open/create    │  │
│  │ - embed  │ │ - extract│ │ - index docs     │  │
│  │          │ │          │ │ - BM25 search    │  │
│  └────┬─────┘ └──────────┘ │ - KNN search     │  │
│       │                    │ - hybrid search   │  │
│       │ embeddings         │ - batch ops       │  │
│       └────────────────────│ - mapping config  │  │
│                            └────────┬─────────┘  │
│                                     │            │
│  ┌──────────┐                       │            │
│  │ database │───────────────────────┘            │
│  │          │  (metadata, chunk records)         │
│  └──────────┘                                   │
│                                                   │
│  ┌──────────┐                                   │
│  │ express  │  (serve web UI for search results) │
│  └──────────┘                                   │
└─────────────────────────────────────────────────┘
```

The goja-bleve module is loaded via `require("bleve")` and provides factory functions that return Go-backed builder objects. Each builder accumulates configuration in Go structs, then produces a Go-backed "product" object (an open index, a search request, a query) when `.build()` is called. This is the same pattern used by `geppetto.engine().inference(settings).build()` and `extract.options().Formats("json").Build()`.

### Go-Side Object Lifecycle

All complex objects (indexes, search requests, queries, batches, search results) are Go structs stored as hidden references on JavaScript wrapper objects, following the geppetto convention:

1. A Go struct (`indexRef`, `queryRef`, `searchRequestRef`, etc.) holds the real data
2. The struct is attached to a JS object via a hidden `__bleve_ref` property (non-enumerable, non-writable, non-configurable)
3. JS methods on the wrapper read/write the Go struct through the hidden reference
4. Clone operations create new Go structs + new JS wrappers (immutable builder pattern)
5. On `.build()`, the accumulated Go struct is consumed to produce the final Go object (e.g., `bleve.Index`, `bleve.SearchRequest`)

This approach gives us:
- **Type safety**: JS code cannot construct invalid query trees or mappings; the Go side validates
- **Performance**: no serialization/deserialization of large vectors or result sets through JSON
- **Ergonomics**: method-chaining builders feel natural in JS
- **Debuggability**: each wrapper has a `.toJSON()` method for inspection

---

## JavaScript API Design

The module is accessed as `require("bleve")`. The top-level exports are namespace objects and factory functions.

### Top-Level Namespace

```
require("bleve") → {
  // ── Index lifecycle ──────────────────────────────
  open(path)                    → Index          // open existing
  create(path)                  → IndexBuilder   // builder for new index
  memory()                      → IndexBuilder   // builder for in-memory index

  // ── Mapping factories ────────────────────────────
  mapping()                     → MappingBuilder
  docMapping()                  → DocMappingBuilder
  field()                       → FieldBuilder

  // ── Query factories ──────────────────────────────
  match(text)                   → MatchQuery
  matchPhrase(text)             → MatchPhraseQuery
  term(text)                    → TermQuery
  prefix(text)                  → PrefixQuery
  fuzzy(text)                   → FuzzyQuery
  wildcard(pattern)             → WildcardQuery
  regexp(pattern)               → RegexpQuery
  queryString(q)                → QueryStringQuery
  matchAll()                    → MatchAllQuery
  matchNone()                   → MatchNoneQuery
  bool()                        → BoolQuery
  conj(...queries)              → ConjQuery
  disj(...queries)              → DisjQuery
  termRange()                   → TermRangeBuilder
  numRange()                    → NumRangeBuilder
  dateRange()                   → DateRangeBuilder
  docIDs(ids)                   → DocIDQuery
  ipRange(cidr)                 → IPRangeQuery

  // ── Search request ───────────────────────────────
  search(query)                 → SearchRequestBuilder

  // ── KNN / Vector search ──────────────────────────
  knn(field, vector, k)         → KNNRequest

  // ── Version ──────────────────────────────────────
  version                       → string
}
```

### Design Decision: Namespace vs. Flat Functions

The geppetto module uses a flat namespace (`gp.engine()`, `gp.schema.string()`, `gp.tool()`) while goja-text uses a nested approach (`extract.options()`, `extract.markdownCodeBlocks()`). For goja-bleve, we adopt a **flat namespace** at the top level (like geppetto) because the concept space is coherent: everything is about index + search. The query factories are all at the top level since they are small, leaf objects that don't need a separate namespace.

---

### Index Lifecycle

#### `bleve.open(path)` → `Index`

Opens an existing bleve index at the given filesystem path. Returns a Go-backed `Index` wrapper object.

- If the path does not exist or is not a valid bleve index, throws a GoError
- The index is open and ready for search/index operations immediately
- Call `.close()` when done to release resources

Pseudocode for the `Index` wrapper object:

```
Index {
  // ── Document operations ────────────────
  index(id, doc)                → void          // index a document (JS object → Go struct via mapping)
  delete(id)                    → void          // delete by ID
  get(id)                       → object | null // retrieve stored document fields
  
  // ── Batch operations ───────────────────
  newBatch()                    → Batch
  batch(b)                      → void          // execute batch
  
  // ── Search ─────────────────────────────
  search(req)                   → SearchResult  // run a SearchRequest
  searchInContext(req)          → SearchResult  // with context awareness
  
  // ── Metadata ──────────────────────────
  docCount()                    → number
  fields()                      → string[]
  mapping()                     → MappingSnapshot  // read-only view of current mapping
  name()                        → string
  
  // ── Lifecycle ─────────────────────────
  close()                       → void
  
  // ── Internal ──────────────────────────
  // Go ref: *bleve.Index (the real open index handle)
  // toJSON() → { name, docCount, fields }
}
```

**Key point**: `index(id, doc)` accepts a plain JavaScript object. Goja converts it to `map[string]any` which bleve's `Index()` method accepts natively—bleve walks the map using the configured mapping to decide what to index. This is the same approach used in the existing `bm25.go` service.

#### `bleve.create(path)` → `IndexBuilder`

Builder for creating a new persistent index. The builder accumulates configuration, then `.build()` creates the index and returns an open `Index`.

```
IndexBuilder {
  mapping(m)                    → IndexBuilder  // set IndexMapping (from MappingBuilder.build())
  using(indexType, kvStore)     → IndexBuilder  // override defaults (advanced)
  build()                       → Index         // create index, return open handle
  
  // If no mapping is set, a default IndexMapping is used (all fields dynamic, text type)
}
```

#### `bleve.memory()` → `IndexBuilder`

Same as `create()` but creates an in-memory-only index. The `.build()` call produces an `Index` that is not persisted to disk. Useful for tests, temporary indexes, and pipelined RAG workflows where the index is built fresh each run.

---

### Mapping API

Mapping configuration defines the schema of the index: what fields exist, what types they are, how they are analyzed, and—for vector fields—what dimensions and similarity metric to use.

#### `bleve.mapping()` → `MappingBuilder`

Top-level index mapping builder. Configures the default document mapping, type-specific mappings, and analysis defaults.

```
MappingBuilder {
  defaultMapping(dm)            → MappingBuilder   // set the default DocMapping
  addTypeMapping(type, dm)      → MappingBuilder   // register a type-specific DocMapping
  defaultAnalyzer(name)         → MappingBuilder   // e.g., "standard", "keyword"
  defaultField(name)            → MappingBuilder   // default field for queries ("_all")
  storeDynamic(bool)            → MappingBuilder   // store unmapped fields?
  indexDynamic(bool)            → MappingBuilder   // index unmapped fields?
  build()                       → IndexMapping      // produces the Go IndexMappingImpl ref
}
```

#### `bleve.docMapping()` → `DocMappingBuilder`

Document mapping builder. Describes the fields and sub-documents for one document type.

```
DocMappingBuilder {
  addField(name, fm)            → DocMappingBuilder  // add a FieldMapping under a name
  addSubDoc(name, dm)           → DocMappingBuilder  // nested object mapping
  dynamic(bool)                 → DocMappingBuilder  // index unmapped sub-fields?
  enabled(bool)                 → DocMappingBuilder  // disable mapping entirely?
  build()                       → DocMapping         // produces Go DocumentMapping ref
}
```

#### `bleve.field()` → `FieldBuilder`

Field mapping builder. The most granular configuration unit. This is where vector field configuration lives.

```
FieldBuilder {
  // ── Type selectors (one required) ────────────────────
  text()                        → FieldBuilder
  keyword()                     → FieldBuilder   // text with keyword analyzer (no tokenization)
  number()                      → FieldBuilder
  datetime()                    → FieldBuilder
  boolean()                     → FieldBuilder
  vector(dims)                  → FieldBuilder   // vector field, dims required
  vectorBase64(dims)            → FieldBuilder   // base64-encoded vector field
  
  // ── Common options ──────────────────────────────────
  name(name)                    → FieldBuilder   // override field name
  analyzer(name)                → FieldBuilder   // analyzer name for text fields
  store(bool)                   → FieldBuilder   // store original value in index?
  index(bool)                   → FieldBuilder   // index for search?
  docValues(bool)               → FieldBuilder   // enable doc values (facets/sort)?
  
  // ── Vector-specific options ──────────────────────────
  similarity(metric)            → FieldBuilder   // "cosine" | "dot_product" | "euclidean"
  optimizedFor(strategy)        → FieldBuilder   // "flat" | "ivf" | "bivf-sq8"
  
  // ── Terminal ────────────────────────────────────────
  build()                       → FieldMapping    // produces Go FieldMapping ref
}
```

**Vector field example** (pseudocode):

```javascript
const bleve = require("bleve")

const embeddingField = bleve.field()
  .vector(1536)              // OpenAI ada-002 dimensions
  .similarity("cosine")
  .optimizedFor("ivf")
  .store(false)
  .build()

const textField = bleve.field()
  .text()
  .analyzer("standard")
  .store(true)
  .build()

const chunkMapping = bleve.docMapping()
  .addField("text", textField)
  .addField("embedding", embeddingField)
  .addField("source_id", bleve.field().keyword().build())
  .addField("chunk_index", bleve.field().number().build())
  .build()

const idxMapping = bleve.mapping()
  .defaultMapping(chunkMapping)
  .build()
```

**Important**: Vector fields require the `vectors` build tag when compiling the Go binary that hosts the goja runtime. Without it, `bleve.field().vector()` will throw an error at `.build()` time. This is a bleve-level constraint, not a goja-bleve constraint.

---

### Query API

Queries are immutable Go-backed objects. Each query factory returns a wrapper whose methods return the same wrapper (mutating in place) or a new wrapper (for compound queries). Queries are consumed by `SearchRequestBuilder.query()`.

#### Leaf Queries

| Factory | Go equivalent | Description |
|---|---|---|
| `bleve.match(text)` | `NewMatchQuery` | Full-text match, analyzer-based |
| `bleve.matchPhrase(text)` | `NewMatchPhraseQuery` | Phrase match, token order matters |
| `bleve.term(text)` | `NewTermQuery` | Exact term match |
| `bleve.prefix(text)` | `NewPrefixQuery` | Prefix match |
| `bleve.fuzzy(text)` | `NewFuzzyQuery` | Levenshtein-distance match |
| `bleve.wildcard(pattern)` | `NewWildcardQuery` | Wildcard pattern (* and ?) |
| `bleve.regexp(pattern)` | `NewRegexpQuery` | Regular expression match |
| `bleve.queryString(q)` | `NewQueryStringQuery` | Lucene-style query string |
| `bleve.matchAll()` | `NewMatchAllQuery` | Matches all documents |
| `bleve.matchNone()` | `NewMatchNoneQuery` | Matches no documents |
| `bleve.ipRange(cidr)` | `NewIPRangeQuery` | IP/CIDR match |

All leaf query wrappers support:

```
.setField(name)       → self     // restrict to a specific field
.setBoost(n)         → self     // boost this query's contribution
```

#### Compound Queries

- **`bleve.bool()`** → `BoolQuery` with methods `.addMust(q)`, `.addShould(q)`, `.addMustNot(q)`. This is the primary way to combine queries. Maps to `NewBooleanQuery`.

- **`bleve.conj(q1, q2, ...)`** → `ConjQuery`. All sub-queries must match. Maps to `NewConjunctionQuery`.

- **`bleve.disj(q1, q2, ...)`** → `DisjQuery`. At least one sub-query must match. Maps to `NewDisjunctionQuery`.

#### Range Queries (Builders)

Range queries use a builder because they have multiple optional parameters:

```
bleve.termRange()
  .min(val)           → self     // lower bound (inclusive by default)
  .max(val)           → self     // upper bound (exclusive by default)
  .field(name)        → self
  .inclusiveMin(bool) → self
  .inclusiveMax(bool) → self
  .build()            → TermRangeQuery

bleve.numRange()
  .min(val)           → self     // float64
  .max(val)           → self     // float64
  .field(name)        → self
  .inclusiveMin(bool) → self
  .inclusiveMax(bool) → self
  .build()            → NumericRangeQuery

bleve.dateRange()
  .start(val)         → self     // string or Date
  .end(val)           → self
  .field(name)        → self
  .inclusiveStart(bool) → self
  .inclusiveEnd(bool)   → self
  .build()            → DateRangeQuery
```

---

### KNN / Vector Search API

This is the centerpiece of the module for RAG use cases.

#### `bleve.knn(field, vector, k)` → `KNNRequest`

Creates a KNN (k-nearest neighbors) search request. The `vector` parameter accepts a JavaScript `Float32Array` or a plain array of numbers (the Go side converts to `[]float32`). `k` is the number of nearest neighbors to return.

```
KNNRequest {
  boost(n)                      → KNNRequest   // boost for score fusion weighting
  filter(query)                 → KNNRequest   // pre-filter: only search docs matching query
  params(obj)                   → KNNRequest   // search params (e.g., IVF nprobe)
  
  // Internal: produces *bleve.KNNRequest for the SearchRequest
}
```

**Integration with geppetto embeddings**:

The typical RAG workflow is: use geppetto to embed the query text, then pass the embedding vector to bleve for KNN search. The goja-bleve module does not call geppetto itself—instead, the JS script wires them:

```javascript
// pseudocode: vector search workflow
const bleve = require("bleve")
const gp = require("geppetto")

const queryText = "what is cosine similarity?"

// 1. Embed the query using geppetto
const embedder = gp.engine()
  .inference(embeddingSettings)
  .build()
const vector = await embedder.embed(queryText)  // returns Float32Array-like

// 2. Open index and search
const idx = bleve.open("data/indexes/my-corpus")
const results = idx.search(
  bleve.search(bleve.matchNone())
    .knn(bleve.knn("embedding", vector, 10))
    .size(10)
    .fields("text", "source_id", "chunk_index")
    .build()
)
```

---

### Search Request API

#### `bleve.search(query)` → `SearchRequestBuilder`

The search request builder is the primary entry point for querying an index. It takes a query (from any of the query factories) and accumulates search parameters.

```
SearchRequestBuilder {
  // ── Query ─────────────────────────────────
  query(q)                      → self    // set/replace the text query
  
  // ── Pagination ────────────────────────────
  size(n)                       → self    // results per page (default 10)
  from(n)                       → self    // offset (default 0)
  
  // ── Fields & highlighting ─────────────────
  fields(...names)              → self    // stored fields to return in hits
  highlight(style?)             → self    // enable highlighting (optional style name)
  highlightFields(...names)     → self    // restrict highlighting to specific fields
  
  // ── KNN ───────────────────────────────────
  knn(request)                  → self    // add a KNNRequest (can call multiple times)
  knnOperator(op)              → self    // "and" | "or" (default "or")
  
  // ── Score fusion ──────────────────────────
  score(mode)                   → self    // "rrf" | "rsf" | "none" | "" (default)
  scoreRankConstant(n)          → self    // RRF k parameter (default 60)
  scoreWindowSize(n)            → self    // fusion window size (default from+size)
  
  // ── Sort ──────────────────────────────────
  sortBy(...fields)             → self    // sort by field names ("-_score" for desc)
  
  // ── Explain ───────────────────────────────
  explain(bool)                 → self    // include score explanation
  
  // ── Facets ────────────────────────────────
  // (v2: facet support, not in initial release)
  
  // ── Terminal ──────────────────────────────
  build()                       → SearchRequest   // Go-backed SearchRequest ref
}
```

**Hybrid search example** (pseudocode):

```javascript
// Hybrid BM25 + KNN with RRF fusion
const results = idx.search(
  bleve.search(
    bleve.disj(
      bleve.match(queryText).setField("text"),
      bleve.match(queryText).setField("title").setBoost(2.0)
    )
  )
  .knn(bleve.knn("embedding", queryVector, 20))
  .knnOperator("or")
  .score("rrf")
  .scoreRankConstant(60)
  .size(10)
  .fields("text", "title", "source_id", "chunk_index")
  .explain(true)
  .build()
)
```

This maps directly to the two-phase search in bleve's `search_knn.go`:

1. **Phase 1 (pre-search)**: Collect KNN hits from the vector index
2. **Phase 2 (main search)**: Run the text query, then merge with KNN hits using the rescoring algorithm (RRF or RSF)

---

### Batch API

#### `idx.newBatch()` → `Batch`

Batches group many index/delete operations for efficient writes.

```
Batch {
  index(id, doc)                → self     // add document to batch
  delete(id)                    → self     // add deletion to batch
  size()                        → number   // number of operations in batch
  reset()                       → self     // clear batch for reuse
}
```

**Batch indexing example** (pseudocode):

```javascript
const batch = idx.newBatch()
for (const chunk of chunks) {
  batch.index(chunk.id, {
    text: chunk.text,
    source_id: chunk.sourceId,
    chunk_index: chunk.index,
    embedding: chunk.vector   // Float32Array from geppetto
  })
  if (batch.size() >= 500) {
    idx.batch(batch)
    batch.reset()
  }
}
if (batch.size() > 0) {
  idx.batch(batch)
}
```

This mirrors the existing batch pattern in `bm25.go` (lines 82–96), but expressed in JavaScript.

---

### Search Result API

When `idx.search(req)` is called, it returns a Go-backed `SearchResult` wrapper:

```
SearchResult {
  total                         → number       // total matching documents
  maxScore                      → number       // highest score
  took                          → number       // duration in milliseconds
  hits                          → Hit[]        // array of hit wrappers
  
  toJSON()                      → object       // serializable snapshot
}

Hit {
  id                            → string       // document ID
  score                         → number       // relevance score
  fields                        → object       // requested stored fields
  fragments                     → object       // highlighted fragments (if requested)
  scoreBreakdown                → number[]     // per-component scores (hybrid)
  explanation                   → object       // score explanation tree (if explain=true)
  index                         → string       // index name
  
  toJSON()                      → object       // serializable snapshot
}
```

The `hits` array and `Hit` objects are Go-backed wrappers that read directly from the bleve `search.DocumentMatch` structs—no intermediate JSON serialization. The `.toJSON()` methods produce plain JS objects for serialization or logging.

---

### Complete RAG Pipeline Example

This example shows a full RAG pipeline in JavaScript, using all the goja modules together:

```javascript
// pseudocode: full RAG evaluation pipeline

const bleve = require("bleve")
const gp = require("geppetto")
const extract = require("extract")
const db = require("database")
const fs = require("fs")
const yaml = require("yaml")

// ── Step 1: Configure embedding engine ──────────────────────
const embeddingSettings = /* ... resolved from gp.inferenceProfiles ... */
const embedder = gp.engine().inference(embeddingSettings).build()

// ── Step 2: Build or open the search index ──────────────────
const indexPath = "data/indexes/corpus-001"
let idx

if (fs.exists(indexPath)) {
  idx = bleve.open(indexPath)
} else {
  const mapping = bleve.mapping()
    .defaultMapping(
      bleve.docMapping()
        .addField("text", bleve.field().text().analyzer("standard").store(true).build())
        .addField("title", bleve.field().text().analyzer("standard").store(true).build())
        .addField("embedding", bleve.field().vector(1536).similarity("cosine").optimizedFor("ivf").build())
        .addField("source_id", bleve.field().keyword().store(true).build())
        .addField("chunk_index", bleve.field().number().store(true).build())
        .build()
    )
    .build()
  
  idx = bleve.create(indexPath).mapping(mapping).build()
  
  // ── Step 3: Load documents, chunk, embed, index ───────────
  const docs = db.query("SELECT id, text, source_id, chunk_index FROM chunks")
  const batch = idx.newBatch()
  
  for (const doc of docs) {
    const vector = await embedder.embed(doc.text)
    batch.index(doc.id, {
      text: doc.text,
      source_id: doc.source_id,
      chunk_index: doc.chunk_index,
      embedding: vector
    })
    if (batch.size() >= 500) {
      idx.batch(batch)
      batch.reset()
    }
  }
  if (batch.size() > 0) idx.batch(batch)
}

// ── Step 4: Search ──────────────────────────────────────────
const queryText = "what is cosine similarity?"
const queryVector = await embedder.embed(queryText)

// BM25-only search
const bm25Results = idx.search(
  bleve.search(
    bleve.disj(
      bleve.match(queryText).setField("text"),
      bleve.match(queryText).setField("title").setBoost(2.0)
    )
  )
  .size(10)
  .fields("text", "title", "source_id", "chunk_index")
  .build()
)

// Vector-only search
const vectorResults = idx.search(
  bleve.search(bleve.matchNone())
  .knn(bleve.knn("embedding", queryVector, 10))
  .size(10)
  .fields("text", "title", "source_id", "chunk_index")
  .build()
)

// Hybrid search (BM25 + KNN with RRF fusion)
const hybridResults = idx.search(
  bleve.search(
    bleve.disj(
      bleve.match(queryText).setField("text"),
      bleve.match(queryText).setField("title").setBoost(2.0)
    )
  )
  .knn(bleve.knn("embedding", queryVector, 20))
  .knnOperator("or")
  .score("rrf")
  .scoreRankConstant(60)
  .size(10)
  .fields("text", "title", "source_id", "chunk_index")
  .explain(true)
  .build()
)

// ── Step 5: Evaluate ───────────────────────────────────────
// Compare retrieved chunks against ground truth
for (const hit of hybridResults.hits) {
  console.log(`${hit.id} score=${hit.score} fields=${JSON.stringify(hit.fields)}`)
}

idx.close()
```

---

## Go Implementation Architecture

### Package Structure

```
goja-bleve/
├── pkg/
│   ├── module.go                  // NativeModule implementation, Loader, init()
│   ├── api_index.go              // Index wrapper, IndexBuilder, open/create/memory
│   ├── api_mapping.go            // MappingBuilder, DocMappingBuilder, FieldBuilder
│   ├── api_query.go              // All query wrappers (match, term, bool, etc.)
│   ├── api_search.go             // SearchRequestBuilder, SearchResult, Hit
│   ├── api_knn.go                // KNNRequest wrapper
│   ├── api_batch.go              // Batch wrapper
│   ├── codec.go                  // JS value → Go conversion helpers
│   ├── provider.go               // providerapi.Module registration
│   ├── doc.go                    // Package documentation
│   └── logcopter.go              // Logging setup
├── cmd/
│   └── goja-bleve/
│       └── main.go               // Standalone REPL for testing
├── go.mod
├── Makefile
└── AGENT.md
```

### Core Go Types (Pseudocode)

The implementation uses Go struct references attached to JS objects, following the geppetto pattern. Here are the key Go types:

```
// moduleRuntime holds per-VM state (like geppetto's moduleRuntime)
moduleRuntime {
    vm           *goja.Runtime
    indexes      map[string]*indexRef    // track open indexes by name
}

// indexRef wraps an open bleve.Index
indexRef {
    api     *moduleRuntime
    index   bleve.Index
}

// mappingRef wraps a built IndexMappingImpl
mappingRef {
    mapping  *mapping.IndexMappingImpl
}

// docMappingRef wraps a DocumentMapping
docMappingRef {
    dm       *mapping.DocumentMapping
}

// fieldMappingRef wraps a FieldMapping
fieldMappingRef {
    fm       *mapping.FieldMapping
}

// queryRef wraps any bleve query.Query
queryRef {
    query    query.Query
}

// searchRequestRef wraps a built SearchRequest
searchRequestRef {
    req      *bleve.SearchRequest
}

// batchRef wraps a bleve.Batch
batchRef {
    api      *moduleRuntime
    index    bleve.Index
    batch    *bleve.Batch
}

// searchResultRef wraps a bleve.SearchResult
searchResultRef {
    result   *bleve.SearchResult
}

// hitRef wraps a search.DocumentMatch
hitRef {
    hit      *search.DocumentMatch
}
```

### Builder State (Pseudocode)

Builders accumulate partial state before `.build()` finalizes it:

```
// indexBuilderRef accumulates index creation config
indexBuilderRef {
    api         *moduleRuntime
    path        string
    memOnly     bool
    mappingRef  *mappingRef       // nil = use default
    indexType   string            // "" = default
    kvStore     string            // "" = default
}

// mappingBuilderRef accumulates IndexMappingImpl config
mappingBuilderRef {
    api              *moduleRuntime
    defaultMapping   *docMappingRef
    typeMappings     map[string]*docMappingRef
    defaultAnalyzer  string
    storeDynamic     *bool
    indexDynamic     *bool
}

// searchRequestBuilderRef accumulates SearchRequest config
searchRequestBuilderRef {
    api               *moduleRuntime
    queryRef          *queryRef
    size              int
    from              int
    fields            []string
    highlight         *highlightConfig
    knnRequests       []*knnRequestRef
    knnOperator       string
    scoreMode         string
    scoreRankConstant int
    scoreWindowSize   int
    sortBy            []string
    explain           bool
}

// knnRequestRef accumulates KNNRequest config
knnRequestRef {
    field       string
    vector      []float32
    k           int64
    boost       *float64
    filterRef   *queryRef
    params      json.RawMessage
}

// fieldBuilderRef accumulates FieldMapping config
fieldBuilderRef {
    fieldType          string       // "text" | "number" | "vector" | etc.
    dims               int          // vector dims
    similarity         string       // vector similarity metric
    optimizedFor       string       // vector index optimization
    name               string
    analyzer           string
    store              *bool
    index              *bool
    docValues          *bool
    // ... other FieldMapping fields
}
```

### The `attachRef` / `getRef` Pattern

This is the critical bridge between Go and JS. It is copied from geppetto (see `geppetto/pkg/js/modules/geppetto/module.go` lines ~245–270):

- **`attachRef(obj, ref)`**: Stores the Go struct pointer on the JS object under a hidden key (`__bleve_ref`). Makes the property non-enumerable, non-writable, non-configurable.
- **`getRef(value)`**: Retrieves the Go struct pointer from a JS object. Used by methods that accept wrappers as arguments (e.g., `searchRequestBuilder.knn(knnRef)` extracts the `knnRequestRef`).

This pattern ensures that JS code cannot forge or tamper with Go-backed objects. A `queryRef` can only be produced by a query factory, not by constructing a plain JS object.

### Vector Conversion: JS Array → Go `[]float32`

The `codec.go` file handles converting JS arrays/Float32Arrays to Go `[]float32` slices. This is critical for KNN search and vector field indexing. The conversion must:

1. Accept `Float32Array` (goja supports typed arrays)
2. Accept plain JS arrays of numbers
3. Validate length against the mapping's `Dims` when indexing
4. Normalize to unit length when similarity is "cosine" (bleve does this internally, but we can validate early)

### Context Propagation

The goja runtime supports context propagation through `runtimebridge.CurrentOwnerContext(vm)`. Search calls should use `idx.searchInContext()` internally to respect cancellation and timeouts, just as the database module uses `QueryContext` (see `database/database.go` line 78).

---

## Implementation Plan

### Phase 1: Core Index + BM25 Search (foundation)

**Goal**: Open/create indexes, index documents, run text queries.

**Files to create**:
- `pkg/module.go` — NativeModule, Loader, init(), moduleRuntime
- `pkg/api_index.go` — indexRef, IndexBuilder, open/create/memory
- `pkg/api_query.go` — query factories (match, term, bool, conj, disj, prefix, fuzzy, matchAll, matchNone)
- `pkg/api_search.go` — SearchRequestBuilder, SearchResult, Hit
- `pkg/api_batch.go` — Batch wrapper
- `pkg/codec.go` — JS↔Go conversion helpers
- `pkg/doc.go`, `pkg/logcopter.go`

**Test**: Create an in-memory index, index some documents with text fields, run a match query, verify hits.

**Key bleve source files to reference**:
- `bleve/index.go` — Index interface, New(), Open(), NewMemOnly()
- `bleve/query.go` — NewMatchQuery, NewBooleanQuery, etc.
- `bleve/search.go` — SearchRequest, SearchResult
- `bleve/mapping.go` — NewIndexMapping, NewTextFieldMapping, etc.

### Phase 2: Mapping Configuration + Field Types

**Goal**: Full mapping builder API, all field types.

**Files to create/modify**:
- `pkg/api_mapping.go` — MappingBuilder, DocMappingBuilder, FieldBuilder

**Test**: Create an index with explicit field mappings (text + keyword + number + datetime), verify that unmapped fields are/aren't indexed depending on dynamic settings.

**Key bleve source files to reference**:
- `bleve/mapping/index.go` — IndexMappingImpl structure
- `bleve/mapping/document.go` — DocumentMapping structure
- `bleve/mapping/field.go` — FieldMapping structure, all field type constructors

### Phase 3: Vector Fields + KNN Search

**Goal**: Vector field mapping, KNN queries, pure vector search.

**Files to create/modify**:
- `pkg/api_knn.go` — KNNRequest wrapper
- `pkg/api_mapping.go` — add vector/vectorBase64 support to FieldBuilder
- `pkg/api_search.go` — add .knn(), .knnOperator() to SearchRequestBuilder
- `pkg/codec.go` — JS array → []float32 conversion

**Build tag**: The Go binary must be compiled with `-tags vectors` to enable bleve's vector support. The module should detect at init time whether vector support is available and set a flag. If vector support is not compiled in, `bleve.field().vector()` should throw a clear error at `.build()` time.

**Test**: Create an in-memory index with a vector field, index documents with embeddings, run KNN search, verify nearest neighbors.

**Key bleve source files to reference**:
- `bleve/mapping_vector.go` — NewVectorFieldMapping, NewVectorBase64FieldMapping
- `bleve/mapping/mapping_vectors.go` — vector field processing, validation, normalization
- `bleve/search_knn.go` — KNNRequest, SearchRequest KNN handling, two-phase search
- `bleve/search/query/knn.go` — KNNQuery struct and Searcher method
- `bleve/document/field_vector.go` — VectorField struct

### Phase 4: Hybrid Search + Score Fusion

**Goal**: Combine BM25 and KNN with RRF/RSF fusion.

**Files to modify**:
- `pkg/api_search.go` — add .score(), .scoreRankConstant(), .scoreWindowSize()

**Test**: Create an index with both text and vector fields, index documents, run hybrid search with RRF fusion, verify that results combine both signals.

**Key bleve source files to reference**:
- `bleve/rescorer.go` — rescorer struct, prepareSearchRequest, rescore, mergeDocs
- `bleve/search.go` — IsScoreFusionRequested, ScoreRRF, ScoreRSF, RequestParams
- `bleve/search_knn.go` — two-phase search orchestration, KNN operator handling

### Phase 5: Provider Registration + Integration

**Goal**: Register goja-bleve as a providerapi package, wire it into the goja engine.

**Files to create**:
- `pkg/provider.go` — providerapi.Module registration, config schema, HostServices interface

**Integration**: Add to go.work, import in the host application's module registration.

**Key source files to reference**:
- `geppetto/pkg/js/modules/geppetto/provider/provider.go` — provider pattern
- `go-go-goja/pkg/xgoja/providerapi/module.go` — Module struct, ModuleFactory
- `goja-text/pkg/xgoja/providers/text/text.go` — simpler provider example

### Phase 6: TypeScript Declarations

**Goal**: Implement `TypeScriptDeclarer` for auto-generated .d.ts files.

**Files to modify**:
- `pkg/module.go` — add TypeScriptModule() method

**Key source files to reference**:
- `go-go-goja/modules/typing.go` — TypeScriptDeclarer interface
- `go-go-goja/modules/database/database.go` — TypeScriptModule() example
- `goja-text/pkg/extract/typescript.go` — another example

---

## Decision Records

### DR-1: Go-backed refs vs. plain JS objects

**Context**: We need to represent bleve objects (indexes, queries, mappings, search results) in JavaScript. Two options:

- **Option A**: Plain JS objects with JSON-serializable state. Go functions accept these objects and parse them on each call.
- **Option B**: Go-backed reference objects attached via hidden properties. Go functions read the hidden ref directly.

**Decision**: Option B (Go-backed refs).

**Rationale**:
- Type safety: JS code cannot construct an invalid query tree; all queries are built through validated Go constructors
- Performance: No serialization/deserialization of large vectors or result sets through JSON
- Consistency: This is the pattern established by geppetto (engineRef, inferenceSettingsRef, schemaBuilderRef, etc.)
- The existing rag-eval search service already uses Go structs throughout; the goja bridge simply makes them accessible from JS

**Consequences**:
- JS code cannot inspect/serialize these objects without explicit `.toJSON()` methods
- Debugging requires using `.toJSON()` or `.debug()` methods (following geppetto's pattern)
- The module must implement `attachRef`/`getRef` infrastructure

### DR-2: Builder pattern vs. constructor options

**Context**: How should users configure complex objects like SearchRequest (which has 10+ optional parameters)?

- **Option A**: Constructor options object (like `{size: 10, fields: ["text"], explain: true}`)
- **Option B**: Fluent builder pattern (like `.size(10).fields("text").explain(true).build()`)

**Decision**: Option B (fluent builder).

**Rationale**:
- Consistency with geppetto and goja-text APIs
- Type safety: each builder method validates its arguments immediately
- Discoverability: IDE autocomplete shows available methods
- Immutability: builders can return new wrappers, preventing accidental mutation

**Consequences**:
- More Go code to write (each builder needs a ref struct and wrapper construction)
- Slightly more verbose JS at call sites (`.build()` terminal call required)
- But the resulting API is more maintainable and harder to misuse

### DR-3: Module name: "bleve" vs. "search"

**Context**: What should the `require()` name be?

- **Option A**: `require("bleve")` — names the module after the library it wraps
- **Option B**: `require("search")` — names the module after the abstract capability

**Decision**: Option A (`require("bleve")`).

**Rationale**:
- Bleve is the specific implementation; "search" is too generic and could conflict with other search modules
- The geppetto module is named after its library, not "llm" or "inference"
- If a future module wraps a different search engine (e.g., Meilisearch), it would get its own name

**Consequences**:
- Scripts must use `require("bleve")` explicitly
- The module is clearly tied to bleve's capabilities and limitations

### DR-4: Vector support detection

**Context**: Bleve's vector features require the `vectors` build tag. How should the module handle missing vector support?

**Decision**: Runtime detection with clear errors.

- At `init()` time, probe whether `bleve.NewVectorFieldMapping` exists (it won't be compiled without the tag)
- Set a `vectorSupportAvailable` boolean on the moduleRuntime
- If a JS script calls `.vector()` on a FieldBuilder when support is unavailable, throw a descriptive error: `"vector fields require the 'vectors' build tag; recompile with -tags vectors"`
- Expose `bleve.vectorSupport` as a read-only boolean so scripts can check at runtime

**Rationale**: Fail clearly at the point of use rather than silently producing broken indexes.

### DR-5: Index ownership and lifecycle

**Context**: Who is responsible for closing indexes? What happens if a script opens an index and never closes it?

**Decision**: The module tracks open indexes and closes them on runtime shutdown.

- `moduleRuntime` maintains a `map[string]*indexRef` of all indexes opened via the module
- On runtime shutdown (via Close callback), all open indexes are closed
- Scripts should still call `idx.close()` explicitly for resource hygiene
- Double-close is a no-op (not an error)

**Rationale**: Prevents resource leaks in long-running goja runtimes. Matches the pattern used by the database module for connection cleanup.

---

## Key File Reference Map

For the implementer, here are the critical source files to study, organized by subsystem:

### Bleve Core API
| File | What to study |
|---|---|
| `bleve/index.go` | Index interface, Batch, New(), Open(), NewMemOnly() |
| `bleve/query.go` | All NewXxxQuery factory functions |
| `bleve/search.go` | SearchRequest, SearchResult, SearchStatus, NewSearchRequestOptions() |
| `bleve/mapping.go` | NewIndexMapping, NewDocumentMapping, NewTextFieldMapping, etc. |
| `bleve/mapping_vector.go` | NewVectorFieldMapping, NewVectorBase64FieldMapping |
| `bleve/rescorer.go` | Hybrid search rescoring logic (RRF, RSF) |
| `bleve/builder.go` | Offline index builder (may be useful for bulk loading) |

### Bleve Internal Details
| File | What to study |
|---|---|
| `bleve/search_knn.go` | KNNRequest, two-phase search, KNN operator |
| `bleve/search/query/knn.go` | KNNQuery struct, Searcher() method |
| `bleve/mapping/field.go` | FieldMapping struct, all fields including vector-specific ones |
| `bleve/mapping/index.go` | IndexMappingImpl, validation, custom analysis |
| `bleve/mapping/document.go` | DocumentMapping, field/subdoc registration |
| `bleve/mapping/mapping_vectors.go` | Vector field processing, normalization, validation |
| `bleve/document/field_vector.go` | VectorField struct, constructor |

### Goja Module Patterns
| File | What to study |
|---|---|
| `go-go-goja/modules/common.go` | NativeModule interface, Registry, Register |
| `go-go-goja/modules/exports.go` | SetExport helper |
| `go-go-goja/modules/typing.go` | TypeScriptDeclarer interface |
| `go-go-goja/modules/database/database.go` | Simple module example (configure/query/exec/close) |
| `geppetto/pkg/js/modules/geppetto/module.go` | Complex module with moduleRuntime, attachRef/getRef |
| `geppetto/pkg/js/modules/geppetto/api_engine_builder.go` | Builder pattern with cloneFor |
| `geppetto/pkg/js/modules/geppetto/api_schema_builders.go` | Recursive builder (schema.object().property().build()) |
| `geppetto/pkg/js/modules/geppetto/api_inference_settings.go` | Go-backed ref with toJSON/clone/debug |
| `goja-text/pkg/extract/module.go` | Simpler builder pattern (options/Build) |
| `goja-text/pkg/extract/types.go` | ExtractOptionsBuilder with validation |

### Provider Integration
| File | What to study |
|---|---|
| `go-go-goja/pkg/xgoja/providerapi/module.go` | Module struct, ModuleFactory, ModuleContext |
| `go-go-goja/pkg/xgoja/providerapi/registry.go` | Package registration |
| `geppetto/pkg/js/modules/geppetto/provider/provider.go` | Full provider example with config schema |
| `goja-text/pkg/xgoja/providers/text/text.go` | Simpler provider (no config) |

### Existing RAG Usage
| File | What to study |
|---|---|
| `2026-05-27--rag-evaluation-system/internal/services/search/bm25.go` | How bleve is used for BM25 indexing and search today |
| `2026-05-27--rag-evaluation-system/internal/services/search/hybrid.go` | RRF fusion implementation (manual, not using bleve's built-in rescoring) |
| `2026-05-27--rag-evaluation-system/internal/services/search/vector.go` | Manual vector search (not using bleve's KNN at all) |

---

## Risks and Open Questions

### Risk: Vector build tag fragmentation

Bleve's vector support is behind a build tag. If the host binary is compiled without `-tags vectors`, the module degrades gracefully (no vector fields, no KNN). However, the rag-evaluation system *needs* vector search, so the binary must always be compiled with the tag in practice. The runtime detection is a safety net, not the expected path.

**Mitigation**: Document the build tag requirement prominently. Consider making it the default in the Makefile.

### Risk: Large vector arrays in JS↔Go bridge

Passing embedding vectors (e.g., 1536 float32 values = 6KB) across the JS↔Go bridge on every KNN query could be a performance concern. Goja's `Float32Array` support may or may not allow zero-copy transfer.

**Mitigation**: Profile early. If needed, explore shared ArrayBuffer or pre-allocated buffer patterns. The existing geppetto module already passes vectors for embedding results, so there is precedent.

### Risk: Index thread safety

Bleve indexes are thread-safe for reads but not for concurrent writes. The goja runtime is single-threaded by default (one VM = one goroutine), so this is naturally safe. However, if the host application creates multiple runtimes that share an index, writes could conflict.

**Mitigation**: Document that an Index wrapper should only be used from the runtime that created it. The module tracks indexes per moduleRuntime.

### Open Question: Should we expose bleve's MultiSearch?

Bleve supports searching across multiple indexes simultaneously via `MultiSearch()`. This is useful for sharded corpora but adds API complexity. **Deferred to v2**.

### Open Question: Should we expose bleve's Alias?

An IndexAlias lets you treat multiple indexes as one virtual index. Useful for rolling index updates. **Deferred to v2**.

### Open Question: Should the module support geppetto embedding integration directly?

Currently, the design expects JS scripts to manually call geppetto for embeddings and pass vectors to bleve. An alternative is a convenience method like `idx.searchWithEmbedding(queryText, embedder)` that handles the embed-then-search round trip internally.

**Decision for v1**: No convenience methods. Keep the modules independent and composable. JS scripts have full control over the pipeline. Convenience wrappers can be added later as a separate helper module.

---

## Testing Strategy

### Unit Tests

Each phase should have corresponding Go test files that create a goja runtime, load the bleve module, run JS code, and verify results. Follow the pattern from `goja-text/pkg/extract/module_test.go`:

- Create a runtime factory with the bleve module enabled
- Run JS code via `vm.RunString()`
- Export the result and assert on Go values

### Integration Tests

- **BM25 round-trip**: Create in-memory index → index 50 documents → search → verify ranking
- **Vector round-trip**: Create in-memory index with vector field → index documents with random embeddings → KNN search → verify top-k ordering
- **Hybrid round-trip**: Create in-memory index with text + vector → index documents → hybrid search with RRF → verify fusion behavior
- **Batch indexing**: Verify batch flush semantics match the existing `bm25.go` service behavior
- **Mapping validation**: Verify that invalid mappings (zero-dim vectors, unknown similarity) produce clear errors

### Smoke Test

A standalone script (in `scripts/` or `examples/`) that exercises the full RAG pipeline: create index → chunk text → embed → index → search → print results. This is the "does it work end-to-end" test.

---

## Glossary

| Term | Definition |
|---|---|
| **bleve** | Go library for full-text and vector search (`github.com/blevesearch/bleve/v2`) |
| **goja** | Go JavaScript runtime (`github.com/dop251/goja`) |
| **goja-bleve** | The native module this document specifies |
| **KNN** | K-nearest neighbors search over vector fields |
| **RRF** | Reciprocal Rank Fusion — hybrid search method that combines rank positions |
| **RSF** | Relative Score Fusion — hybrid search method that combines normalized scores |
| **scorch** | Bleve's default index type (based on upsidedown KV store) |
| **go-faiss** | Go bindings for Facebook's FAISS vector similarity library |
| **NativeModule** | Go interface for goja modules: `Name()`, `Doc()`, `Loader()` |
| **moduleRuntime** | Per-VM state object (like geppetto's pattern) |
| **ref** | Go struct attached to a JS object via hidden property |
| **builder** | Go struct that accumulates config and produces a ref on `.build()` |
| **providerapi** | Go-go-goja's provider registration system for xgoja hosts |
