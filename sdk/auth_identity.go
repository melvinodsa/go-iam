package sdk

type AuthIdentityType string

const (
	AuthIdentityTypeEmail AuthIdentityType = "email"
	AuthIdentityTypePhone AuthIdentityType = "phone"
)

type AuthMetadataType interface {
	UpdateUserDetails(user *User)
}

type AuthIdentity struct {
	Type     AuthIdentityType `json:"type"`
	Value    string           `json:"value"`
	Metadata AuthMetadataType `json:"metadata"`
}
