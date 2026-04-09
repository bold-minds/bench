package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// Package with benchmarks
	pkgA := filepath.Join(root, "internal", "alpha")
	os.MkdirAll(pkgA, 0755)
	os.WriteFile(filepath.Join(pkgA, "alpha.go"), []byte("package alpha\n"), 0644)
	os.WriteFile(filepath.Join(pkgA, "alpha_test.go"), []byte(`package alpha

import "testing"

func BenchmarkAlpha(b *testing.B) {
	for i := 0; i < b.N; i++ {}
}
`), 0644)

	// Package without benchmarks
	pkgB := filepath.Join(root, "internal", "beta")
	os.MkdirAll(pkgB, 0755)
	os.WriteFile(filepath.Join(pkgB, "beta.go"), []byte("package beta\n"), 0644)
	os.WriteFile(filepath.Join(pkgB, "beta_test.go"), []byte(`package beta

import "testing"

func TestBeta(t *testing.T) {}
`), 0644)

	// Nested package with benchmarks
	pkgC := filepath.Join(root, "internal", "gamma", "deep")
	os.MkdirAll(pkgC, 0755)
	os.WriteFile(filepath.Join(pkgC, "deep.go"), []byte("package deep\n"), 0644)
	os.WriteFile(filepath.Join(pkgC, "bench_test.go"), []byte(`package deep

import "testing"

func BenchmarkDeep(b *testing.B) {
	for i := 0; i < b.N; i++ {}
}
`), 0644)

	return root
}

func TestFindPackages_All(t *testing.T) {
	root := setupTestRepo(t)

	pkgs, err := FindPackages(root, model.Target{Kind: model.TargetAll})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages with benchmarks, got %d: %v", len(pkgs), pkgs)
	}

	// Should find alpha and gamma/deep, but not beta
	found := map[string]bool{}
	for _, p := range pkgs {
		found[p] = true
	}
	if !found[filepath.Join("internal", "alpha")] {
		t.Error("expected to find internal/alpha")
	}
	if !found[filepath.Join("internal", "gamma", "deep")] {
		t.Error("expected to find internal/gamma/deep")
	}
	if found[filepath.Join("internal", "beta")] {
		t.Error("should not find internal/beta (no benchmarks)")
	}
}

func TestFindPackages_SpecificPackage(t *testing.T) {
	root := setupTestRepo(t)

	pkgs, err := FindPackages(root, model.Target{
		Kind:  model.TargetPackage,
		Value: filepath.Join("internal", "alpha"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
}

func TestFindPackages_SpecificPackage_NoBenchmarks(t *testing.T) {
	root := setupTestRepo(t)

	_, err := FindPackages(root, model.Target{
		Kind:  model.TargetPackage,
		Value: filepath.Join("internal", "beta"),
	})
	if err == nil {
		t.Fatal("expected error for package without benchmarks")
	}
}

func TestFindPackages_BenchmarkName(t *testing.T) {
	root := setupTestRepo(t)

	pkgs, err := FindPackages(root, model.Target{
		Kind:  model.TargetBenchmark,
		Value: "BenchmarkDeep",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0] != filepath.Join("internal", "gamma", "deep") {
		t.Errorf("expected internal/gamma/deep, got %q", pkgs[0])
	}
}

func TestFindPackages_BenchmarkName_NotFound(t *testing.T) {
	root := setupTestRepo(t)

	_, err := FindPackages(root, model.Target{
		Kind:  model.TargetBenchmark,
		Value: "BenchmarkNonexistent",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent benchmark")
	}
}

func TestFindPackages_Recursive(t *testing.T) {
	root := setupTestRepo(t)

	pkgs, err := FindPackages(root, model.Target{
		Kind:  model.TargetRecursive,
		Value: "internal/gamma/...",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package under gamma, got %d: %v", len(pkgs), pkgs)
	}
}

func TestFindPackages_SkipsHiddenAndVendor(t *testing.T) {
	root := setupTestRepo(t)

	// Create a hidden dir and vendor dir with benchmarks
	for _, name := range []string{".hidden", "vendor"} {
		dir := filepath.Join(root, name)
		os.MkdirAll(dir, 0755)
		os.WriteFile(filepath.Join(dir, "bench_test.go"), []byte(`package x
import "testing"
func BenchmarkHidden(b *testing.B) {}
`), 0644)
	}

	pkgs, err := FindPackages(root, model.Target{Kind: model.TargetAll})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range pkgs {
		if p == ".hidden" || p == "vendor" {
			t.Errorf("should not find package in %s", p)
		}
	}
}
