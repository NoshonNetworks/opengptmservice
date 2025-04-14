package inference

import (
	"time"

	"opengptmservice/pkg/models"
)

type Service struct {
	provider models.Provider
}

func NewService(provider models.Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) GenerateResponse(prompt string, model string) (models.InferenceResponse, error) {
	startTime := time.Now()

	response, err := s.provider.Generate(prompt, model)
	if err != nil {
		return models.InferenceResponse{}, err
	}

	return models.InferenceResponse{
		Response: response,
		Model:    model,
		Time:     time.Since(startTime).Milliseconds(),
	}, nil
}

func (s *Service) ChatCompletion(request models.ChatRequest) (models.ChatResponse, error) {
	startTime := time.Now()

	// Use default model if none specified
	if request.Model == "" {
		request.Model = "llama2" // This should come from config
	}

	response, err := s.provider.ChatCompletion(request)
	if err != nil {
		return models.ChatResponse{}, err
	}

	// Add timing information to the response
	// Note: We're not modifying the response structure, but we could add timing info if needed
	_ = startTime

	return response, nil
}

func (s *Service) ListModels() ([]string, error) {
	return s.provider.ListModels()
}

func (s *Service) GetModelInfo(model string) (models.ModelInfo, error) {
	return s.provider.GetModelInfo(model)
}
