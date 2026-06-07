package pkg

import (
	"os"
	"strings"
	"testing"

	"github.com/go-go-golems/go-go-goja/pkg/tsgen/render"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
)

func TestTypeScriptDeclarationSnapshot(t *testing.T) {
	out, err := render.Bundle(&spec.Bundle{
		HeaderComment: "// goja-bleve TypeScript declaration snapshot",
		Modules:       []*spec.Module{module{}.TypeScriptModule()},
	})
	if err != nil {
		t.Fatalf("render TypeScript declaration: %v", err)
	}
	golden, err := os.ReadFile("testdata/bleve.d.ts.golden")
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(string(golden)) {
		t.Fatalf("TypeScript declaration snapshot mismatch\n--- got ---\n%s\n--- want ---\n%s", out, string(golden))
	}
}

func TestTypeScriptDeclarationMentionsPublicAPISurface(t *testing.T) {
	out, err := render.Bundle(&spec.Bundle{Modules: []*spec.Module{module{}.TypeScriptModule()}})
	if err != nil {
		t.Fatalf("render TypeScript declaration: %v", err)
	}
	for _, needle := range []string{
		`declare module "bleve"`,
		`export interface MappingBuilder`,
		`export interface FieldBuilder`,
		`vector(dims: number): this`,
		`export interface SearchRequestBuilder`,
		`score(mode: ScoreMode): this`,
		`knn(field: string, vector: Vector, k: number, boost?: number): this`,
		`export interface Index`,
		`export interface Batch`,
		`export interface SearchResult`,
	} {
		if !strings.Contains(out, needle) {
			t.Fatalf("declaration missing %q\n%s", needle, out)
		}
	}
}
