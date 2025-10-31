# AGENTS.md

This is a command-line interface for querying the Kagi FastGPT API, built in Go following the KISS (Keep It Simple, Stupid) principle. All code lives in a single `main.go` file for simplicity.

## Project Overview

A fast, production-ready CLI tool that queries Kagi's FastGPT API and returns AI-powered search results with web context. Designed for both human users and AI agents with multiple output formats (text, markdown, JSON).

**Repository:** github.com/grantcarthew/kagi
**Language:** Go 1.25.3
**License:** Mozilla Public License 2.0
**Architecture:** Single-file design (main.go) with comprehensive test coverage

## Setup Commands

```bash
# Clone and build
git clone git@github.com:grantcarthew/kagi.git
cd kagi
go build

# Install dependencies (handled automatically by go build)
go mod download

# Set API key (required for running)
export KAGI_API_KEY='your-api-key-here'
```

## Build and Test Commands

```bash
# Build the binary
go build

# Build with version info
go build -ldflags="-X main.version=1.0.0"

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Verbose test output
go test -v ./...

# Run specific test
go test -run TestValidateQuery

# Check formatting
gofmt -l .

# Format code
gofmt -w .

# Run go vet
go vet ./...
```

## Code Style Guidelines

### Go Conventions

- **Formatting:** Use `gofmt` for all code (standard Go formatting)
- **Linting:** Follow `golint` recommendations
- **Naming:** Use camelCase for unexported, PascalCase for exported identifiers
- **Comments:** Document all exported functions, types, and constants
- **Error handling:** Return errors explicitly, don't panic except for unrecoverable situations

### Project-Specific Standards

1. **Single-file architecture:** All code stays in `main.go` (KISS principle)
2. **Constants at top:** Define all constants in the const block after imports
3. **Type definitions:** Group related types together
4. **Function organization:**
   - Main and CLI setup functions first
   - API client functions
   - Formatting/output functions
   - Utility functions last
5. **Test naming:** `Test<FunctionName>` for unit tests, `Test<Scenario>` for integration tests
6. **No external dependencies** except:
   - `github.com/spf13/cobra` (CLI framework)
   - `golang.org/x/term` (terminal detection)
   - Standard library

### Error Messages

- Start with "Error: " prefix
- Be specific and actionable
- Include suggestions when possible
- Example: `Error: API key required. Set KAGI_API_KEY environment variable or use --api-key flag`

### Exit Codes

- `0`: Success
- `1`: Error (API, network, validation)
- `130`: Interrupted (Ctrl+C)

## Testing Instructions

### Test Coverage

The project has comprehensive test coverage (83 tests in `main_test.go`):

- API request/response handling
- Input validation
- Output formatting (text, markdown, JSON)
- Color detection logic
- Error handling
- CLI flag parsing
- Signal handling

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage report
go test -cover ./...

# Run specific test suite
go test -run TestFormatting

# Run with race detector
go test -race ./...
```

### Adding New Tests

1. Add test functions to `main_test.go`
2. Use table-driven tests for multiple scenarios
3. Mock HTTP responses using `httptest.NewServer`
4. Test both success and error paths
5. Verify output format and content
6. Test edge cases (empty input, timeouts, etc.)

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

### Test Requirements

- All new features must include tests
- Maintain or improve test coverage
- Tests must pass before committing
- Run `go test -v ./...` to verify

## Development Workflow

### Branch Management

- Main branch: `main`
- Feature branches: `feature/description` or descriptive names
- Bug fixes: `fix/description`
- Keep branches short-lived

### Commit Message Format

Follow conventional commit style:

```
type(scope): brief description

Detailed explanation if needed

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

Examples:
- `feat: add markdown output format`
- `fix: handle timeout errors correctly`
- `docs: update installation instructions`
- `test: add coverage for JSON formatting`

### Pull Request Guidelines

1. Run all tests before creating PR
2. Update README.md if adding user-facing features
3. Update docs/design-record.md for architectural changes
4. Keep PRs focused on a single feature/fix
5. Include examples in PR description
6. Ensure CI passes (when configured)

## API Integration

### Kagi FastGPT API

- **Endpoint:** `https://kagi.com/api/v0/fastgpt`
- **Authentication:** Bearer token via `Authorization` header
- **Method:** POST
- **Content-Type:** `application/json`

### Request Structure

```json
{
  "query": "user query text",
  "web_search": true,
  "cache": true
}
```

### Response Structure

```json
{
  "meta": {
    "id": "request-id",
    "node": "server-node",
    "ms": 1234
  },
  "data": {
    "output": "AI response text",
    "references": [
      {
        "title": "Source Title",
        "snippet": "Excerpt from source",
        "url": "https://example.com"
      }
    ]
  }
}
```

### API Key Management

- Environment variable: `KAGI_API_KEY`
- CLI flag override: `--api-key`
- Never log or expose the full API key
- Mask in debug output: show only `***`

### Timeout Handling

