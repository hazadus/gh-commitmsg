// Package llm contains utilities for working with GitHub Models API
package llm

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/auth"
	"gopkg.in/yaml.v3"
)

//go:embed commitmsg.prompt.yml
var standupPromptYAML []byte

// PromptConfig represents the structure of the prompt configuration file
// It includes the model parameters and the messages to be sent to the model.
type PromptConfig struct {
	Name            string          `yaml:"name"`
	Description     string          `yaml:"description"`
	Model           string          `yaml:"model"`
	ModelParameters ModelParameters `yaml:"modelParameters"`
	Messages        []PromptMessage `yaml:"messages"`
}

// ModelParameters defines the parameters for the model
type ModelParameters struct {
	Temperature float64 `yaml:"temperature"`
	TopP        float64 `yaml:"topP"`
}

// PromptMessage represents a single message in the prompt configuration
type PromptMessage struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

// Request represents the structure of the request to the GitHub Models API
type Request struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	Stream      bool      `json:"stream"`
}

// Message represents a single message in the request to the GitHub Models API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents the structure of the response from the GitHub Models API
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Client is a wrapper around the GitHub Models API client
type Client struct {
	token string
}

// NewClient initializes a new GitHub Models API client
// It retrieves the GitHub token from the environment installed by the `gh` CLI tool.
func NewClient() (*Client, error) {
	fmt.Print("  Checking GitHub token... ")

	host, _ := auth.DefaultHost()
	token, _ := auth.TokenForHost(host) // check GH_TOKEN, GITHUB_TOKEN, keychain, etc

	if token == "" {
		fmt.Println("Failed")
		return nil, fmt.Errorf("no GitHub token found, please run 'gh auth login' to authenticate")
	}
	fmt.Println("Done")

	return &Client{token: token}, nil
}

// GenerateCommitMessage generates a commit message based on the provided changes summary
func (c *Client) GenerateCommitMessage(changesSummary string) (string, error) {
	fmt.Print("  Loading prompt configuration... ")
	promptConfig, err := loadPromptConfig()
	if err != nil {
		fmt.Println("Failed")
		return "", err
	}
	fmt.Println("Done")

	selectedModel := promptConfig.Model

	// Build messages from the prompt config, replacing template variables
	messages := make([]Message, len(promptConfig.Messages))
	for i, msg := range promptConfig.Messages {
		content := msg.Content
		// Replace the {{changes}} template variable
		content = strings.ReplaceAll(content, "{{changes}}", changesSummary)

		messages[i] = Message{
			Role:    msg.Role,
			Content: content,
		}
	}

	request := Request{
		Messages:    messages,
		Model:       selectedModel,
		Temperature: promptConfig.ModelParameters.Temperature,
		TopP:        promptConfig.ModelParameters.TopP,
		Stream:      false,
	}

	fmt.Printf("  Calling GitHub Models API (%s)... ", selectedModel)
	response, err := c.callGitHubModels(request)
	if err != nil {
		fmt.Println("Failed")
		return "", err
	}
	fmt.Println("Done")

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response generated from the model")
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

func loadPromptConfig() (*PromptConfig, error) {
	var config PromptConfig
	err := yaml.Unmarshal(standupPromptYAML, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt configuration: %w", err)
	}
	return &config, nil
}

// callGitHubModels makes the API call to GitHub Models
func (c *Client) callGitHubModels(request Request) (*Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://models.github.ai/inference/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
