package pkg

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
)

const (
	// ModuleName is the JavaScript require() name for this native module.
	ModuleName = "bleve"
	// hiddenRefKey stores Go references on JavaScript wrapper objects.
	hiddenRefKey = "__bleve_ref"
)

type module struct{}

var _ modules.NativeModule = (*module)(nil)
var _ modules.TypeScriptDeclarer = (*module)(nil)

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

func (module) TypeScriptModule() *spec.Module {
	return &spec.Module{
		Name:        ModuleName,
		Description: "Bleve full-text, vector, and hybrid search builders for goja.",
		RawDTS: []string{
			"export type Vector = number[] | Float32Array | Float64Array;",
			"export type ScoreMode = 'default' | 'none' | 'rrf' | 'rsf';",
			"export type KNNOperator = 'or' | 'and';",
			"export type VectorSimilarity = 'cosine' | 'dot_product' | 'dot-product' | 'dot' | 'l2_norm' | 'l2-norm' | 'l2' | 'euclidean' | string;",
			"export type VectorOptimization = 'recall' | 'latency' | 'memory-efficient' | 'memory_efficient' | 'memory' | string;",
			"export interface Mapping { readonly type: 'mapping'; }",
			"export interface DocumentMapping { readonly type: 'documentMapping'; }",
			"export interface FieldMapping { readonly type: 'fieldMapping'; }",
			"export interface Query { readonly type: 'query'; field(name: string): this; boost(value: number): this; }",
			"export interface SearchRequest { readonly type: 'searchRequest'; }",
			"export type FieldMappingLike = FieldMapping | FieldBuilder;",
			"export interface MappingBuilder { defaultMapping(mapping: DocumentMapping): this; addTypeMapping(name: string, mapping: DocumentMapping): this; typeMapping(name: string, mapping: DocumentMapping): this; typeField(name: string): this; defaultAnalyzer(name: string): this; defaultField(name: string): this; storeDynamic(enabled: boolean): this; indexDynamic(enabled: boolean): this; docValuesDynamic(enabled: boolean): this; build(): Mapping; }",
			"export interface DocumentMappingBuilder { field(name: string, field: FieldMappingLike): this; addField(name: string, field: FieldMappingLike): this; subDocument(name: string, doc: DocumentMapping): this; addSubDoc(name: string, doc: DocumentMapping): this; dynamic(enabled: boolean): this; enabled(enabled: boolean): this; nested(enabled: boolean): this; defaultAnalyzer(name: string): this; build(): DocumentMapping; }",
			"export interface FieldBuilder { text(): this; keyword(): this; number(): this; datetime(): this; boolean(): this; geoPoint(): this; geoShape(): this; ip(): this; disabled(): this; vector(dims: number): this; vectorBase64(dims: number): this; name(name: string): this; analyzer(name: string): this; store(enabled: boolean): this; index(enabled: boolean): this; docValues(enabled: boolean): this; includeTermVectors(enabled: boolean): this; includeInAll(enabled: boolean): this; dateFormat(name: string): this; similarity(name: VectorSimilarity): this; optimizedFor(name: VectorOptimization): this; build(): FieldMapping; }",
			"export interface BoolQuery extends Query { addMust(...queries: Query[]): this; addShould(...queries: Query[]): this; addMustNot(...queries: Query[]): this; }",
			"export interface SearchRequestBuilder { query(query: Query): this; size(n: number): this; from(n: number): this; fields(names: string[]): this; sort(names: string[]): this; highlight(fields?: string[] | string, style?: string): this; explain(enabled: boolean): this; score(mode: ScoreMode): this; scoreRankConstant(n: number): this; scoreWindowSize(n: number): this; knnOperator(operator: KNNOperator): this; knn(field: string, vector: Vector, k: number, boost?: number): this; build(): SearchRequest; }",
			"export interface IndexBuilder { mapping(mapping: Mapping): this; name(name: string): this; build(): Index; }",
			"export interface Index { index(id: string, doc: unknown): void; delete(id: string): void; search(request: SearchRequest): SearchResult; docCount(): number; newBatch(): Batch; batch(): Batch; close(): void; }",
			"export interface Batch { index(id: string, doc: unknown): this; delete(id: string): this; size(): number; operationCount(): number; reset(): this; execute(): void; }",
			"export interface SearchResult { total: number; maxScore: number; took: string; hits: SearchHit[]; }",
			"export interface SearchHit { id: string; score: number; fields: Record<string, unknown>; fragments?: unknown; locations?: unknown; sort?: unknown[]; explanation?: unknown; scoreBreakdown?: unknown; }",
		},
		Functions: []spec.Function{
			{Name: "mapping", Returns: spec.Named("MappingBuilder")},
			{Name: "indexMapping", Returns: spec.Named("MappingBuilder")},
			{Name: "docMapping", Returns: spec.Named("DocumentMappingBuilder")},
			{Name: "documentMapping", Returns: spec.Named("DocumentMappingBuilder")},
			{Name: "field", Returns: spec.Named("FieldBuilder")},
			{Name: "search", Returns: spec.Named("SearchRequestBuilder")},
			{Name: "searchRequest", Returns: spec.Named("SearchRequestBuilder")},
			{Name: "create", Params: []spec.Param{{Name: "path", Type: spec.String()}}, Returns: spec.Named("IndexBuilder")},
			{Name: "open", Params: []spec.Param{{Name: "path", Type: spec.String()}}, Returns: spec.Named("IndexBuilder")},
			{Name: "memory", Returns: spec.Named("IndexBuilder")},
			{Name: "match", Params: []spec.Param{{Name: "text", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "matchPhrase", Params: []spec.Param{{Name: "text", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "term", Params: []spec.Param{{Name: "term", Type: spec.Any()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "queryString", Params: []spec.Param{{Name: "query", Type: spec.String()}}, Returns: spec.Named("Query")},
			{Name: "prefix", Params: []spec.Param{{Name: "prefix", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "fuzzy", Params: []spec.Param{{Name: "term", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "regexp", Params: []spec.Param{{Name: "pattern", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "wildcard", Params: []spec.Param{{Name: "pattern", Type: spec.String()}}, Returns: spec.Named("QueryBuilder")},
			{Name: "bool", Returns: spec.Named("BoolQuery")},
			{Name: "conj", Params: []spec.Param{{Name: "queries", Type: spec.Named("Query"), Variadic: true}}, Returns: spec.Named("Query")},
			{Name: "conjunction", Params: []spec.Param{{Name: "queries", Type: spec.Named("Query"), Variadic: true}}, Returns: spec.Named("Query")},
			{Name: "disj", Params: []spec.Param{{Name: "queries", Type: spec.Named("Query"), Variadic: true}}, Returns: spec.Named("Query")},
			{Name: "disjunction", Params: []spec.Param{{Name: "queries", Type: spec.Named("Query"), Variadic: true}}, Returns: spec.Named("Query")},
			{Name: "matchAll", Returns: spec.Named("Query")},
			{Name: "matchNone", Returns: spec.Named("Query")},
		},
	}
}

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
