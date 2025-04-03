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
