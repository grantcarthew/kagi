package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Version is set via ldflags during build
var version = "dev"

const (
	// API configuration
	apiEndpoint    = "https://kagi.com/api/v0/fastgpt"
	defaultTimeout = 30 // seconds

	// HTTP headers
	contentTypeJSON  = "application/json"
	authHeaderPrefix = "Bot "

	// Request defaults
	webSearchEnabled = true
	cacheEnabled     = true

	// Exit codes
	exitSuccess   = 0
	exitError     = 1
	exitInterrupt = 130

	// Output formats
	formatText     = "text"
	formatMarkdown = "md"
	formatJSON     = "json"

	// Color modes
	colorAuto   = "auto"
	colorAlways = "always"
	colorNever  = "never"

	// Environment variables
	envAPIKey = "KAGI_API_KEY"
)

const helpTemplate = `USAGE:
  kagi [options] <query...>

DESCRIPTION:
  kagi queries the Kagi FastGPT API and returns AI-powered search results
  with web context. Designed for both human users and AI agents.

  Output formats: text (default), markdown (md), or JSON.

  API key: Set KAGI_API_KEY environment variable or use --api-key flag.

EXAMPLES:
  # Basic query
  kagi golang best practices

  # Different output formats
  kagi -f md golang best practices
  kagi -f json golang concurrency > result.json

  # Using stdin
  echo "explain kubernetes" | kagi

  # With options
  kagi --heading --timeout 60 golang generics
  kagi -q golang channels              # Quiet mode (output body only)

OPTIONS:
  -f, --format string      Output format: text (txt) | md (markdown) | json (default "text")
  -q, --quiet              Output only response body (no heading or references)
      --heading            Include query as heading in text format
  -t, --timeout int        HTTP request timeout in seconds (default 30)
  -c, --color string       Color output: auto | always | never (default "auto")

      --api-key string     Kagi API key (overrides KAGI_API_KEY env var)

      --verbose            Output process information to stderr
      --debug              Output detailed debug information to stderr

  -h, --help               Display this help message
  -v, --version            Display version information
`

type FastGPTRequest struct {
	Query     string `json:"query"`
	WebSearch bool   `json:"web_search"`
	Cache     bool   `json:"cache"`
}

type FastGPTResponse struct {
	Meta struct {
		ID   string `json:"id"`
		Node string `json:"node"`
		MS   int    `json:"ms"`
	} `json:"meta"`
	Data struct {
		Output     string      `json:"output"`
		Tokens     int         `json:"tokens"`
		References []Reference `json:"references"`
	} `json:"data"`
}

type FastGPTError struct {
	Error []struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}

type Reference struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

type Config struct {
	APIKey  string
	Query   string
	Format  string
	Timeout int
	Heading bool
	Quiet   bool
	Color   string
	Verbose bool
	Debug   bool
}

var (
	flagAPIKey  string
	flagFormat  string
	flagTimeout int
	flagHeading bool
	flagQuiet   bool
	flagColor   string
	flagVerbose bool
	flagDebug   bool
	flagVersion bool
)

var rootCmd = &cobra.Command{
	Use:          "kagi [options] <query...>",
	Short:        "Query Kagi FastGPT API from the command line",
	Args:         cobra.ArbitraryArgs,
	RunE:         runCobra,
	SilenceUsage: true,
}

func init() {
	rootCmd.Flags().StringVar(&flagAPIKey, "api-key", "", "Kagi API key (overrides KAGI_API_KEY env var)")
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", formatText, "Output format: text | txt | md | markdown | json")
	rootCmd.Flags().IntVarP(&flagTimeout, "timeout", "t", defaultTimeout, "HTTP request timeout in seconds")
	rootCmd.Flags().BoolVar(&flagHeading, "heading", false, "Include query as heading in text format")
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Output only response body (no heading or references)")
	rootCmd.Flags().StringVarP(&flagColor, "color", "c", colorAuto, "Color output: auto | always | never")
	rootCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Output process information to stderr")
	rootCmd.Flags().BoolVar(&flagDebug, "debug", false, "Output detailed debug information to stderr")
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "Display version information")

	rootCmd.SetHelpTemplate(helpTemplate)
}

