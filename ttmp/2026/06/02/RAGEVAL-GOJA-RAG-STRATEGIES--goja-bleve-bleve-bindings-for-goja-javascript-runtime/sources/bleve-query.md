---
Title: Source: bleve-query.md
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

### Documentation

- User Guide
	- [Getting Started](https://blevesearch.com/docs/Getting%20Started/)
		- [Terminology](https://blevesearch.com/docs/Terminology/)
		- [Text Analysis](https://blevesearch.com/docs/Text-Analysis/)
		- [Introduction to Index Mappings](https://blevesearch.com/docs/Index-Mapping/)
		- Queries
		- [Query String Query](https://blevesearch.com/docs/Query-String-Query/)
		- [Sorting](https://blevesearch.com/docs/Sorting/)
		- [Index Aliases](https://blevesearch.com/docs/IndexAlias/)
		- [Example Applications](https://blevesearch.com/docs/Example-Applications/)
		- [Command-line Utility](https://blevesearch.com/docs/bleve/)
		- How do I...?
		- [Ignore/Disable Sections of 
Documents](https://blevesearch.com/docs/Disabling%20Sections%20of%20Documents/)
				- [Highlight Matches in 
Results](https://blevesearch.com/docs/Highlight%20Matches%20in%20Results/)
				- [Include Facets in 
Results](https://blevesearch.com/docs/Result-Faceting/)
				- [Pronounce Bleve](https://blevesearch.com/docs/Pronunciation/)
		- [Applications Using Bleve](https://blevesearch.com/docs/Applications-using-bleve/)
- Developer Guide
	- [Building](https://blevesearch.com/docs/Building/)
		- [Using blevex packages](https://blevesearch.com/docs/Blevex/)
		- [Contributing](https://blevesearch.com/docs/Contributing/)
		- [Package Structure](https://blevesearch.com/docs/Package-Structure/)

## Queries

Queries supported by bleve:

### Term

A term query is the simplest possible query. It performs an exact match in the index for the 
provided term.

Most of the time users should use a Match Query instead.

### Match

A match query is like a term query, but the input text is analyzed first. An attempt is made to use 
the same analyzer that was used when the field was indexed.

The match query can optionally perform fuzzy matching. If the fuzziness parameter is set to a 
non-zero integer the analyzed text will be matched with the specified level of fuzziness. Also, the 
prefix\_length parameter can be used to require that the term also have the same prefix of the 
specified length.

### Phrase

A phrase query searches for terms occurring in the specified position and offsets.

The phrase query is performing an exact term match for all the phrase constituents. If you want the 
phrase to be analyzed, consider using the Match Phrase Query instead.

### Match Phrase

The match phrase query is like the phrase query, but the input text is analyzed and a phrase query 
is built with the terms resulting from the analysis.

### Prefix

The prefix query finds documents containing terms that start with the provided prefix.

### Fuzzy

A fuzzy query is a term query that matches terms within a specified edit distance (Levenshtein 
distance). Also, you can optionally specify that the term must have a matching prefix of the 
specified length.

### Conjunction

The conjunction query is a compound query. Result documents must satisfy all of the child queries.

### Disjunction

The disjunction query is a compound query. Result documents must satisfy a configurable `min` 
number of child queries. By default this `min` is set to 1.

### Boolean

The boolean query is useful combination of conjunction and disjunction queries. The query takes 
three lists of queries:

- must - result documents must satisfy all of these queries
- should - result documents should satisfy at least `minShould` of these queries
- must not - result documents must not satisfy any of these queries

The `minShould` value is configurable, defaults to 0.

### Numeric Range

The numeric range query finds documents containing a numeric value in the specified field within 
the specified range. You can omit one endpoint, but not both. The `inclusiveMin` and `inclusiveMax` 
properties control whether or not the end points are included or excluded.

### Date Range

The date range query finds documents containing a date value in the specified field within the 
specified range. You can omit one endpoint, but not both. The inclusiveStart and inclusiveEnd 
properties control whether or not the end points are included or excluded.

### Query String

The query language query allows humans to describe complex queries using a simple syntax. There is 
a separate page with the [full query language 
specification](https://blevesearch.com/docs/Query-String-Query/)

### Match All

The match all query will match all documents in the index.

### Match None

The match none query will not match any documents in the index.

### Doc ID Query

The doc ID query will match only documents that contain one of the supplied document identifiers.
