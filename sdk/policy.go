package sdk

import "github.com/melvinodsa/go-iam/utils/goiamuniverse"

type Policy struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Definition  PolicyDefinition `json:"definition"`
}

type PolicyDefinition struct {
	Arguments []PolicyArgument
}

type PolicyArgument struct {
	Name        string
	Description string
	DataType    goiamuniverse.DataType
}

type PolicyResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *Policy `json:"data,omitempty"`
}

type PoliciesResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Policy `json:"data,omitempty"`
}
