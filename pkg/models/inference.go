package models

type InferenceRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Model  string `json:"model,omitempty"`
}

type InferenceResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Time     int64  `json:"time_ms"`
}
