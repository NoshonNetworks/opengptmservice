package services

import (
	"fmt"
	"log"

	"opengptmservice/internal/models"
)

// CVService handles CV generation operations
type CVService struct {
	atomaClient *AtomaClient
}

// NewCVService creates a new CV service instance
func NewCVService() *CVService {
	config := NewAtomaConfig()
	return &CVService{
		atomaClient: NewAtomaClient(config),
	}
}

// GenerateCV generates a CV based on GitHub data
func (s *CVService) GenerateCV(data *models.GitHubData) (string, error) {
	// Create a detailed prompt for CV generation
	prompt := fmt.Sprintf(`Generate a professional CV for a software developer based on their GitHub profile and activity.

User Profile:
- Name: %s
- GitHub: %s
- Bio: %s
- Location: %s
- Company: %s
- Blog: %s
- Twitter: %s
- Public Repositories: %d
- Public Gists: %d
- Followers: %d
- Following: %d
- Member since: %s

Organizations:
%s

Repositories:
%s

Pull Requests:
%s

Please generate a professional CV in markdown format with the following sections:
1. Professional Summary
2. Technical Skills (based on languages and technologies used)
3. Professional Experience (based on organizations and company)
4. Notable Projects (top 5 repositories with descriptions)
5. Open Source Contributions (based on pull requests and contributions)
6. Social Presence (GitHub, Twitter, Blog)
7. Additional Information

Format the CV in a clean, professional style using markdown. Include relevant links to GitHub repositories and pull requests.`,
		data.Profile.Name,
		data.Profile.Login,
		data.Profile.Bio,
		data.Profile.Location,
		data.Profile.Company,
		data.Profile.Blog,
		data.Profile.TwitterUsername,
		data.Profile.PublicRepos,
		data.Profile.PublicGists,
		data.Profile.Followers,
		data.Profile.Following,
		data.Profile.CreatedAt,
		s.formatOrganizations(data.Organizations),
		s.formatRepositories(data.Repositories),
		s.formatPullRequests(data.PullRequests))

	log.Printf("Generating CV with Atoma API")
	return s.atomaClient.GenerateText(prompt)
}

func (s *CVService) formatOrganizations(orgs []models.Organization) string {
	var result string
	for _, org := range orgs {
		result += fmt.Sprintf("- %s: %s\n", org.Login, org.Description)
	}
	return result
}

func (s *CVService) formatRepositories(repos []models.Repository) string {
	var result string
	for _, repo := range repos {
		result += fmt.Sprintf("- %s: %s (Language: %s, Stars: %d, Forks: %d)\n",
			repo.Name,
			repo.Description,
			repo.Language,
			repo.Stars,
			repo.Forks)
	}
	return result
}

func (s *CVService) formatPullRequests(prs []models.PullRequest) string {
	var result string
	for _, pr := range prs {
		result += fmt.Sprintf("- %s (State: %s, Repo: %s)\n",
			pr.Title,
			pr.State,
			pr.Repo)
	}
	return result
}
