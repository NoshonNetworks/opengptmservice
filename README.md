# Developer CV Generator

A Go-based web application that generates professional CVs for developers based on their GitHub profile. The application analyzes your repositories, contributions, and technical skills to create a tailored CV.

## Features

- GitHub OAuth integration
- Repository analysis
- AI-powered CV generation using GPT
- Markdown CV export
- Modern, responsive UI

## Prerequisites

- Go 1.21 or later
- GitHub OAuth App credentials
- OpenAI API key

## Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/developer-cv-generator.git
cd developer-cv-generator
```

2. Create a GitHub OAuth App:

   - Go to GitHub Settings > Developer Settings > OAuth Apps
   - Create a new OAuth App
   - Set the callback URL to `http://localhost:8080/auth/github/callback`
   - Copy the Client ID and Client Secret

3. Create a `config.yaml` file in the project root:

```yaml
server:
  port: 8080

github:
  client_id: "your_github_client_id"
  client_secret: "your_github_client_secret"
  redirect_url: "http://localhost:8080/auth/github/callback"

openai:
  api_key: "your_openai_api_key"
  model: "gpt-3.5-turbo"
```

4. Install dependencies:

```bash
go mod download
```

5. Run the application:

```bash
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`.

## Usage

1. Visit `http://localhost:8080`
2. Click "Login with GitHub"
3. Authorize the application
4. Wait for the CV generation
5. Download your CV in Markdown format

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   └── github.go
│   ├── github/
│   │   └── repo.go
│   └── gpt/
│       └── cv.go
├── web/
│   ├── static/
│   └── templates/
│       └── index.html
├── config.yaml
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
