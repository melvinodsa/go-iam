package sdk

type AuthVerifyCodeResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthCallbackResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    *AuthVerifyCodeResponse `json:"data"`
}

type AuthLoginResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    AuthLoginDataResponse `json:"data"`
}

type AuthLoginDataResponse struct {
	LoginUrl string `json:"login_url"`
}

type AuthRedirectResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type ClientCredentialsRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ClientCredentialsResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    *ClientCredentialsDataResponse `json:"data"`
}

type ClientCredentialsDataResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}