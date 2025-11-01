# Contributing to kagi

Thank you for your interest in contributing to kagi! This document provides guidelines and instructions for contributing.

## Quick Links

- [Issues](https://github.com/grantcarthew/kagi/issues) - Report bugs or request features
- [Pull Requests](https://github.com/grantcarthew/kagi/pulls) - Submit code changes
- [AGENTS.md](AGENTS.md) - Detailed technical documentation for AI agents and developers

## Ways to Contribute

- Report bugs
- Suggest new features or improvements
- Improve documentation
- Submit pull requests with bug fixes or features
- Help answer questions in issues

## Reporting Bugs

When reporting bugs, please:

1. Check [existing issues](https://github.com/grantcarthew/kagi/issues) first
2. Use the bug report template when creating a new issue
3. Provide:
   - kagi version: `kagi --version`
   - Operating system and version
   - Full command that triggered the bug
   - Complete error message or unexpected output
   - Debug output (if applicable): `kagi --debug "your query"`

## Development Setup

### Prerequisites

- Go 1.25.3 or later
- Git
- Kagi API key (for running the tool)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/grantcarthew/kagi.git
cd kagi

# Install dependencies
go mod download

# Build the project
go build

# Set up API key (required for running)
export KAGI_API_KEY='your-api-key-here'

# Verify installation
./kagi --version
```

### Project Structure

kagi follows a simple, single-file architecture (KISS principle):

- `main.go` - All application code (types, API client, CLI, formatting, utilities)
- `main_test.go` - Comprehensive test suite (83 tests)
- `go.mod` - Go module definition
- `go.sum` - Dependency checksums

This intentional simplicity keeps the codebase easy to understand and maintain.

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestValidateQuery

# Run with race detector
go test -race ./...
```

**Test requirements:**
- Tests use `httptest.NewServer` for mocking API responses
- No real API key needed for tests
- All tests should pass before submitting PRs

## Code Style

### Go Conventions

- Follow standard Go formatting: run `gofmt` on all code
- Use `go vet` to catch common mistakes
- Follow `golint` recommendations
- Use descriptive variable and function names
- Document all exported functions, types, and constants

### Project-Specific Standards

**Single-file architecture:**
- All code stays in `main.go` (KISS principle)
- Organize code logically: constants, types, main/CLI, API, formatting, utilities

**Naming conventions:**
- Exported: `PascalCase`
- Unexported: `camelCase`
- Constants: `FormatMarkdown`, `FormatText`, `FormatJSON`

**Error handling:**
- Return errors explicitly, don't panic
- Provide clear, actionable error messages
- Start error messages with "Error: " prefix
- Include suggestions when possible

**Exit codes:**
- `0`: Success
- `1`: Error (API, network, validation)
- `130`: Interrupted (Ctrl+C)

**Dependencies:**
- Minimize external dependencies
- Current dependencies:
  - `github.com/spf13/cobra` (CLI framework)
  - `golang.org/x/term` (terminal detection)
  - Standard library only otherwise

### Code Quality

Before submitting:

```bash
# Format code
gofmt -w .

# Vet code
go vet ./...

# Run tests
go test -v ./...
```

## Branch Management

- Main branch: `main`
- Feature branches: `feature/description`
- Bug fix branches: `fix/description`
- Always work on feature branches, PR to `main`

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) style:

```
feat(api): add caching support for API responses
fix(format): handle empty references in markdown output
docs: update installation instructions
chore(deps): update cobra to v1.8.0
test(api): add timeout handling tests
refactor(output): simplify color detection logic
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `chore` - Maintenance tasks, dependencies
- `test` - Test additions or changes
- `refactor` - Code restructuring without behavior change

**Format:**
```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

## Pull Request Process

1. **Create an issue first** (for significant changes)
2. **Fork the repository** and create a feature branch
3. **Make your changes:**
   - Write clear, focused commits
   - Follow code style guidelines
   - Add/update tests as needed
   - Update documentation if changing CLI interface
4. **Test thoroughly:**
   - Run all tests: `go test -v ./...`
   - Build successfully: `go build`
   - Test manually with various queries
   - Verify all formats work: text, markdown, JSON
5. **Submit pull request:**
   - Use the PR template
   - Reference related issue(s)
   - Describe what changed and why
   - Keep PRs focused on a single feature/fix
6. **Respond to review feedback**

### Pull Request Guidelines

- Run tests and build before submitting
- Ensure code is formatted with `gofmt`
- Update relevant documentation
- Keep PRs focused on a single feature or fix
- Add tests for new functionality
- Maintain or improve test coverage

## Documentation

When changing functionality:

- Update `README.md` for user-facing changes
- Update `AGENTS.md` for technical implementation details
- Update CLI help text if adding/changing flags
- Add design decisions to `docs/design-record.md` for architectural changes

## Testing Philosophy

- **Table-driven tests**: Use for multiple scenarios
- **Mock HTTP responses**: Use `httptest.NewServer` for API tests
- **Test both paths**: Success and error cases
- **Edge cases**: Empty input, timeouts, malformed responses
- Focus on behavior, not implementation details

Example test pattern:

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "expected", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := functionUnderTest(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
            }
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Common Tasks

### Adding a New Output Format

1. Add format constant to const block in main.go
2. Add format to validation in `validateFormat()`
3. Implement formatting function (e.g., `formatYAML()`)
4. Add format to switch statement in main execution
5. Update help text and README.md
6. Add comprehensive tests to main_test.go

### Adding a New CLI Flag

1. Add flag to root command in `init()` or command definition
2. Add corresponding variable (global or in struct)
3. Add validation if needed
4. Update help template
5. Document in README.md
6. Add tests for flag parsing and behavior

### Modifying API Request

1. Update `FastGPTRequest` struct if adding fields
2. Modify request creation in API call function
3. Update tests with new request structure
4. Document changes in docs/design-record.md
5. Test with actual API (requires valid key)

## Questions or Need Help?

- Open an issue for questions
- Check [AGENTS.md](AGENTS.md) for detailed technical documentation
- Review existing code and tests for examples

## License

By contributing to kagi, you agree that your contributions will be licensed under the [Mozilla Public License 2.0](LICENSE).

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). Be respectful, professional, and constructive in all interactions. We welcome contributors of all experience levels.