- Default: 30 seconds
- Configurable via `--timeout` flag
- Cancel on Ctrl+C (signal handling)
- Context-based cancellation for clean shutdown

## Output Formats

### Text Format (default)

- Clean, readable output
- Numbered references at the end
- Optional heading with `--heading`
- Colored output in terminals (auto-detected)

### Markdown Format (`-f md`)

- H1 heading with query
- H2 "References" section
- Markdown links and blockquotes
- Suitable for documentation/README files

### JSON Format (`-f json`)

- Complete API response
- Parseable by `jq` and other tools
- No colors, heading, or formatting
- Ideal for automation and scripting

### Quiet Mode (`-q`)

- Output only the response body
- No heading, no references
- Works with all formats
- Perfect for piping and automation

## Color Output

### Detection Logic

1. Check `--color` flag: `always`, `never`, `auto`
2. Auto mode: detect if stdout is a terminal using `term.IsTerminal()`
3. Disable colors when piped or redirected

### Color Usage

- **Blue:** Headings and section titles
- **Cyan:** URLs and links
- **Yellow:** Reference numbers
- **Green:** Success messages (verbose mode)
- **Red:** Error messages

## Security Considerations

### API Key Protection

- Never commit API keys to version control
- Don't log full API keys (mask in debug output)
- Use environment variables or secure flag input
- Warn users about key exposure in scripts

### Input Validation

- Validate all user input before API calls
- Sanitize query strings (handled by JSON marshaling)
- Validate timeout values (must be positive)
- Check format values against allowed list

### HTTP Security

- Use HTTPS for all API calls
- Validate TLS certificates
- Set reasonable timeouts to prevent hanging
- Handle malformed responses gracefully

### Dependencies

- Minimal external dependencies (only Cobra and term)
- Keep dependencies updated
- Review dependency security advisories
- Use `go mod tidy` to clean unused dependencies

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

## Troubleshooting

### Build Issues

**Issue:** `go build` fails with dependency errors
**Solution:** Run `go mod tidy` then `go build`

**Issue:** Version not set in binary
**Solution:** Use build flag: `go build -ldflags="-X main.version=1.0.0"`

### Test Issues

**Issue:** Tests fail with "connection refused"
**Solution:** Tests use mock HTTP server, ensure no network required. Check test setup.

**Issue:** Tests fail on color output
**Solution:** Color tests may behave differently in CI. Check terminal detection logic.

### Runtime Issues

**Issue:** "API key required" error
**Solution:** Set `export KAGI_API_KEY='your-key'` or use `--api-key` flag

**Issue:** Timeout errors on complex queries
**Solution:** Increase timeout: `kagi --timeout 60 "your query"`

**Issue:** Colors broken in output
**Solution:** Use `--color never`

### Development Issues

**Issue:** Code formatting doesn't match style
**Solution:** Run `gofmt -w .` before committing

**Issue:** Tests pass locally but fail in PR
**Solution:** Ensure no environment-specific dependencies (check API key handling)

## Project Structure

```
kagi/
├── main.go              # All application code (types, API, CLI, formatting)
├── main_test.go         # Comprehensive test suite (83 tests)
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── LICENSE              # Mozilla Public License 2.0
├── README.md            # User-facing documentation
├── CHANGELOG.md         # Version history and changes
├── ROLE.md              # Role definition for AI agents
├── AGENTS.md            # This file
├── docs/
│   ├── design-record.md # Architectural decisions and specifications
│   └── tasks/           # Development task tracking
└── reference/
    └── homebrew-tap/    # Homebrew formula for distribution
```

## Distribution

### Homebrew

The project is distributed via Homebrew tap:

```bash
brew install grantcarthew/tap/kagi
```

Formula location: `reference/homebrew-tap/`

### Go Install

Direct installation via Go:

```bash
go install github.com/grantcarthew/kagi@latest
```

### Building from Source

```bash
git clone git@github.com:grantcarthew/kagi.git
cd kagi
go build
```

## Design Philosophy

This project follows the **KISS principle** (Keep It Simple, Stupid):

1. **Single-file architecture:** All code in `main.go` for simplicity
2. **Minimal dependencies:** Only essential external packages
3. **Clear structure:** Logical organization within the single file
4. **Comprehensive tests:** High coverage in a single test file
5. **Straightforward API:** Simple request/response with Kagi
6. **Clean output:** Multiple formats for different use cases

When making changes, always ask: "Is this adding necessary complexity, or can it be simpler?"

## Quick Reference

```bash
# Build
go build

# Test
go test -v ./...

# Format
gofmt -w .

# Run
./kagi "your query here"

# Run with options
./kagi -f md --heading "golang best practices"

# Debug
./kagi --debug "test query"
```

## Resources

- **Kagi API Documentation:** https://help.kagi.com/kagi/api/fastgpt.html
- **Cobra Documentation:** https://github.com/spf13/cobra
- **Go Style Guide:** https://go.dev/doc/effective_go
- **Project Design Record:** docs/design-record.md
- **GitHub Repository:** https://github.com/grantcarthew/kagi
