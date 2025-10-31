# Kagi CLI - Project Implementation Guide

**Project:** Kagi FastGPT Command Line Interface
**Repository:** github.com/grantcarthew/kagi
**License:** Mozilla Public License 2.0
**Language:** Go 1.22+
**Status:** Planning → Implementation

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Technical Stack](#technical-stack)
3. [Prerequisites](#prerequisites)
4. [Project Structure](#project-structure)
5. [Implementation Phases](#implementation-phases)
6. [Testing Strategy](#testing-strategy)
7. [Quality Checklist](#quality-checklist)
8. [Release Process](#release-process)

---

## Project Overview

### Goal

Build a production-ready CLI tool for querying Kagi's FastGPT API with clean output suitable for both human users and AI agents.

### Success Criteria

- ✅ Simple, intuitive interface
- ✅ Robust error handling
- ✅ Multiple output formats (text, markdown, JSON)
- ✅ TTY-aware color output
- ✅ Comprehensive tests
- ✅ Easy distribution (Homebrew + go install)

### Design Authority

All implementation decisions are documented in `design-record.md`. Refer to that document for:

- Feature specifications
- Flag behavior
- Error handling requirements
- Output format details
- Configuration precedence

---

## Technical Stack

### Language & Version

- **Go:** 1.22+ (tested on Go 1.25.3)
- **Platform:** Cross-platform (Linux, macOS, Windows)

### Dependencies

| Package                | Purpose         | Version       |
| ---------------------- | --------------- | ------------- |
| github.com/spf13/cobra | CLI framework   | Latest stable |
| golang.org/x/term      | TTY detection   | Latest stable |
| Standard library       | HTTP, JSON, I/O | Stdlib        |

### Development Tools

- `go test` - Testing
- `go build` - Building
- `go mod` - Dependency management
- `goreleaser` - Release automation (optional)

---

## Prerequisites

### Required

- Go 1.22 or higher
- Git
- GitHub account (for releases)

### Optional

- Homebrew (for tap testing)
- goreleaser (for automated releases)

### Setup

```bash
# Verify Go version
go version  # Should be 1.22+

# Clone repository (already done)
cd /Users/gcarthew/Projects/kagi

# Initialize Go module (if not done)
go mod init github.com/grantcarthew/kagi

# Install dependencies
go mod tidy
```

---

## Project Structure

### Directory Layout

```
kagi/
├── cmd/
│   └── root.go                 # Root command definition
├── internal/
│   ├── api/
│   │   ├── client.go          # HTTP client & API interaction
│   │   ├── client_test.go     # API client tests
│   │   └── types.go           # API request/response types
│   ├── config/
│   │   ├── config.go          # Configuration management
│   │   └── config_test.go     # Config tests
│   ├── format/
│   │   ├── output.go          # Output formatting
│   │   ├── output_test.go     # Format tests
│   │   └── color.go           # Color/ANSI handling
│   └── input/
│       ├── query.go           # Query parsing (args + stdin)
│       └── query_test.go      # Query parsing tests
├── main.go                     # Entry point
├── go.mod                      # Module definition
├── go.sum                      # Dependency checksums
├── LICENSE                     # MPL 2.0 (already exists)
├── README.md                   # User documentation
├── design-record.md            # Design decisions (already exists)
├── PROJECT.md                  # This file
├── .gitignore                  # Git ignore rules
├── .goreleaser.yml             # Release configuration (optional)
└── dist/                       # Build artifacts (gitignored)
```

### Package Responsibilities

**cmd/root.go**

- Command definition with Cobra
- Flag setup and parsing
- Orchestration of query → API → output flow
- Signal handling (SIGINT)

**internal/api/client.go**

- HTTP client creation
- API request construction
- API response parsing
- Error response handling

**internal/api/types.go**

- Request/response struct definitions
- JSON marshaling tags

**internal/config/config.go**

- Configuration precedence logic
- Environment variable reading
- Validation

**internal/format/output.go**

- Text formatting
- Markdown formatting
- JSON formatting
- Reference section formatting

**internal/format/color.go**

- ANSI color code constants
- TTY detection
- Color stripping

**internal/input/query.go**

- Argument concatenation
- Stdin reading
- Query validation

---

## Implementation Phases

### Phase 1: Project Setup & Scaffolding

**Objective:** Establish project structure and dependencies.

**Tasks:**

- [x] Initialize Git repository (done)
- [x] Add LICENSE file (done)
- [ ] Create `go.mod` with module path
- [ ] Add `.gitignore` for Go projects
- [ ] Create directory structure (cmd/, internal/)
- [ ] Install Cobra dependency: `go get github.com/spf13/cobra`
- [ ] Install term package: `go get golang.org/x/term`
- [ ] Create stub files (main.go, cmd/root.go)
- [ ] Verify project compiles: `go build`

**Deliverable:** Compiling Go project with structure in place.

---

### Phase 2: Core API Client

**Objective:** Implement Kagi API integration.

**Tasks:**

- [ ] Define API request/response types (`internal/api/types.go`)
  - [ ] `FastGPTRequest` struct
  - [ ] `FastGPTResponse` struct
  - [ ] `FastGPTError` struct
  - [ ] `Reference` struct
- [ ] Implement API client (`internal/api/client.go`)
  - [ ] Create HTTP client with timeout
  - [ ] Build POST request with headers
  - [ ] Marshal request JSON
  - [ ] Parse success response
  - [ ] Parse error response
  - [ ] Handle HTTP errors
  - [ ] Handle network errors

**Manual Testing:**
- Test with real Kagi API key
- Verify successful query returns data
- Test timeout behavior
- Test error responses (invalid key, etc.)

**Deliverable:** Working API client package.

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 3: CLI Framework & Flags

**Objective:** Implement command structure and flag parsing.

**Tasks:**

- [ ] Set up Cobra root command (`cmd/root.go`)
  - [ ] Define command description
  - [ ] Configure flag parsing (allow args anywhere)
  - [ ] Define version constant (e.g., `const version = "1.0.0"`)
- [ ] Add flags with validation
  - [ ] `--help, -h` (built-in)
  - [ ] `--version, -v` with version output
  - [ ] `--api-key` (string)
  - [ ] `--format, -f` (string, validate: text/txt/md/markdown/json)
  - [ ] `--heading` (boolean)
  - [ ] `--quiet, -q` (boolean)
  - [ ] `--timeout, -t` (int, validate: positive)
  - [ ] `--color, -c` (string, validate: auto/always/never)
  - [ ] `--verbose` (boolean)
  - [ ] `--debug` (boolean)
- [ ] Implement configuration loading (`internal/config/config.go`)
  - [ ] Read `KAGI_API_KEY` environment variable
  - [ ] Apply flag precedence
  - [ ] Validate configuration
  - [ ] Return errors for missing/invalid config
- [ ] Implement query input (`internal/input/query.go`)
  - [ ] Concatenate args with spaces
  - [ ] Read from stdin if no args
  - [ ] Validate query not empty
  - [ ] Trim whitespace

**Manual Testing:**
- Test all flags: `go run . --help`, `go run . --version`
- Test query parsing: `go run . test query`, `go run . "quoted query"`
- Test stdin: `echo "test" | go run .`
- Test config precedence: `KAGI_API_KEY=env go run . --api-key flag test`
- Test validation: missing API key, empty query, invalid flag values

**Deliverable:** CLI accepting all flags and parsing queries correctly.

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 4: Output Formatting

**Objective:** Implement all output formats with color support.

**Tasks:**

- [ ] Implement color support (`internal/format/color.go`)
  - [ ] Define ANSI color constants
  - [ ] Implement TTY detection using `term.IsTerminal()`
  - [ ] Color application logic (auto/always/never)
  - [ ] Color stripping function
- [ ] Implement text format (`internal/format/output.go`)
  - [ ] Format output body
  - [ ] Format references section (numbered list)
  - [ ] Add heading if `--heading` flag set
  - [ ] Apply colors if enabled
- [ ] Implement markdown format
  - [ ] Format heading as `# query`
  - [ ] Format output body
  - [ ] Format references as markdown links with blockquotes
- [ ] Implement JSON format
  - [ ] Pretty-print full API response
  - [ ] Handle `--quiet` (output field only)
- [ ] Implement quiet mode for all formats
  - [ ] Text: just output body
  - [ ] Markdown: just output body
  - [ ] JSON: just `.data.output` as JSON string

**Manual Testing:**
- Test text format: `go run . test query`
- Test markdown: `go run . -f md test query`
- Test JSON: `go run . -f json test query`
- Test heading: `go run . --heading test query`
- Test quiet: `go run . -q test query`
- Test colors in terminal vs pipe: `go run . test | cat` (no color)
- Test color flags: `--color always`, `--color never`

**Deliverable:** Complete output formatting for all modes.

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 5: Error Handling

**Objective:** Implement comprehensive error handling per design spec.

**Tasks:**

- [ ] Define error types and messages
  - [ ] Missing API key error
  - [ ] Missing query error
  - [ ] Invalid flag value errors
  - [ ] API error responses
  - [ ] Network errors
  - [ ] Timeout errors
  - [ ] JSON parse errors
- [ ] Implement error formatting
  - [ ] Prefix with "Error: "
  - [ ] Include actionable hints
  - [ ] Output to stderr
- [ ] Implement verbosity levels
  - [ ] Default: silent (errors only)
  - [ ] `--verbose`: process info to stderr
  - [ ] `--debug`: verbose + detailed debug info
- [ ] Set appropriate exit codes
  - [ ] 0 for success
  - [ ] 1 for runtime errors
  - [ ] 2 for usage errors
  - [ ] 130 for SIGINT
- [ ] Handle signals
  - [ ] Graceful SIGINT (Ctrl+C) handling
  - [ ] Cleanup on exit

**Manual Testing:**
- Test each error scenario from design-record.md
- Test `--verbose` output
- Test `--debug` output
- Test Ctrl+C handling
- Verify error messages go to stderr: `go run . 2>&1 | grep Error`
- Verify exit codes: `go run . invalid; echo $?`

**Deliverable:** Robust error handling with clear user feedback.

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 6: Main Integration

**Objective:** Wire everything together in main execution flow.

**Tasks:**

- [ ] Implement main execution logic (`cmd/root.go`)
  - [ ] Load configuration
  - [ ] Get query (args or stdin)
  - [ ] Validate inputs
  - [ ] Create API client
  - [ ] Make API request
  - [ ] Format output
  - [ ] Write to stdout
  - [ ] Handle errors appropriately
- [ ] Implement version command
  - [ ] Standard: version + repo + issues URL
  - [ ] Quiet: version number only
- [ ] Implement main.go entry point
  - [ ] Execute root command
  - [ ] Exit with proper code

**Manual Testing:**
- [ ] Test all flag combinations
- [ ] Test with real Kagi API
- [ ] Test piping: `go run . test | less`
- [ ] Test redirects: `go run . test > output.txt`
- [ ] Test stdin: `echo "query" | go run .`
- [ ] Test version: `go run . --version`, `go run . -v -q`
- [ ] Test in different terminals (iTerm, Terminal.app, etc.)
- [ ] Test color auto-detection

**Deliverable:** Fully functional CLI tool.

**Note:** Comprehensive tests will be written in Phase 7.

---

### Phase 7: Testing

**Objective:** Achieve comprehensive test coverage after core implementation is complete.

**Why After Implementation:**
Following KISS principles, we write tests after the core is working to:
- Focus on implementation without context switching
- Understand the full system before testing it
- Avoid testing code that might change during development
- Write better tests with complete understanding of edge cases

**Tasks:**

- [ ] Write unit tests for all packages (`*_test.go` files)
  - [ ] `internal/api/client_test.go`
    - [ ] Successful API calls (mocked HTTP)
    - [ ] API error responses
    - [ ] HTTP errors (404, 500, etc.)
    - [ ] Timeout handling
    - [ ] Invalid JSON responses
    - [ ] Network failures
  - [ ] `internal/config/config_test.go`
    - [ ] Configuration precedence (flags > env > defaults)
    - [ ] API key validation
    - [ ] Flag value validation
  - [ ] `internal/format/output_test.go`
    - [ ] Text format (with/without heading, quiet)
    - [ ] Markdown format (with/without quiet)
    - [ ] JSON format (with/without quiet)
    - [ ] Color application
    - [ ] TTY detection
    - [ ] Reference formatting
    - [ ] Empty references handling
  - [ ] `internal/input/query_test.go`
    - [ ] Argument concatenation
    - [ ] Stdin reading
    - [ ] Args + stdin precedence
    - [ ] Empty query validation
    - [ ] Whitespace handling

- [ ] Write integration tests
  - [ ] Mock API server for end-to-end testing
  - [ ] Test complete execution flows
  - [ ] Test all flag combinations
  - [ ] Test configuration precedence
  - [ ] Test error scenarios

- [ ] Test edge cases
  - [ ] Empty API responses
  - [ ] Missing references in response
  - [ ] Very long queries (>1000 chars)
  - [ ] Special characters in query
  - [ ] Unicode in query and output
  - [ ] Multiple spaces in query args

- [ ] Test error conditions
  - [ ] All 11 error scenarios from design-record.md
  - [ ] Verify error messages match spec
  - [ ] Verify exit codes (0, 1, 2, 130)
  - [ ] Verify stderr output

- [ ] Verify test coverage
  - [ ] Run `go test -cover ./...`
  - [ ] Aim for >80% coverage
  - [ ] Identify untested paths
  - [ ] Add tests for gaps

- [ ] Test on multiple platforms (if possible)
  - [ ] macOS (primary development platform)
  - [ ] Linux (via Docker or CI)
  - [ ] Windows (if accessible)

**Deliverable:** >80% test coverage with all tests passing.

**Note:** This completes the core implementation. Phases 8-9 cover documentation and distribution.

---

### Phase 8: Documentation

**Objective:** Create comprehensive user and developer documentation.

**Tasks:**

- [ ] Write README.md
  - [ ] Project description
  - [ ] Installation instructions (Homebrew + go install)
  - [ ] Quick start guide
  - [ ] Usage examples for all formats
  - [ ] Flag reference
  - [ ] Environment variable documentation
  - [ ] Examples with stdin
  - [ ] Troubleshooting section
  - [ ] Contributing guidelines
  - [ ] License information
- [ ] Add code documentation
  - [ ] Godoc comments for all public functions
  - [ ] Package documentation
  - [ ] Example code snippets
- [ ] Create usage examples
  - [ ] Basic query
  - [ ] With different formats
  - [ ] With flags
  - [ ] Piping and redirects
  - [ ] Error scenarios
- [ ] Update design-record.md if needed
  - [ ] Document any deviations
  - [ ] Add implementation notes

**Deliverable:** Complete documentation for users and developers.

---

### Phase 9: Distribution & Release

**Objective:** Set up automated releases and distribution channels.

**Tasks:**

- [ ] Create GitHub release workflow
  - [ ] Add `.goreleaser.yml` configuration
  - [ ] Configure build targets (Linux, macOS, Windows)
  - [ ] Set up GitHub Actions for releases
  - [ ] Test release process
- [ ] Create Homebrew tap
  - [ ] Create `homebrew-tap` repository
  - [ ] Write formula for kagi
  - [ ] Test installation via brew
  - [ ] Document tap installation
- [ ] Tag first release
  - [ ] Create annotated tag: `git tag -a v1.0.0 -m "Initial release"`
  - [ ] Push tag: `git push origin v1.0.0`
  - [ ] Verify GitHub release created
  - [ ] Download and test binaries
- [ ] Document installation methods
  - [ ] Homebrew: `brew install grantcarthew/tap/kagi`
  - [ ] Go install: `go install github.com/grantcarthew/kagi@latest`
  - [ ] Direct download from releases
- [ ] Announce release
  - [ ] Update README with installation instructions
  - [ ] Share with intended audience

**Deliverable:** v1.0.0 release with multiple distribution methods.

---

## Testing Strategy

### Approach: Build First, Test After

Following KISS principles, testing is done in Phase 7 after core implementation (Phases 1-6) is complete.

**Why This Approach:**

- **Focus:** Build features without context-switching to test writing
- **Understanding:** Better tests when you understand the complete system
- **Efficiency:** Avoid testing code that changes during implementation
- **Simplicity:** No test infrastructure to maintain during development

### During Implementation (Phases 1-6)

**Manual testing only:**

- Use `go run . <args>` to verify features work
- Test with real Kagi API key
- Check error scenarios manually
- Verify output formats in terminal
- No `*_test.go` files created yet

### After Implementation (Phase 7)

#### Unit Tests

**Location:** `*_test.go` files alongside implementation
**Command:** `go test ./...`
**Coverage:** `go test -cover ./...`

**Requirements:**

- Test all public functions
- Test error conditions
- Test edge cases
- Mock external dependencies (HTTP/API calls)
- Aim for >80% coverage

**Packages to test:**

- `internal/api` - API client with mocked HTTP
- `internal/config` - Configuration precedence and validation
- `internal/format` - Output formatting for all formats
- `internal/input` - Query parsing from args and stdin

#### Integration Tests

**Location:** `cmd/root_test.go` or separate `integration_test.go`
**Approach:** End-to-end testing with mocked API

**Scenarios:**

- Successful query with various flag combinations
- Error handling for all 11 error scenarios
- Output format validation (text, markdown, JSON)
- Configuration precedence (flags > env > defaults)
- Exit codes (0, 1, 2, 130)

#### Manual Testing Checklist

**Environment:** Real terminal with Kagi API key

**Final verification before release:**

- [ ] Basic query: `kagi golang best practices`
- [ ] All formats: text, markdown, JSON
- [ ] All flags: heading, quiet, timeout, color, verbose, debug
- [ ] Stdin input: `echo "test" | kagi`
- [ ] Environment variable: `KAGI_API_KEY=xxx kagi test`
- [ ] Flag override: `KAGI_API_KEY=env kagi --api-key flag test`
- [ ] Error: missing API key
- [ ] Error: missing query
- [ ] Error: invalid flag value
- [ ] Error: network timeout
- [ ] Color in terminal (should colorize)
- [ ] Color in pipe (should not colorize): `kagi test | cat`
- [ ] Version output: `kagi --version`
- [ ] Quiet version: `kagi -v -q`
- [ ] Help output: `kagi --help`
- [ ] Interrupt: Ctrl+C during query

---

## Quality Checklist

### Code Quality

- [ ] All code formatted with `gofmt`
- [ ] No linting errors: `golangci-lint run` (if available)
- [ ] All exported functions documented
- [ ] No hardcoded secrets or API keys
- [ ] Proper error wrapping with context
- [ ] No panics in production code

### Functionality

- [ ] All flags work as specified
- [ ] All output formats correct
- [ ] Color output works correctly
- [ ] TTY detection accurate
- [ ] Error messages helpful and actionable
- [ ] Exit codes correct for all scenarios
- [ ] Configuration precedence respected
- [ ] Timeout handling works
- [ ] Signal handling graceful

### Testing

- [ ] All tests pass: `go test ./...`
- [ ] Test coverage >80%
- [ ] Integration tests cover main flows
- [ ] Edge cases tested
- [ ] Error conditions tested

### Documentation

- [ ] README complete with examples
- [ ] All flags documented
- [ ] Installation instructions clear
- [ ] Code comments comprehensive
- [ ] design-record.md up to date

### Distribution

- [ ] Builds on Linux, macOS, Windows
- [ ] Binaries are properly sized (not bloated)
- [ ] GitHub releases working
- [ ] Homebrew formula tested
- [ ] `go install` works

---

## Release Process

### Pre-release Checklist

- [ ] All tests passing
- [ ] Documentation complete
- [ ] CHANGELOG.md updated (if exists)
- [ ] Version number decided (semver)
- [ ] No uncommitted changes

### Release Steps

1. **Update version references**

   ```bash
   # Update version constant in cmd/root.go
   # Change: const version = "1.0.0"
   # To:     const version = "1.1.0"
   ```

2. **Create annotated tag**

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0: Initial release"
   git push origin v1.0.0
   ```

3. **Trigger release** (if using goreleaser)

   ```bash
   goreleaser release --clean
   ```

   Or use GitHub Actions to auto-release on tag push.

4. **Verify release**

   - Check GitHub releases page
   - Download binaries and test
   - Verify checksums

5. **Update Homebrew tap** (if not automated)

   ```bash
   cd ../homebrew-tap
   # Update formula with new version and SHA256
   git commit -am "Update kagi to v1.0.0"
   git push
   ```

6. **Test installation**

   ```bash
   brew uninstall kagi
   brew install grantcarthew/tap/kagi
   kagi --version
   ```

7. **Announce**
   - Update README if needed
   - Notify users/stakeholders

### Version Numbering

**Semantic Versioning (semver):**

- **Major (x.0.0):** Breaking changes
- **Minor (1.x.0):** New features, backward compatible
- **Patch (1.0.x):** Bug fixes, backward compatible

**Examples:**

- v1.0.0 - Initial release
- v1.0.1 - Fix timeout bug
- v1.1.0 - Add new output format
- v2.0.0 - Change flag names (breaking)

---

## Development Workflow

### Daily Development

```bash
# Pull latest
git pull

# Create feature branch
git checkout -b feature/api-client

# Make changes
# ... code ...

# Format code
gofmt -w .

# Run manual tests during Phases 1-6
go run . test query

# Run automated tests (Phase 7+ only)
go test ./...

# Commit
git add .
git commit -m "Implement API client with timeout handling"

# Push
git push origin feature/api-client

# Create PR or merge to main
```

### Quick Commands

**Build:**

```bash
go build -o kagi
./kagi --version
```

**Test:**

```bash
go test ./...                    # All tests
go test -v ./internal/api        # Specific package
go test -cover ./...             # With coverage
go test -race ./...              # Race detection
```

**Run without building:**

```bash
go run . golang best practices
go run . --help
```

**Clean:**

```bash
go clean
rm -rf dist/
```

---

## Troubleshooting

### Common Issues

**Import errors:**

```bash
go mod tidy
```

**Dependency conflicts:**

```bash
go mod tidy
go mod verify
```

**Build failures:**

```bash
go clean
go build -v
```

**Test failures:**

```bash
go test -v ./...  # Verbose output
go test -run TestName  # Run specific test
```

---

## Success Metrics

### v1.0.0 Release Criteria

- [ ] All phases completed
- [ ] All tests passing
- [ ] Documentation complete
- [ ] Manual testing successful
- [ ] Installable via Homebrew
- [ ] Installable via go install
- [ ] Zero known critical bugs

### Post-Release

- Monitor GitHub issues for bugs
- Respond to user feedback
- Plan minor releases for enhancements
- Maintain backward compatibility

---

## Timeline Estimate

**Phase 1:** 1-2 hours (setup)
**Phase 2:** 3-4 hours (API client)
**Phase 3:** 3-4 hours (CLI framework & flags)
**Phase 4:** 3-4 hours (output formatting)
**Phase 5:** 2-3 hours (error handling)
**Phase 6:** 2-3 hours (integration)
**Phase 7:** 6-8 hours (comprehensive testing - all tests written here)
**Phase 8:** 3-4 hours (documentation)
**Phase 9:** 2-3 hours (distribution setup)

**Total: ~25-35 hours** (3-4 days of focused work)

**Note:** Testing is concentrated in Phase 7 following the KISS approach - build first, test after.

---

## Next Steps

1. **Review this document** with stakeholders
2. **Confirm design-record.md** specifications
3. **Begin Phase 1** (project setup)
4. **Track progress** through phases
5. **Iterate** based on testing feedback

---

## Conclusion

This project guide provides a structured approach to implementing the Kagi CLI tool. Follow the phases sequentially, completing all tasks and tests before moving forward. Refer to `design-record.md` for all design decisions and specifications.

**Ready to build!** 🚀
