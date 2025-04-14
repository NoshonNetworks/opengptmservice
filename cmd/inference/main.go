package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opengptmservice/internal/inference"
	"opengptmservice/internal/providers/ollama"
	"opengptmservice/pkg/config"
	"opengptmservice/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		cfg = config.GetDefaultConfig()
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("Starting OpenGPTM Service")

	// Initialize Ollama provider
	ollamaProvider := ollama.NewOllamaProvider(cfg.Ollama.BaseURL)
	ollamaProvider.Client = &http.Client{
		Timeout: 5 * time.Second,
	}

	// Initialize service with the provider
	service := inference.NewService(ollamaProvider)

	// Initialize handler with the service
	handler := inference.NewHandler(service)

	// Set up router
	r := gin.Default()

	// Add middleware for logging
	r.Use(gin.Logger())

	// Register routes
	r.GET("/health", handler.HealthCheck)
	r.POST("/inference", handler.Inference)
	r.POST("/chat", handler.ChatCompletion)
	r.GET("/models", handler.ListModels)
	r.GET("/models/:model", handler.GetModelInfo)

	// Start the server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Info("Server starting", zap.String("address", addr))
		if err := r.Run(addr); err != nil {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
}
