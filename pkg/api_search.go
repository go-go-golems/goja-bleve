package pkg

import "github.com/dop251/goja"

func (m *moduleRuntime) searchRequestBuilder() *goja.Object {
	ref := &searchRequestRef{refBase: refBase{api: m, kind: refKindSearchRequest}}
	return m.newWrapper(ref, refKindSearchRequest)
}
