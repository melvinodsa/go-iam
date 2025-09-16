package sdk

// AuthVerifyCodeResponse represents the response from OAuth2 authorization code verification.
// This contains the access token that can be used to authenticate API requests.
type AuthVerifyCodeResponse struct {
	AccessToken string `json:"access_token"` // JWT access token for API authentication
}

// AuthCallbackResponse represents the response from OAuth2 callback processing.
// This is returned after a user completes the OAuth2 authorization flow.
type AuthCallbackResponse struct {
	Success bool                    `json:"success"` // Indicates if the callback was processed successfully
	Message string                  `json:"message"` // Human-readable message about the operation
	Data    *AuthVerifyCodeResponse `json:"data"`    // Token data (present only on success)
}

// AuthLoginResponse represents the response from initiating an OAuth2 login flow.
// This contains the URL where the user should be redirected to complete authentication.
type AuthLoginResponse struct {
	Success bool                  `json:"success"` // Indicates if the login initiation was successful
	Message string                `json:"message"` // Human-readable message about the operation
	Data    AuthLoginDataResponse `json:"data"`    // Login flow data
}

// AuthLoginDataResponse contains the data needed to continue the OAuth2 login flow.
type AuthLoginDataResponse struct {
	LoginUrl string `json:"login_url"` // URL where the user should be redirected for authentication
}

// AuthRedirectResponse represents a response that requires client-side redirection.
// This is used for OAuth2 flows that need to redirect the user to external providers.
type AuthRedirectResponse struct {
	RedirectUrl string `json:"redirect_url"` // URL where the client should redirect the user
}

// ClientCredentialsRequest represents a request for OAuth2 client credentials flow.
// This is used for server-to-server authentication where no user interaction is required.
type ClientCredentialsRequest struct {
	ClientId     string `json:"client_id"`     // OAuth2 client identifier
	ClientSecret string `json:"client_secret"` // OAuth2 client secret
}

// ClientCredentialsResponse represents the response from OAuth2 client credentials flow.
type ClientCredentialsResponse struct {
	Success bool                    `json:"success"` // Indicates if the credentials were accepted
	Message string                  `json:"message"` // Human-readable message about the operation
	Data    *AuthVerifyCodeResponse `json:"data"`    // Token data (present only on success)
}

// ClientCredentialsDataResponse contains detailed token information from client credentials flow.
// This includes both access and refresh tokens with expiration information.
type ClientCredentialsDataResponse struct {
	AccessToken  string `json:"access_token"`  // JWT access token for API authentication
	RefreshToken string `json:"refresh_token"` // Token used to refresh the access token
	TokenType    string `json:"token_type"`    // Type of token (typically "Bearer")
	ExpiresIn    int64  `json:"expires_in"`    // Number of seconds until the access token expires
}

// RefreshTokenRequest represents a request to refresh an access token.
// This is used to obtain a new access token when the current one expires.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"` // The refresh token obtained from previous authentication
}

// RefreshTokenResponse represents the response from a token refresh operation.
type RefreshTokenResponse struct {
	Success bool                           `json:"success"` // Indicates if the token was refreshed successfully
	Message string                         `json:"message"` // Human-readable message about the operation
	Data    *ClientCredentialsDataResponse `json:"data"`    // New token data (present only on success)
}
