package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/bold-minds/bench/internal/config"
	"github.com/bold-minds/bench/internal/datafile"
	"github.com/bold-minds/bench/internal/discovery"
	"github.com/bold-minds/bench/internal/model"
	"github.com/bold-minds/bench/internal/render"
)

func newRenderCmd(o *opts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render [target]",
		Short: "Regenerate markdown from existing data",
		Long:  "Regenerate markdown reports from existing .benchmarks.yaml files without running benchmarks.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runRender(o, args)
		},
	}
	cmd.Flags().StringVarP(&o.template, "template", "t", "", "custom template file")
	return cmd
}

func runRender(o *opts, args []string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(root, o.config)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	target := ParseTarget(args, o.recursive)

	pkgs, err := discovery.FindPackages(root, target)
	if err != nil {
		return fmt.Errorf("finding packages: %w", err)
	}

	if o.dryRun {
		fmt.Println("Would render markdown for:")
		for _, p := range pkgs {
			fmt.Printf("  ./%s\n", p)
		}
		return nil
	}

	for _, pkg := range pkgs {
		if renderErr := renderPackage(o, root, pkg, cfg); renderErr != nil {
			return renderErr
		}
	}

	return nil
}

func renderPackage(o *opts, root, pkg string, cfg *model.Config) error {
	absDir := filepath.Join(root, pkg)

	pf, err := datafile.Read(absDir)
	if err != nil {
		return fmt.Errorf("reading data for %s: %w", pkg, err)
	}

	if len(pf.Benchmarks) == 0 {
		if o.verbose {
			fmt.Fprintf(os.Stderr, "skipping %s (no benchmark data)\n", pkg)
		}
		return nil
	}

	tmplPath := cfg.DefaultTemplate
	if o.template != "" {
		tmplPath = o.template
	} else if pc, ok := cfg.Packages[pkg]; ok && pc.Template != "" {
		tmplPath = pc.Template
	}

	mdPath := filepath.Join(absDir, "benchmarks.md")

	var buf bytes.Buffer
	if err = render.Render(&buf, pf, tmplPath); err != nil {
		return fmt.Errorf("rendering markdown for %s: %w", pkg, err)
	}

	//nolint:gosec // markdown report needs to be readable
	if err = os.WriteFile(mdPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing markdown for %s: %w", pkg, err)
	}

	fmt.Printf("rendered %s/benchmarks.md\n", pkg)
	return nil
}
