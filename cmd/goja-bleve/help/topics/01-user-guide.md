---
Title: Goja Bleve user guide
Slug: goja-bleve-user-guide
Short: Practical guide to indexing, querying, batching, persistent indexes, vector search, and xgoja integration with the Bleve JavaScript module.
Topics:
- goja-bleve
- bleve
- javascript
- search
- vector-search
- xgoja
Commands:
- goja-bleve
- eval
- run
- repl
- search
- batch
- vector
Flags:
- output
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

This guide explains how to use `goja-bleve` as a JavaScript search toolkit inside an xgoja runtime. It covers the whole workflow: choose an index shape, define mappings, index documents, run queries, use batches, manage persistent indexes, and decide when vector search is appropriate.

Use this guide when you are building an application script, a RAG retrieval prototype, an indexing utility, or a generated xgoja binary that should expose Bleve to JavaScript authors.

## Mental model

`goja-bleve` puts Bleve's Go search engine behind a JavaScript module:

```javascript
const bleve = require("bleve");
```

The JavaScript API is intentionally builder-oriented. Builders let scripts describe Go-backed Bleve objects without exposing raw Go pointers or mutable internals. Terminal `.build()` calls produce opaque wrapper objects that can be passed to other builders and index methods.

The usual data flow is:

```text
JavaScript documents
  -> mapping builders
  -> index builder
  -> index.index(...) or batch.index(...)
  -> query builders
  -> search request builder
  -> index.search(...)
  -> JavaScript result object
```

The same module supports short-lived in-memory indexes and persistent on-disk indexes. It also exposes vector mapping and KNN request builders for binaries compiled with vector support.

## Choose the right index lifecycle

Use an in-memory index when scripts are temporary, tests are deterministic, or documents are loaded fresh every run:

```javascript
const index = bleve.memory()
  .mapping(mapping)
  .name("scratch")
  .build();
```

Use a persistent index when you want to reuse indexed data across runs:

```javascript
const index = bleve.create("./data/articles.bleve")
  .mapping(mapping)
  .build();

// Later:
const sameIndex = bleve.open("./data/articles.bleve").build();
```

Persistent indexes must be closed:

```javascript
index.close();
```

Treat close as part of your script's normal lifecycle, not as optional cleanup. The runtime can clean up some open resources, but explicit close is easier to reason about and reduces file-lock surprises.

## Design mappings from query behavior

Mappings determine how documents are indexed. Start from the questions users will ask.

Use `text()` for human language:

```javascript
const title = bleve.field().text().store(true).build();
const body = bleve.field().text().includeTermVectors(true).build();
```

Use `keyword()` for exact matching and faceted-style values:

```javascript
const status = bleve.field().keyword().docValues(true).build();
```

Use numeric, datetime, boolean, IP, geo, and disabled fields when the document shape requires them:

```javascript
const publishedAt = bleve.field().datetime().dateFormat("dateTimeOptional").build();
const views = bleve.field().number().build();
const archived = bleve.field().boolean().build();
const rawPayload = bleve.field().disabled().build();
```

Then assemble a document mapping:

```javascript
const article = bleve.docMapping()
  .field("title", title)
  .field("body", body)
  .field("status", status)
  .field("publishedAt", publishedAt)
  .build();

const mapping = bleve.mapping()
  .defaultMapping(article)
  .defaultField("body")
  .defaultAnalyzer("standard")
  .build();
```

The current API intentionally does not expose custom analyzer registration. Use built-in analyzer names and let Bleve validate the final mapping. If you need application-specific analysis, design it as a separate provider/host concern rather than hiding it in ad hoc JavaScript.

## Index documents safely

For individual documents:

```javascript
index.index("article-1", {
  title: "Search from JavaScript",
  body: "Bleve runs inside a goja-hosted runtime.",
  status: "published",
});
```

For updates, call `index.index(id, doc)` again with the same ID. For deletes:

```javascript
index.delete("article-1");
```

For bulk indexing, use a batch:

```javascript
const batch = index.batch();

for (const row of rows) {
  if (row.deleted) {
    batch.delete(row.id);
  } else {
    batch.index(row.id, row);
  }
}

console.log(`queued ${batch.operationCount()} operations`);
batch.execute();
```

After `execute()`, do not reuse that batch. Create a new one. This single-use rule prevents subtle bugs where a script accidentally submits old operations twice.

## Build queries intentionally

Use `match` for analyzed full-text terms:

```javascript
bleve.match("privacy preserving search").field("body").boost(2.0)
```

Use `matchPhrase` when order matters:

```javascript
bleve.matchPhrase("vector search").field("title")
```

Use `term`, `prefix`, `wildcard`, `regexp`, and `queryString` for exact or parser-driven searches:

```javascript
bleve.term("published").field("status")
bleve.prefix("goja").field("title")
bleve.queryString("title:goja +status:published")
```

