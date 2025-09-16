package sdk

// Metadata represents contextual information about a request or operation.
// This structure typically contains the authenticated user information
// and the scope of projects they have access to, used for authorization
// and multi-tenant data filtering.
type Metadata struct {
	User       *User    // The authenticated user making the request
	ProjectIds []string // Array of project IDs the user has access to
}
