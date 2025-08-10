# Contributing to ex

Thank you for your interest in contributing to the `ex` Go exception library! 🎉

We welcome contributions that help improve the library while maintaining its focus on simplicity, performance, and Go idioms.

## 🚀 Getting Started

### Prerequisites

- **Go 1.19+** - Required for development
- **Git** - For version control
- **golangci-lint** (optional) - For comprehensive linting

### Development Setup

1. **Fork and clone the repository**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ex.git
   cd ex
   ```

2. **Run the validation pipeline**:
   ```bash
   ./scripts/validate.sh
   ```

## 🎯 What We're Looking For

We welcome contributions in these areas:

### ✅ **Encouraged Contributions**
- **Bug fixes** - Fix issues or edge cases
- **Performance improvements** - Optimize without breaking compatibility
- **Test enhancements** - Add test cases, improve coverage
- **Documentation improvements** - Clarify usage, add examples
- **Validation pipeline improvements** - Enhance CI/CD and tooling

### ⚠️ **Requires Discussion**
- **New error types** - Additional `ExType` constants
- **API changes** - Modifications to public interfaces
- **New dependencies** - Adding external packages
- **Breaking changes** - Changes that affect backward compatibility

### ❌ **Not Accepted**
- **Feature creep** - Complex features that don't align with Go idioms
- **Non-idiomatic Go** - Code that doesn't follow Go conventions
- **Performance regressions** - Changes that significantly slow down the library

## 📋 Contribution Process

### 1. **Create an Issue First**
For significant changes, please create an issue to discuss:
- What problem you're solving
- Your proposed approach
- Any potential breaking changes

### 2. **Development Workflow**

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow the code style guidelines below
   - Add tests for new functionality
   - Update documentation as needed

3. **Validate your changes**:
   ```bash
   # Run the full validation pipeline
   ./scripts/validate.sh
   
   # Or run individual checks
   go fmt ./...
   go vet ./...
   go test -race ./...
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a pull request**:
   ```bash
   git push origin feature/your-feature-name
   ```

### 3. **Pull Request Guidelines**

Your PR should:
- **Have a clear title** describing the change
- **Reference any related issues** using `Fixes #123` or `Closes #123`
- **Include tests** for new functionality
- **Pass all validation checks** (CI will verify this)
- **Maintain backward compatibility** unless discussed otherwise

## 🎨 Code Style Guidelines

### Go Conventions
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Write clear, concise comments for public APIs
- Follow Go's error handling patterns

### Testing Standards
- Write table-driven tests where appropriate
- Test both success and error cases
- Include edge cases (nil values, empty strings, etc.)
- Use `testify/assert` for assertions (already included)

### Example Test Structure
```go
func TestNewException(t *testing.T) {
    tests := []struct {
        name     string
        code     ex.ExType
        id       int
        message  string
        expected string
    }{
        {
            name:     "basic exception",
            code:     ex.ExTypeIncorrectData,
            id:       400,
            message:  "test message",
            expected: "test message",
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            exc := ex.New(tt.code, tt.id, tt.message)
            assert.Equal(t, tt.expected, exc.Error())
        })
    }
}
```

### Documentation Standards
- Document all public functions and types
- Include usage examples in docstrings
- Use clear, concise language
- Follow Go's documentation conventions

## 🧪 Testing Requirements

### Test Coverage
- All new code must have tests
- Aim for high test coverage (we use 80% as minimum)
- Include both unit and integration tests

### Test Types
1. **Unit Tests** - Test individual functions and methods
2. **Integration Tests** - Test `errors.Is` and `errors.As` compatibility
3. **Edge Case Tests** - Test boundary conditions and error cases
4. **Race Tests** - Ensure thread safety (`go test -race`)

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run specific test
go test -run TestNewException
```

## 🔍 Code Review Process

### What We Look For
- **Correctness** - Does the code work as intended?
- **Performance** - Does it maintain or improve performance?
- **Style** - Does it follow Go conventions?
- **Tests** - Are there adequate tests?
- **Documentation** - Is it properly documented?
- **Compatibility** - Does it maintain backward compatibility?

### Review Timeline
- Initial review within 2-3 business days
- Follow-up reviews within 1-2 business days
- Merge after approval and passing CI

## 🚨 Reporting Issues

### Bug Reports
When reporting bugs, please include:
- **Go version** (`go version`)
- **Library version** or commit hash
- **Minimal reproduction case**
- **Expected vs actual behavior**
- **Error messages** (if any)

### Feature Requests
For feature requests, please describe:
- **Use case** - What problem does it solve?
- **Proposed API** - How would it work?
- **Alternatives considered** - Other approaches you've thought of
- **Backward compatibility** - Any breaking changes?

## 📝 Commit Message Guidelines

We follow conventional commits for clear history:

```
type(scope): description

[optional body]

[optional footer]
```

### Types
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Test additions or changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

### Examples
```
feat: add ExType.String() method for debugging

fix: handle nil inner errors correctly in Unwrap()

docs: add error chaining examples to README

test: add edge cases for WithInnerError method
```

## 🏆 Recognition

Contributors will be:
- Listed in release notes for their contributions
- Mentioned in the README contributors section
- Given credit in commit messages and PR descriptions

## 📞 Getting Help

If you need help or have questions:
- **Create an issue** for bugs or feature discussions
- **Check existing issues** to see if your question has been answered
- **Review the README.md** for usage examples and API documentation

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the same MIT License that covers the project.

---

Thank you for helping make `ex` better! 🙏
