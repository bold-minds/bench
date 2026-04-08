package datafile

import (
	"testing"

	"github.com/bold-minds/bench/internal/model"
)

func TestMerge_NewBenchmark(t *testing.T) {
	pf := &model.PackageFile{
		Benchmarks: make(map[string][]model.HistoryEntry),
	}

	results := []model.BenchmarkResult{
		{Name: "BenchmarkNew-8", NsOp: 100.5, BOp: 64, AllocsOp: 2},
	}

	Merge(pf, results, "2025-08-09", 20)

	if pf.LastUpdated != "2025-08-09" {
		t.Errorf("expected last_updated 2025-08-09, got %q", pf.LastUpdated)
	}

	history, ok := pf.Benchmarks["BenchmarkNew-8"]
	if !ok {
		t.Fatal("expected BenchmarkNew-8 in benchmarks")
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(history))
	}
	if history[0].NsOp != 100.5 {
		t.Errorf("expected 100.5, got %f", history[0].NsOp)
	}
}

func TestMerge_AppendToExisting(t *testing.T) {
	pf := &model.PackageFile{
		Benchmarks: map[string][]model.HistoryEntry{
			"BenchmarkFoo-8": {
				{Date: "2025-06-01", NsOp: 200, BOp: 128, AllocsOp: 4},
			},
		},
	}

	results := []model.BenchmarkResult{
		{Name: "BenchmarkFoo-8", NsOp: 150, BOp: 96, AllocsOp: 3},
	}

	Merge(pf, results, "2025-08-09", 20)

	history := pf.Benchmarks["BenchmarkFoo-8"]
	if len(history) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(history))
	}
	if history[1].NsOp != 150 {
		t.Errorf("expected 150, got %f", history[1].NsOp)
	}
}

func TestMerge_HistoryCap(t *testing.T) {
	pf := &model.PackageFile{
		Benchmarks: map[string][]model.HistoryEntry{
			"BenchmarkFoo-8": {
				{Date: "2025-01-01", NsOp: 100},
				{Date: "2025-02-01", NsOp: 90},
				{Date: "2025-03-01", NsOp: 80},
			},
		},
	}

	results := []model.BenchmarkResult{
		{Name: "BenchmarkFoo-8", NsOp: 70},
	}

	Merge(pf, results, "2025-04-01", 3)

	history := pf.Benchmarks["BenchmarkFoo-8"]
	if len(history) != 3 {
		t.Fatalf("expected 3 entries (capped), got %d", len(history))
	}
	// Oldest entry should be trimmed
	if history[0].Date != "2025-02-01" {
		t.Errorf("expected oldest remaining entry to be 2025-02-01, got %q", history[0].Date)
	}
}

func TestMerge_ZeroCap_NoLimit(t *testing.T) {
	pf := &model.PackageFile{
		Benchmarks: map[string][]model.HistoryEntry{
			"BenchmarkFoo-8": make([]model.HistoryEntry, 100),
		},
	}

	Merge(pf, []model.BenchmarkResult{{Name: "BenchmarkFoo-8", NsOp: 50}}, "2025-08-09", 0)

	if len(pf.Benchmarks["BenchmarkFoo-8"]) != 101 {
		t.Errorf("expected 101 entries with no cap, got %d", len(pf.Benchmarks["BenchmarkFoo-8"]))
	}
}
