package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bold-minds/bench/internal/config"
	"github.com/bold-minds/bench/internal/discovery"
	"github.com/bold-minds/bench/internal/model"
	"github.com/bold-minds/bench/internal/runner"
)

func newRunCmd(o *opts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [target]",
		Short: "Run benchmarks without updating files",
		Long:  "Run benchmarks and display results. Does not write any files.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runRun(o, args)
		},
	}
	cmd.Flags().StringVarP(&o.output, "output", "o", "", "output format: table, json, raw")
	return cmd
}

func runRun(o *opts, args []string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(root, o.config)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	outputFmt := cfg.OutputFormat
	if o.output != "" {
		outputFmt = o.output
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
		for _, p := range pkgs {
			fmt.Fprintf(os.Stderr, "  %s\n", p)
		}
	}

	if o.dryRun {
		fmt.Println("Would run benchmarks in:")
		for _, p := range pkgs {
			fmt.Printf("  ./%s\n", p)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	results, rawOutput, err := runner.Run(ctx, root, pkgs, cfg.BenchmarkFlags, benchName)
	if err != nil {
		return fmt.Errorf("running benchmarks: %w", err)
	}

	if outputFmt == "raw" {
		_, err = os.Stdout.WriteString(rawOutput)
		return err
	}

	return FormatResults(os.Stdout, results, outputFmt)
}
