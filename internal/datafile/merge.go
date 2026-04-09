package datafile

import (
	"github.com/bold-minds/bench/internal/model"
)

// Merge appends new benchmark results into an existing PackageFile.
// If historyCap > 0, old entries are trimmed to keep at most historyCap entries.
func Merge(pf *model.PackageFile, results []model.BenchmarkResult, date string, historyCap int) {
	for _, r := range results {
		entry := model.HistoryEntry{
			Date:     date,
			NsOp:     r.NsOp,
			BOp:      r.BOp,
			AllocsOp: r.AllocsOp,
		}

		history := pf.Benchmarks[r.Name]
		history = append(history, entry)

		if historyCap > 0 && len(history) > historyCap {
			history = history[len(history)-historyCap:]
		}

		pf.Benchmarks[r.Name] = history
	}

	pf.LastUpdated = date
}
