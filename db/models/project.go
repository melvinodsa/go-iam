package models

import "time"

type Project struct {
	Id          string     `bson:"id"`
	Name        string     `bson:"name"`
	Tags        []string   `bson:"tags"`
	Description string     `bson:"description"`
	CreatedAt   *time.Time `bson:"created_at"`
	CreatedBy   string     `bson:"created_by"`
	UpdatedAt   *time.Time `bson:"updated_at"`
	UpdatedBy   string     `bson:"updated_by"`
}

type ProjectModel struct {
	iam
	IdKey          string
	NameKey        string
	TagsKey        string
	DescriptionKey string
}

func (p ProjectModel) Name() string {
	return "projects"
}

func GetProjectModel() ProjectModel {
	return ProjectModel{
		IdKey:          "id",
		NameKey:        "name",
		TagsKey:        "tags",
		DescriptionKey: "description",
	}
}