func main() {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		// Clean exit on interrupt
		os.Exit(exitInterrupt)
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(exitError)
	}
}

func runCobra(cmd *cobra.Command, args []string) error {
	if flagVersion {
		if flagQuiet {
			fmt.Println(version)
		} else {
			fmt.Printf("kagi v%s\n", version)
			fmt.Println("Repository: https://github.com/grantcarthew/kagi")
			fmt.Println("Report issues: https://github.com/grantcarthew/kagi/issues/new")
		}
		return nil
	}

	config, err := loadConfig(cmd, args)
	if err != nil {
		return err
	}

	if config.Debug {
		fmt.Fprintf(os.Stderr, "Debug: API Key: ***\n")
		fmt.Fprintf(os.Stderr, "Debug: Query: %s\n", config.Query)
		fmt.Fprintf(os.Stderr, "Debug: Format: %s\n", config.Format)
		fmt.Fprintf(os.Stderr, "Debug: Timeout: %d\n", config.Timeout)
	}

	if config.Verbose || config.Debug {
		fmt.Fprintf(os.Stderr, "Querying Kagi FastGPT API...\n")
	}

	resp, err := queryKagi(config.APIKey, config.Query, config.Timeout)
	if err != nil {
		return err
	}

	if config.Verbose || config.Debug {
		fmt.Fprintf(os.Stderr, "Response received (%dms)\n", resp.Meta.MS)
	}

	output, err := formatOutput(resp, config)
	if err != nil {
		return err
	}
	fmt.Print(output)

	return nil
}

func loadConfig(cmd *cobra.Command, args []string) (*Config, error) {
	// Get API key (flag takes precedence over env var)
	apiKey := flagAPIKey
	if apiKey == "" {
		apiKey = os.Getenv(envAPIKey)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no API key provided\nProvide via --api-key flag or KAGI_API_KEY environment variable")
	}

	query, err := getQuery(args)
	if err != nil {
		return nil, err
	}

	format := normalizeFormat(flagFormat)
	if !isValidFormat(format) {
		return nil, fmt.Errorf("invalid value %q for --format\nValid formats: text, txt, md, markdown, json", flagFormat)
	}

	if flagTimeout <= 0 {
		return nil, fmt.Errorf("invalid timeout value %q\nTimeout must be a positive integer (seconds)", fmt.Sprint(flagTimeout))
	}

	color := strings.ToLower(strings.TrimSpace(flagColor))
	if color != colorAuto && color != colorAlways && color != colorNever {
		return nil, fmt.Errorf("invalid value %q for --color\nValid values: auto, always, never", flagColor)
	}

	// Debug implies verbose
	verbose := flagVerbose
	if flagDebug {
		verbose = true
	}

	return &Config{
		APIKey:  apiKey,
		Query:   query,
		Format:  format,
		Timeout: flagTimeout,
		Heading: flagHeading,
		Quiet:   flagQuiet,
		Color:   color,
		Verbose: verbose,
		Debug:   flagDebug,
	}, nil
}

// getQuery extracts the query from args or stdin
func getQuery(args []string) (string, error) {
	// First, try to get query from args
	if len(args) > 0 {
		query := strings.TrimSpace(strings.Join(args, " "))
		if query != "" {
			return query, nil
		}
	}

	// If no args, try stdin (only if not a terminal)
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
		query := strings.TrimSpace(string(stdinBytes))
		if query != "" {
			return query, nil
		}
	}

	return "", fmt.Errorf("no query provided\nUsage: kagi [flags] <query...>")
}

func normalizeFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))

	switch format {
	case "txt":
		return formatText
	case "markdown":
		return formatMarkdown
	default:
		return format
	}
}

func isValidFormat(format string) bool {
	return format == formatText || format == formatMarkdown || format == formatJSON
}

// ANSI color codes
const (
	ansiReset    = "\033[0m"
	ansiBold     = "\033[1m"
	ansiBlue     = "\033[34m"
	ansiBoldBlue = "\033[1;34m"
	ansiCyan     = "\033[36m"
	ansiYellow   = "\033[33m"
)

func shouldUseColor(config *Config) bool {
	switch config.Color {
	case colorAlways:
		return true
	case colorNever:
		return false
	case colorAuto:
		return term.IsTerminal(int(os.Stdout.Fd()))
	default:
		return false
	}
}

