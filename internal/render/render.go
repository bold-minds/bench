// Package render generates markdown from benchmark data using templates.
package render

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"text/template"

	"github.com/bold-minds/bench/internal/model"
)

// TemplateData is the top-level data passed to the markdown template.
type TemplateData struct {
	Package     string
	LastUpdated string
	Latest      []LatestEntry
	Benchmarks  []BenchmarkHistory
}

// LatestEntry is the most recent result for a benchmark.
type LatestEntry struct {
	Name     string
	NsOp     float64
	BOp      int
	AllocsOp int
}

// BenchmarkHistory is a named benchmark with its full history and deltas.
type BenchmarkHistory struct {
	Name    string
	Entries []HistoryWithDelta
}

// HistoryWithDelta extends a history entry with a computed delta.
type HistoryWithDelta struct {
	Date       string
	NsOp      float64
	BOp        int
	AllocsOp   int
	DeltaNsOp  *float64 // nil for the first entry
}

// Render writes markdown for the given PackageFile. If tmplPath is empty,
// the embedded default template is used.
func Render(w io.Writer, pf *model.PackageFile, tmplPath string) error {
	tmplStr := defaultTemplate
	if tmplPath != "" {
		data, err := os.ReadFile(tmplPath)
		if err != nil {
			return err
		}
		tmplStr = string(data)
	}

	funcMap := template.FuncMap{
		"formatFloat": formatFloat,
		"formatDelta": formatDelta,
	}

	tmpl, err := template.New("benchmarks").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return err
	}

	td := buildTemplateData(pf)
	return tmpl.Execute(w, td)
}

func buildTemplateData(pf *model.PackageFile) TemplateData {
	names := sortedBenchmarkNames(pf.Benchmarks)

	td := TemplateData{
		Package:     pf.Package,
		LastUpdated: pf.LastUpdated,
	}

	for _, name := range names {
		history := pf.Benchmarks[name]
		if len(history) == 0 {
			continue
		}

		// Latest entry
		latest := history[len(history)-1]
		td.Latest = append(td.Latest, LatestEntry{
			Name:     name,
			NsOp:     latest.NsOp,
			BOp:      latest.BOp,
			AllocsOp: latest.AllocsOp,
		})

		// History with deltas
		bh := BenchmarkHistory{Name: name}
		for i, h := range history {
			entry := HistoryWithDelta{
				Date:     h.Date,
				NsOp:     h.NsOp,
				BOp:      h.BOp,
				AllocsOp: h.AllocsOp,
			}
			if i > 0 {
				prev := history[i-1]
				if prev.NsOp != 0 {
					delta := ((h.NsOp - prev.NsOp) / prev.NsOp) * 100
					entry.DeltaNsOp = &delta
				}
			}
			bh.Entries = append(bh.Entries, entry)
		}
		td.Benchmarks = append(td.Benchmarks, bh)
	}

	return td
}

func sortedBenchmarkNames(m map[string][]model.HistoryEntry) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.1f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

func formatDelta(v *float64) string {
	if v == nil {
		return "-"
	}
	d := *v
	if d > 0 {
		return fmt.Sprintf("+%.1f%%", d)
	}
	if d < 0 {
		return fmt.Sprintf("%.1f%%", d)
	}
	return "0.0%"
}
