package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bold-minds/bench/internal/config"
	"github.com/bold-minds/bench/internal/datafile"
	"github.com/bold-minds/bench/internal/discovery"
	"github.com/bold-minds/bench/internal/model"
	"github.com/bold-minds/bench/internal/render"
	"github.com/bold-minds/bench/internal/runner"
)

func newSyncCmd(o *opts) *cobra.Command {
	return &cobra.Command{
		Use:   "sync [target]",
		Short: "Run benchmarks, update data, and regenerate markdown",
		Long:  "Run benchmarks, persist results to .benchmarks.yaml, and regenerate markdown reports.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runSync(o, args)
		},
	}
}

func runSync(o *opts, args []string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(root, o.config)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	target := ParseTarget(args, o.recursive)

	var benchName string
	if target.Kind == model.TargetBenchmark {
		benchName = target.Value
	}

	pkgs, err := discovery.FindPackages(root, target)
	if err != nil {
		return fmt.Errorf("finding packages: %w", err)
	}

	if o.verbose {
		fmt.Fprintf(os.Stderr, "Found %d package(s) with benchmarks\n", len(pkgs))
	}

	if o.dryRun {
		fmt.Println("Would sync benchmarks in:")
		for _, p := range pkgs {
			fmt.Printf("  ./%s\n", p)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	results, _, err := runner.Run(ctx, root, pkgs, cfg.BenchmarkFlags, benchName)
	if err != nil {
		return fmt.Errorf("running benchmarks: %w", err)
	}

	date := time.Now().Format(cfg.DateFormat)
	modPath := readModulePath(root)

	for _, pkgResult := range results {
		pkgDir := findPkgDir(modPath, pkgs, pkgResult.PackagePath)
		if pkgDir == "" {
			fmt.Fprintf(os.Stderr, "warning: could not determine directory for %s, skipping\n", pkgResult.PackagePath)
			continue
		}

		if syncErr := syncPackage(root, pkgDir, pkgResult, cfg, date); syncErr != nil {
			return syncErr
		}
	}

	return nil
}

func syncPackage(root, pkgDir string, pkgResult model.PackageResults, cfg *model.Config, date string) error {
	absDir := filepath.Join(root, pkgDir)

	pf, err := datafile.Read(absDir)
	if err != nil {
		return fmt.Errorf("reading data for %s: %w", pkgDir, err)
	}

	if pf.Package == "" {
		pf.Package = pkgResult.PackagePath
	}

	datafile.Merge(pf, pkgResult.Results, date, cfg.HistoryCap)

	if err = datafile.Write(absDir, pf); err != nil {
		return fmt.Errorf("writing data for %s: %w", pkgDir, err)
	}

	tmplPath := cfg.DefaultTemplate
	if pc, ok := cfg.Packages[pkgDir]; ok && pc.Template != "" {
		tmplPath = pc.Template
	}

	mdPath := filepath.Join(absDir, "benchmarks.md")

	var buf bytes.Buffer
	if err = render.Render(&buf, pf, tmplPath); err != nil {
		return fmt.Errorf("rendering markdown for %s: %w", pkgDir, err)
	}

	//nolint:gosec // markdown report needs to be readable
	if err = os.WriteFile(mdPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing markdown for %s: %w", pkgDir, err)
	}

	fmt.Printf("synced %s (%d benchmarks)\n", pkgDir, len(pkgResult.Results))
	return nil
}

// findPkgDir matches a full package path back to a relative directory.
// The runner returns full import paths (github.com/user/repo/internal/foo)
// but we need the relative dir (internal/foo).
func findPkgDir(modPath string, relPkgs []string, fullPkgPath string) string {
	// Strip module path prefix to get relative path
	if modPath != "" && strings.HasPrefix(fullPkgPath, modPath) {
		rel := fullPkgPath[len(modPath):]
		if len(rel) > 0 && rel[0] == '/' {
			rel = rel[1:]
		}
		if rel == "" {
			rel = "."
		}
		for _, p := range relPkgs {
			if p == rel {
				return p
			}
		}
	}

	// Fallback: match by suffix with path separator boundary
	for _, p := range relPkgs {
		if p == "." {
			if fullPkgPath == modPath {
				return "."
			}
			continue
		}
		if p != "" && strings.HasSuffix(fullPkgPath, "/"+p) {
			return p
		}
	}

	return ""
}

func readModulePath(root string) string {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[7:])
		}
	}
	return ""
}