func colorize(text, colorCode string, useColor bool) string {
	if !useColor {
		return text
	}
	return colorCode + text + ansiReset
}

func formatOutput(resp *FastGPTResponse, config *Config) (string, error) {
	switch config.Format {
	case formatJSON:
		return formatJSON_output(resp, config)
	case formatMarkdown:
		return formatMarkdown_output(resp, config), nil
	default: // formatText
		return formatText_output(resp, config), nil
	}
}

func formatText_output(resp *FastGPTResponse, config *Config) string {
	var output strings.Builder
	useColor := shouldUseColor(config)

	if config.Heading && !config.Quiet {
		heading := "# " + config.Query
		output.WriteString(colorize(heading, ansiBoldBlue, useColor))
		output.WriteString("\n\n")
	}

	output.WriteString(resp.Data.Output)
	output.WriteString("\n")

	if !config.Quiet && len(resp.Data.References) > 0 {
		output.WriteString("\n")
		output.WriteString(colorize("References:", ansiBold, useColor))
		output.WriteString("\n\n")

		for i, ref := range resp.Data.References {
			refNum := fmt.Sprintf("%d. ", i+1)
			output.WriteString(colorize(refNum, ansiYellow, useColor))

			output.WriteString(ref.Title)
			output.WriteString(" - ")

			output.WriteString(colorize(ref.URL, ansiCyan, useColor))

			if ref.Snippet != "" {
				output.WriteString(" - ")
				output.WriteString(ref.Snippet)
			}

			output.WriteString("\n")
		}
	}

	return output.String()
}

func formatMarkdown_output(resp *FastGPTResponse, config *Config) string {
	var output strings.Builder

	if config.Quiet {
		output.WriteString(resp.Data.Output)
		output.WriteString("\n")
		return output.String()
	}

	// Markdown always includes heading
	output.WriteString("# ")
	output.WriteString(config.Query)
	output.WriteString("\n\n")

	output.WriteString(resp.Data.Output)
	output.WriteString("\n")

	if len(resp.Data.References) > 0 {
		output.WriteString("\n## References\n\n")

		for i, ref := range resp.Data.References {
			output.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, ref.Title, ref.URL))

			if ref.Snippet != "" {
				output.WriteString("   > ")
				output.WriteString(ref.Snippet)
				output.WriteString("\n")
			}
		}
	}

	return output.String()
}

func formatJSON_output(resp *FastGPTResponse, config *Config) (string, error) {
	if config.Quiet {
		jsonBytes, err := json.MarshalIndent(resp.Data.Output, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal output to JSON: %w", err)
		}
		return string(jsonBytes) + "\n", nil
	}

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		// Fallback to non-indented if pretty print fails
		jsonBytes, err = json.Marshal(resp)
		if err != nil {
			return "", fmt.Errorf("failed to marshal response to JSON: %w", err)
		}
	}

	return string(jsonBytes) + "\n", nil
}

func queryKagi(apiKey, query string, timeout int) (*FastGPTResponse, error) {
	reqBody := FastGPTRequest{
		Query:     query,
		WebSearch: webSearchEnabled,
		Cache:     cacheEnabled,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", authHeaderPrefix+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("request timeout exceeded (%ds)", timeout)
		}
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError FastGPTError
		if json.Unmarshal(body, &apiError) == nil && len(apiError.Error) > 0 {
			errMsg := apiError.Error[0].Msg
			errCode := apiError.Error[0].Code

			// Provide specific error messages for common status codes
			switch resp.StatusCode {
			case 401, 403:
				return nil, fmt.Errorf("API request failed [%d]: Invalid API key", errCode)
			case 429:
				return nil, fmt.Errorf("API rate limit exceeded, try again later")
			default:
				return nil, fmt.Errorf("API request failed [%d]: %s", errCode, errMsg)
			}
		}

		// Generic HTTP error if we can't parse the error response
		return nil, fmt.Errorf("API returned HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var apiResp FastGPTResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if apiResp.Data.Output == "" {
		return nil, fmt.Errorf("API returned empty response")
	}

	return &apiResp, nil
}
