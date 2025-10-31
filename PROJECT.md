# Kagi CLI - Project Implementation Guide

**Project:** Kagi FastGPT Command Line Interface
**Repository:** github.com/grantcarthew/kagi
**License:** Mozilla Public License 2.0
**Language:** Go 1.22+
**Status:** Core Complete (Phase 6/9) - Production-Ready CLI

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

### Directory Layout (Flat Structure - KISS)

```
kagi/
├── main.go                     # All code (types, API, CLI, formatting)
├── go.mod                      # Module definition
├── go.sum                      # Dependency checksums
├── LICENSE                     # MPL 2.0
├── README.md                   # User documentation (Phase 8)
├── design-record.md            # Design decisions
├── PROJECT.md                  # This file
├── ROLE.md                     # Role definition
├── .gitignore                  # Git ignore rules
├── reference/                  # Reference documentation
│   └── homebrew-tap/          # Homebrew tap example
├── .goreleaser.yml             # Release configuration (Phase 9)
├── kagi                        # Compiled binary (gitignored)
└── dist/                       # Build artifacts (gitignored)
```

**Rationale for Flat Structure:**
- KISS principle: project is <500 lines, no need for package separation
- Easier to navigate and maintain for a simple CLI tool
- All code in main.go: types, constants, API client, CLI, formatting
- Will split to handlers.go only if exceeding 1000 lines

### Code Organization in main.go

**Constants Section:**
- API configuration (endpoint, timeout)
- HTTP headers
- Exit codes
- Output formats
- Color modes
- Environment variables
- ANSI color codes (bold, blue, cyan, yellow, reset)

**Type Definitions:**
- `FastGPTRequest` - API request structure
- `FastGPTResponse` - API response structure
- `FastGPTError` - API error structure
- `Reference` - Reference structure
- `Config` - Application configuration

**Functions:**
- `main()` - Entry point, executes Cobra command
- `runCobra()` - Main command handler
- `loadConfig()` - Configuration loading with precedence
- `getQuery()` - Query extraction from args or stdin
- `queryKagi()` - API client function
- `formatOutput()` - Output format dispatcher
- `formatText_output()` - Text format with references
- `formatMarkdown_output()` - Markdown format with links
- `formatJSON_output()` - JSON format (full or quiet)
- `shouldUseColor()` - TTY detection for color output
- `colorize()` - Apply ANSI color codes
- `normalizeFormat()` - Format alias normalization
- `isValidFormat()` - Format validation

**Cobra Setup:**
- `rootCmd` - Cobra command definition
- `init()` - Flag definitions and help template
- Custom help template (not auto-generated)

---

## Implementation Phases

### Phase 1: Project Setup & Scaffolding ✅ COMPLETED

**Objective:** Establish project structure and dependencies.

**Tasks:**

- [x] Initialize Git repository (done)
- [x] Add LICENSE file (done)
- [x] Create `go.mod` with module path
- [x] Add `.gitignore` for Go projects (already exists)
- [x] ~~Create directory structure (cmd/, internal/)~~ Using flat structure - all code in main.go
- [x] Install Cobra dependency: `go get github.com/spf13/cobra`
- [x] Install term package: `go get golang.org/x/term`
- [x] Create stub files (main.go)
- [x] Verify project compiles: `go build`

**Deliverable:** Compiling Go project with structure in place.

**Implementation Notes:**
- Using flat structure (main.go only) per KISS principles
- Project is <1000 lines, no need for separate packages
- Dependencies installed: Cobra v1.10.1, term v0.36.0

---

### Phase 2: Core API Client ✅ COMPLETED

**Objective:** Implement Kagi API integration.

**Tasks:**

- [x] Define API request/response types (in main.go)
  - [x] `FastGPTRequest` struct
  - [x] `FastGPTResponse` struct
  - [x] `FastGPTError` struct
  - [x] `Reference` struct
- [x] Implement API client (queryKagi function in main.go)
  - [x] Create HTTP client with timeout
  - [x] Build POST request with headers
  - [x] Marshal request JSON
  - [x] Parse success response
  - [x] Parse error response
  - [x] Handle HTTP errors
  - [x] Handle network errors

**Manual Testing:**
- [x] Test with real Kagi API key
- [x] Verify successful query returns data
- [x] Test timeout behavior (context-based)
- [x] Test error responses (401/403 for invalid key, 429 for rate limit)

**Deliverable:** Working API client package.