Use Boolean composition when you need multiple conditions:

```javascript
const query = bleve.bool()
  .addMust(bleve.match("search").field("body"))
  .addShould(bleve.match("goja").field("title"))
  .addMustNot(bleve.term("archived").field("status"));
```

There are also convenience constructors:

```javascript
bleve.conj(q1, q2)      // conjunction
bleve.disj(q1, q2, q3)  // disjunction
bleve.matchAll()
bleve.matchNone()
```

## Shape search requests

A query describes what to match. A search request describes how to execute and return it.

```javascript
const request = bleve.search()
  .query(query)
  .from(0)
  .size(20)
  .fields(["title", "status", "publishedAt"])
  .sort(["-publishedAt", "_score"])
  .highlight(["body"], "html")
  .explain(false)
  .build();

const result = index.search(request);
```

Results are plain JavaScript-friendly objects:

```javascript
{
  total: 12,
  maxScore: 1.74,
  took: "3.2ms",
  hits: [
    { id: "article-1", score: 1.74, fields: { title: "..." } }
  ]
}
```

Request only fields you need. Stored fields cost index space, and returning large fields can dominate script runtime.

## Use vectors and hybrid search when the host supports them

Vector support is build-time gated. The API includes vector builders, but KNN execution requires a binary built with the `vectors` tag and FAISS linked in.

Define a vector field:

```javascript
const embedding = bleve.field()
  .vector(384)
  .similarity("cosine")
  .optimizedFor("recall")
  .build();
```

Add it to a mapping:

```javascript
const chunk = bleve.docMapping()
  .field("text", bleve.field().text())
  .field("embedding", embedding)
  .build();
```

Run KNN:

```javascript
const request = bleve.search()
  .query(bleve.matchNone())
  .knn("embedding", queryVector, 10, 1.0)
  .build();
```

Run hybrid text + vector search:

```javascript
const request = bleve.search()
  .query(bleve.match("database migrations").field("text"))
  .knn("embedding", queryVector, 20, 1.0)
  .score("rrf")
  .scoreRankConstant(60)
  .scoreWindowSize(50)
  .build();
```

Use RRF when text and vector scores are not directly comparable. Use RSF when score scales are meaningful enough for relative score fusion. Keep `scoreWindowSize` at least as large as the result size and large enough to include candidates from each retrieval leg.

## Integrate with xgoja

A generated xgoja spec mounts the provider package and the JavaScript module alias:

```yaml
packages:
  - id: goja-bleve
    import: github.com/go-go-golems/goja-bleve/pkg/xgoja/providers/bleve
modules:
  - package: goja-bleve
    name: bleve
    as: bleve
```

Then scripts use:

```javascript
const bleve = require("bleve");
```

If your runtime also needs file access, database access, LLM calls, or other modules, mount them alongside `bleve`. For RAG scripts, a common shape is:

```javascript
const fs = require("fs");
const bleve = require("bleve");
const geppetto = require("geppetto");
```

Use host modules to read documents and produce embeddings. Use `goja-bleve` to index and retrieve.

## Operational guidelines

- Keep index paths explicit and outside source directories.
- Use in-memory indexes for examples and tests.
- Store only fields that must appear in search results.
- Close persistent indexes.
- Keep vector dimensions fixed per field and validate embeddings before indexing.
- Prefer batches for bulk indexing.
- Treat mapping changes as index-version changes; rebuild indexes when mappings change materially.
- Keep JavaScript wrappers opaque. Do not inspect or modify internal properties such as `__bleve_ref`.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| A field never matches | It may be disabled, not indexed, mapped as a keyword when you expect analyzed text, or queried with the wrong field name. | Inspect the mapping construction and start with a small `matchAll()` request returning stored fields. |
| Returned `fields` are empty | The field was not stored or not requested. | Set `.store(true)` on the field mapping and call `.fields([...])` on the search request. |
| Query builder panics or throws wrapper errors | A builder expected a `Query`, `Mapping`, `DocumentMapping`, or `FieldMapping` wrapper created by this module. | Always pass objects returned by `bleve.*` builders; do not construct wrapper-looking objects by hand. |
| KNN reports dimension mismatch | The query vector length does not match the mapped vector field dimensions. | Check `.vector(dims)` in the mapping and validate embedding lengths before search. |
| Hybrid results look text-only | KNN may not be contributing candidates, vector support may be absent, or the KNN field may be wrong. | Test pure KNN with `.query(bleve.matchNone())`, then add text and fusion once KNN works. |
| Reopening a persistent index fails after mapping changes | Existing index data was built with a different mapping. | Create a new index directory or rebuild the index after mapping changes. |

## See Also

- `goja-bleve-getting-started`
- `goja-bleve-js-api-reference`
