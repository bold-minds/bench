# Contributing to bench

Thank you for your interest in contributing to bench!

## Getting Started

### Prerequisites

- **Go 1.22+**
- **Git**
- **golangci-lint** (optional, for linting)

### Development Setup

1. **Fork and clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/bench.git
   cd bench
   ```

2. **Run tests**:
   ```bash
   go test -race ./...
   ```

3. **Run linter**:
   ```bash
   golangci-lint run
   ```

## What We're Looking For

### Encouraged
- Bug fixes
- Performance improvements
- Test enhancements
- Documentation improvements

### Requires Discussion
- New commands or flags
- API changes
- New dependencies
- Breaking changes

## Development Workflow

1. Create a feature branch
2. Make changes with tests
3. Run `go test -race ./...` and `golangci-lint run`
4. Commit with conventional commit messages
5. Open a pull request

## Commit Messages

```
type(scope): description

feat: add JSON output format for render command
fix: handle missing go.mod gracefully
docs: update CLI usage examples
test: add edge cases for parser
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
