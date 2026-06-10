---
Title: Goja Bleve JavaScript API reference
Slug: goja-bleve-js-api-reference
Short: Reference for the `require("bleve")` module, including mapping builders, query builders, search requests, indexes, batches, vectors, and result objects.
Topics:
- goja-bleve
- bleve
- javascript
- api
- reference
- vector-search
Commands:
- goja-bleve
- eval
- run
- repl
Flags:
- output
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This reference describes the JavaScript surface exposed by `require("bleve")`. The API wraps Bleve's Go types in opaque JavaScript objects and uses builders for mappings, queries, search requests, and indexes.

The reference is written for script authors. It focuses on what each function returns, what the object is used for, and the lifecycle rules that matter when a script runs inside a generated xgoja binary.

## Module import

```javascript
const bleve = require("bleve");
```

All functions below are properties of that module object.

## Top-level factories

| Function | Returns | Purpose |
|---|---|---|
| `bleve.mapping()` / `bleve.indexMapping()` | `MappingBuilder` | Build an index mapping. |
| `bleve.docMapping()` / `bleve.documentMapping()` | `DocumentMappingBuilder` | Build a document mapping. |
| `bleve.field()` | `FieldBuilder` | Build a field mapping. |
| `bleve.memory()` | `IndexBuilder` | Build an in-memory index. |
| `bleve.create(path)` | `IndexBuilder` | Build a new persistent index at `path`. |
| `bleve.open(path)` | `IndexBuilder` | Open an existing persistent index at `path`. |
| `bleve.search()` / `bleve.searchRequest()` | `SearchRequestBuilder` | Build a search request. |
| `bleve.bool()` | `BoolQuery` | Build a mutable Boolean query. |
| `bleve.match(text)` | `QueryBuilder` | Build an analyzed match query. |
| `bleve.matchPhrase(text)` | `QueryBuilder` | Build an analyzed phrase query. |
| `bleve.term(value)` | `QueryBuilder` | Build an exact term query. |
| `bleve.prefix(prefix)` | `QueryBuilder` | Build a prefix query. |
| `bleve.wildcard(pattern)` | `QueryBuilder` | Build a wildcard query. |
| `bleve.regexp(pattern)` | `QueryBuilder` | Build a regular-expression query. |
| `bleve.fuzzy(term)` | `QueryBuilder` | Build a fuzzy term query. |
| `bleve.queryString(query)` | `Query` | Build a query from Bleve query-string syntax. |
| `bleve.matchAll()` | `Query` | Match every document. |
| `bleve.matchNone()` | `Query` | Match no documents; useful for pure KNN. |
| `bleve.conj(...queries)` / `bleve.conjunction(...queries)` | `Query` | Combine queries with AND. |
| `bleve.disj(...queries)` / `bleve.disjunction(...queries)` | `Query` | Combine queries with OR. |

## MappingBuilder

Create with `bleve.mapping()` or `bleve.indexMapping()`.

```javascript
const mapping = bleve.mapping()
  .defaultMapping(docMapping)
  .defaultField("body")
  .defaultAnalyzer("standard")
  .build();
```

Methods:

| Method | Returns | Description |
|---|---|---|
| `.defaultMapping(mapping)` | `this` | Sets the mapping used for documents without a type-specific mapping. |
| `.addTypeMapping(name, mapping)` | `this` | Adds a named type mapping. |
| `.typeMapping(name, mapping)` | `this` | Alias for `.addTypeMapping`. |
| `.typeField(name)` | `this` | Sets the document field used to choose type mappings. |
| `.defaultAnalyzer(name)` | `this` | Sets the analyzer name used when a field does not override it. |
| `.defaultField(name)` | `this` | Sets the field used by queries that do not specify a field. |
| `.storeDynamic(enabled)` | `this` | Controls dynamic field storage. |
| `.indexDynamic(enabled)` | `this` | Controls dynamic field indexing. |
| `.docValuesDynamic(enabled)` | `this` | Controls dynamic doc values. |
| `.build()` | `Mapping` | Produces the opaque mapping wrapper. |

Use `build()` once the mapping is complete. Pass the resulting `Mapping` to an `IndexBuilder`.

## DocumentMappingBuilder

Create with `bleve.docMapping()` or `bleve.documentMapping()`.

```javascript
const doc = bleve.docMapping()
  .field("title", bleve.field().text())
  .field("status", bleve.field().keyword())
  .dynamic(false)
  .build();
```

Methods:

