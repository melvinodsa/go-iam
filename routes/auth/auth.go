package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
)

func Login(c *fiber.Ctx) error {
	log.Debug("received login request")
	pr := providers.GetProviders(c)
	url, err := pr.S.Auth.GetLoginUrl(c.Context(), c.Query("client_id", ""), c.Query("auth_provider", ""), c.Query("state", ""), c.Query("redirect_url", ""))
	if err != nil {
		message := fmt.Errorf("failed to get login url. %w", err).Error()
		log.Errorw("failed to create authprovider", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	return c.Status(http.StatusOK).JSON(sdk.AuthLoginResponse{
		Success: true,
		Message: "Login URL generated successfully",
		Data: sdk.AuthLoginDataResponse{
			LoginUrl: url,
		},
	})
}

func Redirect(c *fiber.Ctx) error {
	log.Debug("received redirect request")
	pr := providers.GetProviders(c)
	code := c.Query("code")
	state := c.Query("state")
	resp, err := pr.S.Auth.Redirect(c.Context(), code, state)
	if err != nil {
		message := fmt.Errorf("failed to redirect. %w", err).Error()
		log.Errorw("failed to redirect", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	log.Debug("redirected successfully")

	return c.Status(http.StatusOK).JSON(sdk.AuthRedirectResponse{
		RedirectUrl: resp.RedirectUrl,
	})
}

func Verify(c *fiber.Ctx) error {
	log.Debug("received callback request")
	pr := providers.GetProviders(c)
	code := c.Query("code")
	resp, err := pr.S.Auth.ClientCallback(c.Context(), code)
	if err != nil {
		message := fmt.Errorf("failed to get callback. %w", err).Error()
		log.Errorw("failed to get callback", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	log.Debug("code verification was successful")

	return c.Status(http.StatusOK).JSON(sdk.AuthCallbackResponse{
		Success: true,
		Message: "Callback successful",
		Data:    resp,
	})
}
