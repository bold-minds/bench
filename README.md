# BenchSync: Go Benchmark Management Tool

## 🎯 Project Overview

**BenchSync** is a standalone CLI tool for managing Go benchmark results across repositories. It provides a robust, repo-agnostic solution for running benchmarks, tracking historical performance data, and generating consistent markdown documentation.

### Key Problems Solved
- Manual benchmark result tracking is error-prone and tedious
- Inconsistent benchmark documentation across repositories
- Complex `go test` command lines for benchmark execution
- No standardized way to track performance deltas over time
- Fragile markdown parsing for historical data extraction

### Solution Approach
- Use committed YAML/JSON files as persistent, canonical data sources
- Template-based markdown generation (no parsing required)
- Single CLI tool installable globally, usable in any Go repository
- Flexible targeting: all packages, specific packages, single tests, recursive discovery

---

## 🚀 CLI Tool: `benchsync`

### Installation & Basic Usage
```bash
# Install once globally
go install github.com/yourusername/benchsync@latest

# Use in any Go repo
cd your-go-repo/
benchsync run all                    # Just run benchmarks, no updates
benchsync sync all                   # Run + update data + regenerate markdown
benchsync sync internal/scope        # Specific package
benchsync sync BenchmarkNewKey       # Single test
benchsync sync --recursive ./...     # Recursive from current dir
```

### Command Structure
```bash
benchsync <command> [target] [flags]

Commands:
  run     Run benchmarks without updating files (simpler than go test)
  sync    Run benchmarks and update data/markdown files
  render  Regenerate markdown from existing data (no benchmark run)
  init    Initialize benchsync in current repo (create config/templates)

Targets:
  all                    All packages with benchmark tests
  <package-path>         Specific package (e.g., internal/scope)
  <benchmark-name>       Single benchmark test
  ./...                  Current directory and subdirectories (with --recursive)

Flags:
  --recursive, -r        Include subdirectories
  --output, -o          Output format for run command (table, json, raw)
  --template, -t        Custom template file
  --config, -c          Custom config file
  --dry-run            Show what would be done without executing
  --verbose, -v         Verbose output
```

### Usage Examples

#### 1. Just Run Benchmarks (No Updates)
```bash
# Simpler than: go test -bench=. -benchmem ./internal/scope
benchsync run internal/scope

# Output formats
benchsync run all --output table    # Pretty table
benchsync run all --output json     # JSON for scripting
benchsync run all --output raw      # Raw go test output
```

#### 2. Run + Update Data + Regenerate Markdown
```bash
benchsync sync all                  # All packages
benchsync sync internal/scope       # Specific package
benchsync sync BenchmarkNewKey      # Single test
benchsync sync --recursive ./...    # Recursive from current dir
```

#### 3. Just Regenerate Markdown (No Benchmark Run)
```bash
benchsync render all                # Useful after manual data edits
benchsync render internal/scope     # Specific package
```

---

## 🏗️ Architecture & Repository Structure

### Standalone Repository Structure
```
benchsync/                           # Standalone repo
├── cmd/benchsync/main.go           # CLI entry point
├── internal/
│   ├── runner/benchmark.go         # Execute go test benchmarks
│   ├── parser/results.go           # Parse benchmark output
│   ├── data/manager.go             # Manage YAML/JSON data files
│   ├── template/renderer.go        # Render markdown from templates
│   ├── discovery/finder.go         # Find benchmark files/packages
│   └── config/config.go            # Configuration management
├── templates/
│   ├── benchmarks.md.tmpl          # Default markdown template
│   └── simple.md.tmpl              # Alternative template
├── examples/                       # Example usage and configs
├── README.md
└── go.mod
```

### Target Repository File Structure
```
your-repo/
├── internal/scope/
│   ├── benchmark_test.go
│   ├── .benchmarks.yaml          # Data file (committed)
│   └── benchmarks.md             # Generated (can be committed or gitignored)
├── internal/query/
│   ├── benchmark_test.go
│   ├── .benchmarks.yaml
│   └── benchmarks.md
└── .benchsync.yaml              # Global config (optional)
```

---

## 💾 Data Structure & Storage

