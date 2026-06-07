---
Title: Source: bleve-index-mapping.md
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
		- Introduction to Index Mappings
		- [Queries](https://blevesearch.com/docs/Query/)
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

## Introduction to Index Mappings

In bleve, the IndexMapping describes how your data model should be indexed.

## Default IndexMapping

To get the default IndexMapping, simply call:

```
indexMapping := bleve.NewIndexMapping()
```

IndexMappings contain DocumentMappings for each of the different types of documents you want to 
support. Further, it contains a DefaultDocumentMapping that will be used for any type which does 
not have an explicit mapping.

## Document Type

How does bleve know what type a document is?

1. If your object implements the interface `bleve.Classifier` then bleve will use string returned 
by its `Type()` method.
2. The IndexMapping has a setting called `TypeField`. You can set this to any document path, and if 
the value at that path is a string, that value will be used as the typed field. If you did not 
customize this setting the default is set to “\_type”.
3. If no type can be determined from 1 or 2, the type is set to the IndexMapping’s `DefaultType`. 
If you did not customize this setting the default is set to “\_default”.

## DocumentMappings

Now that we see how bleve will determine the type, we can provide a customized DocumentMapping for 
each type we’re interested in.

Let’s say we have a document type called `blog`. We can build a DocumentMapping for this type and 
configure the IndexMapping to use it:

```
blogMapping := bleve.NewDocumentMapping()
indexMapping.AddDocumentMapping("blog", blogMapping)
```

We can also set a catch-all mapping that will be used for any type that does not have an explicit 
mapping by setting the DefaultMapping field.

## FieldMappings

Documents are hierarchical and contain named fields. These fields could be values or nested 
sub-documents. We customize the behavior for a named field by setting a DocumentMapping for it. 
Once we have a DocumentMapping for the named field, we can attach 0 or more FieldMappings to it. 
The FieldMappings describe how we want the field to be interpreted and what we want inserted into 
the index.

Let’s say our blog documents have a string field `name` and we want to use the English analyzer 
for this field.

```
nameFieldMapping := bleve.NewTextFieldMapping()
nameFieldMapping.Analyzer = "en"
blogMapping.AddFieldMappingsAt("name", nameFieldMapping)
```

Now let’s say our blog documents have a nested structure describing the `author` field using 
sub-fields `name` and `email`. This time lets say we want index (default) but not store the author 
name. And we want to exclude the email address from the `_all` field.

```
author := bleve.NewDocumentMapping()
authorNameFieldMapping := bleve.NewTextFieldMapping()
authorNameFieldMapping.Store = false
author.AddFieldMappingsAt("name", authorFieldNameMapping)
authorEmailFieldMapping := bleve.NewTextFieldMapping()
authorEmailFieldMapping.IncludeInAll = false
author.AddFieldMappingsAt("email", authorEmailFieldMapping)
blog.AddSubDocumentMapping("author", author)
```

This shows use of a few of the other flags in a FieldMapping. Here is the list:

- Index - index this field, defaults to true
- Store - store this field, defaults to true
- IncludeTermVectors - include term vectors for this field, defaults to true
- IncludeInAll - includes this field in the composite field named `_all`, defaults to true

## Text Field specific options

- Analyzer - the named analyzer to use on this field

There are multiple levels at which you can configure the Default Analyzer if an explicit one 
isn’t specified.

1. Each DocumentMapping has a field `DefaultAnalyzer`. This means you can override the default 
analyzer for each sub-document.
2. The IndexMapping also has a `DefaultAnalyzer`.

The `DefaultAnalyzer` configured with the longest path match to a field will be used.

## Date Field specific options

- DateFormat - the name of a DateTimeParser that will be used to parse a date stored as a string

You can configure a DefaultDateTimeParser in the IndexMapping object.

## Understanding Default Type vs Default Mapping

When Bleve cannot figure out which type a particular document is, it is automatically assigned the 
DefaultType.

Once, Bleve has determined the type, it looks for a DocumentMapping matching this type name. If 
there is no explicitly configured DocumentMapping for this type, then the DefaultMapping is used.

The DefaultType will default to “\_default”, and the DefaultMapping will default to an empty 
default DocumentMapping.

Consider an example from the beer-search sample application. The mapping describes two types 
“beer” and “brewery”. For each of these an explicit DocumentMapping is provided. If you 
attempt to index a document where the type field is missing, it will be assigned the type 
“\_default”. Then Bleve looks to see if there is a mapping configured for “\_default”. 
There is not, so then Bleve proceeds to use the DefaultMapping.
