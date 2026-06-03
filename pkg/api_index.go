package pkg

import "github.com/dop251/goja"

func (m *moduleRuntime) createIndexBuilder(path string) *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "create", path: path}
	return m.newWrapper(ref, refKindIndexBuilder)
}

func (m *moduleRuntime) openIndexBuilder(path string) *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "open", path: path}
	return m.newWrapper(ref, refKindIndexBuilder)
}

func (m *moduleRuntime) memoryIndexBuilder() *goja.Object {
	ref := &indexBuilderRef{refBase: refBase{api: m, kind: refKindIndexBuilder}, mode: "memory"}
	return m.newWrapper(ref, refKindIndexBuilder)
}
