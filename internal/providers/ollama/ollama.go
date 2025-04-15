package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"opengptmservice/pkg/models"
	"time"

	"go.uber.org/zap"
)

// OllamaProvider implements the Provider interface for Ollama
type OllamaProvider struct {
	baseURL      string
	defaultModel string
	client       *http.Client
	log          *zap.Logger
}

// NewProvider creates a new Ollama provider
func NewProvider(baseURL, defaultModel string, log *zap.Logger) *OllamaProvider {
	return &OllamaProvider{
		baseURL:      baseURL,
		defaultModel: defaultModel,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		log: log,
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type OllamaChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Message Message `json:"message"`
}

func (p *OllamaProvider) Generate(prompt string, model string) (string, error) {
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := p.client.Post(
		fmt.Sprintf("%s/api/generate", p.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

func (p *OllamaProvider) ChatCompletion(request models.ChatRequest) (models.ChatResponse, error) {
	// Convert the request to Ollama's format
	ollamaReq := OllamaChatRequest{
		Model:    request.Model,
		Messages: make([]Message, len(request.Messages)),
	}

	for i, msg := range request.Messages {
		ollamaReq.Messages[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := p.client.Post(
		fmt.Sprintf("%s/api/chat", p.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ChatResponse{}, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var ollamaResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return models.ChatResponse{
		Message: models.Message{
			Role:    ollamaResp.Message.Role,
			Content: ollamaResp.Message.Content,
		},
	}, nil
}

func (p *OllamaProvider) ListModels() ([]string, error) {
	url := fmt.Sprintf("%s/api/tags", p.baseURL)
	log.Printf("Making request to: %s", url)

	resp, err := p.client.Get(url)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var models struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &models); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, fmt.Errorf("failed to decode models: %w", err)
	}

	modelNames := make([]string, len(models.Models))
	for i, m := range models.Models {
		modelNames[i] = m.Name
	}

	return modelNames, nil
}

func (p *OllamaProvider) GetModelInfo(model string) (models.ModelInfo, error) {
	// For Ollama, we'll need to implement this based on available information
	// This is a placeholder implementation
	return models.ModelInfo{
		Name:        model,
		Description: fmt.Sprintf("Ollama model: %s", model),
		ContextSize: 4096, // Default context size
		Parameters:  0,    // Unknown by default
	}, nil
}
