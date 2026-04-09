// Package runner executes go test benchmarks and returns parsed results.
package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bold-minds/bench/internal/model"
	"github.com/bold-minds/bench/internal/parser"
)

var validBenchName = regexp.MustCompile(`^Benchmark[A-Za-z0-9_]+$`)

// Run executes benchmarks for the given packages and returns parsed results
// plus the raw output. The benchName filter, if non-empty, restricts which
// benchmarks are run (e.g. "BenchmarkFoo").
func Run(ctx context.Context, root string, packages []string, flags []string, benchName string) ([]model.PackageResults, string, error) {
	benchPattern := "."
	if benchName != "" {
		if !validBenchName.MatchString(benchName) {
			return nil, "", fmt.Errorf("invalid benchmark name %q: must match Benchmark[A-Za-z0-9_]+", benchName)
		}
		benchPattern = "^" + benchName + "$"
	}

	args := []string{"test", "-bench=" + benchPattern, "-run=^$"}
	args = append(args, flags...)

	// Build package list as ./relative/path for each
	for _, pkg := range packages {
		args = append(args, "./"+pkg)
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	rawOutput := stdout.String()

	if err != nil {
		if rawOutput == "" {
			errMsg := stderr.String()
			if errMsg == "" {
				errMsg = err.Error()
			}
			return nil, "", fmt.Errorf("go test failed: %s", strings.TrimSpace(errMsg))
		}
		// Partial output: some packages may have failed. Warn but continue
		// with whatever results were captured.
		fmt.Fprintf(os.Stderr, "warning: go test exited with error (partial results may be returned): %v\n", err)
	}

	results, parseErr := parser.Parse(rawOutput)
	if parseErr != nil {
		return nil, rawOutput, fmt.Errorf("parsing benchmark output: %w", parseErr)
	}

	return results, rawOutput, nil
}
