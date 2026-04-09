// Package model defines the shared types for bench.
package model

// BenchmarkResult holds a single parsed benchmark measurement.
type BenchmarkResult struct {
	Name     string
	NsOp     float64
	BOp      int
	AllocsOp int
}

// PackageResults groups benchmark results by package.
type PackageResults struct {
	PackagePath string
	Results     []BenchmarkResult
}

// HistoryEntry is a single timestamped benchmark measurement stored in YAML.
type HistoryEntry struct {
	Date     string  `yaml:"date"`
	NsOp     float64 `yaml:"ns_op"`
	BOp      int     `yaml:"b_op"`
	AllocsOp int     `yaml:"allocs_op"`
}

// PackageFile is the schema for .benchmarks.yaml files.
type PackageFile struct {
	Package     string                    `yaml:"package"`
	LastUpdated string                    `yaml:"last_updated"`
	Benchmarks  map[string][]HistoryEntry `yaml:"benchmarks"`
}

// TargetKind classifies what the user wants to benchmark.
type TargetKind int

const (
	TargetAll       TargetKind = iota // all packages
	TargetPackage                     // specific package path
	TargetBenchmark                   // single benchmark by name
	TargetRecursive                   // ./... recursive
)

// Target describes what the user wants to benchmark.
type Target struct {
	Kind  TargetKind
	Value string
}

// Config holds merged configuration from .bench.yaml files.
type Config struct {
	BenchmarkFlags  []string                  `yaml:"benchmark_flags"`
	DateFormat      string                    `yaml:"date_format"`
	DefaultTemplate string                    `yaml:"default_template"`
	OutputFormat    string                    `yaml:"output_format"`
	HistoryCap      int                       `yaml:"history_cap"`
	Packages        map[string]PackageConfig  `yaml:"packages"`
}

// PackageConfig holds per-package configuration overrides.
type PackageConfig struct {
	Template string `yaml:"template"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		BenchmarkFlags: []string{"-benchmem", "-count=1"},
		DateFormat:     "2006-01-02",
		OutputFormat:   "table",
		HistoryCap:     20,
	}
}
