package pkg

import (
	"fmt"
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/dop251/goja"
)

func (m *moduleRuntime) mappingBuilder() *goja.Object {
	ref := &mappingBuilderRef{
		refBase: refBase{api: m, kind: refKindMappingBuild},
		mapping: bleve.NewIndexMapping(),
	}
	obj := m.newWrapper(ref, refKindMappingBuild)
	m.mustSet(obj, "defaultMapping", func(dm goja.Value) (*goja.Object, error) {
		docRef, err := getTypedRef[docMappingRef](m, dm, "document mapping")
		if err != nil {
			return nil, err
		}
		ref.mapping.DefaultMapping = docRef.mapping
		return obj, nil
	})
	m.mustSet(obj, "addTypeMapping", func(name string, dm goja.Value) (*goja.Object, error) {
		docRef, err := getTypedRef[docMappingRef](m, dm, "document mapping")
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: type mapping name is required")
		}
		ref.mapping.AddDocumentMapping(name, docRef.mapping)
		return obj, nil
	})
	m.mustSet(obj, "typeMapping", obj.Get("addTypeMapping"))
	m.mustSet(obj, "typeField", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: typeField name is required")
		}
		ref.mapping.TypeField = name
		return obj, nil
	})
	m.mustSet(obj, "defaultAnalyzer", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: defaultAnalyzer name is required")
		}
		ref.mapping.DefaultAnalyzer = name
		return obj, nil
	})
	m.mustSet(obj, "defaultField", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: defaultField name is required")
		}
		ref.mapping.DefaultField = name
		return obj, nil
	})
	m.mustSet(obj, "storeDynamic", func(enabled bool) *goja.Object {
		ref.mapping.StoreDynamic = enabled
		return obj
	})
	m.mustSet(obj, "indexDynamic", func(enabled bool) *goja.Object {
		ref.mapping.IndexDynamic = enabled
		return obj
	})
	m.mustSet(obj, "docValuesDynamic", func(enabled bool) *goja.Object {
		ref.mapping.DocValuesDynamic = enabled
		return obj
	})
	m.mustSet(obj, "build", func() (*goja.Object, error) {
		if err := ref.mapping.Validate(); err != nil {
			return nil, fmt.Errorf("bleve: invalid index mapping: %w", err)
		}
		built := &mappingRef{refBase: refBase{api: m, kind: refKindMapping}, mapping: ref.mapping}
		return m.newWrapper(built, refKindMapping), nil
	})
	return obj
}

func (m *moduleRuntime) docMappingBuilder() *goja.Object {
	ref := &docMappingBuilderRef{
		refBase: refBase{api: m, kind: refKindDocMapBuilder},
		mapping: bleve.NewDocumentMapping(),
	}
	obj := m.newWrapper(ref, refKindDocMapBuilder)
	m.mustSet(obj, "field", func(name string, field goja.Value) (*goja.Object, error) {
		fieldRef, err := getTypedRef[fieldMappingRef](m, field, "field mapping")
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: field name is required")
		}
		ref.mapping.AddFieldMappingsAt(name, fieldRef.mapping)
		return obj, nil
	})
	m.mustSet(obj, "addField", obj.Get("field"))
	m.mustSet(obj, "subDocument", func(name string, doc goja.Value) (*goja.Object, error) {
		docRef, err := getTypedRef[docMappingRef](m, doc, "document mapping")
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: sub-document name is required")
		}
		ref.mapping.AddSubDocumentMapping(name, docRef.mapping)
		return obj, nil
	})
	m.mustSet(obj, "addSubDoc", obj.Get("subDocument"))
	m.mustSet(obj, "dynamic", func(enabled bool) *goja.Object {
		ref.mapping.Dynamic = enabled
		return obj
	})
	m.mustSet(obj, "enabled", func(enabled bool) *goja.Object {
		ref.mapping.Enabled = enabled
		return obj
	})
	m.mustSet(obj, "nested", func(enabled bool) *goja.Object {
		ref.mapping.Nested = enabled
		return obj
	})
	m.mustSet(obj, "defaultAnalyzer", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: defaultAnalyzer name is required")
		}
		ref.mapping.DefaultAnalyzer = name
		return obj, nil
	})
	m.mustSet(obj, "build", func() *goja.Object {
		built := &docMappingRef{refBase: refBase{api: m, kind: refKindDocMapping}, mapping: ref.mapping}
		return m.newWrapper(built, refKindDocMapping)
	})
	return obj
}

