package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	defaultOpenRouterURL   = "https://openrouter.ai/api/v1/chat/completions"
	defaultOpenRouterModel = "openrouter/quasar-alpha"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Default HTTP client
var client HTTPClient = &http.Client{}

// OpenRouter request/response payloads

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// GenerateCommand calls the LLM API and returns the shell command.
func GenerateCommand(nlQuery string) (string, error) {
	if strings.TrimSpace(nlQuery) == "" {
		return "", errors.New("empty query")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", errors.New("missing OPENROUTER_API_KEY")
	}

	openRouterURL := os.Getenv("OPENROUTER_URL")
	if openRouterURL == "" {
		openRouterURL = defaultOpenRouterURL
	}

	openRouterModel := os.Getenv("OPENROUTER_MODEL")
	if openRouterModel == "" {
		openRouterModel = defaultOpenRouterModel
	}

	userPrompt := "Convert the following natural language instruction into a correct shell command only. Do not include explanations.\nInstruction: " + nlQuery

	osType := runtime.GOOS
	systemPrompt := fmt.Sprintf("You are a helpful assistant that translates natural language to shell commands. Generate commands for %s operating system only. If a command is not available in the target OS, suggest an alternative that provides similar functionality.", osType)

	requestBody := ChatRequest{
		Model: openRouterModel,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   100,
		Temperature: 0,
	}

	bodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/rsrini7/aishell") // Replace with your actual repo
	req.Header.Set("X-Title", "AI Shell Command Generator")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	command := chatResp.Choices[0].Message.Content
	// Clean code block formatting if present
	command = stripCodeBlock(command)
	return command, nil
}

func stripCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```bash")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}
