package pkg

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

func (m *moduleRuntime) batchObject(ref *batchRef) *goja.Object {
	obj := m.newWrapper(ref, refKindBatch)
	m.mustSet(obj, "index", func(id string, doc goja.Value) (*goja.Object, error) {
		if err := ref.assertUsable(); err != nil {
			return nil, err
		}
		if strings.TrimSpace(id) == "" {
			return nil, fmt.Errorf("bleve: document id is required")
		}
		if err := ref.batch.Index(id, doc.Export()); err != nil {
			return nil, err
		}
		ref.operation++
		return obj, nil
	})
	m.mustSet(obj, "delete", func(id string) (*goja.Object, error) {
		if err := ref.assertUsable(); err != nil {
			return nil, err
		}
		if strings.TrimSpace(id) == "" {
			return nil, fmt.Errorf("bleve: document id is required")
		}
		ref.batch.Delete(id)
		ref.operation++
		return obj, nil
	})
	m.mustSet(obj, "size", func() (int, error) {
		if ref.batch == nil {
			return 0, fmt.Errorf("bleve: batch is not initialized")
		}
		return ref.batch.Size(), nil
	})
	m.mustSet(obj, "operationCount", func() int { return ref.operation })
	m.mustSet(obj, "reset", func() (*goja.Object, error) {
		if ref.executed {
			return nil, fmt.Errorf("bleve: cannot reset an executed batch")
		}
		if ref.batch == nil {
			return nil, fmt.Errorf("bleve: batch is not initialized")
		}
		ref.batch.Reset()
		ref.operation = 0
		return obj, nil
	})
	m.mustSet(obj, "execute", func() error {
		if err := ref.assertUsable(); err != nil {
			return err
		}
		ref.executed = true
		return ref.index.index.Batch(ref.batch)
	})
	return obj
}

func (r *batchRef) assertUsable() error {
	if r == nil || r.batch == nil {
		return fmt.Errorf("bleve: batch is not initialized")
	}
	if r.executed {
		return fmt.Errorf("bleve: batch has already been executed")
	}
	if r.index == nil || r.index.closed || r.index.index == nil {
		return fmt.Errorf("bleve: batch index is closed")
	}
	return nil
}