### Benchmark Data File (YAML/JSON)
```yaml
# internal/scope/.benchmarks.yaml
package: "github.com/bold-minds/tvzr/internal/scope"
last_updated: "2025-08-09"
benchmarks:
  BenchmarkNewKey_Serial-20:
    history:
      - date: "2025-06-07"
        ns_op: 441.7
        b_op: 208
        allocs_op: 3
      - date: "2025-06-13"
        ns_op: 88.44
        b_op: 48
        allocs_op: 2
  BenchmarkNewKey-20:
    history:
      - date: "2025-06-07"
        ns_op: 120.8
        b_op: 208
        allocs_op: 3
      - date: "2025-06-13"
        ns_op: 40.97
        b_op: 48
        allocs_op: 2

# Manual content (preserved across regenerations)
manual_sections:
  optimization_steps:
    - "Replaced regex validation with manual ASCII validation"
    - "Inlined segment validation logic"
    - "Improved pooling with sync.Pool"
  
  future_improvements:
    - "Investigate further reduction of allocations for deeply nested keys"
    - "Profile for possible inlining"
  
  notes: |
    These optimizations focused on reducing query construction overhead...
```

### Go Data Structures
```go
type BenchmarkData struct {
    Package        string                    `yaml:"package"`
    LastUpdated    string                    `yaml:"last_updated"`
    Benchmarks     map[string]BenchmarkInfo  `yaml:"benchmarks"`
    ManualSections ManualContent             `yaml:"manual_sections"`
}

type BenchmarkInfo struct {
    History []BenchmarkResult `yaml:"history"`
}

type BenchmarkResult struct {
    Date     string  `yaml:"date"`
    NsOp     float64 `yaml:"ns_op"`
    BOp      int     `yaml:"b_op"`
    AllocsOp int     `yaml:"allocs_op"`
}

type ManualContent struct {
    OptimizationSteps   []string `yaml:"optimization_steps"`
    FutureImprovements  []string `yaml:"future_improvements"`
    Notes              string   `yaml:"notes"`
}
```

---

## 🎨 Template System

### Default Markdown Template
```go
// templates/benchmarks.md.tmpl
# {{ .Package }} Benchmarks

## 📊 Latest Benchmarks

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
{{- range .LatestResults }}
| {{ .Name }} | {{ .NsOp }} | {{ .BOp }} | {{ .AllocsOp }} |
{{- end }}

## 🔁 Historical Comparisons

{{- range .Benchmarks }}
### {{ .Name }}
| Date | ns/op | B/op | allocs/op | Δ ns/op | Δ B/op | Δ allocs |
|------|-------|------|-----------|---------|--------|----------|
{{- range .HistoryWithDeltas }}
| {{ .Date }} | {{ .NsOp }} | {{ .BOp }} | {{ .AllocsOp }} | {{ .DeltaNs }} | {{ .DeltaB }} | {{ .DeltaAllocs }} |
{{- end }}
{{- end }}

## ✅ Optimization Steps Applied
{{- range .ManualSections.OptimizationSteps }}
- {{ . }}
{{- end }}

## 🛠 Suggested Future Improvements
{{- range .ManualSections.FutureImprovements }}
- {{ . }}
{{- end }}

## 📋 Notes
{{ .ManualSections.Notes }}

## 🧪 How to Run

To run the {{ .PackageName }} benchmarks:

```sh
go test -benchmem -run=^$ -bench ^Benchmark {{ .Package }}
```

## 📆 Last Updated
{{ .LastUpdated }}
```

---

## ⚙️ Configuration System

### Global Config: `~/.benchsync.yaml`
```yaml
default_template: "benchmarks.md.tmpl"
date_format: "2006-01-02"
benchmark_flags: ["-benchmem", "-count=1"]
output_format: "table"
```

### Per-Repo Config: `.benchsync.yaml`
```yaml
# Override global settings for this repo
packages:
  internal/scope:
    template: "custom-scope.md.tmpl"
    manual_sections:
      optimization_steps:
        - "Replaced regex validation with manual ASCII validation"
      notes: |
        These optimizations focused on reducing query construction...
```

---

## 🔧 Core Implementation Components

### Smart Package Discovery
```go
func FindBenchmarkPackages(target string, recursive bool) ([]string, error) {
    switch {
    case target == "all":
        return findAllBenchmarkPackages(".")
    case strings.HasSuffix(target, "/..."):
        return findBenchmarkPackages(strings.TrimSuffix(target, "/..."), true)
    case strings.HasPrefix(target, "Benchmark"):
        return findPackageForBenchmark(target)
    default:
        return []string{target}, nil
    }
}
```

