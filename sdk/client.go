package sdk

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ErrClientNotFound = errors.New("client not found")

type Client struct {
	Id                    string     `json:"id"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	Secret                string     `json:"secret"`
	Tags                  []string   `json:"tags"`
	RedirectURLs          []string   `json:"redirect_urls"`
	Scopes                []string   `json:"scopes"`
	ProjectId             string     `json:"project_id"`
	DefaultAuthProviderId string     `json:"default_auth_provider_id"`
	GoIamClient           bool       `json:"go_iam_client"` // Indicates if this is a Go-IAM client
	LinkedUserId          string     `json:"linked_user_id"`
	ServiceAccountEmail   string     `json:"service_account_email"`
	Enabled               bool       `json:"enabled"`
	CreatedAt             *time.Time `json:"created_at"`
	CreatedBy             string     `json:"created_by"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedBy             string     `json:"updated_by"`
}

func (c Client) IsServiceAccount() bool {
	return c.HasGoIamAuthProvider() && c.LinkedUserId != ""
}

func (c Client) HasGoIamAuthProvider() bool {
	return c.DefaultAuthProviderId == ""
}

type ClientResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *Client `json:"data,omitempty"`
}

func NewErrorClientResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(ClientResponse{
		Success: false,
		Message: msg,
	})
}

func ClientBadRequest(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusBadRequest, c)
}

func ClientNotFound(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusNotFound, c)
}

func ClientInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusInternalServerError, c)
}

type ClientsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Client `json:"data"`
}

func NewErrorClientsResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(ClientsResponse{
		Success: false,
		Message: msg,
	})
}

func ClientsInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorClientsResponse(msg, http.StatusInternalServerError, c)
}

type ClientQueryParams struct {
	ProjectIds      []string `json:"project_id"`
	GoIamClient     bool     `json:"go_iam_client"`
	SortByUpdatedAt bool     `json:"sort_by_updated_at"`
}
