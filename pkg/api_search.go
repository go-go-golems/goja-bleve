package pkg

import (
	"fmt"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/dop251/goja"
)

func (m *moduleRuntime) searchRequestBuilder() *goja.Object {
	ref := &searchRequestRef{refBase: refBase{api: m, kind: refKindSearchRequest}}
	obj := m.newWrapper(ref, refKindSearchRequest)
	var queryRefValue *queryRef
	var size int
	var from int
	var fields []string

	m.mustSet(obj, "query", func(value goja.Value) (*goja.Object, error) {
		q, err := getTypedRef[queryRef](m, value, "query")
		if err != nil {
			return nil, err
		}
		queryRefValue = q
		return obj, nil
	})
	m.mustSet(obj, "size", func(value int) (*goja.Object, error) {
		if value < 0 {
			return nil, fmt.Errorf("bleve: search size must be non-negative")
		}
		size = value
		return obj, nil
	})
	m.mustSet(obj, "from", func(value int) (*goja.Object, error) {
		if value < 0 {
			return nil, fmt.Errorf("bleve: search from must be non-negative")
		}
		from = value
		return obj, nil
	})
	m.mustSet(obj, "fields", func(value []string) *goja.Object {
		fields = append([]string(nil), value...)
		return obj
	})
	m.mustSet(obj, "build", func() (*goja.Object, error) {
		if queryRefValue == nil || queryRefValue.query == nil {
			return nil, fmt.Errorf("bleve: search query is required before build()")
		}
		request := bleve.NewSearchRequest(queryRefValue.query)
		if size > 0 {
			request.Size = size
		}
		if from > 0 {
			request.From = from
		}
		if len(fields) > 0 {
			request.Fields = append([]string(nil), fields...)
		}
		built := &searchRequestRef{refBase: refBase{api: m, kind: refKindSearchRequest}, request: request}
		return m.newWrapper(built, refKindSearchRequest), nil
	})
	return obj
}

func searchResultToJS(result *bleve.SearchResult) map[string]any {
	out := map[string]any{
		"total":    result.Total,
		"maxScore": result.MaxScore,
		"took":     result.Took.String(),
		"hits":     []map[string]any{},
	}
	hits := make([]map[string]any, 0, len(result.Hits))
	for _, hit := range result.Hits {
		hits = append(hits, searchHitToJS(hit))
	}
	out["hits"] = hits
	return out
}

func searchHitToJS(hit *search.DocumentMatch) map[string]any {
	return map[string]any{
		"id":     hit.ID,
		"score":  hit.Score,
		"fields": hit.Fields,
	}
}
