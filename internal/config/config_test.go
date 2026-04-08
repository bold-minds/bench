package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "table" {
		t.Errorf("expected output format 'table', got %q", cfg.OutputFormat)
	}
	if cfg.HistoryCap != 20 {
		t.Errorf("expected history cap 20, got %d", cfg.HistoryCap)
	}
	if cfg.DateFormat != "2006-01-02" {
		t.Errorf("expected date format '2006-01-02', got %q", cfg.DateFormat)
	}
	if len(cfg.BenchmarkFlags) != 2 {
		t.Errorf("expected 2 default flags, got %d", len(cfg.BenchmarkFlags))
	}
}

func TestLoad_RepoOverridesGlobal(t *testing.T) {
	dir := t.TempDir()

	// Write a repo config that overrides output format
	repoConfig := []byte("output_format: json\nhistory_cap: 50\n")
	if err := os.WriteFile(filepath.Join(dir, ".bench.yaml"), repoConfig, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output format 'json', got %q", cfg.OutputFormat)
	}
	if cfg.HistoryCap != 50 {
		t.Errorf("expected history cap 50, got %d", cfg.HistoryCap)
	}
	// Defaults should still be in place for non-overridden fields
	if cfg.DateFormat != "2006-01-02" {
		t.Errorf("expected default date format, got %q", cfg.DateFormat)
	}
}

func TestLoad_CustomPathOverridesAll(t *testing.T) {
	dir := t.TempDir()

	repoConfig := []byte("output_format: json\n")
	if err := os.WriteFile(filepath.Join(dir, ".bench.yaml"), repoConfig, 0644); err != nil {
		t.Fatal(err)
	}

	customPath := filepath.Join(dir, "custom.yaml")
	customConfig := []byte("output_format: raw\n")
	if err := os.WriteFile(customPath, customConfig, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir, customPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "raw" {
		t.Errorf("expected output format 'raw', got %q", cfg.OutputFormat)
	}
}

func TestLoad_MissingFilesAreOk(t *testing.T) {
	cfg, err := Load("/nonexistent/path", "")
	if err != nil {
		t.Fatalf("expected no error for missing files, got: %v", err)
	}
	if cfg.OutputFormat != "table" {
		t.Errorf("expected default output format, got %q", cfg.OutputFormat)
	}
}

func TestLoad_CustomPathMissingIsError(t *testing.T) {
	_, err := Load("", "/nonexistent/custom.yaml")
	if err == nil {
		t.Fatal("expected error for missing custom config path")
	}
}

func TestLoad_PackageOverrides(t *testing.T) {
	dir := t.TempDir()

	cfg := []byte("packages:\n  internal/scope:\n    template: custom.md.tmpl\n")
	if err := os.WriteFile(filepath.Join(dir, ".bench.yaml"), cfg, 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Load(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Packages == nil {
		t.Fatal("expected packages map to be non-nil")
	}
	pc, ok := result.Packages["internal/scope"]
	if !ok {
		t.Fatal("expected internal/scope in packages")
	}
	if pc.Template != "custom.md.tmpl" {
		t.Errorf("expected template 'custom.md.tmpl', got %q", pc.Template)
	}
}
