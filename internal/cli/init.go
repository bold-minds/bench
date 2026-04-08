package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const defaultConfigContent = `# bench configuration
# See: https://github.com/bold-minds/bench

# Flags passed to go test when running benchmarks
# benchmark_flags:
#   - "-benchmem"
#   - "-count=1"

# Date format for history entries (Go time format)
# date_format: "2006-01-02"

# Default output format for 'bench run' (table, json, raw)
# output_format: table

# Maximum history entries per benchmark (0 = unlimited)
# history_cap: 20

# Per-package overrides
# packages:
#   internal/scope:
#     template: custom.md.tmpl
`

func newInitCmd(o *opts) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize bench in the current repository",
		Long:  "Create a .bench.yaml configuration file with commented defaults.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInit(o)
		},
	}
}

func runInit(o *opts) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(root, ".bench.yaml")

	if _, statErr := os.Stat(path); statErr == nil {
		return fmt.Errorf(".bench.yaml already exists")
	}

	if o.dryRun {
		fmt.Println("Would create .bench.yaml")
		return nil
	}

	//nolint:gosec // config file needs to be readable by the user
	if writeErr := os.WriteFile(path, []byte(defaultConfigContent), 0644); writeErr != nil {
		return writeErr
	}

	fmt.Println("Created .bench.yaml")
	return nil
}
