package pkg

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
)

const (
	// ModuleName is the JavaScript require() name for this native module.
	ModuleName = "bleve"
	// hiddenRefKey stores Go references on JavaScript wrapper objects.
	hiddenRefKey = "__bleve_ref"
)

type module struct{}

var _ modules.NativeModule = (*module)(nil)

// NewLoader returns the native-module loader for direct registration with a
// goja_nodejs require.Registry.
func NewLoader() require.ModuleLoader {
	mod := &module{}
	return mod.Loader
}

// Register registers the bleve native module on a goja_nodejs require registry.
func Register(reg *require.Registry) {
	if reg == nil {
		return
	}
	reg.RegisterNativeModule(ModuleName, NewLoader())
}

func (module) Name() string { return ModuleName }

func (module) Doc() string {
	return `
The bleve module exposes Bleve full-text and vector search to goja.

Current scaffold exports Go-backed builder/reference objects. Phase 0/1 focus
on module loading and hidden Go references. Later phases implement mapping,
index lifecycle, query execution, KNN search, and hybrid score fusion.

Vector/KNN APIs require the host Go binary to be compiled with -tags=vectors and
linked against FAISS. See the FAISS how-to in the RAG evaluation system docs.
`
}

func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	rt := newRuntime(vm)
	exports := moduleObj.Get("exports").(*goja.Object)
	rt.installExports(exports)
}

func init() {
	modules.Register(&module{})
}

type moduleRuntime struct {
	vm            *goja.Runtime
	openIndexes   map[string]*indexRef
	vectorSupport bool
}

func newRuntime(vm *goja.Runtime) *moduleRuntime {
	return &moduleRuntime{
		vm:            vm,
		openIndexes:   map[string]*indexRef{},
		vectorSupport: vectorSupportEnabled,
	}
}

func (m *moduleRuntime) closeAll() error {
	var firstErr error
	for key, ref := range m.openIndexes {
		if ref == nil || ref.closed || ref.index == nil {
			delete(m.openIndexes, key)
			continue
		}
		if err := ref.index.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close index %s: %w", key, err)
		}
		ref.closed = true
		delete(m.openIndexes, key)
	}
	return firstErr
}

func (m *moduleRuntime) installExports(exports *goja.Object) {
	m.mustSet(exports, "version", "0.1.0")
	m.mustSet(exports, "vectorSupport", m.vectorSupport)

	m.mustSet(exports, "mapping", m.mappingBuilder)
	m.mustSet(exports, "indexMapping", m.mappingBuilder)
	m.mustSet(exports, "docMapping", m.docMappingBuilder)
	m.mustSet(exports, "documentMapping", m.docMappingBuilder)
	m.mustSet(exports, "field", m.fieldBuilder)
	m.mustSet(exports, "search", m.searchRequestBuilder)
	m.mustSet(exports, "searchRequest", m.searchRequestBuilder)

	m.mustSet(exports, "create", m.createIndexBuilder)
	m.mustSet(exports, "open", m.openIndexBuilder)
	m.mustSet(exports, "memory", m.memoryIndexBuilder)

	m.mustSet(exports, "match", m.matchQuery)
	m.mustSet(exports, "matchPhrase", m.matchPhraseQuery)
	m.mustSet(exports, "term", m.termQuery)
	m.mustSet(exports, "queryString", m.queryStringQuery)
	m.mustSet(exports, "prefix", m.prefixQuery)
	m.mustSet(exports, "fuzzy", m.fuzzyQuery)
	m.mustSet(exports, "regexp", m.regexpQuery)
	m.mustSet(exports, "wildcard", m.wildcardQuery)
	m.mustSet(exports, "bool", m.boolQuery)
	m.mustSet(exports, "conj", m.conjunctionQuery)
	m.mustSet(exports, "conjunction", m.conjunctionQuery)
	m.mustSet(exports, "disj", m.disjunctionQuery)
	m.mustSet(exports, "disjunction", m.disjunctionQuery)
	m.mustSet(exports, "matchAll", m.matchAllQuery)
	m.mustSet(exports, "matchNone", m.matchNoneQuery)
}

func (m *moduleRuntime) mustSet(o *goja.Object, key string, value any) {
	if err := o.Set(key, value); err != nil {
		panic(m.vm.NewGoError(fmt.Errorf("bleve: set export %s: %w", key, err)))
	}
}
