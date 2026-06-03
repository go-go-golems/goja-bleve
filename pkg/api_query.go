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

func (m *moduleRuntime) termQuery(term string) *goja.Object {
	return m.queryObject(bleve.NewTermQuery(term))
}

func (m *moduleRuntime) queryStringQuery(input string) *goja.Object {
	return m.queryObject(bleve.NewQueryStringQuery(input))
}

func (m *moduleRuntime) matchAllQuery() *goja.Object {
	return m.queryObject(bleve.NewMatchAllQuery())
}

func (m *moduleRuntime) matchNoneQuery() *goja.Object {
	return m.queryObject(bleve.NewMatchNoneQuery())
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
