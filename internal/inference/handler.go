package inference

import (
	"log"
	"net/http"

	"opengptmservice/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Inference(c *gin.Context) {
	var req models.InferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use default model if none specified
	if req.Model == "" {
		req.Model = "llama2" // This should come from config
	}

	resp, err := h.service.GenerateResponse(req.Prompt, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ChatCompletion(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ChatCompletion(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListModels(c *gin.Context) {
	log.Println("ListModels endpoint called")
	models, err := h.service.ListModels()
	if err != nil {
		log.Printf("Error listing models: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Found %d models", len(models))
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (h *Handler) GetModelInfo(c *gin.Context) {
	model := c.Param("model")
	info, err := h.service.GetModelInfo(model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