| Method | Returns | Description |
|---|---|---|
| `.field(name, field)` | `this` | Adds a field mapping. Accepts a `FieldBuilder` or built `FieldMapping`. |
| `.addField(name, field)` | `this` | Alias for `.field`. |
| `.subDocument(name, doc)` | `this` | Adds a nested document mapping. |
| `.addSubDoc(name, doc)` | `this` | Alias for `.subDocument`. |
| `.dynamic(enabled)` | `this` | Enables or disables dynamic mapping for this document. |
| `.enabled(enabled)` | `this` | Enables or disables this document mapping. |
| `.nested(enabled)` | `this` | Marks the document mapping as nested where supported by Bleve. |
| `.defaultAnalyzer(name)` | `this` | Sets the analyzer for fields under this document mapping. |
| `.build()` | `DocumentMapping` | Produces the opaque document mapping wrapper. |

## FieldBuilder

Create with `bleve.field()`.

Field type methods:

| Method | Purpose |
|---|---|
| `.text()` | Analyzed full-text field. |
| `.keyword()` | Exact keyword field. |
| `.number()` | Numeric field. |
| `.datetime()` | Date/time field. |
| `.boolean()` | Boolean field. |
| `.geoPoint()` | Geographic point field. |
| `.geoShape()` | Geographic shape field. |
| `.ip()` | IP address field. |
| `.disabled()` | Field is present in documents but not indexed. |
| `.vector(dims)` | Numeric vector field with `dims` dimensions. |
| `.vectorBase64(dims)` | Base64-encoded vector field with `dims` dimensions. |

Common field options:

| Method | Purpose |
|---|---|
| `.name(name)` | Sets the explicit field name. |
| `.analyzer(name)` | Sets the analyzer for text fields. |
| `.store(enabled)` | Stores field values for result retrieval. |
| `.index(enabled)` | Enables/disables indexing. |
| `.docValues(enabled)` | Enables doc values. |
| `.includeTermVectors(enabled)` | Stores term vectors for highlighting/explanation use cases. |
| `.includeInAll(enabled)` | Includes the field in Bleve's all-field behavior where applicable. |
| `.dateFormat(name)` | Sets date parser/format name for datetime fields. |
| `.similarity(name)` | Sets vector similarity; common values include `cosine`, `dot_product`, and `l2_norm`. |
| `.optimizedFor(name)` | Sets vector optimization preference; common values include `recall`, `latency`, and `memory-efficient`. |
| `.build()` | Produces a `FieldMapping`. |

Vector methods are available in the API, but vector execution requires a binary built with vector support.

## Query and QueryBuilder

Most query factories return a `QueryBuilder`. A query builder supports:

| Method | Returns | Description |
|---|---|---|
| `.field(name)` | `this` | Restricts the query to a field. |
| `.boost(value)` | `this` | Sets query boost. |

Query objects are consumed by `SearchRequestBuilder.query(...)`, Boolean queries, conjunctions, and disjunctions.

Examples:

```javascript
const text = bleve.match("semantic search").field("body").boost(2);
const exact = bleve.term("published").field("status");
const phrase = bleve.matchPhrase("hybrid retrieval").field("title");
const parsed = bleve.queryString("title:bleve +status:published");
```

## BoolQuery

Create with `bleve.bool()`.

| Method | Returns | Description |
|---|---|---|
| `.addMust(...queries)` | `this` | Adds required clauses. |
| `.addShould(...queries)` | `this` | Adds optional/scoring clauses. |
| `.addMustNot(...queries)` | `this` | Adds prohibited clauses. |
| `.field(name)` | `this` | Inherited query field helper where applicable. |
| `.boost(value)` | `this` | Inherited query boost helper where applicable. |

Example:

```javascript
const query = bleve.bool()
  .addMust(bleve.match("search").field("body"))
  .addShould(bleve.match("goja").field("title"))
  .addMustNot(bleve.term("archived").field("status"));
```

## SearchRequestBuilder

Create with `bleve.search()` or `bleve.searchRequest()`.

