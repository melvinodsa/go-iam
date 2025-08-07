package sdk

import (
	"errors"

	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

var ErrPolicyNotFound = errors.New("policy not found")

type Policy struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Definition  PolicyDefinition `json:"definition"`
}

type PolicyDefinition struct {
	Arguments []PolicyArgument `json:"arguments,omitempty"`
}

type PolicyArgument struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	DataType    goiamuniverse.DataType `json:"data_type,omitempty"`
}

type PolicyResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *Policy `json:"data,omitempty"`
}

type PoliciesResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    PolicyList `json:"data,omitempty"`
}

type PolicyList struct {
	Policies []Policy `json:"policies"`
	Total    int      `json:"total"`
	Skip     int64    `json:"skip"`
	Limit    int64    `json:"limit"`
}

type PolicyQuery struct {
	Query string `json:"query,omitempty"`
	Skip  int64  `json:"skip,omitempty"`
	Limit int64  `json:"limit,omitempty"`
}
