package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// AtomaConfig holds the configuration for Atoma API
type AtomaConfig struct {
	APIKey      string
	Model       string
	BaseURL     string
	MaxTokens   int
	Temperature float64
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

// NewAtomaConfig creates a new Atoma configuration from viper
func NewAtomaConfig() *AtomaConfig {
	return &AtomaConfig{
		APIKey:      viper.GetString("atoma.api_key"),
		Model:       viper.GetString("atoma.model"),
		BaseURL:     "https://api.atoma.network/v1/chat/completions",
		MaxTokens:   2000,
		Temperature: 0.7,
		Timeout:     120 * time.Second, // Increased to 2 minutes
		MaxRetries:  3,                 // Maximum number of retries
		RetryDelay:  5 * time.Second,   // Initial delay between retries
	}
}

// AtomaRequest represents the request body for Atoma API
type AtomaRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Message represents a message in the chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AtomaResponse represents the response from Atoma API
type AtomaResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// AtomaClient handles communication with Atoma API
type AtomaClient struct {
	config *AtomaConfig
	client *http.Client
}

// NewAtomaClient creates a new Atoma client
func NewAtomaClient(config *AtomaConfig) *AtomaClient {
	return &AtomaClient{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// GenerateText sends a request to Atoma API and returns the generated text
func (c *AtomaClient) GenerateText(prompt string) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff delay
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * c.config.RetryDelay
			log.Printf("Retry attempt %d/%d after %v delay", attempt, c.config.MaxRetries, delay)
			time.Sleep(delay)
		}

		startTime := time.Now()
		log.Printf("Starting Atoma API request (attempt %d/%d) at %v", attempt+1, c.config.MaxRetries+1, startTime)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
		defer cancel()

		// Prepare the request
		reqBody := AtomaRequest{
			Model: c.config.Model,
			Messages: []Message{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Temperature: c.config.Temperature,
			MaxTokens:   c.config.MaxTokens,
		}

		// Convert request body to JSON
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			lastErr = fmt.Errorf("failed to marshal request: %v", err)
			continue
		}

		// Create HTTP request with context
		req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %v", err)
			continue
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

		// Send request
		log.Printf("Sending request to Atoma API...")
		resp, err := c.client.Do(req)
		if err != nil {
			if err == context.DeadlineExceeded {
				lastErr = fmt.Errorf("request timed out after %v", c.config.Timeout)
				continue
			}
			lastErr = fmt.Errorf("failed to send request: %v", err)
			continue
		}

		// Log response status and timing
		elapsed := time.Since(startTime)
		log.Printf("Atoma API response received after %v", elapsed)
		log.Printf("Atoma API response status: %d", resp.StatusCode)

		// Read response body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close body immediately after reading
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %v", err)
			continue
		}

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			// Don't retry on 4xx errors (client errors)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", lastErr
			}
			continue
		}

		// Parse response
		var atomaResp AtomaResponse
		if err := json.Unmarshal(body, &atomaResp); err != nil {
			lastErr = fmt.Errorf("failed to parse response: %v", err)
			continue
		}

		// Check for errors in response
		if atomaResp.Error != nil {
			lastErr = fmt.Errorf("API returned error: %s", atomaResp.Error.Message)
			continue
		}

		// Return the generated text
		if len(atomaResp.Choices) > 0 {
			totalTime := time.Since(startTime)
			log.Printf("Successfully generated text in %v", totalTime)
			return atomaResp.Choices[0].Message.Content, nil
		}

		lastErr = fmt.Errorf("no response content generated")
	}

	return "", fmt.Errorf("all retry attempts failed: %v", lastErr)
}
