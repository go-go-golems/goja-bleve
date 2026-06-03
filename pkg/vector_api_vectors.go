//go:build vectors

package pkg

import (
	"fmt"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func newVectorFieldMapping(dims int, base64 bool) (*mapping.FieldMapping, error) {
	if dims <= 0 {
		return nil, fmt.Errorf("bleve: vector dims must be positive")
	}
	var field *mapping.FieldMapping
	if base64 {
		field = bleve.NewVectorBase64FieldMapping()
	} else {
		field = bleve.NewVectorFieldMapping()
	}
	field.Dims = dims
	return field, nil
}

func addKNNToSearchRequest(request *bleve.SearchRequest, field string, vector []float32, k int64, boost float64) error {
	request.AddKNN(field, vector, k, boost)
	return nil
}

func setKNNOperator(request *bleve.SearchRequest, operator string) error {
	switch operator {
	case "", "or":
		request.AddKNNOperator("or")
	case "and":
		request.AddKNNOperator("and")
	default:
		return fmt.Errorf("bleve: KNN operator must be 'or' or 'and'")
	}
	return nil
}

func validateKNNAgainstIndexMapping(indexMapping *mapping.IndexMappingImpl, request *bleve.SearchRequest) error {
	if indexMapping == nil || request == nil {
		return nil
	}
	for _, knn := range request.KNN {
		if knn == nil {
			continue
		}
		fieldMapping := indexMapping.FieldMappingForPath(knn.Field)
		if fieldMapping.Type == "" {
			return fmt.Errorf("bleve: kNN field %q is not mapped", knn.Field)
		}
		if fieldMapping.Dims > 0 && len(knn.Vector) != fieldMapping.Dims {
			return fmt.Errorf("bleve: kNN vector dimension mismatch for field %q: got %d, want %d", knn.Field, len(knn.Vector), fieldMapping.Dims)
		}
	}
	return nil
}
