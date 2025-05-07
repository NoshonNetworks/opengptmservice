package models

import "time"

// UserProfile represents a GitHub user profile
type UserProfile struct {
	Login           string `json:"login"`
	Name            string `json:"name"`
	Bio             string `json:"bio"`
	PublicRepos     int    `json:"public_repos"`
	PublicGists     int    `json:"public_gists"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	CreatedAt       string `json:"created_at"`
	Location        string `json:"location"`
	Company         string `json:"company"`
	Blog            string `json:"blog"`
	TwitterUsername string `json:"twitter_username"`
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topics      []string  `json:"topics"`
	Owner       struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"owner"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Title        string    `json:"title"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Repo         string    `json:"repository_url"`
	Number       int       `json:"number"`
	Additions    int       `json:"additions"`
	Deletions    int       `json:"deletions"`
	ChangedFiles int       `json:"changed_files"`
}

// Organization represents a GitHub organization
type Organization struct {
	Login       string `json:"login"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	Type        string `json:"type"`
}

// GitHubData represents all data collected from GitHub
type GitHubData struct {
	Profile       *UserProfile
	Repositories  []Repository
	Organizations []Organization
	PullRequests  []PullRequest
}
