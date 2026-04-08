# bench

CLI tool for managing Go benchmark results. Runs benchmarks, tracks historical performance in YAML, and generates markdown reports.

## Install

```bash
go install github.com/bold-minds/bench/cmd/bench@latest
```

## Quick Start

```bash
cd your-go-repo/

# Initialize (creates .bench.yaml)
bench init

# Run benchmarks and see results
bench run all

# Run + persist to YAML + generate markdown
bench sync all

# Regenerate markdown from existing data
bench render all
```

## Commands

### `bench run [target]`

Run benchmarks and display results. No files are written.

```bash
bench run all                       # All packages with benchmarks
bench run internal/scope            # Specific package
bench run BenchmarkFoo              # Single benchmark
bench run ./... --recursive         # Recursive

bench run all --output table        # Pretty table (default)
bench run all --output json         # JSON
bench run all --output raw          # Raw go test output
```

### `bench sync [target]`

Run benchmarks, update `.benchmarks.yaml`, and regenerate `benchmarks.md`.

```bash
bench sync all
bench sync internal/scope
```

### `bench render [target]`

Regenerate markdown from existing `.benchmarks.yaml` files. Useful after manual edits.

```bash
bench render all
bench render internal/scope --template custom.md.tmpl
```

### `bench init`

Create a `.bench.yaml` config file with commented defaults.

## File Layout

After running `bench sync`, your repo will have:

```
your-repo/
  internal/scope/
    benchmark_test.go       # Your benchmarks
    .benchmarks.yaml        # Historical data (commit this)
    benchmarks.md           # Generated report
  .bench.yaml               # Config (optional)
```

### Data File (`.benchmarks.yaml`)

```yaml
package: github.com/you/repo/internal/scope
last_updated: 2025-08-09
benchmarks:
    BenchmarkFoo-8:
        - date: 2025-06-07
          ns_op: 441.7
          b_op: 208
          allocs_op: 3
        - date: 2025-06-13
          ns_op: 88.44
          b_op: 48
          allocs_op: 2
```

## Configuration

### `.bench.yaml`

```yaml
# Flags passed to go test
benchmark_flags:
  - "-benchmem"
  - "-count=1"

# Output format for 'run' command: table, json, raw
output_format: table

# Max history entries per benchmark (0 = unlimited)
history_cap: 20

# Per-package template overrides
packages:
  internal/scope:
    template: custom.md.tmpl
```

Config is loaded from (in order, later overrides earlier):
1. `~/.bench.yaml` (global defaults)
2. `.bench.yaml` (repo root)
3. `--config` flag (explicit override)

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Custom config file |
| `--verbose` | `-v` | Verbose output |
| `--dry-run` | | Show what would be done |
| `--recursive` | `-r` | Include subdirectories |
| `--output` | `-o` | Output format (run only) |
| `--template` | `-t` | Custom template (render only) |

## License

MIT
