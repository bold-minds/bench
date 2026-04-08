package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func testPackageFile() *model.PackageFile {
	return &model.PackageFile{
		Package:     "github.com/example/repo/internal/scope",
		LastUpdated: "2025-08-09",
		Benchmarks: map[string][]model.HistoryEntry{
			"BenchmarkFoo-8": {
				{Date: "2025-06-07", NsOp: 441.7, BOp: 208, AllocsOp: 3},
				{Date: "2025-06-13", NsOp: 88.44, BOp: 48, AllocsOp: 2},
			},
			"BenchmarkBar-8": {
				{Date: "2025-06-07", NsOp: 120.8, BOp: 208, AllocsOp: 3},
			},
		},
	}
}

func TestRender_DefaultTemplate(t *testing.T) {
	var buf bytes.Buffer
	pf := testPackageFile()

	if err := Render(&buf, pf, ""); err != nil {
		t.Fatalf("render error: %v", err)
	}

	output := buf.String()

	// Should contain package name
	if !strings.Contains(output, "github.com/example/repo/internal/scope") {
		t.Error("expected package name in output")
	}

	// Should contain last updated
	if !strings.Contains(output, "2025-08-09") {
		t.Error("expected last updated date in output")
	}

	// Should contain benchmark names (sorted: Bar before Foo)
	barIdx := strings.Index(output, "BenchmarkBar")
	fooIdx := strings.Index(output, "BenchmarkFoo")
	if barIdx < 0 || fooIdx < 0 {
		t.Fatal("expected both benchmark names in output")
	}
	if barIdx > fooIdx {
		t.Error("expected BenchmarkBar before BenchmarkFoo (sorted)")
	}

	// Should contain the ns/op values
	if !strings.Contains(output, "441.70") {
		t.Error("expected 441.70 ns/op in output")
	}

	// Should contain delta for Foo's second entry
	// (88.44 - 441.7) / 441.7 * 100 = -79.98...
	if !strings.Contains(output, "-80.0%") {
		t.Errorf("expected -80.0%% delta in output, got:\n%s", output)
	}

	// First entry should have "-" for delta
	if !strings.Contains(output, "| - |") {
		t.Error("expected '-' delta for first history entry")
	}
}

func TestRender_EmptyBenchmarks(t *testing.T) {
	var buf bytes.Buffer
	pf := &model.PackageFile{
		Package:     "example",
		LastUpdated: "2025-01-01",
		Benchmarks:  map[string][]model.HistoryEntry{},
	}

	if err := Render(&buf, pf, ""); err != nil {
		t.Fatalf("render error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "example") {
		t.Error("expected package name even with no benchmarks")
	}
}

func TestRender_Deterministic(t *testing.T) {
	pf := testPackageFile()

	var buf1, buf2 bytes.Buffer
	if err := Render(&buf1, pf, ""); err != nil {
		t.Fatal(err)
	}
	if err := Render(&buf2, pf, ""); err != nil {
		t.Fatal(err)
	}

	if buf1.String() != buf2.String() {
		t.Error("render output is not deterministic")
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{100.0, "100.0"},
		{441.7, "441.70"},
		{1.473, "1.47"},
		{0.0, "0.0"},
	}

	for _, tt := range tests {
		got := formatFloat(tt.in)
		if got != tt.want {
			t.Errorf("formatFloat(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFormatDelta(t *testing.T) {
	pos := 12.5
	neg := -33.3
	zero := 0.0

	tests := []struct {
		in   *float64
		want string
	}{
		{nil, "-"},
		{&pos, "+12.5%"},
		{&neg, "-33.3%"},
		{&zero, "0.0%"},
	}

	for _, tt := range tests {
		got := formatDelta(tt.in)
		if got != tt.want {
			t.Errorf("formatDelta(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
