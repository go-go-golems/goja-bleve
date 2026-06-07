package pkg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dop251/goja"
)

func valueToFloat32Vector(vm *goja.Runtime, v goja.Value, expectedDims int) ([]float32, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil, fmt.Errorf("bleve: vector is required")
	}
	obj := v.ToObject(vm)
	lengthValue := obj.Get("length")
	if lengthValue == nil || goja.IsUndefined(lengthValue) {
		return nil, fmt.Errorf("bleve: vector must be array-like")
	}
	length := int(lengthValue.ToInteger())
	if length <= 0 {
		return nil, fmt.Errorf("bleve: vector must not be empty")
	}
	if expectedDims > 0 && length != expectedDims {
		return nil, fmt.Errorf("bleve: vector dimension mismatch: got %d, want %d", length, expectedDims)
	}
	out := make([]float32, length)
	for i := 0; i < length; i++ {
		item := obj.Get(strconv.Itoa(i))
		n := item.ToFloat()
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, fmt.Errorf("bleve: vector[%d] must be a finite number", i)
		}
		out[i] = float32(n)
	}
	return out, nil
}
