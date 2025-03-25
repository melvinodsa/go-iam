package models

type ResourceMap struct {
	ResourcecId string   `bson:"resource_id"`
	RoleId      []string `bson:"role_id"`
}

func (u ResourceMapModel) Name() string {
	return "resouceMap"
}

type ResourceMapModel struct {
	iam
	ResourceIdKey string
	RoleIdKey     []string
}

func GetRoleMap() ResourceMapModel {
	return ResourceMapModel{
		ResourceIdKey: "resource_id",
		RoleIdKey:     []string{"role_id"},
	}
}
