package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"opengptmservice/internal/auth"
	"opengptmservice/internal/services"
)

func main() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Initialize services
	githubService := services.NewGitHubService()
	cvService := services.NewCVService()

	// Initialize router
	r := gin.Default()

	// Set up template functions
	r.SetFuncMap(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*")

	// Serve static files
	r.Static("/static", "./web/static")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Developer CV Generator",
		})
	})

	r.GET("/auth/github", func(c *gin.Context) {
		url := auth.GetGitHubAuthURL()
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/auth/github/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			log.Printf("No code provided in callback")
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		log.Printf("Received code: %s", code)

		// Exchange code for access token
		token, err := auth.ExchangeCodeForToken(code)
		if err != nil {
			log.Printf("Error exchanging code for token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get access token: %v", err)})
			return
		}

		log.Printf("Successfully obtained access token")

		// Get user info
		userInfo, err := auth.GetGitHubUserInfo(token)
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user info: %v", err)})
			return
		}

		log.Printf("Successfully fetched user info for: %s", userInfo["login"])

		// Get GitHub data
		githubData, err := githubService.GetUserData(token)
		if err != nil {
			log.Printf("Error getting GitHub data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get GitHub data: %v", err)})
			return
		}

		log.Printf("Successfully fetched GitHub data")

		// Generate CV
		cv, err := cvService.GenerateCV(githubData)
		if err != nil {
			log.Printf("Error generating CV: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate CV: %v", err)})
			return
		}

		log.Printf("Successfully generated CV")

		// Return the CV
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Your Developer CV",
			"cv":    cv,
			"user":  userInfo,
		})
	})

	// Start server
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
