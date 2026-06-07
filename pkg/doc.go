// Package pkg implements the goja-bleve native module exposed to JavaScript as
// require("bleve").
//
// The current implementation provides the Phase 0/1 module scaffold: native
// module registration, top-level factories, runtime state, build-tag vector
// detection, and Go-backed JavaScript wrapper references. Later phases add
// concrete Bleve mapping, indexing, query, KNN, and hybrid search behavior.
package pkg
