package pkg

import (
	"fmt"
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/dop251/goja"
)

func (m *moduleRuntime) matchQuery(text string) *goja.Object {
	return m.queryObject(bleve.NewMatchQuery(text))
}

func (m *moduleRuntime) matchPhraseQuery(text string) *goja.Object {
	return m.queryObject(bleve.NewMatchPhraseQuery(text))
}

func (m *moduleRuntime) termQuery(term string) *goja.Object {
	return m.queryObject(bleve.NewTermQuery(term))
}

func (m *moduleRuntime) queryStringQuery(input string) *goja.Object {
	return m.queryObject(bleve.NewQueryStringQuery(input))
}

func (m *moduleRuntime) prefixQuery(prefix string) *goja.Object {
	return m.queryObject(bleve.NewPrefixQuery(prefix))
}

func (m *moduleRuntime) fuzzyQuery(term string) *goja.Object {
	return m.queryObject(bleve.NewFuzzyQuery(term))
}

func (m *moduleRuntime) regexpQuery(pattern string) *goja.Object {
	return m.queryObject(bleve.NewRegexpQuery(pattern))
}

func (m *moduleRuntime) wildcardQuery(pattern string) *goja.Object {
	return m.queryObject(bleve.NewWildcardQuery(pattern))
}

func (m *moduleRuntime) matchAllQuery() *goja.Object {
	return m.queryObject(bleve.NewMatchAllQuery())
}

func (m *moduleRuntime) matchNoneQuery() *goja.Object {
	return m.queryObject(bleve.NewMatchNoneQuery())
}

func (m *moduleRuntime) boolQuery() *goja.Object {
	q := bleve.NewBooleanQuery()
	obj := m.queryObject(q)
	m.mustSet(obj, "addMust", func(values ...goja.Value) (*goja.Object, error) {
		queries, err := m.queryRefs(values, "bool.addMust")
		if err != nil {
			return nil, err
		}
		q.AddMust(queries...)
		return obj, nil
	})
	m.mustSet(obj, "addShould", func(values ...goja.Value) (*goja.Object, error) {
		queries, err := m.queryRefs(values, "bool.addShould")
		if err != nil {
			return nil, err
		}
		q.AddShould(queries...)
		return obj, nil
	})
	m.mustSet(obj, "addMustNot", func(values ...goja.Value) (*goja.Object, error) {
		queries, err := m.queryRefs(values, "bool.addMustNot")
		if err != nil {
			return nil, err
		}
		q.AddMustNot(queries...)
		return obj, nil
	})
	return obj
}

func (m *moduleRuntime) conjunctionQuery(values ...goja.Value) (*goja.Object, error) {
	queries, err := m.queryRefs(values, "conjunction")
	if err != nil {
		return nil, err
	}
	return m.queryObject(bleve.NewConjunctionQuery(queries...)), nil
}

func (m *moduleRuntime) disjunctionQuery(values ...goja.Value) (*goja.Object, error) {
	queries, err := m.queryRefs(values, "disjunction")
	if err != nil {
		return nil, err
	}
	return m.queryObject(bleve.NewDisjunctionQuery(queries...)), nil
}

func (m *moduleRuntime) queryRefs(values []goja.Value, label string) ([]query.Query, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("bleve: %s requires at least one query", label)
	}
	queries := make([]query.Query, 0, len(values))
	for _, value := range values {
		ref, err := getTypedRef[queryRef](m, value, "query")
		if err != nil {
			return nil, err
		}
		if ref.query == nil {
			return nil, fmt.Errorf("bleve: %s received empty query", label)
		}
		queries = append(queries, ref.query)
	}
	return queries, nil
}

func (m *moduleRuntime) queryObject(q query.Query) *goja.Object {
	ref := &queryRef{refBase: refBase{api: m, kind: refKindQuery}, query: q}
	obj := m.newWrapper(ref, refKindQuery)
	m.mustSet(obj, "field", func(name string) (*goja.Object, error) {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("bleve: query field name is required")
		}
		if fieldable, ok := ref.query.(interface{ SetField(string) }); ok {
			fieldable.SetField(name)
			return obj, nil
		}
		return nil, fmt.Errorf("bleve: query type %T does not support field()", ref.query)
	})
	m.mustSet(obj, "boost", func(boost float64) (*goja.Object, error) {
		if boost <= 0 {
			return nil, fmt.Errorf("bleve: query boost must be positive")
		}
		if boostable, ok := ref.query.(interface{ SetBoost(float64) }); ok {
			boostable.SetBoost(boost)
			return obj, nil
		}
		return nil, fmt.Errorf("bleve: query type %T does not support boost()", ref.query)
	})
	return obj
}
