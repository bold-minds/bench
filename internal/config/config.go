// Package config loads and merges bench configuration.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/bold-minds/bench/internal/model"
)

// Load reads configuration by merging (in order of increasing priority):
// ~/.bench.yaml, <repoRoot>/.bench.yaml, and customPath (if set).
func Load(repoRoot, customPath string) (*model.Config, error) {
	cfg := model.DefaultConfig()

	// Global config — errors other than "not found" are surfaced
	home, homeErr := os.UserHomeDir()
	if homeErr == nil {
		if err := mergeFile(&cfg, filepath.Join(home, ".bench.yaml")); err != nil {
			return nil, fmt.Errorf("loading global config: %w", err)
		}
	}

	// Repo config — errors other than "not found" are surfaced
	if repoRoot != "" {
		if err := mergeFile(&cfg, filepath.Join(repoRoot, ".bench.yaml")); err != nil {
			return nil, fmt.Errorf("loading repo config: %w", err)
		}
	}

	// Custom config (flag override) — must exist if specified
	if customPath != "" {
		if _, statErr := os.Stat(customPath); statErr != nil {
			return nil, statErr
		}
		if mergeErr := mergeFile(&cfg, customPath); mergeErr != nil {
			return nil, mergeErr
		}
	}

	return &cfg, nil
}

func mergeFile(dst *model.Config, path string) error {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		if errors.Is(readErr, fs.ErrNotExist) {
			return nil
		}
		return readErr
	}

	var overlay model.Config
	if unmarshalErr := yaml.Unmarshal(data, &overlay); unmarshalErr != nil {
		return unmarshalErr
	}

	if len(overlay.BenchmarkFlags) > 0 {
		dst.BenchmarkFlags = overlay.BenchmarkFlags
	}
	if overlay.DateFormat != "" {
		dst.DateFormat = overlay.DateFormat
	}
	if overlay.DefaultTemplate != "" {
		dst.DefaultTemplate = overlay.DefaultTemplate
	}
	if overlay.OutputFormat != "" {
		dst.OutputFormat = overlay.OutputFormat
	}
	if overlay.HistoryCap > 0 {
		dst.HistoryCap = overlay.HistoryCap
	}
	if len(overlay.Packages) > 0 {
		if dst.Packages == nil {
			dst.Packages = make(map[string]model.PackageConfig)
		}
		for k, v := range overlay.Packages {
			dst.Packages[k] = v
		}
	}

	return nil
}