| Method | Returns | Description |
|---|---|---|
| `.query(query)` | `this` | Sets the main Bleve query. |
| `.size(n)` | `this` | Sets the number of hits to return. |
| `.from(n)` | `this` | Sets result offset. |
| `.fields(names)` | `this` | Requests stored fields by name. |
| `.sort(names)` | `this` | Sets sort expressions, for example `["-_score"]` or `["-publishedAt"]`. |
| `.highlight(fields, style)` | `this` | Enables highlighting for fields. |
| `.explain(enabled)` | `this` | Enables score explanations. |
| `.score(mode)` | `this` | Sets score fusion mode: `default`, `none`, `rrf`, or `rsf`. |
| `.scoreRankConstant(n)` | `this` | Sets RRF rank constant. |
| `.scoreWindowSize(n)` | `this` | Sets fusion candidate window size. |
| `.knnOperator(operator)` | `this` | Sets KNN operator: `or` or `and`. |
| `.knn(field, vector, k, boost)` | `this` | Adds a KNN clause. Vector can be `number[]`, `Float32Array`, or `Float64Array`. |
| `.build()` | `SearchRequest` | Produces the opaque search request wrapper. |

Examples:

```javascript
const request = bleve.search()
  .query(bleve.match("privacy").field("body"))
  .size(10)
  .fields(["title", "status"])
  .build();
```

Pure KNN:

```javascript
const request = bleve.search()
  .query(bleve.matchNone())
  .knn("embedding", [0.1, 0.2, 0.3, 0.4], 5, 1.0)
  .build();
```

Hybrid RRF:

```javascript
const request = bleve.search()
  .query(bleve.match("migration").field("text"))
  .knn("embedding", embedding, 20, 1.0)
  .score("rrf")
  .scoreRankConstant(60)
  .scoreWindowSize(50)
  .build();
```

## IndexBuilder

Create with `bleve.memory()`, `bleve.create(path)`, or `bleve.open(path)`.

| Method | Returns | Description |
|---|---|---|
| `.mapping(mapping)` | `this` | Sets the mapping for a new in-memory or persistent index. |
| `.name(name)` | `this` | Sets a diagnostic index name where supported. |
| `.build()` | `Index` | Opens or creates the index. |

For `open(path)`, the mapping is already part of the existing index and normally should not be supplied.

## Index

| Method | Returns | Description |
|---|---|---|
| `.index(id, doc)` | `void` | Adds or replaces a document. |
| `.delete(id)` | `void` | Deletes a document by ID. |
| `.search(request)` | `SearchResult` | Executes a built search request. |
| `.docCount()` | `number` | Returns the indexed document count. |
| `.newBatch()` / `.batch()` | `Batch` | Creates a new batch bound to this index. |
| `.close()` | `void` | Closes the index. |

## Batch

| Method | Returns | Description |
|---|---|---|
| `.index(id, doc)` | `this` | Queues an index operation. |
| `.delete(id)` | `this` | Queues a delete operation. |
| `.size()` | `number` | Returns the approximate batch size. |
| `.operationCount()` | `number` | Returns queued operation count. |
| `.reset()` | `this` | Clears queued operations before execution. |
| `.execute()` | `void` | Submits the batch and marks it executed. |

After `.execute()` succeeds, later mutation or reset calls throw an error. Create a fresh batch for additional work.

## SearchResult and SearchHit

`index.search(request)` returns:

```typescript
interface SearchResult {
  total: number;
  maxScore: number;
  took: string;
  hits: SearchHit[];
}

interface SearchHit {
  id: string;
  score: number;
  fields: Record<string, unknown>;
  fragments?: unknown;
  locations?: unknown;
  sort?: unknown[];
  explanation?: unknown;
  scoreBreakdown?: unknown;
}
```

The exact optional properties depend on the request and Bleve's response. Request stored fields with `.fields([...])` if you need document values in `hit.fields`.

## TypeScript declaration snapshot

The repository tests a TypeScript declaration snapshot in `pkg/testdata/bleve.d.ts.golden`. Treat that file as a compact machine-checkable API summary. Update it whenever public JavaScript APIs change.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| A method is missing from a builder | The object is already built or is the wrong builder type. | Keep builder variables separate from built wrapper variables. |
| `field()` or `mapping()` accepts an object but later fails | A plain JavaScript object was passed where an opaque goja-bleve wrapper was required. | Use objects returned by this module's builders only. |
| Vector arrays are rejected | The value is not a supported vector representation or has the wrong dimensions. | Use `number[]`, `Float32Array`, or `Float64Array`, and match the mapped dimension count. |
| Result fields are absent | Fields were not stored or were not requested. | Set `.store(true)` on the field mapping and call `.fields([...])` in the request. |
| Batch reuse fails | The batch has already been executed. | Create a new batch from the index. |

## See Also

- `goja-bleve-getting-started`
- `goja-bleve-user-guide`
