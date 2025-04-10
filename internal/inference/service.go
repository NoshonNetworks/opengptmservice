package inference

import (
	"log"
)

func NewInferenceService() *InferenceService {
	return &InferenceService{}
}

type InferenceService struct {
	llm *llm.LLM
}

