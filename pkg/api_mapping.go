package pkg

import "github.com/dop251/goja"

func (m *moduleRuntime) mappingBuilder() *goja.Object {
	ref := &mappingBuilderRef{refBase: refBase{api: m, kind: refKindMappingBuild}}
	return m.newWrapper(ref, refKindMappingBuild)
}

func (m *moduleRuntime) docMappingBuilder() *goja.Object {
	ref := &docMappingBuilderRef{refBase: refBase{api: m, kind: refKindDocMapBuilder}}
	return m.newWrapper(ref, refKindDocMapBuilder)
}

func (m *moduleRuntime) fieldBuilder() *goja.Object {
	ref := &fieldBuilderRef{refBase: refBase{api: m, kind: refKindFieldBuilder}}
	return m.newWrapper(ref, refKindFieldBuilder)
}
