package sdk

import "time"

type User struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"project_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	IsEnabled bool       `json:"is_enabled"`
	Expiry    *time.Time `json:"expiry"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
}
