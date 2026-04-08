package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func testResults() []model.PackageResults {
	return []model.PackageResults{
		{
			PackagePath: "github.com/example/repo",
			Results: []model.BenchmarkResult{
				{Name: "BenchmarkFoo-8", NsOp: 441.7, BOp: 208, AllocsOp: 3},
				{Name: "BenchmarkBar-8", NsOp: 120.8, BOp: 64, AllocsOp: 1},
			},
		},
	}
}

func TestFormatResults_Table(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatResults(&buf, testResults(), "table"); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "BenchmarkFoo-8") {
		t.Error("expected BenchmarkFoo-8 in table output")
	}
	if !strings.Contains(output, "441.70") {
		t.Error("expected 441.70 in table output")
	}
}

func TestFormatResults_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatResults(&buf, testResults(), "json"); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Name": "BenchmarkFoo-8"`) {
		t.Errorf("expected JSON with benchmark name, got:\n%s", output)
	}
	if !strings.Contains(output, `"NsOp": 441.7`) {
		t.Error("expected NsOp in JSON output")
	}
}
