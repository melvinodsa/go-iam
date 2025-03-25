package models

type RoleMap struct {
	RoleId string   `bson:"role_id"`
	UserId []string `bson:"user_id"`
}

func (u RoleMapModel) Name() string {
	return "roleMap"
}

type RoleMapModel struct {
	iam
	RoleIdKey string
	UserIdKey []string
}

func GetRoleMap() RoleMapModel {
	return RoleMapModel{
		RoleIdKey: "role_id",
		UserIdKey: []string{"user_id"},
	}
}