**Implementation Notes:**
- All constants defined (API endpoint, headers, exit codes, etc.)
- Context-based timeout handling with proper error messages
- Specific error handling for common HTTP status codes
- Empty response validation
- API key sanitization in debug output (shows "***")

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 3: CLI Framework & Flags ✅ COMPLETED

**Objective:** Implement command structure and flag parsing.

**Tasks:**

- [x] Set up Cobra root command (in main.go)
  - [x] Define command description
  - [x] Configure flag parsing (allow args anywhere)
  - [x] Define version variable (set via ldflags: `var version = "dev"`)
- [x] Add flags with validation
  - [x] `--help, -h` (built-in)
  - [x] `--version, -v` with version output (rich or quiet)
  - [x] `--api-key` (string)
  - [x] `--format, -f` (string, validate: text/txt/md/markdown/json)
  - [x] `--heading` (boolean)
  - [x] `--quiet, -q` (boolean)
  - [x] `--timeout, -t` (int, validate: positive)
  - [x] `--color, -c` (string, validate: auto/always/never)
  - [x] `--verbose` (boolean)
  - [x] `--debug` (boolean, implies verbose)
- [x] Implement configuration loading (loadConfig function in main.go)
  - [x] Read `KAGI_API_KEY` environment variable
  - [x] Apply flag precedence (flags > env > defaults)
  - [x] Validate configuration
  - [x] Return errors for missing/invalid config
- [x] Implement query input (getQuery function in main.go)
  - [x] Concatenate args with spaces
  - [x] Read from stdin if no args (with TTY detection)
  - [x] Validate query not empty
  - [x] Trim whitespace

**Manual Testing:**
- [x] Test all flags: `./kagi --help`, `./kagi --version`, `./kagi -v -q`
- [x] Test query parsing: `./kagi test query`, `./kagi "quoted query"`
- [x] Test stdin: `echo "test" | ./kagi`
- [x] Test config precedence: `KAGI_API_KEY=env ./kagi test`
- [x] Test validation: missing API key, empty query, invalid flag values
- [x] Test format normalization: `-f txt` → text, `-f markdown` → md
- [x] Test verbose/debug flags

**Deliverable:** CLI accepting all flags and parsing queries correctly.

**Implementation Notes:**
- Custom help template (not Cobra auto-generated) following snag pattern
- Version set via ldflags in Homebrew formula: `-X main.version=x.x.x`
- Format normalization with switch statement (txt→text, markdown→md)
- Fail-fast validation for critical errors
- Debug implies verbose (automatic)
- Clean error messages matching design spec
- Working end-to-end: can query Kagi API with all flags

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 4: Output Formatting ✅ COMPLETED

**Objective:** Implement all output formats with color support.

**Tasks:**

- [x] Implement color support (in main.go)
  - [x] Define ANSI color constants (bold, blue, cyan, yellow)
  - [x] Implement TTY detection using `term.IsTerminal()`
  - [x] Color application logic (auto/always/never)
  - [x] `colorize()` helper function
- [x] Implement text format (`formatText_output()` in main.go)
  - [x] Format output body
  - [x] Format references section (numbered list)
  - [x] Add heading if `--heading` flag set
  - [x] Apply colors if enabled
- [x] Implement markdown format (`formatMarkdown_output()`)
  - [x] Format heading as `# query`
  - [x] Format output body
  - [x] Format references as markdown links with blockquotes
- [x] Implement JSON format (`formatJSON_output()`)
  - [x] Pretty-print full API response
  - [x] Handle `--quiet` (output field only)
- [x] Implement quiet mode for all formats
  - [x] Text: just output body
  - [x] Markdown: just output body
  - [x] JSON: just `.data.output` as JSON string

**Manual Testing:**
- [x] Test text format: default output with references
- [x] Test markdown: proper heading and reference links
- [x] Test JSON: full response and quiet mode
- [x] Test heading: `--heading` shows "# query"
- [x] Test quiet: `-q` outputs only body
- [x] Test colors in pipe: `| cat` (no color, auto-detected)
- [x] Test color flags: `--color always` (ANSI codes), `--color never` (plain)

**Deliverable:** Complete output formatting for all modes.

