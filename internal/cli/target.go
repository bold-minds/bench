package cli

import (
	"strings"

	"github.com/bold-minds/bench/internal/model"
)

// ParseTarget interprets the user's positional argument into a Target.
func ParseTarget(args []string, recursive bool) model.Target {
	if len(args) == 0 || args[0] == "all" {
		return model.Target{Kind: model.TargetAll}
	}

	value := args[0]

	if recursive || strings.HasSuffix(value, "/...") {
		return model.Target{Kind: model.TargetRecursive, Value: value}
	}

	if strings.HasPrefix(value, "Benchmark") {
		return model.Target{Kind: model.TargetBenchmark, Value: value}
	}

	return model.Target{Kind: model.TargetPackage, Value: value}
}
