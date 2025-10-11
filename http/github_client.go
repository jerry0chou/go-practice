package http

import (
	"fmt"
	"time"
)

type GitHubUser struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GitHubRepo struct {
	ID               int        `json:"id"`
	NodeID           string     `json:"node_id"`
	Name             string     `json:"name"`
	FullName         string     `json:"full_name"`
	Private          bool       `json:"private"`
	Owner            GitHubUser `json:"owner"`
	HTMLURL          string     `json:"html_url"`
	Description      string     `json:"description"`
	Fork             bool       `json:"fork"`
	URL              string     `json:"url"`
	ForksURL         string     `json:"forks_url"`
	KeysURL          string     `json:"keys_url"`
	CollaboratorsURL string     `json:"collaborators_url"`
	TeamsURL         string     `json:"teams_url"`
	HooksURL         string     `json:"hooks_url"`
	IssueEventsURL   string     `json:"issue_events_url"`
	EventsURL        string     `json:"events_url"`
	AssigneesURL     string     `json:"assignees_url"`
	BranchesURL      string     `json:"branches_url"`
	TagsURL          string     `json:"tags_url"`
	BlobsURL         string     `json:"blobs_url"`
	GitTagsURL       string     `json:"git_tags_url"`
	GitRefsURL       string     `json:"git_refs_url"`
	TreesURL         string     `json:"trees_url"`
	StatusesURL      string     `json:"statuses_url"`
	LanguagesURL     string     `json:"languages_url"`
	StargazersCount  int        `json:"stargazers_count"`
	WatchersCount    int        `json:"watchers_count"`
	Language         string     `json:"language"`
	ForksCount       int        `json:"forks_count"`
	Archived         bool       `json:"archived"`
	Disabled         bool       `json:"disabled"`
	OpenIssuesCount  int        `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	AllowForking  bool      `json:"allow_forking"`
	IsTemplate    bool      `json:"is_template"`
	Topics        []string  `json:"topics"`
	Visibility    string    `json:"visibility"`
	Forks         int       `json:"forks"`
	OpenIssues    int       `json:"open_issues"`
	Watchers      int       `json:"watchers"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PushedAt      time.Time `json:"pushed_at"`
}

type GitHubClient struct {
	client  *HTTPClient
	baseURL string
	token   string
}

func NewGitHubClient(token string) *GitHubClient {
	client := NewHTTPClient("https://api.github.com")

	client.SetHeaders(map[string]string{
		"Accept":        "application/vnd.github.v3+json",
		"User-Agent":    "Go-GitHub-Client/1.0",
		"Authorization": fmt.Sprintf("token %s", token),
	})

	return &GitHubClient{
		client:  client,
		baseURL: "https://api.github.com",
		token:   token,
	}
}

func NewGitHubClientWithoutAuth() *GitHubClient {
	client := NewHTTPClient("https://api.github.com")

	client.SetHeaders(map[string]string{
		"Accept":     "application/vnd.github.v3+json",
		"User-Agent": "Go-GitHub-Client/1.0",
	})

	return &GitHubClient{
		client:  client,
		baseURL: "https://api.github.com",
		token:   "",
	}
}

func (gc *GitHubClient) GetUser(username string) (*GitHubUser, error) {
	var user GitHubUser
	err := gc.client.GetJSON(fmt.Sprintf("/users/%s", username), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", username, err)
	}
	return &user, nil
}

func (gc *GitHubClient) GetUserRepos(username string) ([]GitHubRepo, error) {
	var repos []GitHubRepo
	err := gc.client.GetJSON(fmt.Sprintf("/users/%s/repos", username), &repos)
	if err != nil {
		return nil, fmt.Errorf("failed to get repos for user %s: %w", username, err)
	}
	return repos, nil
}

func (gc *GitHubClient) GetRepo(owner, repo string) (*GitHubRepo, error) {
	var repository GitHubRepo
	err := gc.client.GetJSON(fmt.Sprintf("/repos/%s/%s", owner, repo), &repository)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo %s/%s: %w", owner, repo, err)
	}
	return &repository, nil
}