**Implementation Notes:**
- ANSI colors: References (bold), numbers (yellow), URLs (cyan), headings (bold blue)
- Color auto-detection: Uses `term.IsTerminal(int(os.Stdout.Fd()))`
- Text format: Numbered references with title - URL - snippet
- Markdown format: `[title](url)` links with `> snippet` blockquotes
- JSON format: Pretty-printed with 2-space indent
- Quiet mode: Consistent across all formats (output only)
- All formatting functions: `formatText_output()`, `formatMarkdown_output()`, `formatJSON_output()`
- Helper functions: `shouldUseColor()`, `colorize()`
- Current line count: 564 lines (well under 1000 threshold)

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 5: Error Handling ✅ COMPLETED

**Objective:** Implement comprehensive error handling per design spec.

**Tasks:**

- [x] Define error types and messages
  - [x] Missing API key error
  - [x] Missing query error
  - [x] Invalid flag value errors
  - [x] API error responses
  - [x] Network errors
  - [x] Timeout errors
  - [x] JSON parse errors
- [x] Implement error formatting
  - [x] Prefix with "Error: "
  - [x] Include actionable hints
  - [x] Output to stderr
- [x] Implement verbosity levels
  - [x] Default: silent (errors only)
  - [x] `--verbose`: process info to stderr
  - [x] `--debug`: verbose + detailed debug info
- [x] Set appropriate exit codes
  - [x] 0 for success
  - [x] 1 for errors (simplified from spec)
  - [x] 130 for SIGINT
- [x] Handle signals
  - [x] Graceful SIGINT (Ctrl+C) handling
  - [x] Clean exit on signal

**Manual Testing:**
- [x] Test each error scenario from design-record.md
- [x] Test `--verbose` output
- [x] Test `--debug` output
- [x] Test Ctrl+C handling (exit code 130)
- [x] Verify error messages go to stderr
- [x] Verify exit codes work correctly

**Deliverable:** Robust error handling with clear user feedback.

**Implementation Notes:**
- All 11 error scenarios from design-record.md verified and working
- Signal handler in `main()` catches `os.Interrupt` and `syscall.SIGTERM`
- Exit codes: 0 (success), 1 (all errors), 130 (interrupt)
- Simplified exit codes (1 for all errors vs 1/2 split) for KISS
- All errors prefixed with "Error: " and output to stderr
- Verbose mode shows "Querying Kagi FastGPT API..." and "Response received (Xms)"
- Debug mode shows API key (sanitized as "***"), query, format, timeout
- Debug implies verbose automatically
- Stdin read error handling with proper error messages
- Timeout errors distinguish between request timeout and network timeout
- API errors show specific messages for 401/403 (invalid key) and 429 (rate limit)
- Current line count: 576 lines

**Note:** Unit tests will be written in Phase 7 after core implementation is complete.

---

### Phase 6: Main Integration ✅ COMPLETED

**Objective:** Wire everything together in main execution flow.

**Tasks:**

- [x] Implement main execution logic (`runCobra()` in main.go)
  - [x] Load configuration
  - [x] Get query (args or stdin)
  - [x] Validate inputs
  - [x] Create API client
  - [x] Make API request
  - [x] Format output
  - [x] Write to stdout
  - [x] Handle errors appropriately
- [x] Implement version command
  - [x] Standard: version + repo + issues URL
  - [x] Quiet: version number only
- [x] Implement main.go entry point
  - [x] Execute root command
  - [x] Signal handling for graceful shutdown
  - [x] Exit with proper code

**Manual Testing:**
- [x] Test all flag combinations
- [x] Test with real Kagi API
- [x] Test piping: `./kagi test | less`
- [x] Test redirects: `./kagi test > output.txt`
- [x] Test stdin: `echo "query" | ./kagi`
- [x] Test version: `./kagi --version`, `./kagi -v -q`
- [x] Test color auto-detection (terminal vs pipe)

**Deliverable:** Fully functional CLI tool.

**Implementation Notes:**
- All components fully integrated and working end-to-end
- `main()` function sets up signal handling and executes Cobra command
- `runCobra()` orchestrates: config → query → API → format → output
- Version flag shows rich output by default, version number only with `-q`
- All flags working correctly in combination
- Query input from args or stdin with proper TTY detection
- Color auto-detection working (TTY = colored, pipe = plain)
- Verbose and debug modes provide appropriate context
- All error paths tested and working
- Exit codes properly set for all scenarios
- Clean shutdown on Ctrl+C (SIGINT)
- Current line count: 576 lines (well under threshold)
- **CLI is production-ready** - all core functionality complete

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
