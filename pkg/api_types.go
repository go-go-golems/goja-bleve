package pkg

import (
	"fmt"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/dop251/goja"
)

type refKind string

const (
	refKindIndex         refKind = "index"
	refKindIndexBuilder  refKind = "indexBuilder"
	refKindMapping       refKind = "mapping"
	refKindMappingBuild  refKind = "mappingBuilder"
	refKindDocMapping    refKind = "docMapping"
	refKindDocMapBuilder refKind = "docMappingBuilder"
	refKindFieldMapping  refKind = "fieldMapping"
	refKindFieldBuilder  refKind = "fieldBuilder"
	refKindQuery         refKind = "query"
	refKindSearchRequest refKind = "searchRequest"
	refKindBatch         refKind = "batch"
	refKindKNN           refKind = "knn"
)

type refBase struct {
	api    *moduleRuntime
	kind   refKind
	closed bool
}

type indexRef struct {
	refBase
	name  string
	path  string
	index bleve.Index
}

type indexBuilderRef struct {
	refBase
	mode string
	path string
}

type mappingRef struct {
	refBase
	mapping any
}

type mappingBuilderRef struct {
	refBase
}

type docMappingRef struct {
	refBase
	mapping any
}

type docMappingBuilderRef struct {
	refBase
}

type fieldMappingRef struct {
	refBase
	mapping any
}

type fieldBuilderRef struct {
	refBase
	fieldType string
}

type queryRef struct {
	refBase
	query any
}

type searchRequestRef struct {
	refBase
	request *bleve.SearchRequest
}

type batchRef struct {
	refBase
	index    *indexRef
	executed bool
}

type knnRef struct {
	refBase
	field  string
	vector []float32
	k      int
	boost  float64
}

func (m *moduleRuntime) attachRef(o *goja.Object, ref any) {
	_ = o.Set(hiddenRefKey, ref)
	_ = o.DefineDataProperty(
		hiddenRefKey,
		o.Get(hiddenRefKey),
		goja.FLAG_FALSE, // writable
		goja.FLAG_FALSE, // enumerable
		goja.FLAG_FALSE, // configurable
	)
}

func (m *moduleRuntime) getRef(v goja.Value) any {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil
	}
	obj, ok := v.(*goja.Object)
	if !ok {
		return nil
	}
	raw := obj.Get(hiddenRefKey)
	if raw == nil || goja.IsUndefined(raw) || goja.IsNull(raw) {
		return nil
	}
	return raw.Export()
}

func getTypedRef[T any](m *moduleRuntime, v goja.Value, expected string) (*T, error) {
	ref := m.getRef(v)
	if ref == nil {
		return nil, fmt.Errorf("bleve: expected %s wrapper, got value without Go reference", expected)
	}
	typed, ok := ref.(*T)
	if !ok {
		return nil, fmt.Errorf("bleve: expected %s wrapper, got %T", expected, ref)
	}
	return typed, nil
}

func (m *moduleRuntime) newWrapper(ref any, kind refKind) *goja.Object {
	obj := m.vm.NewObject()
	m.mustSet(obj, "type", string(kind))
	m.attachRef(obj, ref)
	return obj
}

func (r refBase) assertOpen(label string) error {
	if r.closed {
		return fmt.Errorf("bleve: %s is closed", label)
	}
	return nil
}
