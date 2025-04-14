# OpenGPTM Service

A Go-based inference service that provides a REST API for interacting with local LLM models through Ollama.

## Features

- REST API for model inference and chat completion
- Support for Ollama provider
- Model management (listing, info retrieval)
- Configurable logging
- Health checks
- Graceful shutdown

## Prerequisites

- Go 1.21 or later
- Ollama (for local model hosting)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/NoshonNetworks/opengptmservice.git
cd opengptmservice
```

2. Install dependencies:

```bash
go mod tidy
```

3. Install Ollama:

```bash
# For macOS
brew install ollama

# For Linux
curl -fsSL https://ollama.com/install.sh | sh
```

4. Start Ollama:

```bash
# Start Ollama in a separate terminal
ollama serve
```

5. Pull a model:

```bash
# Pull the model you want to use (e.g., llama3.2)
ollama pull llama3.2
```

## Configuration

The service can be configured using the `config/config.yaml` file. Here's the default configuration:

```yaml
server:
  port: 8081
  host: "0.0.0.0"

ollama:
  base_url: "http://localhost:11434"
  default_model: "llama3.2:latest"

logging:
  level: "info"
  format: "json"
```

## Running the Service

1. Make sure Ollama is running:

```bash
# In a separate terminal
ollama serve
```

2. Start the inference service:

```bash
go run cmd/inference/main.go
```

The service will start on the configured port (default: 8081).

## API Endpoints

### Health Check

```bash
curl http://localhost:8081/health
```

Response:

```json
{ "status": "ok" }
```

### List Available Models

```bash
curl http://localhost:8081/models
```

Response:

```json
{ "models": ["llama3.2:latest"] }
```

### Get Model Info

```bash
curl http://localhost:8081/models/llama3.2:latest
```

### Text Generation

```bash
curl -X POST http://localhost:8081/inference \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, how are you?", "model": "llama3.2:latest"}'
```

### Chat Completion

```bash
curl -X POST http://localhost:8081/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:latest",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }'
```

## Troubleshooting

### Port Conflicts

If you encounter port conflicts:

1. Check which process is using the port:

```bash
lsof -i :8081  # For inference service
lsof -i :11434 # For Ollama
```

2. Stop the conflicting process or update the port in `config/config.yaml`

### Ollama Not Found

If `ollama` command is not found:

1. Verify Ollama installation:

```bash
which ollama
```

2. Add Ollama to your PATH if not found
3. Restart your terminal

### Model Not Found

If you get 404 errors for model requests:

1. Verify the model is available:

```bash
curl http://localhost:8081/models
```

2. Pull the model if not available:

```bash
ollama pull llama3.2
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o opengptmservice cmd/inference/main.go
```

### Adding New Providers

To add a new provider:

1. Create a new package under `internal/providers/`
2. Implement the `models.Provider` interface
3. Update the main.go file to use your new provider

## License

MIT License
