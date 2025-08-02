package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// LoginRoute registers the login route
func LoginRoute(router fiber.Router, basePath string) {
	routePath := "/login"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Login",
		Description: "Login to the application",
		Tags:        routeTags,
		RequestBody: nil,
		Response: &docs.ApiResponse{
			Description: "Login URL generated successfully",
			Content:     new(sdk.AuthLoginResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "client_id",
				In:          "query",
				Description: "The client ID",
				Required:    true,
			},
			{
				Name:        "auth_provider",
				In:          "query",
				Description: "The authentication provider",
				Required:    false,
			},
			{
				Name:        "state",
				In:          "query",
				Description: "State parameter for CSRF protection",
				Required:    false,
			},
			{
				Name:        "redirect_url",
				In:          "query",
				Description: "The URL to redirect to after login",
				Required:    false,
			},
		},
		UnAuthenticated:      true,
		ProjectIDNotRequired: true,
	})
	router.Get(routePath, Login)
}

func Login(c *fiber.Ctx) error {
	log.Debug("received login request")
	pr := providers.GetProviders(c)

	url, err := pr.S.Auth.GetLoginUrl(c.Context(), c.Query("client_id", ""), c.Query("auth_provider", ""), c.Query("state", ""), c.Query("redirect_url", ""))
	if err != nil {
		message := fmt.Errorf("failed to get login url. %w", err).Error()
		log.Errorw("failed to get login url", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}

	postBack := c.Query("postback", "false")
	if postBack == "true" {
		return c.Status(http.StatusOK).JSON(sdk.AuthLoginResponse{
			Success: true,
			Message: "Login URL generated successfully",
			Data: sdk.AuthLoginDataResponse{
				LoginUrl: url,
			},
		})
	}
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

// RedirectRoute registers the redirect route
func RedirectRoute(router fiber.Router, basePath string) {
	routePath := "/authp-callback"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Redirect",
		Description: "Redirect to the authentication provider",
		Tags:        routeTags,
		RequestBody: nil,
		Response: &docs.ApiResponse{
			Description: "Redirect URL generated successfully",
			Content:     new(sdk.AuthRedirectResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "code",
				In:          "query",
				Description: "The authentication code",
				Required:    true,
			},
			{
				Name:        "state",
				In:          "query",
				Description: "State parameter for CSRF protection",
				Required:    false,
			},
			{
				Name:        "postback",
				In:          "query",
				Description: "Whether to return the redirect URL in the response",
				Required:    false,
			},
		},
		UnAuthenticated:      true,
		ProjectIDNotRequired: true,
	})
	router.Get(routePath, Redirect)
}

func Redirect(c *fiber.Ctx) error {
	log.Debug("received redirect request")
	pr := providers.GetProviders(c)
	code := c.Query("code")
	state := c.Query("state")
	postback := c.Query("postback", "false")
	resp, err := pr.S.Auth.Redirect(c.Context(), code, state)
	if err != nil {
		message := fmt.Errorf("failed to redirect. %w", err).Error()
		log.Errorw("failed to redirect", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	log.Debug("redirected successfully")
	if postback == "true" {
		return c.Status(http.StatusOK).JSON(sdk.AuthRedirectResponse{
			RedirectUrl: resp.RedirectUrl,
		})
	}
	return c.Redirect(resp.RedirectUrl, http.StatusTemporaryRedirect)
}

// VerifyRoute registers the verify route
func VerifyRoute(router fiber.Router, basePath string) {
	routePath := "/verify"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Verify",
		Description: "Verify the authentication code",
		Tags:        routeTags,
		RequestBody: nil,
		Response: &docs.ApiResponse{
			Description: "Verification successful",
			Content:     new(sdk.AuthCallbackResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "code",
				In:          "query",
				Description: "The authentication code",
				Required:    true,
			},
		},
		UnAuthenticated:      true,
		ProjectIDNotRequired: true,
	})
	router.Get(routePath, Verify)
}

func Verify(c *fiber.Ctx) error {
	log.Debug("received callback request")
	pr := providers.GetProviders(c)
	code := c.Query("code")
	resp, err := pr.S.Auth.ClientCallback(c.Context(), code)
	if err != nil {
		message := fmt.Errorf("failed to get callback. %w", err).Error()
		return sdk.AuthProviderInternalServerError(message, c)
	}
	log.Debug("code verification was successful")

	return c.Status(http.StatusOK).JSON(sdk.AuthCallbackResponse{
		Success: true,
		Message: "Callback successful",
		Data:    resp,
	})
}
