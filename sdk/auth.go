package sdk

type AuthRedirectResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type AuthVerifyCodeResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthCallbackResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    *AuthVerifyCodeResponse `json:"data"`
}
