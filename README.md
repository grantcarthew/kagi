# kagi

A fast, simple command-line interface for [Kagi FastGPT](https://kagi.com/fastgpt) - AI-powered search with web context.

## Features

- **Simple Interface**: Just type your query and get AI-powered answers
- **Multiple Output Formats**: Plain text, Markdown, or JSON
- **Smart Color Output**: Automatically detects terminals and pipes
- **Web References**: Includes sources with every response
- **Flexible Input**: Accept queries from arguments or stdin
- **Production Ready**: Comprehensive error handling and timeout controls
- **AI Agent Friendly**: Clean output formats for automation

## Installation

### Homebrew (macOS/Linux)

```bash
brew install grantcarthew/tap/kagi
```

### Go Install

```bash
go install github.com/grantcarthew/kagi@latest
```

### From Source

```bash
git clone https://github.com/grantcarthew/kagi.git
cd kagi
go build
```

## Quick Start

### 1. Get a Kagi API Key

1. Sign up at [kagi.com](https://kagi.com)
2. Visit [Account Settings > API](https://kagi.com/settings?p=api)
3. Generate a FastGPT API key

### 2. Set Your API Key

```bash
export KAGI_API_KEY='your-api-key-here'
```

Add this to your `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish` to make it permanent.

### 3. Ask Questions

```bash
kagi golang best practices
```

## Usage

```
kagi [options] <query...>
```

### Basic Examples

```bash
# Simple query
kagi what is kubernetes

# Multiple words automatically joined
kagi golang concurrency patterns

# With quotes (optional)
kagi "explain docker containers"
```

### Output Formats

#### Text Format (Default)

Clean, readable output with numbered references:

```bash
kagi golang generics
```

Output:

```
Docker containers are lightweight, standalone executable packages...

References:
  1. What are Containers? - https://www.docker.com/resources/what-container
     Docker containers are a standardized unit of software that packages...

  2. Container Technology - https://kubernetes.io/docs/concepts/containers/
     Containers decouple applications from underlying host infrastructure...
```

#### Markdown Format

Perfect for documentation or README files:

```bash
kagi -f md golang generics > generics-guide.md
kagi -f markdown how to use rust > rust-guide.md
```

Output includes markdown links and blockquotes:

```markdown
# golang generics

Go generics allow you to write flexible, reusable code...

## References

1. [Go Generics Tutorial](https://go.dev/doc/tutorial/generics)

   > Learn how to add generic functions and types to your Go code...

2. [Type Parameters Proposal](https://go.dev/blog/generics)
   > The Go team introduces type parameters to the language...
```

#### JSON Format

For scripting and automation:

```bash
kag -f json golang channels | jq '.data.output'
```

Full API response:

```json
{
  "meta": {
    "id": "...",
    "node": "...",
    "ms": 1234
  },
  "data": {
    "output": "Go channels are a typed conduit...",
    "references": [
      {
        "title": "Go by Example: Channels",
        "snippet": "Channels are the pipes that connect...",
        "url": "https://gobyexample.com/channels"
      }
    ]
  }
}
```

### Using Stdin

Read queries from pipes or redirects:

```bash
# From echo
echo "explain kubernetes pods" | kagi

# From file
cat questions.txt | kagi

# From here-doc
kagi << EOF
What are the benefits of using Go for
backend development?
EOF

# In scripts
QUERY="golang vs rust performance"
echo "$QUERY" | kagi -f md > comparison.md
```

### Common Options

```bash
# Quiet mode (output only, no references)
kagi -q golang channels

# With heading
kagi --heading golang best practices

# Markdown with quiet mode
kagi -f md -q golang testing > testing.md

# Increase timeout for complex queries
kagi --timeout 60 "comprehensive guide to kubernetes"

# Force color output in pipes
kagi --color always golang patterns | less -R

# Disable colors
kagi --color never golang patterns

# Debug mode
kagi --debug golang generics
```

### Examples for Automation

```bash
# Save JSON response
kagi -f json "golang concurrency" > response.json

# Extract just the answer
kagi -q "what is docker" > answer.txt

# Generate markdown documentation
kagi -f md "golang error handling best practices" > errors.md

# Batch processing
cat topics.txt | while read topic; do
  kagi -f md -q "$topic" > "docs/${topic}.md"
done

# CI/CD usage
if kagi -q "is golang 1.23 stable"; then
  echo "Proceeding with upgrade"
fi
```

## Command Reference

### Options

| Flag        | Short | Default         | Description                                            |
| ----------- | ----- | --------------- | ------------------------------------------------------ |
| `--api-key` |       | `$KAGI_API_KEY` | Kagi API key (overrides environment variable)          |
| `--format`  | `-f`  | `text`          | Output format: `text`, `txt`, `md`, `markdown`, `json` |
| `--quiet`   | `-q`  | `false`         | Output only response body (no heading or references)   |
| `--heading` |       | `false`         | Include query as heading in text format                |
| `--timeout` | `-t`  | `30`            | HTTP request timeout in seconds                        |
| `--color`   | `-c`  | `auto`          | Color output: `auto`, `always`, `never`                |
| `--verbose` |       | `false`         | Output process information to stderr                   |
| `--debug`   |       | `false`         | Output detailed debug information to stderr            |
| `--version` | `-v`  |                 | Display version information                            |
| `--help`    | `-h`  |                 | Display help message                                   |

### Environment Variables

| Variable       | Description                                           |
| -------------- | ----------------------------------------------------- |
| `KAGI_API_KEY` | Your Kagi API key (required unless using `--api-key`) |
| `NO_COLOR`     | Disable color output (standard)                       |

### Exit Codes

| Code  | Meaning                                            |
| ----- | -------------------------------------------------- |
| `0`   | Success                                            |
| `1`   | Error (API error, network error, validation error) |
| `130` | Interrupted (Ctrl+C)                               |

## Color Output

The CLI automatically detects whether output is going to a terminal or pipe:

```bash
# Colored output (terminal)
kagi golang patterns

# No colors (piped)
kagi golang patterns | less

# Force colors in pipe
kagi --color always golang patterns | less -R

# Disable colors
kagi --color never golang patterns
```

Colors are also disabled if the `NO_COLOR` environment variable is set.

## Error Handling

### Common Errors

**Missing API Key:**

```bash
$ kagi test
Error: API key required. Set KAGI_API_KEY environment variable or use --api-key flag
```

**Empty Query:**

```bash
$ kagi
Error: query cannot be empty
```

**Invalid Format:**

```bash
$ kagi -f xml test
Error: invalid format 'xml'. Valid formats: text, txt, md, markdown, json
```

**API Error:**

```bash
$ kagi test
Error: API error (401): Invalid or missing API key
```

**Network Timeout:**

```bash
$ kagi --timeout 1 test
Error: request timeout after 1s
```

### Verbose Output

Use `--verbose` to see what's happening:

```bash
$ kagi --verbose golang channels
Querying Kagi FastGPT API...
Response received (1234ms)
[output...]
```

### Debug Output

Use `--debug` for detailed information:

```bash
$ kagi --debug golang channels
API Key: ***
Query: golang channels
Format: text
Timeout: 30s
Querying Kagi FastGPT API...
Response received (1234ms)
[output...]
```

## Integration Examples

### Shell Scripts

```bash
#!/bin/bash
# ask-kagi.sh - Interactive Kagi queries

if [ -z "$KAGI_API_KEY" ]; then
  echo "Error: Set KAGI_API_KEY environment variable"
  exit 1
fi

while true; do
  echo -n "Ask Kagi: "
  read -r query
  [ -z "$query" ] && break

  kagi "$query"
  echo ""
done
```

### Vim Integration

Add to your `.vimrc`:

```vim
" Query selected text with Kagi
vnoremap <leader>k :!kagi -q<CR>

" Query word under cursor
nnoremap <leader>k :!kagi -q <cword><CR>
```

### Fish Shell Function

Add to your `~/.config/fish/functions/ask.fish`:

```fish
function ask --description "Quick Kagi query"
    kagi $argv
end
```

Then use:

```bash
ask golang channels
```

### Alfred/Raycast Workflow

```bash
#!/bin/bash
# Kagi query for Alfred/Raycast
query="$1"
kagi -q "$query" | pbcopy
echo "Copied to clipboard!"
```

## Troubleshooting

### API Key Not Found

Make sure your API key is set:

```bash
echo $KAGI_API_KEY
```

If empty, add to your shell config:

```bash
# Bash/Zsh
echo 'export KAGI_API_KEY="your-key"' >> ~/.bashrc

# Fish
set -Ux KAGI_API_KEY "your-key"
```

### Timeout Errors

Increase timeout for complex queries:

```bash
kagi --timeout 60 "your long query"
```

### Color Output Issues

If colors appear broken:

```bash
# Disable colors
kagi --color never query

# Or set environment variable
export NO_COLOR=1
```

### Rate Limiting

Kagi API has rate limits. If you receive a 429 error:

```
Error: API error (429): Rate limit exceeded. Please try again later.
```

Wait a few moments before retrying.

## Development

### Building

```bash
git clone https://github.com/grantcarthew/kagi.git
cd kagi
go build
```

### Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Project Structure

```
kagi/
├── main.go          # All code (types, API client, CLI, formatting)
├── main_test.go     # Comprehensive test suite (83 tests)
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── LICENSE          # Mozilla Public License 2.0
├── README.md        # This file
├── PROJECT.md       # Development guide
└── ROLE.md          # Role definition
```

The project follows the KISS principle with a flat structure - all code in a single `main.go` file.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Guidelines

1. Follow Go conventions (`gofmt`, `golint`)
2. Add tests for new features
3. Update documentation
4. Keep it simple (KISS principle)

### Reporting Issues

Please include:

- Operating system and version
- Go version (`go version`)
- Full command and error output
- Expected vs actual behavior

Report issues at: https://github.com/grantcarthew/kagi/issues

## License

Mozilla Public License 2.0 (MPL-2.0)

See [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Powered by [Kagi FastGPT](https://kagi.com/fastgpt)

## Links

- **Repository**: https://github.com/grantcarthew/kagi
- **Issues**: https://github.com/grantcarthew/kagi/issues
- **Kagi**: https://kagi.com
- **Kagi API Docs**: https://help.kagi.com/kagi/api/fastgpt.html

---

**Note**: This is an unofficial tool and is not affiliated with or endorsed by Kagi.
