package project

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromSdkToModel(project sdk.Project) models.Project {
	return models.Project{
		Id:          project.Id,
		Name:        project.Name,
		Tags:        project.Tags,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		CreatedBy:   project.CreatedBy,
		UpdatedAt:   project.UpdatedAt,
		UpdatedBy:   project.UpdatedBy,
	}
}

func fromModelToSdk(project *models.Project) *sdk.Project {
	return &sdk.Project{
		Id:          project.Id,
		Name:        project.Name,
		Tags:        project.Tags,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		CreatedBy:   project.CreatedBy,
		UpdatedAt:   project.UpdatedAt,
		UpdatedBy:   project.UpdatedBy,
	}
}

func fromModelListToSdk(projects []models.Project) []sdk.Project {
	var sdkProjects []sdk.Project
	for i := range projects {
		sdkProjects = append(sdkProjects, *fromModelToSdk(&projects[i]))
	}
	return sdkProjects
}
