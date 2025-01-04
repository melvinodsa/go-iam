package middlewares

import (
	"github.com/melvinodsa/go-iam/services/project"
)

type Middlewares struct {
	projectSvc project.Service
}

func NewMiddlewares(projectSvc project.Service) *Middlewares {
	return &Middlewares{
		projectSvc: projectSvc,
	}
}
