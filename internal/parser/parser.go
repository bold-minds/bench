// Package parser extracts structured benchmark results from go test output.
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bold-minds/bench/internal/model"
)

var benchLine = regexp.MustCompile(
	`^(Benchmark\S+)\s+\d+\s+([\d.]+)\s+ns/op` +
		`(?:\s+(\d+)\s+B/op)?` +
		`(?:\s+(\d+)\s+allocs/op)?`,
)

var pkgLine = regexp.MustCompile(`^pkg:\s+(\S+)`)

// Parse extracts benchmark results from raw go test -bench output.
// Results are grouped by package path. If no pkg: line is found, results
// are grouped under an empty string key.
func Parse(output string) ([]model.PackageResults, error) {
	var (
		currentPkg string
		byPkg      = make(map[string][]model.BenchmarkResult)
		order      []string
	)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)

		if m := pkgLine.FindStringSubmatch(line); m != nil {
			currentPkg = m[1]
			if _, seen := byPkg[currentPkg]; !seen {
				order = append(order, currentPkg)
				byPkg[currentPkg] = nil
			}
			continue
		}

		m := benchLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		nsOp, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			continue
		}

		var bOp, allocsOp int
		if m[3] != "" {
			bOp, _ = strconv.Atoi(m[3])
		}
		if m[4] != "" {
			allocsOp, _ = strconv.Atoi(m[4])
		}

		if _, seen := byPkg[currentPkg]; !seen {
			order = append(order, currentPkg)
		}

		byPkg[currentPkg] = append(byPkg[currentPkg], model.BenchmarkResult{
			Name:     m[1],
			NsOp:     nsOp,
			BOp:      bOp,
			AllocsOp: allocsOp,
		})
	}

	results := make([]model.PackageResults, 0, len(order))
	for _, pkg := range order {
		results = append(results, model.PackageResults{
			PackagePath: pkg,
			Results:     byPkg[pkg],
		})
	}
	return results, nil
}
