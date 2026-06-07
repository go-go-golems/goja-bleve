//go:build !vectors

package pkg

import (
	"fmt"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func newVectorFieldMapping(_ int, _ bool) (*mapping.FieldMapping, error) {
	return nil, fmt.Errorf("bleve: vector fields require building the host with -tags=vectors")
}

func addKNNToSearchRequest(_ *bleve.SearchRequest, _ string, _ []float32, _ int64, _ float64) error {
	return fmt.Errorf("bleve: KNN search requires building the host with -tags=vectors")
}

func setKNNOperator(_ *bleve.SearchRequest, _ string) error {
	return fmt.Errorf("bleve: KNN search requires building the host with -tags=vectors")
}

func validateKNNAgainstIndexMapping(_ mapping.IndexMapping, _ *bleve.SearchRequest) error {
	return nil
}