### Flexible Benchmark Runner
```go
type BenchmarkRunner struct {
    UpdateData bool  // false for "run", true for "sync"
    OutputFormat string
    Flags []string
}

func (br *BenchmarkRunner) Execute(packages []string) (*Results, error) {
    results := br.runBenchmarks(packages)
    
    if br.UpdateData {
        br.updateDataFiles(results)
        br.regenerateMarkdown()
    }
    
    return results, nil
}
```

### Data Management
```go
func (bd *BenchmarkData) AddResults(newResults []BenchmarkResult) {
    today := time.Now().Format("2006-01-02")
    for _, result := range newResults {
        if benchmark, exists := bd.Benchmarks[result.Name]; exists {
            benchmark.History = append(benchmark.History, result)
        } else {
            bd.Benchmarks[result.Name] = BenchmarkInfo{
                History: []BenchmarkResult{result},
            }
        }
    }
    bd.LastUpdated = today
}
```

### Template Rendering
```go
func RenderMarkdown(data *BenchmarkData, templatePath string) (string, error) {
    // Calculate deltas for template
    enrichedData := data.WithDeltas()
    
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, enrichedData)
    return buf.String(), err
}
```

---

## 🎯 Key Benefits

### 1. **Robust Data Persistence**
- ✅ Structured YAML/JSON files (no fragile markdown parsing)
- ✅ Version-controlled with the repository
- ✅ Clean git diffs when benchmark data changes
- ✅ Easy to edit manually if needed

### 2. **Template Flexibility**
- ✅ Easy to change markdown format without touching data
- ✅ Consistent formatting across all repositories
- ✅ Custom templates for specific needs
- ✅ Manual content preservation (notes, optimization steps)

### 3. **Simple Delta Calculations**
- ✅ Compare last two entries in history array
- ✅ Automatic percentage calculations
- ✅ Visual indicators (↑/↓ arrows)
- ✅ No complex parsing required

### 4. **Repo Portability**
- ✅ Single source of truth (separate repository)
- ✅ Install once, use everywhere
- ✅ No Makefiles or complex setup
- ✅ Works with any Go repository structure

### 5. **User-Friendly Commands**
- ✅ Simpler than complex `go test` invocations
- ✅ Flexible targeting (all/package/single test/recursive)
- ✅ Separation of concerns (run vs sync vs render)
- ✅ Multiple output formats

---

## 🚀 Development Roadmap

### Phase 1: MVP
- [ ] Basic CLI structure with cobra
- [ ] Simple benchmark runner (go test execution)
- [ ] YAML data file management
- [ ] Basic template rendering
- [ ] Core commands: run, sync, render

### Phase 2: Discovery & Targeting
- [ ] Smart package discovery (all/package/test/recursive)
- [ ] Benchmark test file detection
- [ ] Target validation and error handling
- [ ] Dry-run functionality

### Phase 3: Templates & Configuration
- [ ] Template system with default templates
- [ ] Global and per-repo configuration
- [ ] Custom template support
- [ ] Manual content preservation

### Phase 4: Polish & Features
- [ ] Multiple output formats (table, json, raw)
- [ ] Verbose logging and debugging
- [ ] Error handling and validation
- [ ] Performance optimizations

### Phase 5: Documentation & Examples
- [ ] Comprehensive README
- [ ] Usage examples and best practices
- [ ] Template customization guide
- [ ] Integration with CI/CD workflows

---

## 📋 Requirements Summary

### Core Requirements
- ✅ **No AI integration** - Pure algorithmic approach
- ✅ **Repo-agnostic** - Works with any Go repository
- ✅ **Single source of truth** - Separate repository, install globally
- ✅ **No Makefiles** - Pure Go workflow
- ✅ **Persistent data** - Committed YAML/JSON files
- ✅ **Template-based** - No fragile markdown parsing

### CLI Requirements
- ✅ **Flexible targeting**: all/package/single test/recursive
- ✅ **Run without updates** - Simpler benchmark execution
- ✅ **Separation of concerns** - run vs sync vs render
- ✅ **Multiple output formats** - table, json, raw
- ✅ **Configuration support** - Global and per-repo

### Data Requirements
- ✅ **Historical tracking** - Performance over time
- ✅ **Delta calculations** - Automatic performance comparisons
- ✅ **Manual content preservation** - Notes, optimization steps
- ✅ **Git-friendly** - Clean diffs, version control

---

This specification provides a complete blueprint for building a robust, maintainable benchmark management tool that solves the real-world problems of Go benchmark documentation and tracking.
