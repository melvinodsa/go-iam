package sdk

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrClientNotFound is returned when a requested OAuth2 client cannot be found.
var ErrClientNotFound = errors.New("client not found")

// Client represents an OAuth2 client in the Go IAM system.
// Clients can be external applications that integrate with the IAM system
// or internal service accounts used for server-to-server communication.
type Client struct {
	Id                    string     `json:"id"`                       // Unique identifier for the client
	Name                  string     `json:"name"`                     // Display name of the client
	Description           string     `json:"description"`              // Description of the client's purpose
	Secret                string     `json:"secret"`                   // Client secret for authentication
	Tags                  []string   `json:"tags"`                     // Tags for categorizing the client
	RedirectURLs          []string   `json:"redirect_urls"`            // Allowed redirect URLs for OAuth2 flows
	Scopes                []string   `json:"scopes"`                   // OAuth2 scopes this client can request
	ProjectId             string     `json:"project_id"`               // ID of the project this client belongs to
	DefaultAuthProviderId string     `json:"default_auth_provider_id"` // Default auth provider for this client
	GoIamClient           bool       `json:"go_iam_client"`            // Indicates if this is a Go-IAM internal client
	LinkedUserId          string     `json:"linked_user_id"`           // Associated user ID for service accounts
	ServiceAccountEmail   string     `json:"service_account_email"`    // Email address for service account clients
	Enabled               bool       `json:"enabled"`                  // Whether the client is active
	CreatedAt             *time.Time `json:"created_at"`               // Timestamp when client was created
	CreatedBy             string     `json:"created_by"`               // ID of the user who created this client
	UpdatedAt             *time.Time `json:"updated_at"`               // Timestamp when client was last updated
	UpdatedBy             string     `json:"updated_by"`               // ID of the user who last updated this client
}

// IsServiceAccount returns true if this client represents a service account.
// A service account is a Go-IAM client that has an associated user account.
func (c Client) IsServiceAccount() bool {
	return c.HasGoIamAuthProvider() && c.LinkedUserId != ""
}

// HasGoIamAuthProvider returns true if this client uses Go-IAM's internal authentication.
// Clients without a DefaultAuthProviderId use Go-IAM's built-in authentication system.
func (c Client) HasGoIamAuthProvider() bool {
	return c.DefaultAuthProviderId == ""
}

// ClientResponse represents an API response containing a single OAuth2 client.
type ClientResponse struct {
	Success bool    `json:"success"`        // Indicates if the operation was successful
	Message string  `json:"message"`        // Human-readable message about the operation
	Data    *Client `json:"data,omitempty"` // The client data (present only on success)
}

// NewErrorClientResponse creates a new error response for client operations.
func NewErrorClientResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(ClientResponse{
		Success: false,
		Message: msg,
	})
}

// ClientBadRequest returns a 400 Bad Request error response for client operations.
func ClientBadRequest(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusBadRequest, c)
}

// ClientNotFound returns a 404 Not Found error response for client operations.
func ClientNotFound(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusNotFound, c)
}

// ClientInternalServerError returns a 500 Internal Server Error response for client operations.
func ClientInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorClientResponse(msg, http.StatusInternalServerError, c)
}

// ClientsResponse represents an API response containing a list of OAuth2 clients.
type ClientsResponse struct {
	Success bool     `json:"success"` // Indicates if the operation was successful
	Message string   `json:"message"` // Human-readable message about the operation
	Data    []Client `json:"data"`    // Array of client data
}

// NewErrorClientsResponse creates a new error response for clients list operations.
func NewErrorClientsResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(ClientsResponse{
		Success: false,
		Message: msg,
	})
}

// ClientsInternalServerError returns a 500 Internal Server Error response for clients list operations.
func ClientsInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorClientsResponse(msg, http.StatusInternalServerError, c)
}

// ClientQueryParams represents query parameters for filtering and sorting OAuth2 clients.
type ClientQueryParams struct {
	ProjectIds      []string `json:"project_id"`         // Filter clients by project IDs
	GoIamClient     bool     `json:"go_iam_client"`      // Filter by Go-IAM internal clients
	SortByUpdatedAt bool     `json:"sort_by_updated_at"` // Sort results by update timestamp
}
