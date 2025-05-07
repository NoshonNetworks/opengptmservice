package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"opengptmservice/internal/models"
)

// GitHubService handles all GitHub-related operations
type GitHubService struct {
	client *http.Client
}

// NewGitHubService creates a new GitHub service instance
func NewGitHubService() *GitHubService {
	return &GitHubService{
		client: &http.Client{},
	}
}

// GetUserData fetches all user data from GitHub
func (s *GitHubService) GetUserData(accessToken string) (*models.GitHubData, error) {
	profile, err := s.getUserProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}

	orgs, err := s.getUserOrganizations(accessToken)
	if err != nil {
		log.Printf("Warning: Failed to get organizations: %v", err)
	}

	prs, err := s.getUserPullRequests(accessToken)
	if err != nil {
		log.Printf("Warning: Failed to get pull requests: %v", err)
	}

	repos, err := s.getRepositories(accessToken)
	if err != nil {
		return nil, err
	}

	return &models.GitHubData{
		Profile:       profile,
		Repositories:  repos,
		Organizations: orgs,
		PullRequests:  prs,
	}, nil
}

func (s *GitHubService) getUserProfile(accessToken string) (*models.UserProfile, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %s", resp.Status)
	}

	var profile models.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *GitHubService) getUserOrganizations(accessToken string) ([]models.Organization, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/orgs", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organizations: %s", resp.Status)
	}

	var orgs []models.Organization
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, err
	}

	return orgs, nil
}

func (s *GitHubService) getUserPullRequests(accessToken string) ([]models.PullRequest, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/issues?filter=all&state=all", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pull requests: %s", resp.Status)
	}

	var issues []models.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	// Filter only pull requests
	var prs []models.PullRequest
	for _, issue := range issues {
		if issue.Repo != "" { // This indicates it's a PR
			prs = append(prs, issue)
		}
	}

	return prs, nil
}

func (s *GitHubService) getRepositories(accessToken string) ([]models.Repository, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get repositories: %s, body: %s", resp.Status, string(body))
	}

	var repos []models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	// Sort by stars
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})

	return repos, nil
}
