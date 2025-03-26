package models

type ResourceMap struct {
	ResourceId string   `bson:"resource_id"`
	RoleId     []string `bson:"role_id"`
}

func (u ResourceMapModel) Name() string {
	return "resourceMap"
}

type ResourceMapModel struct {
	iam
	ResourceIdKey string
	RoleIdKey     string
}

func GetResourceMap() ResourceMapModel {
	return ResourceMapModel{
		ResourceIdKey: "resource_id",
		RoleIdKey:     "role_id",
	}
}
