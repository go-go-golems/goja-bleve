package pkg

import "testing"

func TestPhase3ErrorPathsAreClear(t *testing.T) {
	tests := []struct {
		name string
		js   string
		want []string
	}{
		{
			name: "missing create path",
			js:   `const bleve = require("bleve"); bleve.create("").build();`,
			want: []string{"create index path is required"},
		},
		{
			name: "invalid mapping object",
			js:   `const bleve = require("bleve"); bleve.memory().mapping({}).build();`,
			want: []string{"expected index mapping wrapper", "without Go reference"},
		},
		{
			name: "invalid document id",
			js:   `const bleve = require("bleve"); const idx = bleve.memory().build(); idx.index("", { text: "x" });`,
			want: []string{"document id is required"},
		},
		{
			name: "invalid query object",
			js:   `const bleve = require("bleve"); bleve.search().query({}).build();`,
			want: []string{"expected query wrapper", "without Go reference"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := newBleveTestVM()
			_, err := vm.RunString(tt.js)
			if err == nil {
				t.Fatalf("expected error")
			}
			if got := err.Error(); !containsAll(got, tt.want) {
				t.Fatalf("error = %q, want parts %#v", got, tt.want)
			}
		})
	}
}
