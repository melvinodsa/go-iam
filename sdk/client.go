package sdk

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Client struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Secret       string     `json:"secret"`
	Tags         []string   `json:"tags"`
	RedirectURLs []string   `json:"redirect_urls"`
	Scopes       []string   `json:"scopes"`
	Enabled      bool       `json:"enabled"`
	CreatedAt    *time.Time `json:"created_at"`
	CreatedBy    string     `json:"created_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
	UpdatedBy    string     `json:"updated_by"`
}

type ClientResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *Client `json:"data,omitempty"`
}

func NewErrorClientResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(ClientResponse{
		Success: false,
		Message: msg,
	})
}

func ClientBadRequest(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusBadRequest, c)
}

func ClientInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusInternalServerError, c)
}

type ClientsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Client `json:"data,omitempty"`
}
