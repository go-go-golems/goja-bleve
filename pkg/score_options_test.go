package pkg

import "testing"

func TestSearchScoreOptionsValidateInputs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		script   string
		contains []string
	}{
		{
			name:     "invalid score",
			script:   `const bleve = require("bleve"); bleve.search().score("made-up");`,
			contains: []string{"score", "rrf", "rsf"},
		},
		{
			name:     "invalid rank constant",
			script:   `const bleve = require("bleve"); bleve.search().scoreRankConstant(0);`,
			contains: []string{"rank constant", "positive"},
		},
		{
			name:     "invalid window size",
			script:   `const bleve = require("bleve"); bleve.search().scoreWindowSize(0);`,
			contains: []string{"window size", "positive"},
		},
		{
			name:     "window smaller than size",
			script:   `const bleve = require("bleve"); bleve.search().query(bleve.matchAll()).size(5).score("rrf").scoreWindowSize(2).build();`,
			contains: []string{"score window size", "Size (5)"},
		},
	} {
		vm := newBleveTestVM()
		_, err := vm.RunString(tc.script)
		if err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
		if got := err.Error(); !containsAll(got, tc.contains) {
			t.Fatalf("%s: error = %q", tc.name, got)
		}
	}
}
