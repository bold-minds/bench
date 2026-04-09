// Package discovery locates Go packages that contain benchmark tests.
package discovery

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bold-minds/bench/internal/model"
)

// FindPackages returns package directory paths (relative to root) that match
// the given target. Paths are relative to root using "./" prefix.
func FindPackages(root string, target model.Target) ([]string, error) {
	switch target.Kind {
	case model.TargetAll:
		return findAll(root)
	case model.TargetRecursive:
		base := target.Value
		if base == "" || base == "./..." {
			base = "."
		} else {
			base = strings.TrimSuffix(base, "/...")
		}
		return findUnder(root, base)
	case model.TargetPackage:
		dir := filepath.Join(root, target.Value)
		ok, err := hasBenchmarks(dir)
		if err != nil {
			return nil, fmt.Errorf("checking %s: %w", target.Value, err)
		}
		if !ok {
			return nil, fmt.Errorf("no benchmark tests found in %s", target.Value)
		}
		return []string{target.Value}, nil
	case model.TargetBenchmark:
		return findByBenchmarkName(root, target.Value)
	default:
		return nil, fmt.Errorf("unknown target kind: %d", target.Kind)
	}
}

func findAll(root string) ([]string, error) {
	return findUnder(root, ".")
}

func findUnder(root, base string) ([]string, error) {
	searchRoot := filepath.Join(root, base)
	var pkgs []string

	err := filepath.WalkDir(searchRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		// Skip hidden dirs and vendor
		name := d.Name()
		if strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata" {
			return filepath.SkipDir
		}

		hasBench, benchErr := hasBenchmarks(path)
		if benchErr != nil {
			return benchErr
		}
		if hasBench {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			pkgs = append(pkgs, rel)
		}
		return nil
	})

	return pkgs, err
}

func findByBenchmarkName(root, name string) ([]string, error) {
	seen := make(map[string]struct{})
	var pkgs []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			n := d.Name()
			if strings.HasPrefix(n, ".") || n == "vendor" || n == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}

		found, findErr := fileContainsBenchmark(path, name)
		if findErr != nil || !found {
			return findErr
		}

		dir := filepath.Dir(path)
		rel, relErr := filepath.Rel(root, dir)
		if relErr != nil {
			return relErr
		}
		if _, ok := seen[rel]; !ok {
			seen[rel] = struct{}{}
			pkgs = append(pkgs, rel)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("benchmark %q not found in any package", name)
	}
	return pkgs, nil
}

// hasBenchmarks reports whether the directory contains any _test.go file
// with a func Benchmark signature.
func hasBenchmarks(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		found, findErr := fileContainsBenchmark(filepath.Join(dir, e.Name()), "")
		if findErr != nil {
			return false, findErr
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

// fileContainsBenchmark scans a file for func Benchmark lines.
// If name is non-empty, it looks for that specific benchmark function.
func fileContainsBenchmark(path, name string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if name != "" {
			if strings.HasPrefix(trimmed, "func "+name+"(") {
				return true, nil
			}
		} else {
			if strings.HasPrefix(strings.TrimSpace(line), "func Benchmark") {
				return true, nil
			}
		}
	}
	return false, scanner.Err()
}
