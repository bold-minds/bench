// Package datafile reads and writes .benchmarks.yaml files.
package datafile

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bold-minds/bench/internal/model"
)

const FileName = ".benchmarks.yaml"

// Read loads a PackageFile from the given directory.
// Returns a zero-value PackageFile if the file does not exist.
func Read(dir string) (*model.PackageFile, error) {
	path := filepath.Join(dir, FileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &model.PackageFile{
				Benchmarks: make(map[string][]model.HistoryEntry),
			}, nil
		}
		return nil, err
	}

	var pf model.PackageFile
	if unmarshalErr := yaml.Unmarshal(data, &pf); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	if pf.Benchmarks == nil {
		pf.Benchmarks = make(map[string][]model.HistoryEntry)
	}
	return &pf, nil
}

// Write saves a PackageFile to the given directory with sorted keys
// for deterministic, diff-friendly output.
func Write(dir string, pf *model.PackageFile) error {
	// Build a sorted representation for clean YAML output
	doc := sortedPackageFile{
		Package:     pf.Package,
		LastUpdated: pf.LastUpdated,
	}

	names := make([]string, 0, len(pf.Benchmarks))
	for name := range pf.Benchmarks {
		names = append(names, name)
	}
	sort.Strings(names)

	doc.Benchmarks = make([]sortedBenchmark, len(names))
	for i, name := range names {
		doc.Benchmarks[i] = sortedBenchmark{
			Name:    name,
			History: pf.Benchmarks[name],
		}
	}

	data, err := marshalSorted(doc)
	if err != nil {
		return err
	}

	//nolint:gosec // benchmark data needs to be readable
	return os.WriteFile(filepath.Join(dir, FileName), data, 0644)
}

type sortedPackageFile struct {
	Package     string
	LastUpdated string
	Benchmarks  []sortedBenchmark
}

type sortedBenchmark struct {
	Name    string
	History []model.HistoryEntry
}

// marshalSorted produces YAML with benchmark names as sorted map keys.
func marshalSorted(doc sortedPackageFile) ([]byte, error) {
	// Build the YAML node tree manually for sorted key order
	root := &yaml.Node{Kind: yaml.MappingNode}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "package"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: doc.Package},
		&yaml.Node{Kind: yaml.ScalarNode, Value: "last_updated"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: doc.LastUpdated},
	)

	benchmarksMap := &yaml.Node{Kind: yaml.MappingNode}
	for _, b := range doc.Benchmarks {
		nameNode := &yaml.Node{Kind: yaml.ScalarNode, Value: b.Name}

		historySeq := &yaml.Node{Kind: yaml.SequenceNode}
		for _, h := range b.History {
			entry := &yaml.Node{Kind: yaml.MappingNode}
			entry.Content = append(entry.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "date"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: h.Date, Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "ns_op"},
				floatNode(h.NsOp),
				&yaml.Node{Kind: yaml.ScalarNode, Value: "b_op"},
				intNode(h.BOp),
				&yaml.Node{Kind: yaml.ScalarNode, Value: "allocs_op"},
				intNode(h.AllocsOp),
			)
			historySeq.Content = append(historySeq.Content, entry)
		}

		benchmarksMap.Content = append(benchmarksMap.Content, nameNode, historySeq)
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "benchmarks"},
		benchmarksMap,
	)

	doc2 := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	return yaml.Marshal(doc2)
}

func floatNode(v float64) *yaml.Node {
	n := &yaml.Node{Kind: yaml.ScalarNode}
	// Format with explicit decimal so YAML infers float, not int
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	n.Value = s
	return n
}

func intNode(v int) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: strconv.Itoa(v),
	}
}
