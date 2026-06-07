---
Title: Source: blevesearch-homepage.md
Ticket: RAGEVAL-GOJA-RAG-STRATEGIES
Status: active
Topics: [bleve]
DocType: reference
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Source material fetched for goja-bleve design investigation"
LastUpdated: 2026-06-02T21:30:00.000000-04:00
WhatFor: "Reference source material"
WhenToUse: "When investigating bleve documentation and APIs"
---

```
import "github.com/blevesearch/bleve/v2"

func main() {
    // open a new index
    mapping := bleve.NewIndexMapping()
    index, err := bleve.New("example.bleve", mapping)

    // index some data
    err = index.Index(identifier, your_data)

    // search for some text
    query := bleve.NewMatchQuery("text")
    search := bleve.NewSearchRequest(query)
    searchResults, err := index.Search(search)
}
```

## Bleve modern indexing & search for Go Simple top-level API Index any object in your data model 
Override default mapping to customize behavior Rich set of interfaces for extending the capabilities

### Upcoming Events (See All)

---

#### Simple

Import one package, build an index with three lines of code, query for documents with another three 
lines.

#### Text Analysis

Bleve includes general-purpose analyzers as well as pre-built text analyzers for the following 
languages:

- Arabic (ar), Bulgarian (bg), Catalan (ca), Chinese-Japanese-Korean (cjk), Kurdish (ckb), Danish 
(da), German (de), Greek (el), English (en), Spanish - Castilian (es), Basque (eu), Persian (fa), 
Finnish (fi), French (fr), Gaelic (ga), Spanish - Galician (gl), Hindi (hi), Croatian (hr), 
Hungarian (hu), Armenian (hy), Indonesian (id, in), Italian (it), Dutch (nl), Norwegian (no), 
Polish (pl), Portuguese (pt), Romanian (ro), Russian (ru), Swedish (sv), Turkish (tr)

#### Faceting

Support for aggregating facet information across search results. Supported facet types:

- Terms Facet
- Numeric Range Facet
- Date Range Facet

#### Powerful

By indexing your data with bleve you gain the ability to compose the following query types:

- Term, Phrase, Match, Match Phrase, Prefix, Fuzzy
- Conjunction, Disjunction, Boolean
- Numeric and Date Ranges
- Query String (see [Syntax](https://blevesearch.com/docs/Query-String-Query/))
- Approximate k-nearest-neighbors over vector content

#### Scoring

Industry standard [tf-idf](http://en.wikipedia.org/wiki/Tf%E2%80%93idf), 
[bm25](https://en.wikipedia.org/wiki/Okapi_BM25) scoring with query time boosting.

#### Result Highlighting

Includes support for

```
highlighting matching text
```
within document fragments.

#### Text analysis wizard

Here's a [playground](https://bleveanalysis.couchbase.com/) where you can demo various text 
analysis components.

#### Open Source

The [bleve source](http://github.com/blevesearch/bleve) is available on github and distributed 
under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).
