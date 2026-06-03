package pkg

import (
	"fmt"
	"math"

	"github.com/dop251/goja"
)

func requireFiniteNumber(v goja.Value, label string) (float64, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return 0, fmt.Errorf("bleve: %s is required", label)
	}
	n := v.ToFloat()
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("bleve: %s must be a finite number", label)
	}
	return n, nil
}
