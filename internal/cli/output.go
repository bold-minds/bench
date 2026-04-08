package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bold-minds/bench/internal/model"
)

// FormatResults writes benchmark results in the specified format.
func FormatResults(w io.Writer, results []model.PackageResults, format string) error {
	switch format {
	case "json":
		return formatJSON(w, results)
	case "raw":
		// Raw is handled by the caller writing raw output directly
		return nil
	default:
		return formatTable(w, results)
	}
}

func formatTable(w io.Writer, results []model.PackageResults) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.AlignRight)

	for _, pkg := range results {
		if pkg.PackagePath != "" {
			fmt.Fprintf(tw, "\n%s\n", pkg.PackagePath)
		}
		fmt.Fprintln(tw, "Benchmark\tns/op\tB/op\tallocs/op\t")
		fmt.Fprintln(tw, "---------\t-----\t----\t---------\t")

		for _, r := range pkg.Results {
			fmt.Fprintf(tw, "%s\t%.2f\t%d\t%d\t\n", r.Name, r.NsOp, r.BOp, r.AllocsOp)
		}
	}

	return tw.Flush()
}

func formatJSON(w io.Writer, results []model.PackageResults) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
