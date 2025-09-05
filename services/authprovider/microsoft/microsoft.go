package microsoft

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

type authProvider struct {
	cnf oauth2.Config
}

func NewAuthProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	oauthConfig := oauth2.Config{
		ClientID:     p.GetParam("@MICROSOFT/CLIENT_ID"),
		ClientSecret: p.GetParam("@MICROSOFT/CLIENT_SECRET"),
		RedirectURL:  p.GetParam("@MICROSOFT/REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	}
	return authProvider{cnf: oauthConfig}
}

func (m authProvider) HasRefreshTokenFlow() bool {
	return true
}

func (m authProvider) GetAuthCodeUrl(state string) string {
	return m.cnf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (m authProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	token, err := m.cnf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error verifying the code with microsoft exchange. %w", err)
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

func (m authProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	urlStr := m.cnf.Endpoint.TokenURL
	data := url.Values{}
	data.Set("client_id", m.cnf.ClientID)
	data.Set("client_secret", m.cnf.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("scope", "openid profile email")

	resp, err := http.PostForm(urlStr, data)
	if err != nil {
		return nil, fmt.Errorf("error refreshing the token. %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body: %w", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response. %w", err)
	}

	// Check for HTTP error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error refreshing the token, status: %d, response: %s", resp.StatusCode, string(body))
	}

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

type MicrosoftIdentityEmail struct {
	Email string `json:"email"`
}

func (m MicrosoftIdentityEmail) UpdateUserDetails(user *sdk.User) {
	user.Email = m.Email
}

type MicrosoftIdentityName struct {
	Name string `json:"name"`
}

func (m MicrosoftIdentityName) UpdateUserDetails(user *sdk.User) {
	user.Name = m.Name
}

type MicrosoftIdentityProfilePic struct {
	ProfilePic string `json:"picture"`
}

func (m MicrosoftIdentityProfilePic) UpdateUserDetails(user *sdk.User) {
	user.ProfilePic = m.ProfilePic
}

func (m authProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/oidc/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request. %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

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

	var userInfo struct {
		Email      string `json:"email"`
		FirstName  string `json:"givenname"`
		LastName   string `json:"familyname"`
		ProfilePic string `json:"picture"`
	}
	err = json.Unmarshal(respBytes, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
	}
	return []sdk.AuthIdentity{
		{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityEmail{Email: userInfo.Email}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityName{Name: fmt.Sprintf("%s %s", userInfo.FirstName, userInfo.LastName)}},
		{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityProfilePic{ProfilePic: userInfo.ProfilePic}},
	}, nil
}
