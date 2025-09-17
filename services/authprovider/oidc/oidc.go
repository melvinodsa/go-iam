package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"golang.org/x/oauth2"
)

// authProvider implements the SDK ServiceProvider interface for generic OpenID Connect providers
type authProvider struct {
	cnf          oauth2.Config
	userInfoURL  string
	issuer       string
	providerName string
}

// NewAuthProvider creates a new generic OIDC provider instance
// Required parameters in the AuthProvider configuration:
// - @OIDC/CLIENT_ID: OAuth2 client ID
// - @OIDC/CLIENT_SECRET: OAuth2 client secret
// - @OIDC/REDIRECT_URL: OAuth2 redirect URL
// - @OIDC/AUTHORIZATION_URL: OIDC authorization endpoint
// - @OIDC/TOKEN_URL: OIDC token endpoint
// - @OIDC/USERINFO_URL: OIDC userinfo endpoint
// - @OIDC/SCOPES: Space-separated list of OAuth2 scopes (optional, defaults to "openid profile email")
func NewAuthProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	// Get configuration parameters
	clientID := p.GetParam("@OIDC/CLIENT_ID")
	clientSecret := p.GetParam("@OIDC/CLIENT_SECRET")
	redirectURL := p.GetParam("@OIDC/REDIRECT_URL")
	authURL := p.GetParam("@OIDC/AUTHORIZATION_URL")
	tokenURL := p.GetParam("@OIDC/TOKEN_URL")
	userInfoURL := p.GetParam("@OIDC/USERINFO_URL")

	// Default scopes if not specified
	scopes := []string{"openid", "profile", "email"}

	oauthConfig := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	return authProvider{
		cnf:          oauthConfig,
		userInfoURL:  userInfoURL,
		providerName: p.Name,
	}
}

// HasRefreshTokenFlow indicates that OIDC providers typically support refresh tokens
func (o authProvider) HasRefreshTokenFlow() bool {
	return true
}

// GetAuthCodeUrl returns the authorization URL where users should be redirected for authentication
func (o authProvider) GetAuthCodeUrl(state string) string {
	return o.cnf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// VerifyCode exchanges an authorization code for access and refresh tokens
func (o authProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	token, err := o.cnf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error verifying code with OIDC provider %s: %w", o.providerName, err)
	}

	return &sdk.AuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

// TokenResponse represents the OAuth2 token response structure
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// RefreshToken uses a refresh token to obtain a new access token
func (o authProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	data := url.Values{}
	data.Set("client_id", o.cnf.ClientID)
	data.Set("client_secret", o.cnf.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := http.PostForm(o.cnf.Endpoint.TokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("error refreshing token with OIDC provider %s: %w", o.providerName, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading token refresh response: %w", err)
	}

	// Check for HTTP error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error refreshing token, status: %d, response: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling token response: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return &sdk.AuthToken{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// OIDCIdentityEmail handles email identity information
type OIDCIdentityEmail struct {
	Email string `json:"email"`
}

func (o OIDCIdentityEmail) UpdateUserDetails(user *sdk.User) {
	user.Email = o.Email
}

// OIDCIdentityName handles name identity information
type OIDCIdentityName struct {
	Name string `json:"name"`
}

func (o OIDCIdentityName) UpdateUserDetails(user *sdk.User) {
	user.Name = o.Name
}

// OIDCIdentityProfilePic handles profile picture identity information
type OIDCIdentityProfilePic struct {
	ProfilePic string `json:"picture"`
}

func (o OIDCIdentityProfilePic) UpdateUserDetails(user *sdk.User) {
	user.ProfilePic = o.ProfilePic
}

// UserInfoResponse represents the standard OIDC UserInfo response
type UserInfoResponse struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Picture           string `json:"picture"`
	Locale            string `json:"locale"`
}

// GetIdentity retrieves user identity information using an access token
func (o authProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	req, err := http.NewRequest("GET", o.userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching identity from OIDC provider %s: %w", o.providerName, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading userinfo response: %w", err)
	}

	// Check for HTTP error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching identity, status: %d, response: %s", resp.StatusCode, string(respBytes))
	}

	var userInfo UserInfoResponse
	err = json.Unmarshal(respBytes, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling userinfo response: %s - %w", string(respBytes), err)
	}

	var identities []sdk.AuthIdentity

	// Add email identity if available
	if userInfo.Email != "" {
		identities = append(identities, sdk.AuthIdentity{
			Type:     sdk.AuthIdentityTypeEmail,
			Metadata: OIDCIdentityEmail{Email: userInfo.Email},
		})
	}

	// Add name identity (prefer full name, fallback to given name, then preferred username)
	name := userInfo.Name
	if name == "" && userInfo.GivenName != "" {
		if userInfo.FamilyName != "" {
			name = fmt.Sprintf("%s %s", userInfo.GivenName, userInfo.FamilyName)
		} else {
			name = userInfo.GivenName
		}
	}
	if name == "" {
		name = userInfo.PreferredUsername
	}
	if name != "" {
		identities = append(identities, sdk.AuthIdentity{
			Type:     sdk.AuthIdentityTypeEmail, // Using email type as the general identity type
			Metadata: OIDCIdentityName{Name: name},
		})
	}

	// Add profile picture identity if available
	if userInfo.Picture != "" {
		identities = append(identities, sdk.AuthIdentity{
			Type:     sdk.AuthIdentityTypeEmail, // Using email type as the general identity type
			Metadata: OIDCIdentityProfilePic{ProfilePic: userInfo.Picture},
		})
	}

	return identities, nil
}
