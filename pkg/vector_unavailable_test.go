//go:build !vectors

package pkg

import "testing"

func TestVectorAPIsReportUnavailableWithoutVectorBuildTag(t *testing.T) {
	for _, script := range []string{
		`const bleve = require("bleve"); bleve.field().vector(3);`,
		`const bleve = require("bleve"); bleve.search().query(bleve.matchAll()).knn("embedding", [1, 0, 0], 2);`,
	} {
		vm := newBleveTestVM()
		_, err := vm.RunString(script)
		if err == nil {
			t.Fatalf("expected vector unavailable error for %s", script)
		}
		if got := err.Error(); !containsAll(got, []string{"require", "building the host with -tags=vectors"}) {
			t.Fatalf("vector unavailable error = %q", got)
		}
	}
}
