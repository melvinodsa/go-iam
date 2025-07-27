package sdk

import "time"

type PolicyBeta struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Roles       map[string]string `json:"roles"`
	Description string            `json:"description"`
	CreatedAt   *time.Time        `json:"created_at"`
	CreatedBy   string            `json:"created_by"`
}

type PolicyResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *PolicyBeta `json:"data,omitempty"`
}

type PoliciesResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []PolicyBeta `json:"data,omitempty"`
}