func (m *moduleRuntime) fieldBuilder() *goja.Object {
	ref := &fieldBuilderRef{refBase: refBase{api: m, kind: refKindFieldBuilder}}
	obj := m.newWrapper(ref, refKindFieldBuilder)
	m.mustSet(obj, "text", func() *goja.Object {
		ref.mapping = bleve.NewTextFieldMapping()
		return obj
	})
	m.mustSet(obj, "keyword", func() *goja.Object {
		ref.mapping = bleve.NewKeywordFieldMapping()
		return obj
	})
	m.mustSet(obj, "number", func() *goja.Object {
		ref.mapping = bleve.NewNumericFieldMapping()
		return obj
	})
	m.mustSet(obj, "datetime", func() *goja.Object {
		ref.mapping = bleve.NewDateTimeFieldMapping()
		return obj
	})
	m.mustSet(obj, "boolean", func() *goja.Object {
		ref.mapping = bleve.NewBooleanFieldMapping()
		return obj
	})
	m.mustSet(obj, "geoPoint", func() *goja.Object {
		ref.mapping = bleve.NewGeoPointFieldMapping()
		return obj
	})
	m.mustSet(obj, "geoShape", func() *goja.Object {
		ref.mapping = bleve.NewGeoShapeFieldMapping()
		return obj
	})
	m.mustSet(obj, "ip", func() *goja.Object {
		ref.mapping = bleve.NewIPFieldMapping()
		return obj
	})
	m.mustSet(obj, "disabled", func() *goja.Object {
		ref.mapping = &mapping.FieldMapping{Index: false, Store: false, DocValues: false}
		return obj
	})
	m.mustSet(obj, "vector", func(dims int) (*goja.Object, error) {
		field, err := newVectorFieldMapping(dims, false)
		if err != nil {
			return nil, err
		}
		ref.mapping = field
		return obj, nil
	})
	m.mustSet(obj, "vectorBase64", func(dims int) (*goja.Object, error) {
		field, err := newVectorFieldMapping(dims, true)
		if err != nil {
			return nil, err
		}
		ref.mapping = field
		return obj, nil
	})
	m.installFieldOptions(obj, ref)
	m.mustSet(obj, "build", func() (*goja.Object, error) {
		if ref.mapping == nil {
			return nil, fmt.Errorf("bleve: field type is required before build()")
		}
		built := &fieldMappingRef{refBase: refBase{api: m, kind: refKindFieldMapping}, mapping: ref.mapping}
		return m.newWrapper(built, refKindFieldMapping), nil
	})
	return obj
}

func (m *moduleRuntime) installFieldOptions(obj *goja.Object, ref *fieldBuilderRef) {
	ensure := func() *mapping.FieldMapping {
		if ref.mapping == nil {
			ref.mapping = bleve.NewTextFieldMapping()
		}
		return ref.mapping
	}
	m.mustSet(obj, "name", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: field mapping name is required")
		}
		ensure().Name = name
		return obj, nil
	})
	m.mustSet(obj, "analyzer", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: analyzer name is required")
		}
		ensure().Analyzer = name
		return obj, nil
	})
	m.mustSet(obj, "store", func(enabled bool) *goja.Object {
		ensure().Store = enabled
		return obj
	})
	m.mustSet(obj, "index", func(enabled bool) *goja.Object {
		ensure().Index = enabled
		return obj
	})
	m.mustSet(obj, "docValues", func(enabled bool) *goja.Object {
		ensure().DocValues = enabled
		return obj
	})
	m.mustSet(obj, "includeTermVectors", func(enabled bool) *goja.Object {
		ensure().IncludeTermVectors = enabled
		return obj
	})
	m.mustSet(obj, "includeInAll", func(enabled bool) *goja.Object {
		ensure().IncludeInAll = enabled
		return obj
	})
	m.mustSet(obj, "dateFormat", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: date format name is required")
		}
		ensure().DateFormat = name
		return obj, nil
	})
	m.mustSet(obj, "similarity", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: vector similarity is required")
		}
		ensure().Similarity = normalizeVectorSimilarity(name)
		return obj, nil
	})
	m.mustSet(obj, "optimizedFor", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: vector optimization target is required")
		}
		ensure().VectorIndexOptimizedFor = normalizeVectorOptimization(name)
		return obj, nil
	})
}

func normalizeVectorSimilarity(name string) string {
	switch strings.TrimSpace(strings.ToLower(name)) {
	case "cosine":
		return "cosine"
	case "dot", "dot_product", "dot-product":
		return "dot_product"
	case "l2", "l2_norm", "l2-norm", "euclidean":
		return "l2_norm"
	default:
		return name
	}
}

func normalizeVectorOptimization(name string) string {
	switch strings.TrimSpace(strings.ToLower(name)) {
	case "recall", "flat", "ivf":
		return "recall"
	case "latency":
		return "latency"
	case "memory", "memory-efficient", "memory_efficient":
		return "memory-efficient"
	default:
		return name
	}
}
