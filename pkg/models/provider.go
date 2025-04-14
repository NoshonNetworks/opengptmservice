package models

// Message represents a single message in a chat conversation
type Message struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Message Message `json:"message"`
}

// Provider defines the interface that all LLM providers must implement
type Provider interface {
	// Generate generates a response for the given prompt
	Generate(prompt string, model string) (string, error)

	// ChatCompletion generates a response for a chat conversation
	ChatCompletion(request ChatRequest) (ChatResponse, error)

	// ListModels returns a list of available models
	ListModels() ([]string, error)

	// GetModelInfo returns information about a specific model
	GetModelInfo(model string) (ModelInfo, error)
}

// ModelInfo contains information about a specific model
type ModelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ContextSize int    `json:"context_size"`
	Parameters  int    `json:"parameters"`
}
