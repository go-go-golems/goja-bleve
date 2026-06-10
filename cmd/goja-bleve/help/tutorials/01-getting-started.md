---
Title: Goja Bleve getting started
Slug: goja-bleve-getting-started
Short: Build your first in-memory Bleve index from JavaScript and run BM25 searches from a generated xgoja binary.
Topics:
- goja-bleve
- bleve
- javascript
- xgoja
- getting-started
Commands:
- goja-bleve
- eval
- run
- search
Flags:
- output
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial gets you from an empty checkout to a working JavaScript search script. It uses the generated `goja-bleve` xgoja binary, loads the native `bleve` module with `require("bleve")`, builds an in-memory index, inserts a few documents, and runs a text search.

The important idea is that `goja-bleve` is not a separate HTTP service. It is a native Go module exposed inside a goja JavaScript runtime. Your script creates Bleve mappings, indexes documents, builds search requests, and receives plain JavaScript objects back.

## What you need

Use this path when you want to try the module quickly before embedding it into a larger runtime.

You need:

- Go installed with the version required by this repository.
- A checkout of `github.com/go-go-golems/goja-bleve`.
- The generated command under `cmd/goja-bleve`.

From the repository root, verify the command can run:

```bash
cd cmd/goja-bleve
GOWORK=off go run . modules
```

The output should include a provider package named `goja-bleve` and a module alias named `bleve`. That means JavaScript can call:

```javascript
const bleve = require("bleve");
```

## Create a first script

Create `first-search.js`:

```javascript
const bleve = require("bleve");

const mapping = bleve.mapping()
  .defaultMapping(
    bleve.docMapping()
      .field("title", bleve.field().text())
      .field("body", bleve.field().text())
      .field("category", bleve.field().keyword())
      .build()
  )
  .defaultField("body")
  .build();

const index = bleve.memory()
  .mapping(mapping)
  .name("articles")
  .build();

index.index("1", {
  title: "Introducing goja-bleve",
  body: "Goja scripts can create Bleve indexes and run full-text search.",
  category: "guide",
});

index.index("2", {
  title: "Hybrid retrieval",
  body: "Bleve can combine text search with vector KNN when built with vector support.",
  category: "reference",
});

index.index("3", {
  title: "Operational notes",
  body: "Close persistent indexes when a script is finished with them.",
  category: "ops",
});

const request = bleve.search()
  .query(bleve.match("goja search").field("body"))
  .size(5)
  .fields(["title", "category"])
  .build();

const result = index.search(request);

console.log(JSON.stringify({
  total: result.total,
  hits: result.hits.map(hit => ({
    id: hit.id,
    score: hit.score,
    title: hit.fields.title,
    category: hit.fields.category,
  })),
}, null, 2));

index.close();
```

Run it:

```bash
cd cmd/goja-bleve
GOWORK=off go run . run ../../first-search.js
```

If you put the script somewhere else, adjust the path passed to `run`.

## Understand the script

The script has four phases.

First, it creates a mapping. The mapping tells Bleve how to interpret fields. `text()` fields are tokenized and analyzed for full-text search. `keyword()` fields are indexed as exact values, which is usually what you want for categories, IDs, tags, and status fields.

Second, it creates an in-memory index. `bleve.memory()` is the best choice for tests, short-lived scripts, and examples because it leaves no files behind. Use `bleve.create(path)` and `bleve.open(path)` only when you intentionally want a persistent on-disk index.

Third, it indexes documents. JavaScript objects are converted into Go values and passed to Bleve. Document IDs are strings. Reusing an ID replaces that document.

Fourth, it builds a search request and executes it. Query builders such as `bleve.match(...)`, `bleve.term(...)`, and `bleve.queryString(...)` create query objects. `bleve.search()` builds the request envelope with pagination, fields, sorting, highlighting, explanation, KNN, and score-fusion options.

## Use a batch when indexing many documents

For more than a handful of documents, use `index.newBatch()`:

```javascript
const batch = index.newBatch();

for (const doc of documents) {
  batch.index(doc.id, doc);
}

batch.execute();
```

A batch is single-use after `execute()`. If you need another batch, call `index.newBatch()` again. This lifecycle makes indexing behavior explicit and avoids ambiguity about whether queued operations were retained after submission.

## Use a persistent index

Use `bleve.create(path)` for a new on-disk index and `bleve.open(path)` for an existing one:

```javascript
const index = bleve.create("./tmp/articles.bleve")
  .mapping(mapping)
  .build();

// ... index and search ...

index.close();
```

Persistent indexes are ordinary filesystem directories managed by Bleve. The host application decides path policy. The standalone generated command does not sandbox these paths, so use explicit directories and avoid writing into source-controlled paths by accident.

## Where vectors fit

The default non-vector build supports text search, mapping builders, batches, and request builders. Vector field mappings and KNN request builders are present in the JavaScript API, but actual vector search requires a host binary built with the `vectors` Go build tag and FAISS available at link/runtime.

Use text search first. Add vector search only when you have embeddings, dimensions, similarity choices, and deployment support for FAISS.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `Cannot find module "bleve"` | The generated xgoja binary did not mount the `goja-bleve` provider as `bleve`. | Run `goja-bleve modules` and check the spec; the module should be listed with alias `bleve`. |
| Search returns no hits | The query targets the wrong field, the field is not indexed, or the analyzer does not produce the terms you expect. | Start with `.defaultField("body")`, request `.fields([...])`, and search a known indexed word. |
| Persistent index cannot be reopened | The process did not close the index or the path points at a non-Bleve directory. | Call `index.close()` and keep index paths separate from other application data. |
| KNN search reports vector support errors | The binary was not built with the `vectors` tag or FAISS is unavailable. | Use text search in the default build, or build the vector target described in the vector playbook. |
| Batch mutation fails after execute | Batches are single-use after successful execution. | Create a new batch with `index.newBatch()` for the next group of operations. |

## See Also

- `goja-bleve-user-guide`
- `goja-bleve-js-api-reference`
