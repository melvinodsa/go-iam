package google

import (
	"bytes"
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
	"golang.org/x/oauth2/google"
	gauth "google.golang.org/api/oauth2/v2"
)

type authProvider struct {
	cnf oauth2.Config
}

func NewAuthProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	oauthConfig := oauth2.Config{
		ClientID:     p.GetParam("@GOOGLE/CLIENT_ID"),
		ClientSecret: p.GetParam("@GOOGLE/CLIENT_SECRET"), // Set this in your environment
		RedirectURL:  p.GetParam("@GOOGLE/REDIRECT_URL"),
		Scopes:       []string{gauth.UserinfoEmailScope, gauth.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
	return authProvider{cnf: oauthConfig}
}

func (g authProvider) HasRefreshTokenFlow() bool {
	return true
}

func (g authProvider) GetAuthCodeUrl(state string) string {
	return g.cnf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}
func (g authProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	token, err := g.cnf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error verifying the code with google exchange. %w", err)
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
}

func (g authProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	// Token endpoint
	url := "https://oauth2.googleapis.com/token"

	// Prepare the request body
	data := map[string]string{
		"client_id":     g.cnf.ClientID,
		"client_secret": g.cnf.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling the data. %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error refreshing the token. %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response. %w", err)
	}

	// Parse the JSON response
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the response. %w", err)
	}
	return &sdk.AuthToken{
		AccessToken:  tokenResponse.AccessToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		RefreshToken: refreshToken,
	}, nil
}

type GoogleIdentityEmail struct {
	Email string `json:"email"`
}

func (g GoogleIdentityEmail) UpdateUserDetails(user *sdk.User) {
	user.Email = g.Email
}

type GoogleIdentityName struct {
	FirstName string `json:"given_name"`
}

func (g GoogleIdentityName) UpdateUserDetails(user *sdk.User) {
	user.Name = g.FirstName
}

type GoogleIdentityProfilePic struct {
	ProfilePic string `json:"profile_pic"`
}

func (g GoogleIdentityProfilePic) UpdateUserDetails(user *sdk.User) {
	user.ProfilePic = g.ProfilePic
}

func (g authProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token))
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

	var tokenResponse gauth.Userinfo
	err = json.Unmarshal(respBytes, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
	}
	return []sdk.AuthIdentity{
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GoogleIdentityEmail{Email: tokenResponse.Email}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GoogleIdentityName{FirstName: tokenResponse.GivenName}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: GoogleIdentityProfilePic{ProfilePic: tokenResponse.Picture}},
	}, nil
}
