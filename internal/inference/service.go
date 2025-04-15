package inference

import (
	"time"

	"opengptmservice/pkg/models"

	"go.uber.org/zap"
)

// Service handles the business logic for inference operations
type Service struct {
	provider models.Provider
	log      *zap.Logger
}

// NewService creates a new inference service
func NewService(provider models.Provider, log *zap.Logger) *Service {
	return &Service{
		provider: provider,
		log:      log,
	}
}

// GenerateText generates text using the specified model
func (s *Service) GenerateText(prompt string, model string) (string, error) {
	start := time.Now()
	s.log.Info("Generating text",
		zap.String("model", model),
		zap.String("prompt", prompt))

	result, err := s.provider.Generate(prompt, model)
	if err != nil {
		s.log.Error("Failed to generate text",
			zap.Error(err),
			zap.String("model", model))
		return "", err
	}

	s.log.Info("Text generation completed",
		zap.String("model", model),
		zap.Duration("duration", time.Since(start)))
	return result, nil
}

// ChatCompletion handles chat completion requests
func (s *Service) ChatCompletion(messages []models.Message, model string) (models.ChatResponse, error) {
	start := time.Now()
	s.log.Info("Processing chat completion",
		zap.String("model", model),
		zap.Int("message_count", len(messages)))

	request := models.ChatRequest{
		Model:    model,
		Messages: messages,
	}

	result, err := s.provider.ChatCompletion(request)
	if err != nil {
		s.log.Error("Failed to process chat completion",
			zap.Error(err),
			zap.String("model", model))
		return models.ChatResponse{}, err
	}

	s.log.Info("Chat completion completed",
		zap.String("model", model),
		zap.Duration("duration", time.Since(start)))
	return result, nil
}

// ListModels returns a list of available models
func (s *Service) ListModels() ([]string, error) {
	s.log.Info("Listing available models")
	return s.provider.ListModels()
}

// GetModelInfo returns information about a specific model
func (s *Service) GetModelInfo(model string) (models.ModelInfo, error) {
	s.log.Info("Getting model info", zap.String("model", model))
	return s.provider.GetModelInfo(model)
}
