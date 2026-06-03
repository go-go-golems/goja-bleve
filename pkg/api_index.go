package pkg

import (
	"fmt"
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/dop251/goja"
)

func (m *moduleRuntime) createIndexBuilder(path string) *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "create", path: path}
	return m.installIndexBuilderMethods(ref)
}

func (m *moduleRuntime) openIndexBuilder(path string) *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "open", path: path}
	return m.installIndexBuilderMethods(ref)
}

func (m *moduleRuntime) memoryIndexBuilder() *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "memory"}
	return m.installIndexBuilderMethods(ref)
}

func (m *moduleRuntime) installIndexBuilderMethods(ref *indexBuilderRef) *goja.Object {
	obj := m.newWrapper(ref, refKindIndexBuilder)
	m.mustSet(obj, "mapping", func(value goja.Value) (*goja.Object, error) {
		mappingRef, err := getTypedRef[mappingRef](m, value, "index mapping")
		if err != nil {
			return nil, err
		}
		ref.mapping = mappingRef.mapping
		return obj, nil
	})
	m.mustSet(obj, "build", func() (*goja.Object, error) {
		return m.buildIndex(ref)
	})
	return obj
}

func (m *moduleRuntime) buildIndex(builder *indexBuilderRef) (*goja.Object, error) {
	mapping := builder.mapping
	if mapping == nil {
		mapping = bleve.NewIndexMapping()
	}

	var (
		idx bleve.Index
		err error
	)
	switch builder.mode {
	case "memory":
		idx, err = bleve.NewMemOnly(mapping)
	case "create":
		if strings.TrimSpace(builder.path) == "" {
			return nil, fmt.Errorf("bleve: create index path is required")
		}
		idx, err = bleve.New(builder.path, mapping)
	case "open":
		if strings.TrimSpace(builder.path) == "" {
			return nil, fmt.Errorf("bleve: open index path is required")
		}
		idx, err = bleve.Open(builder.path)
	default:
		return nil, fmt.Errorf("bleve: unknown index builder mode %q", builder.mode)
	}
	if err != nil {
		return nil, fmt.Errorf("bleve: %s index: %w", builder.mode, err)
	}

	name := builder.path
	if name == "" {
		name = "memory"
	}
	ref := &indexRef{refBase: refBase{api: m, kind: refKindIndex}, name: name, path: builder.path, index: idx}
	m.openIndexes[name] = ref
	return m.indexObject(ref), nil
}

func (m *moduleRuntime) indexObject(ref *indexRef) *goja.Object {
	obj := m.newWrapper(ref, refKindIndex)
	m.mustSet(obj, "name", func() string { return ref.name })
	m.mustSet(obj, "docCount", func() (uint64, error) {
		if err := ref.assertOpen("index"); err != nil {
			return 0, err
		}
		return ref.index.DocCount()
	})
	m.mustSet(obj, "index", func(id string, doc goja.Value) error {
		if err := ref.assertOpen("index"); err != nil {
			return err
		}
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("bleve: document id is required")
		}
		return ref.index.Index(id, doc.Export())
	})
	m.mustSet(obj, "delete", func(id string) error {
		if err := ref.assertOpen("index"); err != nil {
			return err
		}
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("bleve: document id is required")
		}
		return ref.index.Delete(id)
	})
	m.mustSet(obj, "newBatch", func() (*goja.Object, error) {
		if err := ref.assertOpen("index"); err != nil {
			return nil, err
		}
		return m.batchObject(&batchRef{
			refBase: refBase{api: m, kind: refKindBatch},
			index:   ref,
			batch:   ref.index.NewBatch(),
		}), nil
	})
	m.mustSet(obj, "batch", obj.Get("newBatch"))
	m.mustSet(obj, "search", func(requestValue goja.Value) (map[string]any, error) {
		if err := ref.assertOpen("index"); err != nil {
			return nil, err
		}
		requestRef, err := getTypedRef[searchRequestRef](m, requestValue, "search request")
		if err != nil {
			return nil, err
		}
		if requestRef.request == nil {
			return nil, fmt.Errorf("bleve: search request is not built")
		}
		result, err := ref.index.Search(requestRef.request)
		if err != nil {
			return nil, err
		}
		return searchResultToJS(result), nil
	})
	m.mustSet(obj, "close", func() error {
		if ref.closed {
			return nil
		}
		ref.closed = true
		delete(m.openIndexes, ref.name)
		if ref.index == nil {
			return nil
		}
		return ref.index.Close()
	})
	return obj
}
