package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

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
			{
				Name:        "code_challenge",
				In:          "query",
				Description: "Code challenge for PKCE. This required for public clients. For security reasons, only S256 is supported.",
				Required:    false,
			},
			{
				Name:        "code_verifier",
				In:          "query",
				Description: "Code verifier for PKCE. This required for public clients",
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

	codeChallenge := c.Query("code_challenge", "")
	// Might have to revisit this when the standards change.
	if len(codeChallenge) != 0 && strings.Compare(codeChallenge, "S256") != 0 {
		log.Debugw("invalid code challenge", "code_challenge", codeChallenge)
		return sdk.AuthProviderBadRequest("invalid code challenge. Only S256 is supported", c)
	}
	url, err := pr.S.Auth.GetLoginUrl(c.Context(), c.Query("client_id", ""), c.Query("auth_provider", ""), c.Query("state", ""), c.Query("redirect_url", ""), c.Query("code_challenge", ""), c.Query("code_verifier", ""))
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
			{
				Name:        "code_verifier",
				In:          "query",
				Description: "The code verifier",
				Required:    false,
			},
			{
				Name:        "client_id",
				In:          "query",
				Description: "The client ID to be provided if code verifier is provided",
				Required:    false,
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
	var clientId, clientSecret string
	// get code verifier from query params
	codeVerifier := c.Query("code_verifier")
	clientId = c.Query("client_id")

	if len(codeVerifier) == 0 || len(clientId) == 0 {
		// get client id and secret from authorization header with basic auth
		clId, clSec, ok := getClientDetails(c)
		if !ok {
			return sdk.AuthProviderBadRequest("missing or invalid authorization header", c)
		}
		clientId = clId
		clientSecret = clSec
	}
	resp, err := pr.S.Auth.ClientCallback(c.Context(), code, codeVerifier, clientId, clientSecret)
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

func getClientDetails(c *fiber.Ctx) (string, string, bool) {
	headers := c.GetReqHeaders()
	authHeaders := headers["Authorization"]
	if len(authHeaders) == 0 {
		return "", "", false
	}
	// extract client id and secret from basic auth
	parts := strings.SplitN(authHeaders[0], " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}
	credentials, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}
	creds := strings.SplitN(string(credentials), ":", 2)
	if len(creds) != 2 {
		return "", "", false
	}
	return creds[0], creds[1], true
}
