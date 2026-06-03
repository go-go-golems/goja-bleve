package pkg

import (
	"fmt"
	"strings"

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
	var sort []string
	var highlightFields []string
	var highlightStyle string
	var explain bool
	var score string
	var scoreRankConstant int
	var scoreWindowSize int
	var knnClauses []*knnRef
	var knnOperator string

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
	m.mustSet(obj, "sort", func(value []string) *goja.Object {
		sort = append([]string(nil), value...)
		return obj
	})
	m.mustSet(obj, "highlight", func(value ...goja.Value) (*goja.Object, error) {
		if len(value) > 0 && !goja.IsUndefined(value[0]) && !goja.IsNull(value[0]) {
			highlightFields = exportedStringSlice(value[0].Export())
		}
		if len(value) > 1 && !goja.IsUndefined(value[1]) && !goja.IsNull(value[1]) {
			highlightStyle = value[1].String()
		}
		return obj, nil
	})
	m.mustSet(obj, "explain", func(enabled bool) *goja.Object {
		explain = enabled
		return obj
	})
	m.mustSet(obj, "score", func(value string) (*goja.Object, error) {
		normalized := strings.TrimSpace(strings.ToLower(value))
		switch normalized {
		case "", "default":
			score = bleve.ScoreDefault
		case "none":
			score = bleve.ScoreNone
		case "rrf":
			score = bleve.ScoreRRF
		case "rsf":
			score = bleve.ScoreRSF
		default:
			return nil, fmt.Errorf("bleve: search score must be one of default, none, rrf, or rsf")
		}
		return obj, nil
	})
	m.mustSet(obj, "scoreRankConstant", func(value int) (*goja.Object, error) {
		if value <= 0 {
			return nil, fmt.Errorf("bleve: score rank constant must be positive")
		}
		scoreRankConstant = value
		return obj, nil
	})
	m.mustSet(obj, "scoreWindowSize", func(value int) (*goja.Object, error) {
		if value <= 0 {
			return nil, fmt.Errorf("bleve: score window size must be positive")
		}
		scoreWindowSize = value
		return obj, nil
	})
	m.mustSet(obj, "knnOperator", func(operator string) (*goja.Object, error) {
		normalized := strings.TrimSpace(strings.ToLower(operator))
		if normalized != "or" && normalized != "and" {
			return nil, fmt.Errorf("bleve: KNN operator must be 'or' or 'and'")
		}
		knnOperator = normalized
		return obj, nil
	})
	m.mustSet(obj, "knn", func(field string, vectorValue goja.Value, k int, boostValue ...float64) (*goja.Object, error) {
		if !m.vectorSupport {
			return nil, fmt.Errorf("bleve: KNN search requires building the host with -tags=vectors")
		}
		if field == "" {
			return nil, fmt.Errorf("bleve: KNN field is required")
		}
		if k <= 0 {
			return nil, fmt.Errorf("bleve: KNN k must be positive")
		}
		vector, err := valueToFloat32Vector(m.vm, vectorValue, 0)
		if err != nil {
			return nil, err
		}
		boost := 1.0
		if len(boostValue) > 0 {
			boost = boostValue[0]
		}
		if boost <= 0 {
			return nil, fmt.Errorf("bleve: KNN boost must be positive")
		}
		knnClauses = append(knnClauses, &knnRef{field: field, vector: vector, k: k, boost: boost})
		return obj, nil
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
		if len(sort) > 0 {
			request.SortBy(sort)
		}
		if score != "" {
			request.Score = score
		}
		if scoreRankConstant > 0 || scoreWindowSize > 0 {
			request.Params = bleve.NewDefaultParams(request.From, request.Size)
			if scoreRankConstant > 0 {
				request.Params.ScoreRankConstant = scoreRankConstant
			}
			if scoreWindowSize > 0 {
				request.Params.ScoreWindowSize = scoreWindowSize
			}
			if err := request.Params.Validate(request.Size); err != nil {
				return nil, fmt.Errorf("bleve: invalid score params: %w", err)
			}
		}
		if len(highlightFields) > 0 || highlightStyle != "" {
			if highlightStyle != "" {
				request.Highlight = bleve.NewHighlightWithStyle(highlightStyle)
			} else {
				request.Highlight = bleve.NewHighlight()
			}
			for _, field := range highlightFields {
				request.Highlight.AddField(field)
			}
		}
		if knnOperator != "" {
			if err := setKNNOperator(request, knnOperator); err != nil {
				return nil, err
			}
		}
		for _, knn := range knnClauses {
			if err := addKNNToSearchRequest(request, knn.field, knn.vector, int64(knn.k), knn.boost); err != nil {
				return nil, err
			}
		}
		request.Explain = explain
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
	out := map[string]any{
		"id":     hit.ID,
		"score":  hit.Score,
		"fields": hit.Fields,
	}
	if len(hit.Fragments) > 0 {
		out["fragments"] = hit.Fragments
	}
	if len(hit.Locations) > 0 {
		out["locations"] = hit.Locations
	}
	if len(hit.Sort) > 0 {
		out["sort"] = hit.Sort
	}
	if hit.Expl != nil {
		out["explanation"] = hit.Expl
	}
	if len(hit.ScoreBreakdown) > 0 {
		out["scoreBreakdown"] = hit.ScoreBreakdown
	}
	return out
}

func exportedStringSlice(v any) []string {
	switch vv := v.(type) {
	case []string:
		return append([]string(nil), vv...)
	case []any:
		out := make([]string, 0, len(vv))
		for _, item := range vv {
			out = append(out, fmt.Sprint(item))
		}
		return out
	default:
		if v == nil {
			return nil
		}
		return []string{fmt.Sprint(v)}
	}
}
