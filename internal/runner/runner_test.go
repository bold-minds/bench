package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Integration(t *testing.T) {
	// Create a minimal Go module with a benchmark
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module testmod\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x.go"), []byte("package testmod\nfunc Add(a, b int) int { return a + b }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x_test.go"), []byte(`package testmod

import "testing"

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(1, 2)
	}
}
`), 0644); err != nil {
		t.Fatal(err)
	}

	results, raw, err := Run(context.Background(), root, []string{"."}, []string{"-benchmem", "-count=1"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v\nraw output: %s", err, raw)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 package result")
	}
	if len(results[0].Results) != 1 {
		t.Fatalf("expected 1 benchmark result, got %d", len(results[0].Results))
	}
	name := results[0].Results[0].Name
	if name != "BenchmarkAdd" && !strings.HasPrefix(name, "BenchmarkAdd-") {
		t.Errorf("expected BenchmarkAdd[-N], got %q", name)
	}
	if raw == "" {
		t.Error("expected non-empty raw output")
	}
}

func TestRun_BenchNameFilter(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module testmod\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x.go"), []byte("package testmod\nfunc Add(a, b int) int { return a + b }\nfunc Sub(a, b int) int { return a - b }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x_test.go"), []byte(`package testmod

import "testing"

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ { Add(1, 2) }
}

func BenchmarkSub(b *testing.B) {
	for i := 0; i < b.N; i++ { Sub(1, 2) }
}
`), 0644); err != nil {
		t.Fatal(err)
	}

	results, _, err := Run(context.Background(), root, []string{"."}, []string{"-benchmem", "-count=1"}, "BenchmarkAdd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || len(results[0].Results) != 1 {
		t.Fatalf("expected exactly 1 result for filtered benchmark, got %v", results)
	}
	name := results[0].Results[0].Name
	if name != "BenchmarkAdd" && !strings.HasPrefix(name, "BenchmarkAdd-") {
		t.Errorf("expected BenchmarkAdd[-N], got %q", name)
	}
}

func TestRun_InvalidPackage(t *testing.T) {
	root := t.TempDir()

	// No go.mod at all — go test should fail
	_, _, err := Run(context.Background(), root, []string{"."}, nil, "")
	if err == nil {
		t.Fatal("expected error for missing go.mod")
	}
}
