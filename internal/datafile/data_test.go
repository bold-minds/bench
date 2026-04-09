package datafile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func TestReadWrite_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	pf := &model.PackageFile{
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

	if err := Write(dir, pf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Verify file exists
	data, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	// Read back
	got, err := Read(dir)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if got.Package != pf.Package {
		t.Errorf("package: got %q, want %q", got.Package, pf.Package)
	}
	if got.LastUpdated != pf.LastUpdated {
		t.Errorf("last_updated: got %q, want %q", got.LastUpdated, pf.LastUpdated)
	}
	if len(got.Benchmarks) != 2 {
		t.Fatalf("expected 2 benchmarks, got %d", len(got.Benchmarks))
	}

	fooHistory := got.Benchmarks["BenchmarkFoo-8"]
	if len(fooHistory) != 2 {
		t.Fatalf("expected 2 history entries for Foo, got %d", len(fooHistory))
	}
	if fooHistory[0].NsOp != 441.7 {
		t.Errorf("expected 441.7 ns/op, got %f", fooHistory[0].NsOp)
	}

	// Verify sorted key order in YAML
	content := string(data)
	barIdx := strings.Index(content, "BenchmarkBar")
	fooIdx := strings.Index(content, "BenchmarkFoo")
	if barIdx > fooIdx {
		t.Error("expected BenchmarkBar before BenchmarkFoo in sorted output")
	}
}

func TestRead_MissingFile(t *testing.T) {
	dir := t.TempDir()

	pf, err := Read(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pf.Benchmarks == nil {
		t.Error("expected non-nil Benchmarks map")
	}
	if len(pf.Benchmarks) != 0 {
		t.Errorf("expected empty Benchmarks, got %d", len(pf.Benchmarks))
	}
}

func TestWrite_DeterministicOutput(t *testing.T) {
	dir := t.TempDir()

	pf := &model.PackageFile{
		Package:     "example",
		LastUpdated: "2025-01-01",
		Benchmarks: map[string][]model.HistoryEntry{
			"BenchmarkZzz-8": {{Date: "2025-01-01", NsOp: 100, BOp: 10, AllocsOp: 1}},
			"BenchmarkAaa-8": {{Date: "2025-01-01", NsOp: 200, BOp: 20, AllocsOp: 2}},
			"BenchmarkMmm-8": {{Date: "2025-01-01", NsOp: 300, BOp: 30, AllocsOp: 3}},
		},
	}

	// Write twice and compare
	if err := Write(dir, pf); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(filepath.Join(dir, FileName))

	if err := Write(dir, pf); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(filepath.Join(dir, FileName))

	if string(first) != string(second) {
		t.Error("output is not deterministic across writes")
	}

	// Verify alphabetical order
	content := string(first)
	aIdx := strings.Index(content, "BenchmarkAaa")
	mIdx := strings.Index(content, "BenchmarkMmm")
	zIdx := strings.Index(content, "BenchmarkZzz")
	if aIdx >= mIdx || mIdx >= zIdx {
		t.Error("benchmarks not in alphabetical order")
	}
}
