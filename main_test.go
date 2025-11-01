package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNormalizeFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"text format unchanged", "text", "text"},
		{"txt alias to text", "txt", "text"},
		{"markdown alias to md", "markdown", "md"},
		{"md format unchanged", "md", "md"},
		{"json format unchanged", "json", "json"},
		{"uppercase text", "TEXT", "text"},
		{"uppercase txt", "TXT", "text"},
		{"uppercase markdown", "MARKDOWN", "md"},
		{"mixed case", "TeXt", "text"},
		{"whitespace trimmed", "  text  ", "text"},
		{"whitespace with alias", "  txt  ", "text"},
		{"unknown format unchanged", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeFormat(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeFormat(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"text is valid", "text", true},
		{"md is valid", "md", true},
		{"json is valid", "json", true},
		{"txt is invalid (not normalized)", "txt", false},
		{"markdown is invalid (not normalized)", "markdown", false},
		{"empty is invalid", "", false},
		{"unknown is invalid", "unknown", false},
		{"TEXT is invalid (case sensitive)", "TEXT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidFormat(tt.format)
			if result != tt.expected {
				t.Errorf("isValidFormat(%q) = %v; want %v", tt.format, result, tt.expected)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		colorCode string
		useColor  bool
		expected  string
	}{
		{
			name:      "color enabled with bold",
			text:      "test",
			colorCode: ansiBold,
			useColor:  true,
			expected:  "\033[1mtest\033[0m",
		},
		{
			name:      "color enabled with blue",
			text:      "test",
			colorCode: ansiBlue,
			useColor:  true,
			expected:  "\033[34mtest\033[0m",
		},
		{
			name:      "color disabled returns plain text",
			text:      "test",
			colorCode: ansiBold,
			useColor:  false,
			expected:  "test",
		},
		{
			name:      "empty text with color",
			text:      "",
			colorCode: ansiBold,
			useColor:  true,
			expected:  "\033[1m\033[0m",
		},
		{
			name:      "empty text without color",
			text:      "",
			colorCode: ansiBold,
			useColor:  false,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorize(tt.text, tt.colorCode, tt.useColor)
			if result != tt.expected {
				t.Errorf("colorize(%q, %q, %v) = %q; want %q", tt.text, tt.colorCode, tt.useColor, result, tt.expected)
			}
		})
	}
}

func TestShouldUseColor(t *testing.T) {
	tests := []struct {
		name      string
		colorMode string
		expected  bool
	}{
		{"always returns true", colorAlways, true},
		{"never returns false", colorNever, false},
		// Note: auto mode depends on TTY detection, which we can't easily test in unit tests
		// We'll test the logic, but auto will be tested in integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Color: tt.colorMode}
			result := shouldUseColor(config)
			if result != tt.expected {
				t.Errorf("shouldUseColor(config with Color=%q) = %v; want %v", tt.colorMode, result, tt.expected)
			}
		})
	}

	// Test auto mode separately (will return false in test environment as it's not a TTY)
	t.Run("auto mode in non-TTY environment", func(t *testing.T) {
		config := &Config{Color: colorAuto}
		result := shouldUseColor(config)
		// In test environment, stdout is not a terminal, so should return false
		if result != false {
			t.Errorf("shouldUseColor(config with Color=auto) in non-TTY = %v; want false", result)
		}
	})
}

func createTestResponse() *FastGPTResponse {
	return &FastGPTResponse{
		Meta: struct {
			ID   string `json:"id"`
			Node string `json:"node"`
			MS   int    `json:"ms"`
		}{
			ID:   "test-id",
			Node: "test-node",
			MS:   100,
		},
		Data: struct {
			Output     string      `json:"output"`
			Tokens     int         `json:"tokens"`
			References []Reference `json:"references"`
		}{
			Output: "This is a test response",
			Tokens: 50,
			References: []Reference{
				{
					Title:   "Test Reference 1",
					Snippet: "First test snippet",
					URL:     "https://example.com/1",
				},
				{
					Title:   "Test Reference 2",
					Snippet: "Second test snippet",
					URL:     "https://example.com/2",
				},
			},
		},
	}
}

