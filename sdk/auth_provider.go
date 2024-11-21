package sdk

type AuthProvider struct {
	Id        string `json:"id"`
	Provider  string `json:"provider"`
	IsEnabled bool   `json:"is_enabled"`
}
