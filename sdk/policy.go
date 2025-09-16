package sdk

import (
	"errors"

	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

// ErrPolicyNotFound is returned when a requested policy cannot be found.
var ErrPolicyNotFound = errors.New("policy not found")

// Policy represents a policy in the Go IAM system.
// Policies define fine-grained access control rules that can be applied
// to users and resources. They support parameterization through arguments
// for dynamic policy evaluation.
type Policy struct {
	Id          string           `json:"id"`          // Unique identifier for the policy
	Name        string           `json:"name"`        // Display name of the policy
	Description string           `json:"description"` // Description of what this policy does
	Definition  PolicyDefinition `json:"definition"`  // Policy definition containing logic and arguments
}

// PolicyDefinition contains the structure and parameters of a policy.
// This defines what arguments the policy accepts for dynamic evaluation.
type PolicyDefinition struct {
	Arguments []PolicyArgument `json:"arguments,omitempty"` // Array of argument definitions for the policy
}

// PolicyArgument represents a parameter that can be passed to a policy.
// These arguments allow policies to be customized with different values
// when applied to users or resources.
type PolicyArgument struct {
	Name        string                 `json:"name,omitempty"`        // Name of the argument
	Description string                 `json:"description,omitempty"` // Description of what this argument represents
	DataType    goiamuniverse.DataType `json:"data_type,omitempty"`   // Data type of the argument value
}

// PolicyResponse represents an API response containing a single policy.
type PolicyResponse struct {
	Success bool    `json:"success"`        // Indicates if the operation was successful
	Message string  `json:"message"`        // Human-readable message about the operation
	Data    *Policy `json:"data,omitempty"` // The policy data (present only on success)
}

// PoliciesResponse represents an API response containing a list of policies.
type PoliciesResponse struct {
	Success bool       `json:"success"`        // Indicates if the operation was successful
	Message string     `json:"message"`        // Human-readable message about the operation
	Data    PolicyList `json:"data,omitempty"` // The paginated policy list data
}

// PolicyList represents a paginated list of policies with metadata.
type PolicyList struct {
	Policies []Policy `json:"policies"` // Array of policy objects
	Total    int      `json:"total"`    // Total number of policies matching the query (before pagination)
	Skip     int64    `json:"skip"`     // Number of records skipped
	Limit    int64    `json:"limit"`    // Maximum number of records returned
}

// PolicyQuery represents search and filtering criteria for policy queries.
// This is used for listing policies with text search and pagination.
type PolicyQuery struct {
	Query string `json:"query,omitempty"` // Text search query across policy fields
	Skip  int64  `json:"skip,omitempty"`  // Number of records to skip (pagination)
	Limit int64  `json:"limit,omitempty"` // Maximum number of records to return
}
