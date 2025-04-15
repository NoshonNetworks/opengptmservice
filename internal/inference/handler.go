package inference

import (
	"net/http"

	"opengptmservice/pkg/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for the inference service
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler creates a new inference handler
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	h.log.Info("Health check request received")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Inference handles text generation requests
func (h *Handler) Inference(c *gin.Context) {
	var request models.InferenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Inference request received",
		zap.String("model", request.Model),
		zap.String("prompt", request.Prompt))

	result, err := h.service.GenerateText(request.Prompt, request.Model)
	if err != nil {
		h.log.Error("Failed to generate text", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.InferenceResponse{
		Response: result,
		Model:    request.Model,
	})
}

// ChatCompletion handles chat completion requests
func (h *Handler) ChatCompletion(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Chat completion request received",
		zap.String("model", request.Model),
		zap.Int("message_count", len(request.Messages)))

	response, err := h.service.ChatCompletion(request.Messages, request.Model)
	if err != nil {
		h.log.Error("Failed to process chat completion", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListModels handles requests to list available models
func (h *Handler) ListModels(c *gin.Context) {
	h.log.Info("List models request received")
	models, err := h.service.ListModels()
	if err != nil {
		h.log.Error("Failed to list models", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// GetModelInfo handles requests to get information about a specific model
func (h *Handler) GetModelInfo(c *gin.Context) {
	model := c.Param("model")
	h.log.Info("Get model info request received", zap.String("model", model))

	info, err := h.service.GetModelInfo(model)
	if err != nil {
		h.log.Error("Failed to get model info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