func (gc *GitHubClient) SearchUsers(query string) ([]GitHubUser, error) {
	var result struct {
		TotalCount int          `json:"total_count"`
		Items      []GitHubUser `json:"items"`
	}

	err := gc.client.GetJSON(fmt.Sprintf("/search/users?q=%s", query), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return result.Items, nil
}

func (gc *GitHubClient) SearchRepos(query string) ([]GitHubRepo, error) {
	var result struct {
		TotalCount int          `json:"total_count"`
		Items      []GitHubRepo `json:"items"`
	}

	err := gc.client.GetJSON(fmt.Sprintf("/search/repositories?q=%s", query), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to search repos: %w", err)
	}

	return result.Items, nil
}

func (gc *GitHubClient) GetRateLimit() (map[string]any, error) {
	var rateLimit map[string]any
	err := gc.client.GetJSON("/rate_limit", &rateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit: %w", err)
	}
	return rateLimit, nil
}

func ExampleGitHubAPI() {
	fmt.Println("=== GitHub API Examples ===")

	client := NewGitHubClientWithoutAuth()
	fmt.Println("\n1. Get user information:")
	user, err := client.GetUser("octocat")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("User: %s (%s)\n", user.Name, user.Login)
		fmt.Printf("Bio: %s\n", user.Bio)
		fmt.Printf("Public Repos: %d\n", user.PublicRepos)
		fmt.Printf("Followers: %d, Following: %d\n", user.Followers, user.Following)
		fmt.Printf("Profile: %s\n", user.HTMLURL)
	}

	fmt.Println("\n2. Get user repositories:")
	repos, err := client.GetUserRepos("octocat")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d repositories:\n", len(repos))
		for i, repo := range repos {
			if i >= 5 { // Show only first 5
				fmt.Printf("... and %d more\n", len(repos)-5)
				break
			}
			fmt.Printf("  - %s (%s) - %d stars\n",
				repo.Name, repo.Language, repo.StargazersCount)
		}
	}

	fmt.Println("\n3. Search for repositories:")
	searchRepos, err := client.SearchRepos("golang language:go")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d repositories:\n", len(searchRepos))
		for i, repo := range searchRepos {
			if i >= 3 { // Show only first 3
				fmt.Printf("... and %d more\n", len(searchRepos)-3)
				break
			}
			fmt.Printf("  - %s/%s - %d stars\n",
				repo.Owner.Login, repo.Name, repo.StargazersCount)
		}
	}

	fmt.Println("\n4. Check rate limit:")
	rateLimit, err := client.GetRateLimit()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Rate limit info: %+v\n", rateLimit)
	}
}

func ExampleGitHubWithAuth(token string) {
	if token == "" {
		fmt.Println("No GitHub token provided, skipping authenticated examples")
		return
	}

	fmt.Println("\n=== GitHub API with Authentication ===")

	client := NewGitHubClient(token)
	fmt.Println("\n1. Get authenticated user:")
	user, err := client.GetUser("") // Empty string gets authenticated user
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Authenticated as: %s (%s)\n", user.Name, user.Login)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Public Repos: %d\n", user.PublicRepos)
	}

	fmt.Println("\n2. Check rate limit with auth:")
	rateLimit, err := client.GetRateLimit()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Rate limit info: %+v\n", rateLimit)
	}
}

func SimpleGitHubUserInfo(username string) {
	fmt.Printf("Getting info for GitHub user: %s\n", username)

	client := NewGitHubClientWithoutAuth()

	user, err := client.GetUser(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Login: %s\n", user.Login)
	fmt.Printf("Bio: %s\n", user.Bio)
	fmt.Printf("Location: %s\n", user.Location)
	fmt.Printf("Public Repos: %d\n", user.PublicRepos)
	fmt.Printf("Followers: %d\n", user.Followers)
	fmt.Printf("Following: %d\n", user.Following)
	fmt.Printf("Profile URL: %s\n", user.HTMLURL)

	if user.CreatedAt != (time.Time{}) {
		fmt.Printf("Joined: %s\n", user.CreatedAt.Format("2006-01-02"))
	}
}
