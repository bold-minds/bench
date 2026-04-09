// Package cli implements the bench command-line interface.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

// opts holds all flag values, scoped to a single command execution.
type opts struct {
	config    string
	verbose   bool
	dryRun    bool
	recursive bool
	output    string
	template  string
}

func newRootCmd() *cobra.Command {
	o := &opts{}

	root := &cobra.Command{
		Use:   "bench",
		Short: "Manage Go benchmark results",
		Long:  "bench runs Go benchmarks, tracks historical performance data, and generates markdown reports.",
	}

	root.PersistentFlags().StringVarP(&o.config, "config", "c", "", "custom config file")
	root.PersistentFlags().BoolVarP(&o.verbose, "verbose", "v", false, "verbose output")
	root.PersistentFlags().BoolVar(&o.dryRun, "dry-run", false, "show what would be done without executing")
	root.PersistentFlags().BoolVarP(&o.recursive, "recursive", "r", false, "include subdirectories")

	root.AddCommand(newRunCmd(o))
	root.AddCommand(newSyncCmd(o))
	root.AddCommand(newRenderCmd(o))
	root.AddCommand(newInitCmd(o))
	root.AddCommand(newVersionCmd())

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(_ *cobra.Command, _ []string) {
			os.Stdout.WriteString("bench " + version + "\n")
		},
	}
}

// Execute runs the root command. Returns 0 on success, 1 on error.
func Execute() int {
	if err := newRootCmd().Execute(); err != nil {
		return 1
	}
	return 0
}
