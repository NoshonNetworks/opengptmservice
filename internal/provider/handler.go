package provider

import (
	"log"
)


func NewProviderHandler(inferenceService *inference.InferenceService) *ProviderHandler {
	return &ProviderHandler{
		inferenceService: inferenceService,
	}
}

type ProviderHandler struct {
	inferenceService *inference.InferenceService
}