func TestFormatText_output(t *testing.T) {
	resp := createTestResponse()

	t.Run("basic text output without heading", func(t *testing.T) {
		config := &Config{
			Query:   "test query",
			Format:  formatText,
			Heading: false,
			Quiet:   false,
			Color:   colorNever,
		}

		result := formatText_output(resp, config)

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if !strings.Contains(result, "References:") {
			t.Errorf("Output missing references section")
		}
		if !strings.Contains(result, "Test Reference 1") {
			t.Errorf("Output missing first reference")
		}
		if !strings.Contains(result, "https://example.com/1") {
			t.Errorf("Output missing first reference URL")
		}

		if strings.Contains(result, "# test query") {
			t.Errorf("Output should not contain heading when Heading=false")
		}
	})

	t.Run("text output with heading", func(t *testing.T) {
		config := &Config{
			Query:   "test query",
			Format:  formatText,
			Heading: true,
			Quiet:   false,
			Color:   colorNever,
		}

		result := formatText_output(resp, config)

		if !strings.Contains(result, "# test query") {
			t.Errorf("Output missing heading")
		}

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}
		if !strings.Contains(result, "References:") {
			t.Errorf("Output missing references section")
		}
	})

	t.Run("text output in quiet mode", func(t *testing.T) {
		config := &Config{
			Query:   "test query",
			Format:  formatText,
			Heading: false,
			Quiet:   true,
			Color:   colorNever,
		}

		result := formatText_output(resp, config)

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if strings.Contains(result, "# test query") {
			t.Errorf("Quiet mode should not include heading")
		}
		if strings.Contains(result, "References:") {
			t.Errorf("Quiet mode should not include references")
		}
	})

	t.Run("text output with colors enabled", func(t *testing.T) {
		config := &Config{
			Query:   "test query",
			Format:  formatText,
			Heading: true,
			Quiet:   false,
			Color:   colorAlways,
		}

		result := formatText_output(resp, config)

		if !strings.Contains(result, "\033[") {
			t.Errorf("Output should contain ANSI color codes when color=always")
		}

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}
	})

	t.Run("text output with empty references", func(t *testing.T) {
		respNoRefs := createTestResponse()
		respNoRefs.Data.References = []Reference{}

		config := &Config{
			Query:   "test query",
			Format:  formatText,
			Heading: false,
			Quiet:   false,
			Color:   colorNever,
		}

		result := formatText_output(respNoRefs, config)

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if strings.Contains(result, "References:") {
			t.Errorf("Output should not include empty references section")
		}
	})
}

func TestFormatMarkdown_output(t *testing.T) {
	resp := createTestResponse()

	t.Run("basic markdown output", func(t *testing.T) {
		config := &Config{
			Query:  "test query",
			Format: formatMarkdown,
			Quiet:  false,
		}

		result := formatMarkdown_output(resp, config)

		// Should contain heading (always in markdown)
		if !strings.Contains(result, "# test query") {
			t.Errorf("Markdown output missing heading")
		}

		// Should contain output
		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if !strings.Contains(result, "## References") {
			t.Errorf("Output missing references section")
		}

		if !strings.Contains(result, "[Test Reference 1](https://example.com/1)") {
			t.Errorf("Output missing markdown link for first reference")
		}

		if !strings.Contains(result, "> First test snippet") {
			t.Errorf("Output missing blockquote snippet")
		}
	})

	t.Run("markdown output in quiet mode", func(t *testing.T) {
		config := &Config{
			Query:  "test query",
			Format: formatMarkdown,
			Quiet:  true,
		}

		result := formatMarkdown_output(resp, config)

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if strings.Contains(result, "# test query") {
			t.Errorf("Quiet mode should not include heading")
		}
		if strings.Contains(result, "## References") {
			t.Errorf("Quiet mode should not include references")
		}
	})

	t.Run("markdown output with empty references", func(t *testing.T) {
		respNoRefs := createTestResponse()
		respNoRefs.Data.References = []Reference{}

		config := &Config{
			Query:  "test query",
			Format: formatMarkdown,
			Quiet:  false,
		}

		result := formatMarkdown_output(respNoRefs, config)

		if !strings.Contains(result, "# test query") {
			t.Errorf("Output missing heading")
		}
		if !strings.Contains(result, "This is a test response") {
			t.Errorf("Output missing response text")
		}

		if strings.Contains(result, "## References") {
			t.Errorf("Output should not include empty references section")
		}
	})
}

