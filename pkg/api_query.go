package pkg

import "github.com/dop251/goja"

func (m *moduleRuntime) matchAllQuery() *goja.Object {
	ref := &queryRef{refBase: refBase{api: m, kind: refKindQuery}, query: "matchAll"}
	return m.newWrapper(ref, refKindQuery)
}

func (m *moduleRuntime) matchNoneQuery() *goja.Object {
	ref := &queryRef{refBase: refBase{api: m, kind: refKindQuery}, query: "matchNone"}
	return m.newWrapper(ref, refKindQuery)
}
