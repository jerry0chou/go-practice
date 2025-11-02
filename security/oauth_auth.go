package security

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthProvider represents different OAuth providers
type OAuthProvider string

const (
	GoogleProvider   OAuthProvider = "google"
	GitHubProvider   OAuthProvider = "github"
	FacebookProvider OAuthProvider = "facebook"
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuthUserInfo represents user information from OAuth provider
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

// OAuthAuth handles OAuth authentication
type OAuthAuth struct {
	configs map[OAuthProvider]OAuthConfig
	client  *http.Client
}

// NewOAuthAuth creates a new OAuth authentication instance
func NewOAuthAuth() *OAuthAuth {
	return &OAuthAuth{
		configs: make(map[OAuthProvider]OAuthConfig),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AddProvider adds an OAuth provider configuration
func (o *OAuthAuth) AddProvider(provider OAuthProvider, config OAuthConfig) {
	o.configs[provider] = config
}

// GetAuthURL generates the authorization URL for OAuth flow
func (o *OAuthAuth) GetAuthURL(provider OAuthProvider, state string) (string, error) {
	config, exists := o.configs[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	baseURL := o.getProviderAuthURL(provider)
	params := url.Values{
		"client_id":     {config.ClientID},
		"redirect_uri":  {config.RedirectURL},
		"scope":         {strings.Join(config.Scopes, " ")},
		"response_type": {"code"},
		"state":         {state},
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (o *OAuthAuth) ExchangeCodeForToken(provider OAuthProvider, code string) (string, error) {
	config, exists := o.configs[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	tokenURL := o.getProviderTokenURL(provider)
	data := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {config.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	resp, err := o.client.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("OAuth error: %s", tokenResp.Error)
	}

	return tokenResp.AccessToken, nil
}

// GetUserInfo retrieves user information using access token
func (o *OAuthAuth) GetUserInfo(provider OAuthProvider, accessToken string) (*OAuthUserInfo, error) {
	userInfoURL := o.getProviderUserInfoURL(provider)

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	userInfo.Provider = string(provider)
	return &userInfo, nil
}

// getProviderAuthURL returns the authorization URL for each provider
func (o *OAuthAuth) getProviderAuthURL(provider OAuthProvider) string {
	switch provider {
	case GoogleProvider:
		return "https://accounts.google.com/o/oauth2/v2/auth"
	case GitHubProvider:
		return "https://github.com/login/oauth/authorize"
	case FacebookProvider:
		return "https://www.facebook.com/v18.0/dialog/oauth"
	default:
		return ""
	}
}

// getProviderTokenURL returns the token exchange URL for each provider
func (o *OAuthAuth) getProviderTokenURL(provider OAuthProvider) string {
	switch provider {
	case GoogleProvider:
		return "https://oauth2.googleapis.com/token"
	case GitHubProvider:
		return "https://github.com/login/oauth/access_token"
	case FacebookProvider:
		return "https://graph.facebook.com/v18.0/oauth/access_token"
	default:
		return ""
	}
}

// getProviderUserInfoURL returns the user info URL for each provider
func (o *OAuthAuth) getProviderUserInfoURL(provider OAuthProvider) string {
	switch provider {
	case GoogleProvider:
		return "https://www.googleapis.com/oauth2/v2/userinfo"
	case GitHubProvider:
		return "https://api.github.com/user"
	case FacebookProvider:
		return "https://graph.facebook.com/me?fields=id,name,email,picture"
	default:
		return ""
	}
}

// ValidateState validates the state parameter to prevent CSRF attacks
func (o *OAuthAuth) ValidateState(expectedState, receivedState string) error {
	if expectedState != receivedState {
		return fmt.Errorf("invalid state parameter")
	}
	return nil
}
