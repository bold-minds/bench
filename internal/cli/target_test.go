package cli

import (
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		recursive bool
		wantKind  model.TargetKind
		wantValue string
	}{
		{"no args", nil, false, model.TargetAll, ""},
		{"all", []string{"all"}, false, model.TargetAll, ""},
		{"package", []string{"internal/scope"}, false, model.TargetPackage, "internal/scope"},
		{"benchmark name", []string{"BenchmarkFoo"}, false, model.TargetBenchmark, "BenchmarkFoo"},
		{"recursive suffix", []string{"./..."}, false, model.TargetRecursive, "./..."},
		{"recursive flag", []string{"internal/scope"}, true, model.TargetRecursive, "internal/scope"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTarget(tt.args, tt.recursive)
			if got.Kind != tt.wantKind {
				t.Errorf("kind: got %d, want %d", got.Kind, tt.wantKind)
			}
			if got.Value != tt.wantValue {
				t.Errorf("value: got %q, want %q", got.Value, tt.wantValue)
			}
		})
	}
}
