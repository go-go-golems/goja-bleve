---
Title: Source: bleve-pkg-go-dev.md
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

## README

### bleve

[![Tests](https://github.com/blevesearch/bleve/actions/workflows/tests.yml/badge.svg?branch=master&e
vent=push)](https://github.com/blevesearch/bleve/actions/workflows/tests.yml?query=event%3Apush+bran
ch%3Amaster) [![Coverage 
Status](https://coveralls.io/repos/github/blevesearch/bleve/badge.svg)](https://coveralls.io/github/
blevesearch/bleve) [![Go 
Reference](https://pkg.go.dev/badge/github.com/blevesearch/bleve/v2.svg)](https://pkg.go.dev/github.
com/blevesearch/bleve/v2) [![Join the 
chat](https://badges.gitter.im/join_chat.svg)](https://app.gitter.im/#/room/%23blevesearch_bleve:git
ter.im) [![Go Report 
Card](https://goreportcard.com/badge/github.com/blevesearch/bleve/v2)](https://goreportcard.com/repo
rt/github.com/blevesearch/bleve/v2) 
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/lice
nses/Apache-2.0)

A modern indexing + search library in GO

#### Features

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
Fusion)](https://github.com/blevesearch/bleve/blob/v2.6.0/docs/score_fusion.md)
- [Result pagination](https://github.com/blevesearch/bleve/blob/master/docs/pagination.md)
- Query time boosting
- Search result match highlighting with document fragments
- Aggregations/faceting support:
	- terms facet
		- numeric range facet
		- date range facet

#### Indexing

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

#### Querying

```
index, _ := bleve.Open("example.bleve")
query := bleve.NewQueryStringQuery("bleve")
searchRequest := bleve.NewSearchRequest(query)
searchResult, _ := index.Search(searchRequest)
```

#### Command Line Interface

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

#### Text Analysis

Bleve includes general-purpose analyzers (customizable) as well as pre-built text analyzers for the 
following languages:

Arabic (ar), Bulgarian (bg), Catalan (ca), Chinese-Japanese-Korean (cjk), Kurdish (ckb), Danish 
(da), German (de), Greek (el), English (en), Spanish - Castilian (es), Basque (eu), Persian (fa), 
Finnish (fi), French (fr), Gaelic (ga), Spanish - Galician (gl), Hindi (hi), Croatian (hr), 
Hungarian (hu), Armenian (hy), Indonesian (id, in), Italian (it), Dutch (nl), Norwegian (no), 
Polish (pl), Portuguese (pt), Romanian (ro), Russian (ru), Swedish (sv), Turkish (tr)

#### Text Analysis Wizard

[bleveanalysis.couchbase.com](https://bleveanalysis.couchbase.com/)

#### Discussion/Issues

Discuss usage/development of bleve and/or report issues here:

- [Github issues](https://github.com/blevesearch/bleve/issues)
- [Google group](https://groups.google.com/forum/#!forum/bleve)

#### License

Apache License Version 2.0

## Documentation

### Overview

Package bleve is a library for indexing and searching text.

Example Opening New Index, Indexing Data

```
message := struct{
    Id:   "example"
    From: "xyz@couchbase.com",
    Body: "bleve indexing is easy",
}

mapping := bleve.NewIndexMapping()
index, _ := bleve.New("example.bleve", mapping)
index.Index(message.Id, message)
```

Example Opening Existing Index, Searching Data

```
index, _ := bleve.Open("example.bleve")
query := bleve.NewQueryStringQuery("bleve")
searchRequest := bleve.NewSearchRequest(query)
searchResult, _ := index.Search(searchRequest)
```

### Index

### Examples

### Constants

```
const (
    SearchQueryStartCallbackKey search.ContextKey = "_search_query_start_callback_key"
    SearchQueryEndCallbackKey   search.ContextKey = "_search_query_end_callback_key"
)
```

```
const (
    ScoreDefault = ""
    ScoreNone    = "none"
    ScoreRRF     = "rrf"
    ScoreRSF     = "rsf"
)
```

```
const (
    DefaultScoreRankConstant = 60
)
```

```
const NestedDocumentKey = "_$nested"
```

### Variables

```
var AllowedFusionSort = search.SortOrder{&search.SortScore{Desc: true}}
```

```
var Config *configuration
```

Config contains library level configuration

### Functions

#### added in v2.5.4

```
func DeletedFields(ori, upd *mapping.IndexMappingImpl) (map[string]*index.UpdateFieldInfo, error)
```

Compare two index mappings to identify all of the updatable changes

#### added in v2.5.4

```
func IsScoreFusionRequested(req *SearchRequest) bool
```

Checks if the request is hybrid search. Currently supports: RRF, RSF.

#### added in v2.6.0

```
func LoadAndHighlightAllFields(
    root *search.DocumentMatch,
    req *SearchRequest,
    indexName string,
    r index.IndexReader,
    highlighter highlight.Highlighter,
) (error, uint64)
```

LoadAndHighlightAllFields loads stored fields + highlights for root and its descendants. All 
descendant documents are collected into a \_$nested array in the root DocumentMatch.

#### func LoadAndHighlightFields ¶

```
func LoadAndHighlightFields(hit *search.DocumentMatch, req *SearchRequest,
    indexName string, r index.IndexReader,
    highlighter highlight.Highlighter,
) (error, uint64)
```

#### func MemoryNeededForSearchResult ¶

```
func MemoryNeededForSearchResult(req *SearchRequest) uint64
```

MemoryNeededForSearchResult is an exported helper function to determine the RAM needed to 
accommodate the results for a given search request.

#### func NewBoolFieldQuery ¶

```
func NewBoolFieldQuery(val bool) *query.BoolFieldQuery
```

NewBoolFieldQuery creates a new Query for boolean fields

#### func NewBooleanFieldMapping ¶

```
func NewBooleanFieldMapping() *mapping.FieldMapping
```

NewBooleanFieldMapping returns a default field mapping for booleans

#### func NewBooleanQuery ¶

```
func NewBooleanQuery() *query.BooleanQuery
```

NewBooleanQuery creates a compound Query composed of several other Query objects. These other query 
objects are added using the AddMust() AddShould() and AddMustNot() methods. Result documents must 
satisfy ALL of the must Queries. Result documents must satisfy NONE of the must not Queries. Result 
documents that ALSO satisfy any of the should Queries will score higher.

Example [¶](#example-NewBooleanQuery "Go to Example")

```
must := NewMatchQuery("one")
mustNot := NewMatchQuery("great")
query := NewBooleanQuery()
query.AddMust(must)
query.AddMustNot(mustNot)
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 1
```

#### func NewConjunctionQuery ¶

```
func NewConjunctionQuery(conjuncts ...query.Query) *query.ConjunctionQuery
```

NewConjunctionQuery creates a new compound Query. Result documents must satisfy all of the queries.

Example [¶](#example-NewConjunctionQuery "Go to Example")

```
conjunct1 := NewMatchQuery("great")
conjunct2 := NewMatchQuery("one")
query := NewConjunctionQuery(conjunct1, conjunct2)
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 2
```

#### func NewDateRangeInclusiveQuery ¶

```
func NewDateRangeInclusiveQuery(start, end time.Time, startInclusive, endInclusive *bool) 
*query.DateRangeQuery
```

NewDateRangeInclusiveQuery creates a new Query for ranges of date values. Date strings are parsed 
using the DateTimeParser configured in the

```
top-level config.QueryDateTimeParser
```

Either, but not both endpoints can be nil. startInclusive and endInclusive control inclusion of the 
endpoints.

#### added in v2.3.10

```
func NewDateRangeInclusiveStringQuery(start, end string, startInclusive, endInclusive *bool) 
*query.DateRangeStringQuery
```

NewDateRangeInclusiveStringQuery creates a new Query for ranges of date values. Date strings are 
parsed using the DateTimeParser set using

```
the DateRangeStringQuery.SetDateTimeParser() method.
```

this DateTimeParser is a custom date time parser defined in the index mapping, using 
AddCustomDateTimeParser() method. If no DateTimeParser is set, then the

```
top-level config.QueryDateTimeParser
```

is used. Either, but not both endpoints can be nil. startInclusive and endInclusive control 
inclusion of the endpoints.

#### func NewDateRangeQuery ¶

```
func NewDateRangeQuery(start, end time.Time) *query.DateRangeQuery
```

NewDateRangeQuery creates a new Query for ranges of date values. Date strings are parsed using the 
DateTimeParser configured in the

```
top-level config.QueryDateTimeParser
```

Either, but not both endpoints can be nil.

#### added in v2.3.10

```
func NewDateRangeStringQuery(start, end string) *query.DateRangeStringQuery
```

NewDateRangeStringQuery creates a new Query for ranges of date values. Date strings are parsed 
using the DateTimeParser set using

```
the DateRangeStringQuery.SetDateTimeParser() method.
```

If no DateTimeParser is set, then the

```
top-level config.QueryDateTimeParser
```

is used.

#### func NewDateTimeFieldMapping ¶

```
func NewDateTimeFieldMapping() *mapping.FieldMapping
```

NewDateTimeFieldMapping returns a default field mapping for dates

#### func NewDisjunctionQuery ¶

```
func NewDisjunctionQuery(disjuncts ...query.Query) *query.DisjunctionQuery
```

NewDisjunctionQuery creates a new compound Query. Result documents satisfy at least one Query.

Example [¶](#example-NewDisjunctionQuery "Go to Example")

```
disjunct1 := NewMatchQuery("great")
disjunct2 := NewMatchQuery("named")
query := NewDisjunctionQuery(disjunct1, disjunct2)
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(len(searchResults.Hits))
```
```
Output:
2
```

#### func NewDocIDQuery ¶

```
func NewDocIDQuery(ids []string) *query.DocIDQuery
```

NewDocIDQuery creates a new Query object returning indexed documents among the specified set. 
Combine it with ConjunctionQuery to restrict the scope of other queries output.

#### func NewDocumentDisabledMapping ¶

```
func NewDocumentDisabledMapping() *mapping.DocumentMapping
```

NewDocumentDisabledMapping returns a new document mapping that will not perform any indexing.

#### func NewDocumentMapping ¶

```
func NewDocumentMapping() *mapping.DocumentMapping
```

NewDocumentMapping returns a new document mapping with all the default values.

#### func NewDocumentStaticMapping ¶

```
func NewDocumentStaticMapping() *mapping.DocumentMapping
```

NewDocumentStaticMapping returns a new document mapping that will not automatically index parts of 
a document without an explicit mapping.

#### func NewFuzzyQuery ¶

```
func NewFuzzyQuery(term string) *query.FuzzyQuery
```

NewFuzzyQuery creates a new Query which finds documents containing terms within a specific 
fuzziness of the specified term. The default fuzziness is 1.

The current implementation uses Levenshtein edit distance as the fuzziness metric.

#### func NewGeoBoundingBoxQuery ¶

```
func NewGeoBoundingBoxQuery(topLeftLon, topLeftLat, bottomRightLon, bottomRightLat float64) 
*query.GeoBoundingBoxQuery
```

NewGeoBoundingBoxQuery creates a new Query for performing geo bounding box searches. The arguments 
describe the position of the box and documents which have an indexed geo point inside the box will 
be returned.

#### func NewGeoDistanceQuery ¶

```
func NewGeoDistanceQuery(lon, lat float64, distance string) *query.GeoDistanceQuery
```

NewGeoDistanceQuery creates a new Query for performing geo distance searches. The arguments 
describe a position and a distance. Documents which have an indexed geo point which is less than or 
equal to the provided distance from the given position will be returned.

#### func NewGeoPointFieldMapping ¶

```
func NewGeoPointFieldMapping() *mapping.FieldMapping
```

#### added in v2.3.9

```
func NewGeoShapeCircleQuery(coordinates []float64, radius, relation string) (*query.GeoShapeQuery, 
error)
```

NewGeoShapeCircleQuery creates a new query for a geoshape that is a circle given center point and 
the radius. Radius formats supported: "5in" "5inch" "7yd" "7yards" "9ft" "9feet" "11km" 
"11kilometers" "3nm" "3nauticalmiles" "13mm" "13millimeters" "15cm" "15centimeters" "17mi" 
"17miles" "19m" "19meters" If the unit cannot be determined, the entire string is parsed and the 
unit of meters is assumed.

#### added in v2.3.5

```
func NewGeoShapeFieldMapping() *mapping.FieldMapping
```

#### added in v2.3.9

```
func NewGeoShapeQuery(coordinates [][][][]float64, typ, relation string) (*query.GeoShapeQuery, 
error)
```

NewGeoShapeQuery creates a new Query for matching the given geo shape. This method can be used for 
creating geoshape queries for shape types like: point, linestring, polygon, multipoint, 
multilinestring, multipolygon and envelope.

#### added in v2.3.9

```
func NewGeometryCollectionQuery(coordinates [][][][][]float64, types []string, relation string) 
(*query.GeoShapeQuery, error)
```

NewGeometryCollectionQuery creates a new query for the provided geometrycollection coordinates and 
types, which could contain multiple geo shapes.

#### added in v2.3.0

```
func NewIPFieldMapping() *mapping.FieldMapping
```

#### added in v2.3.0

```
func NewIPRangeQuery(cidr string) *query.IPRangeQuery
```

NewIPRangeQuery creates a new Query for matching IP addresses. If the argument is in CIDR format, 
then the query will match all IP addresses in the network specified. If the argument is an IP 
address, then the query will return documents which contain that IP. Both ipv4 and ipv6 are 
supported.

#### func NewIndexAlias ¶

```
func NewIndexAlias(indexes ...Index) *indexAliasImpl
```

NewIndexAlias creates a new IndexAlias over the provided Index objects.

#### func NewIndexMapping ¶

```
func NewIndexMapping() *mapping.IndexMappingImpl
```

NewIndexMapping creates a new IndexMapping that will use all the default indexing rules

#### added in v2.3.0

```
func NewKeywordFieldMapping() *mapping.FieldMapping
```

NewKeywordFieldMapping returns a field mapping for text using the keyword analyzer, which 
essentially doesn't apply any specific text analysis.

#### func NewMatchAllQuery ¶

```
func NewMatchAllQuery() *query.MatchAllQuery
```

NewMatchAllQuery creates a Query which will match all documents in the index.

Example [¶](#example-NewMatchAllQuery "Go to Example")

```
// finds all documents in the index
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(len(searchResults.Hits))
```
```
Output:
2
```

#### func NewMatchNoneQuery ¶

```
func NewMatchNoneQuery() *query.MatchNoneQuery
```

NewMatchNoneQuery creates a Query which will not match any documents in the index.

Example [¶](#example-NewMatchNoneQuery "Go to Example")

```
// matches no documents in the index
query := NewMatchNoneQuery()
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(len(searchResults.Hits))
```
```
Output:
0
```

#### func NewMatchPhraseQuery ¶

```
func NewMatchPhraseQuery(matchPhrase string) *query.MatchPhraseQuery
```

NewMatchPhraseQuery creates a new Query object for matching phrases in the index. An Analyzer is 
chosen based on the field. Input text is analyzed using this analyzer. Token terms resulting from 
this analysis are used to build a search phrase. Result documents must match this phrase. Queried 
field must have been indexed with IncludeTermVectors set to true.

Example [¶](#example-NewMatchPhraseQuery "Go to Example")

```
// finds all documents with the given phrase in the index
query := NewMatchPhraseQuery("nameless one")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 2
```

#### func NewMatchQuery ¶

```
func NewMatchQuery(match string) *query.MatchQuery
```

NewMatchQuery creates a Query for matching text. An Analyzer is chosen based on the field. Input 
text is analyzed using this analyzer. Token terms resulting from this analysis are used to perform 
term searches. Result documents must satisfy at least one of these term searches.

Example [¶](#example-NewMatchQuery "Go to Example")

```
// finds documents with fields fully matching the given query text
query := NewMatchQuery("named one")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 1
```

#### added in v2.6.0

```
func NewNestedDocumentMapping() *mapping.DocumentMapping
```

NewNestedDocumentMapping returns a new document mapping that will treat all objects as nested 
documents.

#### added in v2.6.0

```
func NewNestedDocumentStaticMapping() *mapping.DocumentMapping
```

NewNestedDocumentStaticMapping returns a new document mapping that will treat all objects as nested 
documents and will not automatically index parts of a nested document without an explicit mapping.

#### func NewNumericFieldMapping ¶

```
func NewNumericFieldMapping() *mapping.FieldMapping
```

NewNumericFieldMapping returns a default field mapping for numbers

#### func NewNumericRangeInclusiveQuery ¶

```
func NewNumericRangeInclusiveQuery(min, max *float64, minInclusive, maxInclusive *bool) 
*query.NumericRangeQuery
```

NewNumericRangeInclusiveQuery creates a new Query for ranges of numeric values. Either, but not 
both endpoints can be nil. Control endpoint inclusion with inclusiveMin, inclusiveMax.

Example [¶](#example-NewNumericRangeInclusiveQuery "Go to Example")

```
value1 := float64(10)
value2 := float64(100)
v1incl := false
v2incl := false

query := NewNumericRangeInclusiveQuery(&value1, &value2, &v1incl, &v2incl)
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 3
```

#### func NewNumericRangeQuery ¶

```
func NewNumericRangeQuery(min, max *float64) *query.NumericRangeQuery
```

NewNumericRangeQuery creates a new Query for ranges of numeric values. Either, but not both 
endpoints can be nil. The minimum value is inclusive. The maximum value is exclusive.

Example [¶](#example-NewNumericRangeQuery "Go to Example")

```
value1 := float64(11)
value2 := float64(100)
data := struct{ Priority float64 }{Priority: float64(15)}
data2 := struct{ Priority float64 }{Priority: float64(10)}

err = exampleIndex.Index("document id 3", data)
if err != nil {
    panic(err)
}
err = exampleIndex.Index("document id 4", data2)
if err != nil {
    panic(err)
}

query := NewNumericRangeQuery(&value1, &value2)
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 3
```

#### func NewPhraseQuery ¶

```
func NewPhraseQuery(terms []string, field string) *query.PhraseQuery
```

NewPhraseQuery creates a new Query for finding exact term phrases in the index. The provided terms 
must exist in the correct order, at the correct index offsets, in the specified field. Queried 
field must have been indexed with IncludeTermVectors set to true.

Example [¶](#example-NewPhraseQuery "Go to Example")

```
// finds all documents with the given phrases in the given field in the index
query := NewPhraseQuery([]string{"nameless", "one"}, "Name")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 2
```

#### func NewPrefixQuery ¶

```
func NewPrefixQuery(prefix string) *query.PrefixQuery
```

NewPrefixQuery creates a new Query which finds documents containing terms that start with the 
specified prefix.

Example [¶](#example-NewPrefixQuery "Go to Example")

```
// finds all documents with terms having the given prefix in the index
query := NewPrefixQuery("name")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(len(searchResults.Hits))
```
```
Output:
2
```

#### func NewQueryStringQuery ¶

```
func NewQueryStringQuery(q string) *query.QueryStringQuery
```

NewQueryStringQuery creates a new Query used for finding documents that satisfy a query string. The 
query string is a small query language for humans.

Example [¶](#example-NewQueryStringQuery "Go to Example")

```
query := NewQueryStringQuery("+one -great")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 1
```

#### func NewRegexpQuery ¶

```
func NewRegexpQuery(regexp string) *query.RegexpQuery
```

NewRegexpQuery creates a new Query which finds documents containing terms that match the specified 
regular expression.

#### func NewTermQuery ¶

```
func NewTermQuery(term string) *query.TermQuery
```

NewTermQuery creates a new Query for finding an exact term match in the index.

Example [¶](#example-NewTermQuery "Go to Example")

```
query := NewTermQuery("great")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 2
```

#### func NewTermRangeInclusiveQuery ¶

```
func NewTermRangeInclusiveQuery(min, max string, minInclusive, maxInclusive *bool) 
*query.TermRangeQuery
```

NewTermRangeInclusiveQuery creates a new Query for ranges of text terms. Either, but not both 
endpoints can be "". Control endpoint inclusion with inclusiveMin, inclusiveMax.

#### func NewTermRangeQuery ¶

```
func NewTermRangeQuery(min, max string) *query.TermRangeQuery
```

NewTermRangeQuery creates a new Query for ranges of text terms. Either, but not both endpoints can 
be "". The minimum value is inclusive. The maximum value is exclusive.

#### func NewTextFieldMapping ¶

```
func NewTextFieldMapping() *mapping.FieldMapping
```

NewTextFieldMapping returns a default field mapping for text

#### func NewWildcardQuery ¶

```
func NewWildcardQuery(wildcard string) *query.WildcardQuery
```

NewWildcardQuery creates a new Query which finds documents containing terms that match the 
specified wildcard. In the wildcard pattern '\*' will match any sequence of 0 or more characters, 
and '?' will match any single character.

#### func SetLog ¶

```
func SetLog(l *log.Logger)
```

SetLog sets the logger used for logging by default log messages are sent to io.Discard

### Types

#### type Batch ¶

```
type Batch struct {
    // contains filtered or unexported fields
}
```

A Batch groups together multiple Index and Delete operations you would like performed at the same 
time. The Batch structure is NOT thread-safe. You should only perform operations on a batch from a 
single thread at a time. Once batch execution has started, you may not modify it.

#### func (\*Batch) Delete ¶

```
func (b *Batch) Delete(id string)
```

Delete adds the specified delete operation to the batch. NOTE: the bleve Index is not updated until 
the batch is executed.

#### func (\*Batch) DeleteInternal ¶

```
func (b *Batch) DeleteInternal(key []byte)
```

DeleteInternal adds the specified delete internal operation to the batch. NOTE: the bleve Index is 
not updated until the batch is executed.

#### func (\*Batch) Index ¶

```
func (b *Batch) Index(id string, data interface{}) error
```

Index adds the specified index operation to the batch. NOTE: the bleve Index is not updated until 
the batch is executed.

#### func (\*Batch) IndexAdvanced ¶

```
func (b *Batch) IndexAdvanced(doc *document.Document) (err error)
```

IndexAdvanced adds the specified index operation to the batch which skips the mapping. NOTE: the 
bleve Index is not updated until the batch is executed.

#### added in v2.5.0

```
func (b *Batch) IndexSynonym(id string, collection string, definition *SynonymDefinition) error
```

#### func (\*Batch) LastDocSize ¶

```
func (b *Batch) LastDocSize() uint64
```

#### func (\*Batch) Merge ¶

```
func (b *Batch) Merge(o *Batch)
```

#### func (\*Batch) PersistedCallback ¶

```
func (b *Batch) PersistedCallback() index.BatchCallback
```

#### func (\*Batch) Reset ¶

```
func (b *Batch) Reset()
```

Reset returns a Batch to the empty state so that it can be reused in the future.

#### func (\*Batch) SetInternal ¶

```
func (b *Batch) SetInternal(key, val []byte)
```

SetInternal adds the specified set internal operation to the batch. NOTE: the bleve Index is not 
updated until the batch is executed.

#### func (\*Batch) SetPersistedCallback ¶

```
func (b *Batch) SetPersistedCallback(f index.BatchCallback)
```

#### func (\*Batch) Size ¶

```
func (b *Batch) Size() int
```

Size returns the total number of operations inside the batch including normal index operations and 
internal operations.

#### func (\*Batch) String ¶

```
func (b *Batch) String() string
```

String prints a user friendly string representation of what is inside this batch.

#### func (\*Batch) TotalDocsSize ¶

```
func (b *Batch) TotalDocsSize() uint64
```

#### type Builder ¶

```
type Builder interface {
    Index(id string, data interface{}) error
    Close() error
}
```

Builder is a limited interface, used to build indexes in an offline mode. Items cannot be updated 
or deleted, and the caller MUST ensure a document is indexed only once.

#### func NewBuilder ¶

```
func NewBuilder(path string, mapping mapping.IndexMapping, config map[string]interface{}) (Builder, 
error)
```

NewBuilder creates a builder, which will build an index at the specified path, using the specified 
mapping and options.

#### type Error ¶

```
type Error int
```

Error represents a more strongly typed bleve error for detecting and handling specific types of 
errors.

```
const (
    ErrorIndexPathExists Error = iota
    ErrorIndexPathDoesNotExist
    ErrorIndexMetaMissing
    ErrorIndexMetaCorrupt
    ErrorIndexClosed
    ErrorAliasMulti
    ErrorAliasEmpty
    ErrorUnknownIndexType
    ErrorEmptyID
    ErrorIndexReadInconsistency
    ErrorTwoPhaseSearchInconsistency
    ErrorSynonymSearchNotSupported
    ErrorTrainingNotSupported
)
```

Constant Error values which can be compared to determine the type of error

#### func (Error) Error ¶

```
func (e Error) Error() string
```

#### type FacetRequest ¶

```
type FacetRequest struct {
    Size           int              \`json:"size"\`
    Field          string           \`json:"field"\`
    TermPrefix     string           \`json:"term_prefix,omitempty"\`
    TermPattern    string           \`json:"term_pattern,omitempty"\`
    NumericRanges  []*numericRange  \`json:"numeric_ranges,omitempty"\`
    DateTimeRanges []*dateTimeRange \`json:"date_ranges,omitempty"\`
    // contains filtered or unexported fields
}
```

A FacetRequest describes a facet or aggregation of the result document set you would like to be 
built.

#### func NewFacetRequest ¶

```
func NewFacetRequest(field string, size int) *FacetRequest
```

NewFacetRequest creates a facet on the specified field that limits the number of entries to the 
specified size.

Example [¶](#example-NewFacetRequest "Go to Example")

```
facet := NewFacetRequest("Name", 1)
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchRequest.AddFacet("facet name", facet)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

// total number of terms
fmt.Println(searchResults.Facets["facet name"].Total)
// number of docs with no value for this field
fmt.Println(searchResults.Facets["facet name"].Missing)
// term with highest occurrences in field name
fmt.Println(searchResults.Facets["facet name"].Terms.Terms()[0].Term)
```
```
Output:
5
2
one
```

#### func (\*FacetRequest) AddDateTimeRange ¶

```
func (fr *FacetRequest) AddDateTimeRange(name string, start, end time.Time)
```

AddDateTimeRange adds a bucket to a field containing date values. Documents with a date value 
falling into this range are tabulated as part of this bucket/range.

Example [¶](#example-FacetRequest.AddDateTimeRange "Go to Example")

```
facet := NewFacetRequest("Created", 1)
facet.AddDateTimeRange("range name", time.Unix(0, 0), time.Now())
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchRequest.AddFacet("facet name", facet)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

// dates in field Created since starting of unix time till now
fmt.Println(searchResults.Facets["facet name"].DateRanges[0].Count)
```
```
Output:
2
```

#### func (\*FacetRequest) AddDateTimeRangeString ¶

```
func (fr *FacetRequest) AddDateTimeRangeString(name string, start, end *string)
```

AddDateTimeRangeString adds a bucket to a field containing date values. Uses defaultDateTimeParser 
to parse the date strings.

#### added in v2.3.10

```
func (fr *FacetRequest) AddDateTimeRangeStringWithParser(name string, start, end *string, parser 
string)
```

AddDateTimeRangeString adds a bucket to a field containing date values. Uses the specified parser 
to parse the date strings. provided the parser is registered in the index mapping.

#### func (\*FacetRequest) AddNumericRange ¶

```
func (fr *FacetRequest) AddNumericRange(name string, min, max *float64)
```

AddNumericRange adds a bucket to a field containing numeric values. Documents with a numeric value 
falling into this range are tabulated as part of this bucket/range.

Example [¶](#example-FacetRequest.AddNumericRange "Go to Example")

```
value1 := float64(11)

facet := NewFacetRequest("Priority", 1)
facet.AddNumericRange("range name", &value1, nil)
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchRequest.AddFacet("facet name", facet)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

// number documents with field Priority in the given range
fmt.Println(searchResults.Facets["facet name"].NumericRanges[0].Count)
```
```
Output:
1
```

#### added in v2.5.6

```
func (fr *FacetRequest) SetPrefixFilter(prefix string)
```

SetPrefixFilter sets the prefix filter for term facets.

#### added in v2.5.6

```
func (fr *FacetRequest) SetRegexFilter(pattern string)
```

SetRegexFilter sets the regex pattern filter for term facets.

#### func (\*FacetRequest) Validate ¶

```
func (fr *FacetRequest) Validate() error
```

#### type FacetsRequest ¶

```
type FacetsRequest map[string]*FacetRequest
```

FacetsRequest groups together all the FacetRequest objects for a single query.

#### func (FacetsRequest) Validate ¶

```
func (fr FacetsRequest) Validate() error
```

#### added in v2.1.0

```
type FileSystemDirectory string
```

FileSystemDirectory is the default implementation for the index.Directory interface.

#### added in v2.1.0

```
func (f FileSystemDirectory) GetWriter(filePath string) (io.WriteCloser,
    error,
)
```

#### type HighlightRequest ¶

```
type HighlightRequest struct {
    Style  *string  \`json:"style"\`
    Fields []string \`json:"fields"\`
}
```

HighlightRequest describes how field matches should be highlighted.

#### func NewHighlight ¶

```
func NewHighlight() *HighlightRequest
```

NewHighlight creates a default HighlightRequest.

Example [¶](#example-NewHighlight "Go to Example")

```
query := NewMatchQuery("nameless")
searchRequest := NewSearchRequest(query)
searchRequest.Highlight = NewHighlight()
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].Fragments["Name"][0])
```
```
Output:
great <mark>nameless</mark> one
```

#### func NewHighlightWithStyle ¶

```
func NewHighlightWithStyle(style string) *HighlightRequest
```

NewHighlightWithStyle creates a HighlightRequest with an alternate style.

Example [¶](#example-NewHighlightWithStyle "Go to Example")

```
query := NewMatchQuery("nameless")
searchRequest := NewSearchRequest(query)
searchRequest.Highlight = NewHighlightWithStyle(ansi.Name)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].Fragments["Name"][0])
```
```
Output:
great �[43mnameless�[0m one
```

#### func (\*HighlightRequest) AddField ¶

```
func (h *HighlightRequest) AddField(field string)
```

#### type Index ¶

```
type Index interface {
    // Index analyzes, indexes or stores mapped data fields. Supplied
    // identifier is bound to analyzed data and will be retrieved by search
    // requests. See Index interface documentation for details about mapping
    // rules.
    Index(id string, data interface{}) error
    Delete(id string) error

    NewBatch() *Batch
    Batch(b *Batch) error

    // Document returns specified document or nil if the document is not
    // indexed or stored.
    Document(id string) (index.Document, error)
    // DocCount returns the number of documents in the index.
    DocCount() (uint64, error)

    Search(req *SearchRequest) (*SearchResult, error)
    SearchInContext(ctx context.Context, req *SearchRequest) (*SearchResult, error)

    Fields() ([]string, error)

    FieldDict(field string) (index.FieldDict, error)
    FieldDictRange(field string, startTerm []byte, endTerm []byte) (index.FieldDict, error)
    FieldDictPrefix(field string, termPrefix []byte) (index.FieldDict, error)

    Close() error

    Mapping() mapping.IndexMapping

    Stats() *IndexStat
    StatsMap() map[string]interface{}

    GetInternal(key []byte) ([]byte, error)
    SetInternal(key, val []byte) error
    DeleteInternal(key []byte) error

    // Name returns the name of the index (by default this is the path)
    Name() string
    // SetName lets you assign your own logical name to this index
    SetName(string)

    // Advanced returns the internal index implementation
    Advanced() (index.Index, error)
}
```

An Index implements all the indexing and searching capabilities of bleve. An Index can be created 
using the New() and Open() methods.

Index() takes an input value, deduces a DocumentMapping for its type, assigns string paths to its 
fields or values then applies field mappings on them.

The DocumentMapping used to index a value is deduced by the following rules:

1. If value implements mapping.bleveClassifier interface, resolve the mapping from BleveType().
2. If value implements mapping.Classifier interface, resolve the mapping from Type().
3. If value has a string field or value at IndexMapping.TypeField.

(defaulting to "\_type"), use it to resolve the mapping. Fields addressing is described below. 4) 
If IndexMapping.DefaultType is registered, return it. 5) Return IndexMapping.DefaultMapping.

Each field or nested field of the value is identified by a string path, then mapped to one or 
several FieldMappings which extract the result for analysis.

Struct values fields are identified by their "json:" tag, or by their name. Nested fields are 
identified by prefixing with their parent identifier, separated by a dot.

Map values entries are identified by their string key. Entries not indexed by strings are ignored. 
Entry values are identified recursively like struct fields.

Slice and array values are identified by their field name. Their elements are processed 
sequentially with the same FieldMapping.

String, float64 and time.Time values are identified by their field name. Other types are ignored.

Each value identifier is decomposed in its parts and recursively address SubDocumentMappings in the 
tree starting at the root DocumentMapping. If a mapping is found, all its FieldMappings are applied 
to the value. If no mapping is found and the root DocumentMapping is dynamic, default mappings are 
used based on value type and IndexMapping default configurations.

Finally, mapped values are analyzed, indexed or stored. See FieldMapping.Analyzer to know how an 
analyzer is resolved for a given field.

Examples:

```
type Date struct {
  Day string \`json:"day"\`
  Month string
  Year string
}

type Person struct {
  FirstName string \`json:"first_name"\`
  LastName string
  BirthDate Date \`json:"birth_date"\`
}
```

A Person value FirstName is mapped by the SubDocumentMapping at "first\_name". Its LastName is 
mapped by the one at "LastName". The day of BirthDate is mapped to the SubDocumentMapping "day" of 
the root SubDocumentMapping "birth\_date". It will appear as the "birth\_date.day" field in the 
index. The month is mapped to "birth\_date.Month".

Example (Indexing) [¶](#example-Index-Indexing "Go to Example (Indexing)")

```
data := struct {
    Name    string
    Created time.Time
    Age     int
}{Name: "named one", Created: time.Now(), Age: 50}
data2 := struct {
    Name    string
    Created time.Time
    Age     int
}{Name: "great nameless one", Created: time.Now(), Age: 25}

// index some data
err = exampleIndex.Index("document id 1", data)
if err != nil {
    panic(err)
}
err = exampleIndex.Index("document id 2", data2)
if err != nil {
    panic(err)
}

// 2 documents have been indexed
count, err := exampleIndex.DocCount()
if err != nil {
    panic(err)
}

fmt.Println(count)
```
```
Output:
2
```

#### func New ¶

```
func New(path string, mapping mapping.IndexMapping) (Index, error)
```

New index at the specified path, must not exist. The provided mapping will be used for all 
Index/Search operations.

Example [¶](#example-New "Go to Example")

```
indexMapping = NewIndexMapping()
exampleIndex, err = New("path_to_index", indexMapping)
if err != nil {
    panic(err)
}
count, err := exampleIndex.DocCount()
if err != nil {
    panic(err)
}

fmt.Println(count)
```
```
Output:
0
```

#### func NewMemOnly ¶

```
func NewMemOnly(mapping mapping.IndexMapping) (Index, error)
```

NewMemOnly creates a memory-only index. The contents of the index is NOT persisted, and will be 
lost once closed. The provided mapping will be used for all Index/Search operations.

#### func NewUsing ¶

```
func NewUsing(path string, mapping mapping.IndexMapping, indexType string, kvstore string, kvconfig 
map[string]interface{}) (Index, error)
```

NewUsing creates index at the specified path, which must not already exist. The provided mapping 
will be used for all Index/Search operations. The specified index type will be used. The specified 
kvstore implementation will be used and the provided kvconfig will be passed to its constructor. 
Note that currently the values of kvconfig must be able to be marshaled and unmarshaled using the 
encoding/json library (used when reading/writing the index metadata file).

#### func Open ¶

```
func Open(path string) (Index, error)
```

Open index at the specified path, must exist. The mapping used when it was created will be used for 
all Index/Search operations.

#### func OpenUsing ¶

```
func OpenUsing(path string, runtimeConfig map[string]interface{}) (Index, error)
```

OpenUsing opens index at the specified path, must exist. The mapping used when it was created will 
be used for all Index/Search operations. The provided runtimeConfig can override settings persisted 
when the kvstore was created. If runtimeConfig has updated mapping, then an index update is 
attempted Throws an error without any changes to the index if an unupdatable mapping is provided

#### type IndexAlias ¶

```
type IndexAlias interface {
    Index

    Add(i ...Index)
    Remove(i ...Index)
    Swap(in, out []Index)
}
```

An IndexAlias is a wrapper around one or more Index objects. It has two distinct modes of 
operation. 1. When it points to a single index, ALL index operations are valid and will be passed 
through to the underlying index. 2. When it points to more than one index, the only valid operation 
is Search. In this case the search will be performed across all the underlying indexes and the 
results merged. Calls to Add/Remove/Swap the underlying indexes are atomic, so you can safely 
change the underlying Index objects while other components are performing operations.

#### added in v2.1.0

```
type IndexCopyable interface {
    // CopyTo creates a fully functional copy of the index at the
    // specified destination directory implementation.
    CopyTo(d index.Directory) error
}
```

IndexCopyable is an index which supports an online copy operation of the index.

#### type IndexErrMap ¶

```
type IndexErrMap map[string]error
```

IndexErrMap tracks errors with the name of the index where it occurred

#### func (IndexErrMap) MarshalJSON ¶

```
func (iem IndexErrMap) MarshalJSON() ([]byte, error)
```

MarshalJSON serializes the error into a string for JSON consumption

#### func (IndexErrMap) UnmarshalJSON ¶

```
func (iem IndexErrMap) UnmarshalJSON(data []byte) error
```

#### added in v2.6.0

```
type IndexFileCopyable interface {
    SetPathInBolt(key []byte, value []byte) error       //dest index
    CopyFile(file string, d index.IndexDirectory) error // source index
}
```

#### type IndexStat ¶

```
type IndexStat struct {
    // contains filtered or unexported fields
}
```

#### func (\*IndexStat) MarshalJSON ¶

```
func (is *IndexStat) MarshalJSON() ([]byte, error)
```

#### type IndexStats ¶

```
type IndexStats struct {
    // contains filtered or unexported fields
}
```

#### func NewIndexStats ¶

```
func NewIndexStats() *IndexStats
```

#### func (\*IndexStats) Register ¶

```
func (i *IndexStats) Register(index Index)
```

#### func (\*IndexStats) String ¶

```
func (i *IndexStats) String() string
```

#### func (\*IndexStats) UnRegister ¶

```
func (i *IndexStats) UnRegister(index Index)
```

#### added in v2.6.0

```
type IndexWithCallbacks interface {
    FileWriterIDsInUse() (map[string]struct{}, error)
    DropFileWriterIDs(ids map[string]struct{}) error
}
```

#### added in v2.5.5

```
type InsightsIndex interface {
    Index
    // TermFrequencies returns the tokens ordered by frequencies for the field index.
    TermFrequencies(field string, limit int, descending bool) ([]index.TermFreq, error)
    // CentroidCardinalities returns the centroids (clusters) from IVF indexes ordered by data 
density.
    CentroidCardinalities(field string, limit int, desceding bool) ([]index.CentroidCardinality, 
error)
}
```

#### added in v2.6.0

```
type OptionalRawMessage json.RawMessage
```

OptionalRawMessage is a wrapper around json.RawMessage that treats empty or \`null\` JSON as nil.

#### added in v2.6.0

```
func (n OptionalRawMessage) MarshalJSON() ([]byte, error)
```

#### added in v2.6.0

```
func (n *OptionalRawMessage) UnmarshalJSON(data []byte) error
```

#### added in v2.5.4

```
type RequestParams struct {
    ScoreRankConstant int \`json:"score_rank_constant,omitempty"\`
    ScoreWindowSize   int \`json:"score_window_size,omitempty"\`
}
```

Additional parameters in the search request. Currently only being used for score fusion parameters.

#### added in v2.5.4

```
func NewDefaultParams(from, size int) *RequestParams
```

#### added in v2.5.4

```
func ParseParams(r *SearchRequest, input []byte) (*RequestParams, error)
```

#### added in v2.5.4

```
func (p *RequestParams) UnmarshalJSON(input []byte) error
```

#### added in v2.5.4

```
func (p *RequestParams) Validate(size int) error
```

#### type SearchQueryEndCallbackFn ¶

```
type SearchQueryEndCallbackFn func(size uint64) error
```

#### type SearchQueryStartCallbackFn ¶

```
type SearchQueryStartCallbackFn func(size uint64) error
```

#### type SearchRequest ¶

```
type SearchRequest struct {
    ClientContextID  string            \`json:"client_context_id,omitempty"\`
    Query            query.Query       \`json:"query"\`
    Size             int               \`json:"size"\`
    From             int               \`json:"from"\`
    Highlight        *HighlightRequest \`json:"highlight,omitempty"\`
    Fields           []string          \`json:"fields,omitempty"\`
    Facets           FacetsRequest     \`json:"facets,omitempty"\`
    Explain          bool              \`json:"explain"\`
    Sort             search.SortOrder  \`json:"sort"\`
    IncludeLocations bool              \`json:"includeLocations"\`
    Score            string            \`json:"score,omitempty"\`
    SearchAfter      []string          \`json:"search_after,omitempty"\`
    SearchBefore     []string          \`json:"search_before,omitempty"\`

    PreSearchData map[string]interface{} \`json:"pre_search_data,omitempty"\`

    Params *RequestParams \`json:"params,omitempty"\`
    // contains filtered or unexported fields
}
```

A SearchRequest describes all the parameters needed to search the index. Query is required. 
Size/From describe how much and which part of the result set to return. Highlight describes 
optional search result highlighting. Fields describes a list of field values which should be 
retrieved for result documents, provided they were stored while indexing. Facets describe the set 
of facets to be computed. Explain triggers inclusion of additional search result score 
explanations. Sort describes the desired order for the results to be returned. Score controls the 
kind of scoring performed SearchAfter supports deep paging by providing a minimum sort key 
SearchBefore supports deep paging by providing a maximum sort key sortFunc specifies the sort 
implementation to use for sorting results.

A special field named "\*" can be used to return all fields.

#### func NewSearchRequest ¶

```
func NewSearchRequest(q query.Query) *SearchRequest
```

NewSearchRequest creates a new SearchRequest for the Query, using default values for all other 
search parameters.

Example [¶](#example-NewSearchRequest "Go to Example")

```
// finds documents with fields fully matching the given query text
query := NewMatchQuery("named one")
searchRequest := NewSearchRequest(query)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
```
```
Output:
document id 1
```

#### func NewSearchRequestOptions ¶

```
func NewSearchRequestOptions(q query.Query, size, from int, explain bool) *SearchRequest
```

NewSearchRequestOptions creates a new SearchRequest for the Query, with the requested size, from 
and explanation search parameters. By default results are ordered by score, descending.

#### func (\*SearchRequest) AddFacet ¶

```
func (r *SearchRequest) AddFacet(facetName string, f *FacetRequest)
```

AddFacet adds a FacetRequest to this SearchRequest

Example [¶](#example-SearchRequest.AddFacet "Go to Example")

```
facet := NewFacetRequest("Name", 1)
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchRequest.AddFacet("facet name", facet)
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

// total number of terms
fmt.Println(searchResults.Facets["facet name"].Total)
// number of docs with no value for this field
fmt.Println(searchResults.Facets["facet name"].Missing)
// term with highest occurrences in field name
fmt.Println(searchResults.Facets["facet name"].Terms.Terms()[0].Term)
```
```
Output:
5
2
one
```

#### added in v2.5.4

```
func (r *SearchRequest) AddParams(params RequestParams)
```

AddParams adds a RequestParams field to the search request

#### func (\*SearchRequest) SetSearchAfter ¶

```
func (r *SearchRequest) SetSearchAfter(after []string)
```

SetSearchAfter sets the request to skip over hits with a sort value less than the provided sort 
after key

#### func (\*SearchRequest) SetSearchBefore ¶

```
func (r *SearchRequest) SetSearchBefore(before []string)
```

SetSearchBefore sets the request to skip over hits with a sort value greater than the provided sort 
before key

#### func (\*SearchRequest) SetSortFunc ¶

```
func (r *SearchRequest) SetSortFunc(s func(sort.Interface))
```

SetSortFunc sets the sort implementation to use when sorting hits.

SearchRequests can specify a custom sort implementation to meet their needs. For instance, by 
specifying a parallel sort that uses all available cores.

#### func (\*SearchRequest) SortBy ¶

```
func (r *SearchRequest) SortBy(order []string)
```

SortBy changes the request to use the requested sort order this form uses the simplified syntax 
with an array of strings each string can either be a field name or the magic value \_id and \_score 
which refer to the doc id and search score any of these values can optionally be prefixed with - to 
reverse the order

Example [¶](#example-SearchRequest.SortBy "Go to Example")

```
// find docs containing "one", order by Age instead of score
query := NewMatchQuery("one")
searchRequest := NewSearchRequest(query)
searchRequest.SortBy([]string{"Age"})
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
fmt.Println(searchResults.Hits[1].ID)
```
```
Output:
document id 2
document id 1
```

#### func (\*SearchRequest) SortByCustom ¶

```
func (r *SearchRequest) SortByCustom(order search.SortOrder)
```

SortByCustom changes the request to use the requested sort order

Example [¶](#example-SearchRequest.SortByCustom "Go to Example")

```
// find all docs, order by Age, with docs missing Age field first
query := NewMatchAllQuery()
searchRequest := NewSearchRequest(query)
searchRequest.SortByCustom(search.SortOrder{
    &search.SortField{
        Field:   "Age",
        Missing: search.SortFieldMissingFirst,
    },
    &search.SortDocID{},
})
searchResults, err := exampleIndex.Search(searchRequest)
if err != nil {
    panic(err)
}

fmt.Println(searchResults.Hits[0].ID)
fmt.Println(searchResults.Hits[1].ID)
fmt.Println(searchResults.Hits[2].ID)
fmt.Println(searchResults.Hits[3].ID)
```
```
Output:
document id 3
document id 4
document id 2
document id 1
```

#### func (\*SearchRequest) SortFunc ¶

```
func (r *SearchRequest) SortFunc() func(data sort.Interface)
```

SortFunc returns the sort implementation to use when sorting hits. Defaults to sort.Sort.

#### func (\*SearchRequest) UnmarshalJSON ¶

```
func (r *SearchRequest) UnmarshalJSON(input []byte) error
```

UnmarshalJSON deserializes a JSON representation of a SearchRequest

#### func (\*SearchRequest) Validate ¶

```
func (r *SearchRequest) Validate() error
```

#### type SearchResult ¶

```
type SearchResult struct {
    Status   *SearchStatus                  \`json:"status"\`
    Request  *SearchRequest                 \`json:"request,omitempty"\`
    Hits     search.DocumentMatchCollection \`json:"hits"\`
    Total    uint64                         \`json:"total_hits"\`
    Cost     uint64                         \`json:"cost"\`
    MaxScore float64                        \`json:"max_score"\`
    Took     time.Duration                  \`json:"took"\`
    Facets   search.FacetResults            \`json:"facets"\`
    // special fields that are applicable only for search
    // results that are obtained from a presearch
    SynonymResult search.FieldTermSynonymMap \`json:"synonym_result,omitempty"\`

    // The following fields are applicable to BM25 preSearch
    BM25Stats *search.BM25Stats \`json:"bm25_stats,omitempty"\`
}
```

A SearchResult describes the results of executing a SearchRequest.

Status - Whether the search was executed on the underlying indexes successfully or failed, and the 
corresponding errors. Request - The SearchRequest that was executed. Hits - The list of documents 
that matched the query and their corresponding scores, score explanation, location info and so on. 
Total - The total number of documents that matched the query. Cost - indicates how expensive was 
the query with respect to bytes read from the mapped index files. MaxScore - The maximum score seen 
across all document hits seen for this query. Took - The time taken to execute the search. Facets - 
The facet results for the search.

#### func MultiSearch ¶

```
func MultiSearch(ctx context.Context, req *SearchRequest, params *multiSearchParams, indexes 
...Index) (*SearchResult, error)
```

MultiSearch executes a SearchRequest across multiple Index objects, then merges the results. The 
indexes must honor any ctx deadline.

#### func (\*SearchResult) Merge ¶

```
func (sr *SearchResult) Merge(other *SearchResult)
```

Merge will merge together multiple SearchResults during a MultiSearch

#### func (\*SearchResult) Size ¶

```
func (sr *SearchResult) Size() int
```

#### func (\*SearchResult) String ¶

```
func (sr *SearchResult) String() string
```

#### type SearchStatus ¶

```
type SearchStatus struct {
    Total      int         \`json:"total"\`
    Failed     int         \`json:"failed"\`
    Successful int         \`json:"successful"\`
    Errors     IndexErrMap \`json:"errors,omitempty"\`
}
```

SearchStatus is a section in the SearchResult reporting how many underlying indexes were queried, 
how many were successful/failed and a map of any errors that were encountered

#### func (\*SearchStatus) Merge ¶

```
func (ss *SearchStatus) Merge(other *SearchStatus)
```

Merge will merge together multiple SearchStatuses during a MultiSearch

#### added in v2.5.0

```
type SynonymDefinition struct {
    // Input is an optional list of terms for unidirectional synonym mapping.
    // When terms are specified in Input, they will map to the terms in Synonyms,
    // making the relationship unidirectional (each Input maps to all Synonyms).
    // If Input is omitted, the relationship is bidirectional among all Synonyms.
    Input []string \`json:"input,omitempty"\`

    // Synonyms is a list of terms that are considered equivalent.
    // If Input is specified, each term in Input will map to each term in Synonyms.
    // If Input is not specified, the Synonyms list will be treated bidirectionally,
    // meaning each term in Synonyms is treated as synonymous with all others.
    Synonyms []string \`json:"synonyms"\`
}
```

SynonymDefinition represents a synonym mapping in Bleve. Each instance associates one or more input 
terms with a list of synonyms, defining how terms are treated as equivalent in searches.

#### added in v2.5.0

```
func (sd *SynonymDefinition) Validate() error
```

#### added in v2.5.0

```
type SynonymIndex interface {
    Index
    // IndexSynonym indexes a synonym definition, with the specified id and belonging to the 
specified collection.
    IndexSynonym(id string, collection string, definition *SynonymDefinition) error
}
```

SynonymIndex supports indexing synonym definitions alongside regular documents. Synonyms, grouped 
by collection name, define term relationships for query expansion in searches.

#### added in v2.6.0

```
type TrainableIndex interface {
    Index
    Train(*Batch) error
}
```

## Directories

| Path | Synopsis |
| --- | --- |
| [analysis](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis) |  |
| 
[analyzer/custom](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/analyzer/custom
) |  |
| 
[analyzer/keyword](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/analyzer/keywo
rd) |  |
| 
[analyzer/simple](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/analyzer/simple
) |  |
| 
[analyzer/standard](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/analyzer/stan
dard) |  |
| [analyzer/web](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/analyzer/web) | 
 |
| 
[char/asciifolding](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/char/asciifol
ding) |  |
| [char/html](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/char/html) |  |
| [char/regexp](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/char/regexp) |  |
| 
[char/zerowidthnonjoiner](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/char/ze
rowidthnonjoiner) |  |
| 
[datetime/flexible](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/datetime/flex
ible) |  |
| [datetime/iso](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/datetime/iso) | 
 |
| 
[datetime/optional](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/datetime/opti
onal) |  |
| 
[datetime/percent](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/datetime/perce
nt) |  |
| 
[datetime/sanitized](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/datetime/san
itized) |  |
| 
[datetime/timestamp/microseconds](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis
/datetime/timestamp/microseconds) |  |
| 
[datetime/timestamp/milliseconds](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis
/datetime/timestamp/milliseconds) |  |
| 
[datetime/timestamp/nanoseconds](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/
datetime/timestamp/nanoseconds) |  |
| 
[datetime/timestamp/seconds](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/date
time/timestamp/seconds) |  |
| [lang/ar](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ar) |  |
| [lang/bg](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/bg) |  |
| [lang/ca](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ca) |  |
| [lang/cjk](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/cjk) |  |
| [lang/ckb](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ckb) |  |
| [lang/cs](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/cs) |  |
| [lang/da](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/da) |  |
| [lang/de](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/de) |  |
| [lang/el](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/el) |  |
| [lang/en](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/en)  Package en 
implements an analyzer with reasonable defaults for processing English text. | Package en 
implements an analyzer with reasonable defaults for processing English text. |
| [lang/es](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/es) |  |
| [lang/eu](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/eu) |  |
| [lang/fa](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/fa) |  |
| [lang/fi](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/fi) |  |
| [lang/fr](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/fr) |  |
| [lang/ga](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ga) |  |
| [lang/gl](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/gl) |  |
| [lang/hi](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/hi) |  |
| [lang/hr](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/hr) |  |
| [lang/hu](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/hu) |  |
| [lang/hy](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/hy) |  |
| [lang/id](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/id) |  |
| [lang/in](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/in) |  |
| [lang/it](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/it) |  |
| [lang/nl](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/nl) |  |
| [lang/no](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/no) |  |
| [lang/pl](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/pl) |  |
| 
[lang/pl/stempel](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/pl/stempel
) |  |
| 
[lang/pl/stempel/javadata](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/p
l/stempel/javadata) |  |
| [lang/pt](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/pt) |  |
| [lang/ro](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ro) |  |
| [lang/ru](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/ru) |  |
| [lang/sv](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/sv) |  |
| [lang/tr](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/lang/tr) |  |
| 
[token/apostrophe](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/apostrop
he) |  |
| 
[token/camelcase](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/camelcase
) |  |
| 
[token/compound](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/compound) 
|  |
| 
[token/edgengram](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/edgengram
) |  |
| [token/elision](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/elision) 
|  |
| 
[token/hierarchy](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/hierarchy
) |  |
| [token/keyword](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/keyword) 
|  |
| [token/length](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/length) | 
 |
| 
[token/lowercase](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/lowercase
)  Package lowercase implements a TokenFilter which converts tokens to lower case according to 
unicode rules. | Package lowercase implements a TokenFilter which converts tokens to lower case 
according to unicode rules. |
| [token/ngram](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/ngram) |  |
| [token/porter](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/porter) | 
 |
| [token/reverse](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/reverse) 
|  |
| [token/shingle](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/shingle) 
|  |
| 
[token/snowball](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/snowball) 
|  |
| [token/stop](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/stop)  
Package stop implements a TokenFilter removing tokens found in a TokenMap. | Package stop 
implements a TokenFilter removing tokens found in a TokenMap. |
| 
[token/truncate](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/truncate) 
|  |
| 
[token/unicodenorm](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/unicode
norm) |  |
| [token/unique](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/token/unique) | 
 |
| 
[tokenizer/character](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/c
haracter) |  |
| 
[tokenizer/exception](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/e
xception)  package exception implements a Tokenizer which extracts pieces matched by a regular 
expression from the input data, delegates the rest to another tokenizer, then insert back extracted 
parts in the token stream. | package exception implements a Tokenizer which extracts pieces matched 
by a regular expression from the input data, delegates the rest to another tokenizer, then insert 
back extracted parts in the token stream. |
| 
[tokenizer/letter](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/lett
er) |  |
| 
[tokenizer/regexp](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/rege
xp) |  |
| 
[tokenizer/single](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/sing
le) |  |
| 
[tokenizer/unicode](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/uni
code) |  |
| [tokenizer/web](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/web) 
|  |
| 
[tokenizer/whitespace](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenizer/
whitespace) |  |
| [tokenmap](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/analysis/tokenmap)  package 
token\_map implements a generic TokenMap, often used in conjunction with filters to remove or 
process specific tokens. | package token\_map implements a generic TokenMap, often used in 
conjunction with filters to remove or process specific tokens. |
| cmd |  |
| [bleve](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/cmd/bleve) command |  |
| [bleve/cmd](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/cmd/bleve/cmd) |  |
| 
[bleve/cmd/scorch](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/cmd/bleve/cmd/scorch) 
|  |
| [config](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/config) |  |
| [document](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/document) |  |
| [fusion](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/fusion) |  |
| [geo](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/geo) |  |
| index |  |
| [scorch](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/scorch) |  |
| 
[scorch/mergeplan](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/scorch/mergeplan)
  Package mergeplan provides a segment merge planning approach that's inspired by Lucene's 
TieredMergePolicy.java and descriptions like 
http://blog.mikemccandless.com/2011/02/visualizing-lucenes-segment-merges.html | Package mergeplan 
provides a segment merge planning approach that's inspired by Lucene's TieredMergePolicy.java and 
descriptions like http://blog.mikemccandless.com/2011/02/visualizing-lucenes-segment-merges.html |
| [upsidedown](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedown) |  |
| 
[upsidedown/store/boltdb](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedown
/store/boltdb)  Package boltdb implements a store.KVStore on top of BoltDB. | Package boltdb 
implements a store.KVStore on top of BoltDB. |
| 
[upsidedown/store/goleveldb](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsided
own/store/goleveldb) |  |
| 
[upsidedown/store/gtreap](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedown
/store/gtreap)  Package gtreap provides an in-memory implementation of the KVStore interfaces using 
the gtreap balanced-binary treap, copy-on-write data structure. | Package gtreap provides an 
in-memory implementation of the KVStore interfaces using the gtreap balanced-binary treap, 
copy-on-write data structure. |
| 
[upsidedown/store/metrics](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedow
n/store/metrics)  Package metrics provides a bleve.store.KVStore implementation that wraps another, 
real KVStore implementation, and uses go-metrics to track runtime performance metrics. | Package 
metrics provides a bleve.store.KVStore implementation that wraps another, real KVStore 
implementation, and uses go-metrics to track runtime performance metrics. |
| 
[upsidedown/store/moss](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedown/s
tore/moss) |  |
| 
[upsidedown/store/null](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/index/upsidedown/s
tore/null) |  |
| [mapping](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/mapping) |  |
| [numeric](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/numeric) |  |
| [registry](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/registry) |  |
| [search](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search) |  |
| [collector](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/collector) |  |
| [facet](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/facet) |  |
| [highlight](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highlight) |  |
| 
[highlight/format/ansi](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highlight/f
ormat/ansi) |  |
| 
[highlight/format/html](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highlight/f
ormat/html) |  |
| 
[highlight/format/plain](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highlight/
format/plain) |  |
| 
[highlight/fragmenter/simple](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highl
ight/fragmenter/simple) |  |
| 
[highlight/highlighter/ansi](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highli
ght/highlighter/ansi) |  |
| 
[highlight/highlighter/html](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/highli
ght/highlighter/html) |  |
| 
[highlight/highlighter/simple](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/high
light/highlighter/simple) |  |
| [query](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/query) |  |
| [scorer](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/scorer) |  |
| [searcher](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/search/searcher) |  |
| [size](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/size) |  |
| [test](https://pkg.go.dev/github.com/blevesearch/bleve/v2@v2.6.0/test) |  |
|  |  |
