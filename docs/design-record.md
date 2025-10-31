# Kagi CLI - Design Record

**Project:** Kagi FastGPT Command Line Interface
**Repository:** github.com/grantcarthew/kagi
**License:** Mozilla Public License 2.0
**Language:** Go 1.22+
**Date:** 2025-10-31

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Design Philosophy](#design-philosophy)
3. [Core Design Decisions](#core-design-decisions)
4. [Feature Specifications](#feature-specifications)
5. [Technical Specifications](#technical-specifications)
6. [Configuration Precedence](#configuration-precedence)
7. [Error Handling](#error-handling)
8. [Output Formats](#output-formats)
9. [Flag Reference](#flag-reference)
10. [API Integration](#api-integration)
11. [Excluded Features](#excluded-features)

---

## Project Overview

### Purpose

A command-line interface for querying the Kagi FastGPT API, designed for both human users and AI agents. The tool provides a simple, efficient way to get AI-powered search results with web context directly in the terminal.

### Target Audience

- **Primary:** AI agents requiring clean, parseable output
- **Secondary:** Human users wanting quick terminal access to Kagi FastGPT
- **Distribution:** Homebrew and `go install`

### Existing Implementation

A Bash version exists at `~/bin/scripts/kagi` with basic functionality. This Go implementation will be a fully-featured, public-facing replacement with enhanced capabilities and better error handling.

---

## Design Philosophy

### KISS Principle (Keep It Simple, Stupid)

Every design decision prioritizes simplicity over features:

- **Single API focus:** FastGPT only, no Universal Summarizer or other APIs
- **No configuration files:** Flags and environment variables only
- **No caching:** Rely on Kagi's server-side caching
- **No interactive mode:** One query, one response
- **No batch processing:** Single query per invocation
- **No advanced features:** No retries, rate limiting, update checks, etc.

### Rationale

Complex features add maintenance burden and cognitive overhead. The tool should do one thing exceptionally well: query Kagi FastGPT and return results cleanly.

---

## Core Design Decisions

### 1. Query Argument Handling

**Decision:** All non-flag arguments are concatenated (space-separated) to form the query.

**Rationale:**

- Natural CLI experience: `kagi golang best practices`
- No need for quotes for simple queries (though quotes still work)
- Cobra handles flag/arg separation automatically
- Flags can appear anywhere: `kagi --format json golang best practices`

**Examples:**

```bash
kagi golang best practices              # Query: "golang best practices"
kagi "golang best practices"            # Query: "golang best practices"
kagi --timeout 60 golang best practices # Query: "golang best practices"
```

### 2. Standard Input Support

**Decision:** Accept queries from stdin when no arguments provided.

**Rationale:**

- Enables piping: `echo "query" | kagi`
- Supports file input: `kagi < query.txt`
- Arguments take precedence over stdin if both exist

### 3. Output Format Strategy

**Decision:** Three distinct formats with different default behaviors:

- **text (default):** Clean output + references, no heading by default
- **markdown:** Structured document with heading + output + references
- **json:** Raw API response

**Rationale:**

- Text format optimized for quick terminal queries and AI agents
- Markdown format for documentation/archival purposes
- JSON format for programmatic consumption
- `--heading` flag only affects text format (markdown always has structure)

### 4. Color Output

**Decision:** Auto-detect TTY and colorize appropriately.

**Rationale:**

- Human users benefit from color-coded sections
- AI agents/pipes get clean output without ANSI codes
- `--color auto` (default) uses Go's `term.IsTerminal()` detection
- Override with `always` or `never` when needed

### 5. Verbosity Levels

**Decision:** Three levels implemented as flags, not stackable:

1. **Default:** Silent (results only to stdout)
2. **--verbose:** Process information to stderr
3. **--debug:** Verbose + detailed debug info to stderr (implies --verbose)

**Rationale:**

- Most users want quiet operation
- `--verbose` for transparency during execution
- `--debug` for troubleshooting without separate flags
- Levels prevent confusion (not `--verbose --verbose --verbose`)

### 6. Short Flags

**Decision:** Provide short flags for common options only:

- `-h` (--help)
- `-v` (--version)
- `-f` (--format)
- `-q` (--quiet)
- `-t` (--timeout)
- `-c` (--color)

**Rationale:**

- Frequently used flags deserve shortcuts
- API key, debug, verbose are infrequent → long form only
- Prevents alphabet soup of flags
- Common convention: version = `-v`, help = `-h`

### 7. No Configuration File

**Decision:** Configuration via flags and environment variables only.

**Rationale:**

- KISS: no config file location debates
- No config file format decisions (YAML, TOML, JSON?)
- No config file validation/error handling
- API key naturally lives in environment
- Other options rarely need persistence

### 8. Timeout Handling

**Decision:** Integer timeout in seconds (default: 30).

**Rationale:**

- Simple integer instead of duration strings ("30s", "1m")
- 30 seconds balances patience vs. responsiveness
- User can override when needed
- Kagi FastGPT typically responds in <10s

### 9. References Section

**Decision:** Always include references in output (except with `--quiet`).

**Rationale:**

- References provide source attribution
- Critical for verifying AI-generated information
- Kagi API provides them, we should surface them
- `--quiet` flag available for users who don't want them

### 10. Version Information

**Decision:** Rich version output by default, minimal with `--quiet`.

**Standard output:**

```
kagi v1.0.0
Repository: https://github.com/grantcarthew/kagi
Report issues: https://github.com/grantcarthew/kagi/issues/new
```

**Quiet output:**

```
1.0.0
```

**Rationale:**

- Helpful information for bug reports (repo, issues link)
- `--quiet` provides machine-readable version
- Follows conventions of tools like `git --version`

---

## Feature Specifications

### Query Processing

1. **Arguments:** Join all non-flag args with spaces
2. **Stdin:** Read from stdin if no args provided
3. **Precedence:** Args override stdin if both exist
4. **Validation:** Error if query is empty or whitespace-only
5. **Encoding:** URL-encode for API transmission

### API Key Management

1. **Sources:** `--api-key` flag or `KAGI_API_KEY` environment variable
2. **Precedence:** Flag overrides environment variable
3. **Validation:** Error if no key provided from either source
4. **Security:** Never log or print the API key (even in debug mode, truncate to "\*\*\*")

### HTTP Client

1. **Timeout:** Configurable via `--timeout` (default: 30 seconds)
2. **Headers:**
   - `Authorization: Bot <api_key>`
   - `Content-Type: application/json`
3. **Method:** POST to `https://kagi.com/api/v0/fastgpt`
4. **Body:** JSON with `{"query": "<query>", "web_search": true}`

### Response Processing

1. **Success:** Extract `.data.output` and `.data.references`
2. **API Errors:** Parse `.error[]` array and display
3. **HTTP Errors:** Report status code and message
4. **Timeout:** Report timeout exceeded
5. **JSON Parse Errors:** Report parse failure

---

## Technical Specifications

### Go Version

- **Minimum:** Go 1.22
- **Tested on:** Go 1.25.3
- **Rationale:** Modern Go features, stable stdlib

### Dependencies

- **Cobra:** CLI framework (github.com/spf13/cobra)
- **Standard library:** net/http, encoding/json, os, io, etc.
- **Terminal detection:** golang.org/x/term for TTY detection

### Project Structure

**Design:**

```
kagi/
├── cmd/
│   └── root.go          # Root command + flags
├── internal/
│   ├── api/
│   │   └── client.go    # Kagi API client
│   ├── format/
│   │   └── output.go    # Output formatting
│   └── config/
│       └── config.go    # Configuration handling
├── main.go              # Entry point
...
```

**Actual Implementation (KISS - Flat Structure):**

```
kagi/
├── main.go              # All code (types, API, CLI, formatting)
├── main_test.go         # All tests
├── go.mod
├── go.sum
├── LICENSE              # MPL 2.0
├── README.md            # User documentation
├── PROJECT.md           # Implementation guide
├── ROLE.md              # Role definition
└── docs/
    └── design-record.md # This file
```

**Rationale for Deviation:**
Per KISS principles, the entire implementation is ~580 lines in a single `main.go` file. The planned package structure would add unnecessary complexity for this scale. Code will only be split if it exceeds 1000 lines.

### Exit Codes

**Design:**

- `0`: Success
- `1`: General error (API, network, parsing)
- `2`: Invalid arguments/usage
- `130`: SIGINT (Ctrl+C)

**Actual Implementation (Simplified):**

- `0`: Success
- `1`: All errors (usage, API, network, parsing)
- `130`: SIGINT (Ctrl+C)

**Rationale for Deviation:**
Following KISS, exit code 1 is used for all error conditions. Distinguishing usage errors (2) from other errors adds minimal value but complicates error handling logic.

### Output Streams

- **stdout:** Query results only
- **stderr:** Errors, verbose/debug messages

---

## Configuration Precedence

### Priority Order (highest to lowest)

1. **CLI flags** (highest priority)
2. **Environment variables** (KAGI_API_KEY only)
3. **Defaults** (lowest priority)

### Examples

**API Key Resolution:**

```bash
# Uses env_key
KAGI_API_KEY=env_key kagi query

# Uses cli_key (flag wins)
KAGI_API_KEY=env_key kagi --api-key cli_key query
```

**Timeout Resolution:**

```bash
# Uses 30 (default)
kagi query

# Uses 60 (flag overrides default)
kagi -t 60 query
```

### Notes

- No configuration file support
- Only `KAGI_API_KEY` environment variable is recognized
- All other configuration via flags only

---

## Error Handling

### Principles

- Prefix all errors with "Error: "
- Be specific and actionable
- Include helpful hints where appropriate
- All errors output to stderr
- Exit with appropriate code

### Error Scenarios

#### 1. Missing API Key

```
Error: no API key provided
Provide via --api-key flag or KAGI_API_KEY environment variable
```

**Exit code:** 1

#### 2. Missing Query

```
Error: query cannot be empty
```

**Exit code:** 1

#### 3. Invalid Flag Value

```
Error: invalid format 'xml'. Valid formats: text, txt, md, markdown, json
```

**Exit code:** 1

#### 4. Invalid Timeout

```
Error: timeout must be a positive integer
```

**Exit code:** 1

#### 5. API Error Response

```
Error: API request failed [403]: Invalid API key
```

**Exit code:** 1

#### 6. Network/Connection Error

```
Error: network request failed: connection timeout
```

**Exit code:** 1

#### 7. HTTP Error (non-2xx)

```
Error: API returned HTTP 500: Internal Server Error
```

**Exit code:** 1

#### 8. Invalid JSON Response

```
Error: failed to parse API response
```

**Exit code:** 1

#### 9. Request Timeout

```
Error: request timeout exceeded (30s)
```

**Exit code:** 1

#### 10. Stdin Read Error

```
Error: failed to read from stdin: <reason>
```

**Exit code:** 1

#### 11. Interrupt (Ctrl+C)

```
(no message, just clean exit)
```

**Exit code:** 130

### Enhanced Error Context

**With --verbose:**

```
Querying Kagi FastGPT API...
Error: API returned HTTP 500: Internal Server Error
```

**With --debug:**

```
Querying Kagi FastGPT API...
POST https://kagi.com/api/v0/fastgpt
Request: {"query":"test query","web_search":true}
Response: HTTP 500 Internal Server Error
Body: {"error":[{"code":500,"msg":"Internal error"}]}
Error: API returned HTTP 500: Internal Server Error
```

---

## Output Formats

### Text Format (default)

**Without --heading:**

```
Python 3.11 was released in 2021 and introduced several new features...

References:

1. What's New In Python 3.11 — Python 3.11.3 documentation - https://docs.python.org/3/whatsnew/3.11.html - ...
2. Introducing the New Features in Python 3.11 - https://earthly.dev/blog/python-3.11-new-features/ - ...
```

**With --heading:**

```
# python 3.11

Python 3.11 was released in 2021 and introduced several new features...

References:

1. What's New In Python 3.11 — Python 3.11.3 documentation - https://docs.python.org/3/whatsnew/3.11.html - ...
2. Introducing the New Features in Python 3.11 - https://earthly.dev/blog/python-3.11-new-features/ - ...
```

**With --quiet:**

```
Python 3.11 was released in 2021 and introduced several new features...
```

### Markdown Format

**Always includes heading and references:**

```markdown
# python 3.11

Python 3.11 was released in 2021 and introduced several new features...

## References

1. [What's New In Python 3.11 — Python 3.11.3 documentation](https://docs.python.org/3/whatsnew/3.11.html)
   > ...
2. [Introducing the New Features in Python 3.11](https://earthly.dev/blog/python-3.11-new-features/)
   > ...
```

**With --quiet:**

```
Python 3.11 was released in 2021 and introduced several new features...
```

### JSON Format

**Raw API response:**

```json
{
  "meta": {
    "id": "120145af-f057-466d-9e6d-7829ac902adc",
    "node": "us-east",
    "ms": 7943
  },
  "data": {
    "output": "Python 3.11 was released in 2021...",
    "tokens": 757,
    "references": [
      {
        "title": "What's New In Python 3.11",
        "snippet": "...",
        "url": "https://docs.python.org/3/whatsnew/3.11.html"
      }
    ]
  }
}
```

**With --quiet:**

```json
"Python 3.11 was released in 2021..."
```

### Color Output

**When enabled (--color always or auto with TTY):**

- Headings: Bold Blue
- Section headers (References): Bold
- URLs: Cyan
- Reference numbers: Yellow

**Implementation:** Use ANSI escape codes, easily stripped for `--color never` or non-TTY.

---

## Flag Reference

### --help, -h

**Type:** Boolean
**Default:** false
**Description:** Display help information and exit.

### --version, -v

**Type:** Boolean
**Default:** false
**Description:** Display version information and exit.
**Behavior:**

- Default: Show version, repo URL, issue URL
- With `--quiet`: Show version number only

### --api-key

**Type:** String
**Default:** "" (reads from KAGI_API_KEY env var)
**Description:** Kagi API key for authentication.
**Precedence:** Overrides KAGI_API_KEY environment variable.

### --format, -f

**Type:** String
**Default:** "text"
**Valid values:** text, txt, md, markdown, json
**Description:** Output format for results.
**Notes:** "txt" is alias for "text", "markdown" is alias for "md"

### --heading

**Type:** Boolean
**Default:** false
**Description:** Include query as heading in output.
**Applies to:** Text format only (ignored for markdown/json)

### --quiet, -q

**Type:** Boolean
**Default:** false
**Description:** Output only the response body, no heading or references.
**Applies to:** All formats

### --timeout, -t

**Type:** Integer
**Default:** 30
**Description:** HTTP request timeout in seconds.
**Validation:** Must be positive integer.

### --color, -c

**Type:** String
**Default:** "auto"
**Valid values:** auto, always, never
**Description:** Control color output.

- `auto`: Color if stdout is a TTY
- `always`: Always use color codes
- `never`: Never use color codes

### --verbose

**Type:** Boolean
**Default:** false
**Description:** Output query process information to stderr.
**Output examples:**

- "Querying Kagi FastGPT API..."
- "Response received (1234ms)"

### --debug

**Type:** Boolean
**Default:** false
**Description:** Output detailed debug information to stderr.
**Implies:** --verbose
**Output includes:**

- HTTP request details
- Request/response headers
- Full response body
- Timing information
- Error stack traces

---

## API Integration

### Kagi FastGPT API

**Endpoint:** `https://kagi.com/api/v0/fastgpt`
**Method:** POST
**Authentication:** `Authorization: Bot <api_key>`

### Request Format

```json
{
  "query": "search query here",
  "web_search": true,
  "cache": true
}
```

**Parameters:**

- `query` (string, required): The search query
- `web_search` (boolean, optional): Enable web search (default: true, always true for now)
- `cache` (boolean, optional): Allow cached responses (default: true)

### Response Format

```json
{
  "meta": {
    "id": "request-id",
    "node": "us-east",
    "ms": 7943
  },
  "data": {
    "output": "AI-generated response",
    "tokens": 757,
    "references": [
      {
        "title": "Reference title",
        "snippet": "Reference snippet",
        "url": "https://example.com"
      }
    ]
  }
}
```

### Error Response Format

```json
{
  "error": [
    {
      "code": 403,
      "msg": "Invalid API key"
    }
  ]
}
```

### Pricing

- **With web_search:** $0.015 per query ($15 per 1000 queries)
- **Cached responses:** Free

### Implementation Notes

1. Always set `web_search: true` (no flag to disable - KISS)
2. Always set `cache: true` (rely on Kagi's caching)
3. Handle both success and error response formats
4. Parse references array even if empty
5. Token count available for future usage tracking

---

## Excluded Features

These features were explicitly considered and rejected to maintain simplicity:

### Configuration File

**Rejected:** No `~/.kagi.yaml` or similar
**Rationale:** Adds complexity, most users only need API key in environment

### Multiple API Support

**Rejected:** No Universal Summarizer, Search API, etc.
**Rationale:** FastGPT is the primary use case, other APIs would complicate interface

### Caching

**Rejected:** No local response caching
**Rationale:** Kagi provides server-side caching, local cache adds state management

### Interactive Mode

**Rejected:** No REPL-style `kagi --interactive`
**Rationale:** Outside scope of simple CLI tool, shell history works fine

### Batch Processing

**Rejected:** No `--file queries.txt` to process multiple queries
**Rationale:** Shell loops handle this: `while read q; do kagi "$q"; done < queries.txt`

### History/Session Management

**Rejected:** No `kagi history` or session tracking
**Rationale:** Adds persistent state, shell history suffices

### Cost Tracking

**Rejected:** No API usage/cost tracking
**Rationale:** Kagi dashboard provides this, adds complexity

### Retry Logic

**Rejected:** No automatic retries with exponential backoff
**Rationale:** User can re-run command, adds complexity for rare failures

### Rate Limiting

**Rejected:** No client-side rate limit handling
**Rationale:** Kagi enforces server-side, adding client logic is premature

### Update Notifications

**Rejected:** No "new version available" warnings
**Rationale:** Homebrew handles updates, notifications annoy AI agents

### Shell Completion

**Rejected:** No bash/zsh/fish completion generation
**Rationale:** Limited value for simple flag set, adds maintenance

### Custom User-Agent

**Rejected:** No version in User-Agent header
**Rationale:** Not required by API, adds version coordination

### Streaming Support

**Rejected:** No streaming responses
**Rationale:** Kagi API doesn't support streaming (flat-rate pricing per query)

### Web Search Toggle

**Rejected:** No `--no-web-search` flag
**Rationale:** Feature currently disabled in API, may be removed entirely

---

## Implementation Notes

### Error Recovery

- Graceful handling of Ctrl+C (SIGINT)
- Proper cleanup of HTTP connections
- No partial output on errors (atomic success/failure)

### Testing Strategy

Following KISS principles, tests are written after core implementation is complete (see PROJECT.md Phase 7):

- Unit tests for all packages (api, config, format, input)
- Integration tests with mock API responses
- Error case coverage for all 11 error scenarios
- Flag parsing validation
- Manual testing during implementation (Phases 1-6)
- Comprehensive automated testing in Phase 7
- Target: >80% test coverage

### Documentation Requirements

- Comprehensive README.md with examples
- Man page (optional, future)
- Inline code documentation
- This design record for maintainers

### Distribution

1. **go install:** `go install github.com/grantcarthew/kagi@latest`
2. **Homebrew:** Custom tap with formula
3. **GitHub Releases:** Binary releases for major platforms

---

## Versioning

**Scheme:** Semantic Versioning (semver)
**Format:** vMAJOR.MINOR.PATCH

**Examples:**

- v1.0.0 - Initial release
- v1.0.1 - Bug fix
- v1.1.0 - New feature (backward compatible)
- v2.0.0 - Breaking change

---

## Future Considerations

If KISS principle is relaxed in future versions, consider:

1. **Config file support** - If users request persistent settings
2. **Additional Kagi APIs** - If demand exists for Summarizer, etc.
3. **Response caching** - If network/cost becomes issue
4. **Structured logging** - If integration with logging systems needed
5. **Plugin system** - If extensibility becomes valuable

**Note:** These should only be added with strong user demand and careful consideration of complexity cost.

---

## Implementation Deviations

The following deviations from the original design were made during implementation, all following the KISS principle:

### 1. Flat Project Structure

**Design:** Separate packages (cmd/, internal/api/, internal/format/, internal/config/)
**Actual:** Single main.go file (~580 lines)
**Rationale:** Project complexity doesn't justify package separation. All code fits comfortably in one file.

### 2. Simplified Exit Codes

**Design:** Exit code 2 for usage errors, 1 for other errors
**Actual:** Exit code 1 for all errors
**Rationale:** Distinguishing usage vs runtime errors adds minimal value but complicates error handling.

### 3. Testing Approach

**Design:** TDD with tests written alongside implementation
**Actual:** Tests written in Phase 7 after implementation complete
**Rationale:** Faster development, better understanding of edge cases, avoids test churn.

### 4. Test Coverage

**Achieved:** 48.3% overall, 100% on business logic
**Analysis:** Untested code is integration/glue code (main, runCobra, loadConfig, queryKagi)
**Rationale:** Higher coverage requires dependency injection, contradicts KISS for this architecture.

All deviations maintain or improve simplicity while preserving functionality.

---

## Conclusion

This design record captures all major decisions made during the design phase of the Kagi CLI tool. It serves as the authoritative reference for implementation and future maintenance. Any deviation from these specifications should be documented with rationale in git commit messages and potentially as amendments to this document.

**Design Approved:** 2025-10-31
**Implementation Completed:** 2025-10-31 (Phases 1-7)
**Status:** Production-ready, documentation phase (Phase 8)
**Next Step:** See PROJECT.md for remaining phases (documentation and distribution)
