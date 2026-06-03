---
Title: Source: bleve-github-readme.md
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

## bleve

[![Tests](https://github.com/blevesearch/bleve/actions/workflows/tests.yml/badge.svg?branch=master&e
vent=push)](https://github.com/blevesearch/bleve/actions/workflows/tests.yml?query=event%3Apush+bran
ch%3Amaster) [![Coverage 
Status](https://camo.githubusercontent.com/072ed4491aa545f915e2010fd80aa7d72c9d1249f8b002d387fa8f043
79b1df3/68747470733a2f2f636f766572616c6c732e696f2f7265706f732f6769746875622f626c6576657365617263682f
626c6576652f62616467652e737667)](https://coveralls.io/github/blevesearch/bleve) [![Go 
Reference](https://camo.githubusercontent.com/eab50b58f98fbeb8c9c48cf7595ffdf9f02c0b8dadf071bbaf69f6
6374881f6d/68747470733a2f2f706b672e676f2e6465762f62616467652f6769746875622e636f6d2f626c6576657365617
263682f626c6576652f76322e737667)](https://pkg.go.dev/github.com/blevesearch/bleve/v2) [![Join the 
chat](https://camo.githubusercontent.com/0b8f43b6491f1f8c4fa2c807b26578b4c0fa72dc62e6a82fe5e3013f719
f1353/68747470733a2f2f6261646765732e6769747465722e696d2f6a6f696e5f636861742e737667)](https://app.git
ter.im/#/room/#blevesearch_bleve:gitter.im) [![Go Report 
Card](https://camo.githubusercontent.com/63f3b71694f12bf267eb4620e707ae1fe5140a039bc4f001f646b1f81fb
3d761/68747470733a2f2f676f7265706f7274636172642e636f6d2f62616467652f6769746875622e636f6d2f626c657665
7365617263682f626c6576652f7632)](https://goreportcard.com/report/github.com/blevesearch/bleve/v2) 
[![Sourcegraph](https://camo.githubusercontent.com/a74f60b2c25009065994ea1ade069f45c653de808abab0f4a
4bd9a2fee308daf/68747470733a2f2f736f7572636567726170682e636f6d2f6769746875622e636f6d2f626c6576657365
617263682f626c6576652f2d2f62616467652e737667)](https://sourcegraph.com/github.com/blevesearch/bleve?
badge) 
[![License](https://camo.githubusercontent.com/a549a7a30bacba7bfceebdc207a8e86c3f2c02995a2527640dca3
0048fd2b64e/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f4c6963656e73652d41706163686525
3230322e302d626c75652e737667)](https://opensource.org/licenses/Apache-2.0)

A modern indexing + search library in GO

## Features

- Index any GO data structure or JSON
- Intelligent defaults backed up by powerful configuration 
([scorch](https://github.com/blevesearch/bleve/blob/master/index/scorch/README.md))
- Supported field types:
	- `text`, `number`, `datetime`, `boolean`, `geopoint`, `geoshape`, `IP`, `vector`
- Supported query types:
	- `term`, `phrase`, `match`, `match_phrase`, `prefix`, `regexp`, `wildcard`, `fuzzy`
		- term range, numeric range, date range, boolean field
		- compound queries: `conjuncts`, `disjuncts`, boolean (`must` / `should` / 
`must_not`)
		- [query string syntax](http://www.blevesearch.com/docs/Query-String-Query/)
		- [geo spatial 
search](https://github.com/blevesearch/bleve/blob/master/geo/README.md)
		- approximate k-nearest neighbors via [vector 
search](https://github.com/blevesearch/bleve/blob/master/docs/vectors.md)
		- [synonym 
search](https://github.com/blevesearch/bleve/blob/master/docs/synonyms.md)
		- [hierarchical nested 
search](https://github.com/blevesearch/bleve/blob/master/docs/hierarchy.md)
- [tf-idf](https://github.com/blevesearch/bleve/blob/master/docs/scoring.md#tf-idf) / 
[bm25](https://github.com/blevesearch/bleve/blob/master/docs/scoring.md#bm25) scoring models
- Hybrid search: exact + semantic
	- Supports [RRF (Reciprocal Rank Fusion) and RSF (Relative Score 
Fusion)](https://github.com/blevesearch/bleve/blob/master/docs/score_fusion.md)
- [Result pagination](https://github.com/blevesearch/bleve/blob/master/docs/pagination.md)
- Query time boosting
- Search result match highlighting with document fragments
- Aggregations/faceting support:
	- terms facet
		- numeric range facet
		- date range facet

## Indexing

```
message := struct {
    Id   string
    From string
    Body string
}{
    Id:   "example",
    From: "xyz@couchbase.com",
    Body: "bleve indexing is easy",
}

mapping := bleve.NewIndexMapping()
index, err := bleve.New("example.bleve", mapping)
if err != nil {
    panic(err)
}
index.Index(message.Id, message)
```

## Querying

```
index, _ := bleve.Open("example.bleve")
query := bleve.NewQueryStringQuery("bleve")
searchRequest := bleve.NewSearchRequest(query)
searchResult, _ := index.Search(searchRequest)
```

## Command Line Interface

To install the CLI for the latest release of bleve, run:

```
go install github.com/blevesearch/bleve/v2/cmd/bleve@latest
```

```
$ bleve --help
Bleve is a command-line tool to interact with a bleve index.

Usage:
  bleve [command]

Available Commands:
  bulk        bulk loads from newline delimited JSON files
  check       checks the contents of the index
  count       counts the number documents in the index
  create      creates a new index
  dictionary  prints the term dictionary for the specified field in the index
  dump        dumps the contents of the index
  fields      lists the fields in this index
  help        Help about any command
  index       adds the files to the index
  mapping     prints the mapping used for this index
  query       queries the index
  registry    registry lists the bleve components compiled into this executable
  scorch      command-line tool to interact with a scorch index

Flags:
  -h, --help   help for bleve

Use "bleve [command] --help" for more information about a command.
```

## Text Analysis

Bleve includes general-purpose analyzers (customizable) as well as pre-built text analyzers for the 
following languages:

Arabic (ar), Bulgarian (bg), Catalan (ca), Chinese-Japanese-Korean (cjk), Kurdish (ckb), Danish 
(da), German (de), Greek (el), English (en), Spanish - Castilian (es), Basque (eu), Persian (fa), 
Finnish (fi), French (fr), Gaelic (ga), Spanish - Galician (gl), Hindi (hi), Croatian (hr), 
Hungarian (hu), Armenian (hy), Indonesian (id, in), Italian (it), Dutch (nl), Norwegian (no), 
Polish (pl), Portuguese (pt), Romanian (ro), Russian (ru), Swedish (sv), Turkish (tr)
