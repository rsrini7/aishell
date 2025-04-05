package llm

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// mockHTTPClient replaces the default HTTP client during tests
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGenerateCommand(t *testing.T) {
	// Setup test environment
	originalAPI := os.Getenv("OPENROUTER_API_KEY")
	defer os.Setenv("OPENROUTER_API_KEY", originalAPI)

	// Create mock response
	mockResponse := `{
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "ls -la"
				}
			}
		]
	}`

	tests := []struct {
		name        string
		query       string
		apiKey      string
		response    string
		expectError bool
		errorMsg    string
		wantCommand string
	}{
		{
			name:        "Valid query with API key",
			query:       "list files",
			apiKey:      "test-api-key",
			response:    mockResponse,
			expectError: false,
			wantCommand: "ls -la",
		},
		{
			name:        "Empty query",
			query:       "",
			apiKey:      "test-api-key",
			expectError: true,
			errorMsg:    "empty query",
		},
		{
			name:        "Missing API key",
			query:       "list files",
			apiKey:      "",
			expectError: true,
			errorMsg:    "missing OPENROUTER_API_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment for this test
			os.Setenv("OPENROUTER_API_KEY", tt.apiKey)

			// Create a mock client
			mockClient := &mockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Return mock response for valid cases
					if tt.response != "" {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(tt.response)),
						}, nil
					}
					return nil, nil
				},
			}

			// Store the original client and restore it after the test
			originalClient := client
			client = mockClient
			defer func() { client = originalClient }()

			cmd, err := GenerateCommand(tt.query)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q but got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cmd != tt.wantCommand {
					t.Errorf("expected command %q but got %q", tt.wantCommand, cmd)
				}
			}
		})
	}
}
