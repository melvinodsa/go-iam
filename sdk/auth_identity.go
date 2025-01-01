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
	Metadata AuthMetadataType `json:"metadata"`
}

func (a AuthIdentity) UpdateUserDetails(user *User) {
	a.Metadata.UpdateUserDetails(user)
}
