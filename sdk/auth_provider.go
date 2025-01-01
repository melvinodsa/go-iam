package sdk

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthProviderType string

const (
	AuthProviderTypeGoogle AuthProviderType = "GOOGLE"
)

type AuthProvider struct {
	Id        string              `json:"id"`
	Name      string              `json:"name"`
	Icon      string              `json:"icon"`
	Provider  AuthProviderType    `json:"provider"`
	Params    []AuthProviderParam `json:"params"`
	Enabled   bool                `json:"enabled"`
	CreatedAt *time.Time          `json:"created_at"`
	UpdatedAt *time.Time          `json:"updated_at"`
	CreatedBy string              `json:"created_by"`
	UpdatedBy string              `json:"updated_by"`
}

func (a AuthProvider) GetParam(key string) string {
	for _, p := range a.Params {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

type AuthProviderParam struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Key      string `json:"key"`
	IsSecret bool   `json:"is_secret"`
}

type AuthProviderResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *AuthProvider `json:"data"`
}

func NewErrorAuthProviderResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(ClientResponse{
		Success: false,
		Message: msg,
	})
}

func AuthProviderBadRequest(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProviderResponse(msg, http.StatusBadRequest, c)
}

func AuthProviderInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProviderResponse(msg, http.StatusInternalServerError, c)
}

type AuthProvidersResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []AuthProvider `json:"data"`
}

func NewErrorAuthProvidersResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(AuthProvidersResponse{
		Success: false,
		Message: msg,
	})
}

func AuthProvidersInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProvidersResponse(msg, http.StatusInternalServerError, c)
}

type ServiceProvider interface {
	GetAuthCodeUrl(state string) string
	VerifyCode(ctx context.Context, code string) (*AuthToken, error)
	RefreshToken(refreshToken string) (*AuthToken, error)
	GetIdentity(token string) ([]AuthIdentity, error)
}
