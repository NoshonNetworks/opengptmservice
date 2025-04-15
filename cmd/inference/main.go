package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opengptmservice/internal/inference"
	"opengptmservice/internal/providers/ollama"
	"opengptmservice/pkg/logger"
	"opengptmservice/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	// Initialize logger
	log := logger.NewLogger(viper.GetString("logging.level"), viper.GetString("logging.format"))
	defer log.Sync()

	log.Info("Starting OpenGPTM Service")

	// Create Ollama provider
	ollamaProvider := ollama.NewProvider(
		viper.GetString("ollama.base_url"),
		viper.GetString("ollama.default_model"),
		log,
	)

	// Create inference service
	service := inference.NewService(ollamaProvider, log)

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RateLimitMiddleware(log))

	// Create handler
	handler := inference.NewHandler(service, log)

	// Register routes
	router.GET("/health", handler.HealthCheck)
	router.POST("/inference", handler.Inference)
	router.POST("/chat", handler.ChatCompletion)
	router.GET("/models", handler.ListModels)
	router.GET("/models/:model", handler.GetModelInfo)

	// Create server
	srv := &http.Server{
		Addr:    viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}