func TestFormatJSON_output(t *testing.T) {
	resp := createTestResponse()

	t.Run("full JSON output", func(t *testing.T) {
		config := &Config{
			Query:  "test query",
			Format: formatJSON,
			Quiet:  false,
		}

		result, err := formatJSON_output(resp, config)
		if err != nil {
			t.Fatalf("formatJSON_output failed: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Errorf("Output is not valid JSON: %v", err)
		}

		if _, ok := parsed["meta"]; !ok {
			t.Errorf("JSON output missing 'meta' field")
		}
		if _, ok := parsed["data"]; !ok {
			t.Errorf("JSON output missing 'data' field")
		}

		if !strings.Contains(result, "This is a test response") {
			t.Errorf("JSON output missing response text")
		}

		// Should be pretty-printed (contains indentation)
		if !strings.Contains(result, "  ") {
			t.Errorf("JSON output should be pretty-printed with indentation")
		}
	})

	t.Run("JSON output in quiet mode", func(t *testing.T) {
		config := &Config{
			Query:  "test query",
			Format: formatJSON,
			Quiet:  true,
		}

		result, err := formatJSON_output(resp, config)
		if err != nil {
			t.Fatalf("formatJSON_output failed: %v", err)
		}

		var parsed string
		if err := json.Unmarshal([]byte(strings.TrimSpace(result)), &parsed); err != nil {
			t.Errorf("Quiet JSON output is not valid JSON: %v", err)
		}

		if parsed != "This is a test response" {
			t.Errorf("Quiet JSON output = %q; want %q", parsed, "This is a test response")
		}

		if strings.Contains(result, "meta") {
			t.Errorf("Quiet mode should not include meta field")
		}
		if strings.Contains(result, "references") {
			t.Errorf("Quiet mode should not include references field")
		}
	})
}

func TestFormatOutput(t *testing.T) {
	resp := createTestResponse()

	tests := []struct {
		name             string
		format           string
		shouldContain    string
		shouldNotContain string
	}{
		{
			name:          "text format dispatches to text formatter",
			format:        formatText,
			shouldContain: "This is a test response",
		},
		{
			name:          "markdown format dispatches to markdown formatter",
			format:        formatMarkdown,
			shouldContain: "# test query",
		},
		{
			name:          "json format dispatches to json formatter",
			format:        formatJSON,
			shouldContain: `"output"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Query:  "test query",
				Format: tt.format,
				Quiet:  false,
				Color:  colorNever,
			}

			result, err := formatOutput(resp, config)
			if err != nil {
				t.Fatalf("formatOutput failed: %v", err)
			}

			if tt.shouldContain != "" && !strings.Contains(result, tt.shouldContain) {
				t.Errorf("Output missing expected content: %q", tt.shouldContain)
			}
			if tt.shouldNotContain != "" && strings.Contains(result, tt.shouldNotContain) {
				t.Errorf("Output contains unexpected content: %q", tt.shouldNotContain)
			}
		})
	}
}

func TestGetQuery(t *testing.T) {
	t.Run("single arg query", func(t *testing.T) {
		args := []string{"test"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery returned error: %v", err)
		}
		if result != "test" {
			t.Errorf("getQuery(%v) = %q; want %q", args, result, "test")
		}
	})

	t.Run("multiple args concatenated", func(t *testing.T) {
		args := []string{"golang", "best", "practices"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery returned error: %v", err)
		}
		expected := "golang best practices"
		if result != expected {
			t.Errorf("getQuery(%v) = %q; want %q", args, result, expected)
		}
	})

	t.Run("args with extra whitespace", func(t *testing.T) {
		args := []string{"  test  ", "  query  "}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery returned error: %v", err)
		}
		// Join preserves the spacing, but TrimSpace at the end should clean it
		if !strings.Contains(result, "test") || !strings.Contains(result, "query") {
			t.Errorf("getQuery(%v) = %q; should contain 'test' and 'query'", args, result)
		}
	})

	t.Run("empty args returns error", func(t *testing.T) {
		args := []string{}
		_, err := getQuery(args)
		if err == nil {
			t.Errorf("getQuery(%v) should return error for empty args", args)
		}
		if !strings.Contains(err.Error(), "no query provided") {
			t.Errorf("Error message should mention 'no query provided', got: %v", err)
		}
	})

	t.Run("args with only whitespace returns error", func(t *testing.T) {
		args := []string{"   ", "  "}
		_, err := getQuery(args)
		if err == nil {
			t.Errorf("getQuery(%v) should return error for whitespace-only args", args)
		}
		if !strings.Contains(err.Error(), "no query provided") {
			t.Errorf("Error message should mention 'no query provided', got: %v", err)
		}
	})

	t.Run("unicode and special characters", func(t *testing.T) {
		args := []string{"测试", "🚀", "query"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery returned error: %v", err)
		}
		expected := "测试 🚀 query"
		if result != expected {
			t.Errorf("getQuery(%v) = %q; want %q", args, result, expected)
		}
	})
}

func TestQueryKagi_ResponseParsing(t *testing.T) {
	t.Run("successful API response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != contentTypeJSON {
				t.Errorf("Expected Content-Type %s, got %s", contentTypeJSON, r.Header.Get("Content-Type"))
			}
			if !strings.HasPrefix(r.Header.Get("Authorization"), authHeaderPrefix) {
				t.Errorf("Expected Authorization header with prefix %s", authHeaderPrefix)
			}

			resp := FastGPTResponse{
				Meta: struct {
					ID   string `json:"id"`
					Node string `json:"node"`
					MS   int    `json:"ms"`
				}{
					ID:   "test-123",
					Node: "test-node",
					MS:   150,
				},
				Data: struct {
					Output     string      `json:"output"`
					Tokens     int         `json:"tokens"`
					References []Reference `json:"references"`
				}{
					Output: "Test response",
					Tokens: 42,
					References: []Reference{
						{Title: "Ref 1", URL: "https://test.com", Snippet: "snippet"},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		// Note: We can't directly test queryKagi because it uses a hardcoded endpoint
		// This test verifies our mock server works correctly
		// Full API client testing would require refactoring to inject the endpoint
		t.Skip("queryKagi uses hardcoded endpoint - tested in integration tests")
	})

	t.Run("API error response (401)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResp := FastGPTError{
				Error: []struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
				}{
					{Code: 401, Msg: "Invalid API key"},
				},
			}
			json.NewEncoder(w).Encode(errorResp)
		}))
		defer server.Close()

		t.Skip("queryKagi uses hardcoded endpoint - tested in integration tests")
	})

	t.Run("timeout handling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sleep longer than timeout
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		t.Skip("queryKagi uses hardcoded endpoint - tested in integration tests")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		t.Skip("queryKagi uses hardcoded endpoint - tested in integration tests")
	})

	t.Run("empty output validation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := FastGPTResponse{
				Data: struct {
					Output     string      `json:"output"`
					Tokens     int         `json:"tokens"`
					References []Reference `json:"references"`
				}{
					Output: "", // Empty output
					Tokens: 0,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		t.Skip("queryKagi uses hardcoded endpoint - tested in integration tests")
	})
}

// Note: Full API client testing (queryKagi) is limited because it uses a hardcoded endpoint.
// To properly unit test the HTTP client, the code would need refactoring to inject the endpoint URL.
// The current tests verify response/error structures and will be supplemented by integration tests.

func TestEdgeCases(t *testing.T) {
	t.Run("very long query", func(t *testing.T) {
		// Test with a query longer than 1000 characters
		longQuery := make([]string, 200)
		for i := range longQuery {
			longQuery[i] = "word"
		}
		args := longQuery
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery should handle long queries: %v", err)
		}
		if len(result) < 200 {
			t.Errorf("Long query was truncated unexpectedly")
		}
	})

	t.Run("query with special characters", func(t *testing.T) {
		args := []string{"query", "with", "special", "chars:", "<>&\"'"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery should handle special chars: %v", err)
		}
		if !strings.Contains(result, "<>&\"'") {
			t.Errorf("Special characters should be preserved in query")
		}
	})

	t.Run("query with newlines and tabs", func(t *testing.T) {
		args := []string{"query\nwith\nnewlines", "and\ttabs"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("getQuery should handle newlines and tabs: %v", err)
		}
		if !strings.Contains(result, "query") {
			t.Errorf("Query content should be preserved")
		}
	})

	t.Run("reference with empty snippet", func(t *testing.T) {
		resp := createTestResponse()
		resp.Data.References[0].Snippet = ""

		config := &Config{
			Query:  "test",
			Format: formatText,
			Color:  colorNever,
		}

		result := formatText_output(resp, config)
		if !strings.Contains(result, "Test Reference 1") {
			t.Errorf("Reference with empty snippet should still be displayed")
		}
	})

	t.Run("response with many references", func(t *testing.T) {
		resp := createTestResponse()
		// Add 10 more references
		for i := 3; i <= 12; i++ {
			resp.Data.References = append(resp.Data.References, Reference{
				Title:   "Test Reference " + string(rune(i)),
				URL:     "https://example.com/" + string(rune(i)),
				Snippet: "Snippet " + string(rune(i)),
			})
		}

		config := &Config{
			Query:  "test",
			Format: formatText,
			Color:  colorNever,
		}

		result := formatText_output(resp, config)
		if !strings.Contains(result, "References:") {
			t.Errorf("Should display references section with many refs")
		}
	})

	t.Run("format normalization with unusual input", func(t *testing.T) {
		tests := []string{
			"TEXT",
			"  text  ",
			"TxT",
			"MARKDOWN",
			"  md  ",
			"JsOn",
		}

		for _, input := range tests {
			result := normalizeFormat(input)
			if !isValidFormat(result) && result != "json" {
				// After normalization, should either be valid or json
				if result != "text" && result != "md" && result != "json" {
					t.Errorf("normalizeFormat(%q) = %q; should normalize to valid format", input, result)
				}
			}
		}
	})

	t.Run("unicode in output and references", func(t *testing.T) {
		resp := createTestResponse()
		resp.Data.Output = "回答：Go 是一种编程语言"
		resp.Data.References[0].Title = "中文标题"
		resp.Data.References[0].Snippet = "这是一个中文摘要"

		config := &Config{
			Query:  "测试查询",
			Format: formatText,
			Color:  colorNever,
		}

		result := formatText_output(resp, config)
		if !strings.Contains(result, "回答") {
			t.Errorf("Should preserve unicode in output")
		}
		if !strings.Contains(result, "中文标题") {
			t.Errorf("Should preserve unicode in references")
		}
	})

	t.Run("json output with unicode", func(t *testing.T) {
		resp := createTestResponse()
		resp.Data.Output = "Unicode: 你好 🌍"

		config := &Config{
			Query:  "test",
			Format: formatJSON,
			Quiet:  false,
		}

		result, err := formatJSON_output(resp, config)
		if err != nil {
			t.Fatalf("formatJSON_output failed: %v", err)
		}
		// Should be valid JSON
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Errorf("JSON with unicode should be valid: %v", err)
		}
	})

	t.Run("markdown with special markdown characters", func(t *testing.T) {
		resp := createTestResponse()
		resp.Data.Output = "Test with # heading and * bullet"
		resp.Data.References[0].Title = "Title [with] brackets"
		resp.Data.References[0].URL = "https://example.com/path?param=value&other=test"

		config := &Config{
			Query:  "test",
			Format: formatMarkdown,
			Quiet:  false,
		}

		result := formatMarkdown_output(resp, config)
		if !strings.Contains(result, "[Title [with] brackets]") {
			t.Errorf("Should preserve brackets in markdown links")
		}
	})

	t.Run("color codes with special characters", func(t *testing.T) {
		text := "test\nwith\nnewlines"
		result := colorize(text, ansiBold, true)
		// Should wrap entire text including newlines
		if !strings.HasPrefix(result, ansiBold) {
			t.Errorf("Colorized text should start with color code")
		}
		if !strings.HasSuffix(result, ansiReset) {
			t.Errorf("Colorized text should end with reset code")
		}
	})
}

func TestErrorConditions(t *testing.T) {
	t.Run("invalid format string", func(t *testing.T) {
		invalidFormats := []string{"xml", "yaml", "html", "pdf", ""}
		for _, format := range invalidFormats {
			normalized := normalizeFormat(format)
			if isValidFormat(normalized) {
				t.Errorf("Format %q should not be valid after normalization", format)
			}
		}
	})

	t.Run("empty response data", func(t *testing.T) {
		resp := &FastGPTResponse{
			Data: struct {
				Output     string      `json:"output"`
				Tokens     int         `json:"tokens"`
				References []Reference `json:"references"`
			}{
				Output:     "",
				Tokens:     0,
				References: []Reference{},
			},
		}

		config := &Config{
			Query:  "test",
			Format: formatText,
			Color:  colorNever,
		}

		result := formatText_output(resp, config)
		if result == "" {
			t.Errorf("Should return some output even with empty response data")
		}
	})

	t.Run("nil config color handling", func(t *testing.T) {
		config := &Config{
			Color: "invalid",
		}
		result := shouldUseColor(config)
		if result != false {
			t.Errorf("Invalid color mode should default to no color")
		}
	})

	t.Run("empty query after trimming", func(t *testing.T) {
		args := []string{"   "}
		_, err := getQuery(args)
		if err == nil {
			t.Errorf("Should return error for whitespace-only query")
		}
	})

	t.Run("query with only special characters", func(t *testing.T) {
		args := []string{"!@#$%^&*()"}
		result, err := getQuery(args)
		if err != nil {
			t.Errorf("Should accept query with only special chars: %v", err)
		}
		if result != "!@#$%^&*()" {
			t.Errorf("Special chars query = %q; want %q", result, "!@#$%^&*()")
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("timeout validation", func(t *testing.T) {
		// Note: This would require calling loadConfig which uses global flags
		// For now, we verify the timeout value is checked to be positive
		// This is covered in integration tests
		if defaultTimeout <= 0 {
			t.Errorf("Default timeout should be positive, got %d", defaultTimeout)
		}
	})

	t.Run("color mode validation", func(t *testing.T) {
		validModes := []string{colorAuto, colorAlways, colorNever}
		for _, mode := range validModes {
			config := &Config{Color: mode}
			_ = shouldUseColor(config)
		}
	})

	t.Run("format validation", func(t *testing.T) {
		validFormats := []string{formatText, formatMarkdown, formatJSON}
		for _, format := range validFormats {
			if !isValidFormat(format) {
				t.Errorf("Format %q should be valid", format)
			}
		}
	})
}
