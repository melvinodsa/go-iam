package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"golang.org/x/oauth2"
)

type authProvider struct {
	cnf oauth2.Config
}

func NewAuthProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	oauthConfig := oauth2.Config{
		ClientID:     p.GetParam("@GITHUB/CLIENT_ID"),
		ClientSecret: p.GetParam("@GITHUB/CLIENT_SECRET"),
		RedirectURL:  p.GetParam("@GITHUB/REDIRECT_URL"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	return authProvider{cnf: oauthConfig}
}

func (g authProvider) HasRefreshTokenFlow() bool {
	return false
}

func (g authProvider) GetAuthCodeUrl(state string) string {
	return g.cnf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g authProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	token, err := g.cnf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error verifying the code with github exchange. %w", err)
	}
	return &sdk.AuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (g authProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	return nil, nil
}

type GitHubIdentityEmail struct {
	Email string `json:"email"`
}

func (g GitHubIdentityEmail) UpdateUserDetails(user *sdk.User) {
	user.Email = g.Email
}

type GitHubIdentityName struct {
	Name string `json:"name"`
}

func (g GitHubIdentityName) UpdateUserDetails(user *sdk.User) {
	user.Name = g.Name
}

type GitHubIdentityProfilePic struct {
	ProfilePic string `json:"avatar_url"`
}

func (g GitHubIdentityProfilePic) UpdateUserDetails(user *sdk.User) {
	user.ProfilePic = g.ProfilePic
}

type GitHubIdentityUsername struct {
	Username string `json:"login"`
}

func (g GitHubIdentityUsername) UpdateUserDetails(user *sdk.User) {
	// We can store the GitHub username in a custom field or use it as display name
	if user.Name == "" {
		user.Name = g.Username
	}
}

func (g authProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	// Get user info from GitHub API
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request. %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching the identity. %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response. %w", err)
	}

	fmt.Println(string(respBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(respBytes))
	}

	var userInfo struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Login     string `json:"login"`
	}
	err = json.Unmarshal(respBytes, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
	}

	fmt.Println("GitHub User Info:", userInfo)

	identities := []sdk.AuthIdentity{
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: userInfo.Email}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityName{Name: userInfo.Name}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityProfilePic{ProfilePic: userInfo.AvatarURL}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityUsername{Username: userInfo.Login}},
	}

	// GitHub API might not return email in the user endpoint if it's private
	// In that case, we need to fetch emails separately
	if userInfo.Email == "" {
		emailIdentities, err := g.getGitHubEmails(token)
		if err == nil && len(emailIdentities) > 0 {
			// Replace the first identity with the primary email
			identities[0] = emailIdentities[0]
		}
	}

	return identities, nil
}

func (g authProvider) getGitHubEmails(token string) ([]sdk.AuthIdentity, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating emails request. %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching emails. %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching emails, status: %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading emails response. %w", err)
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	err = json.Unmarshal(respBytes, &emails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling emails response. %w", err)
	}

	// Find primary email
	for _, email := range emails {
		if email.Primary {
			return []sdk.AuthIdentity{
				{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: email.Email}},
			}, nil
		}
	}

	// If no primary email found, return the first one
	if len(emails) > 0 {
		return []sdk.AuthIdentity{
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: emails[0].Email}},
		}, nil
	}

	return nil, fmt.Errorf("no emails found")
}
