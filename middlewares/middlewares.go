package middlewares

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/project"
)

type Middlewares struct {
	projectSvc project.Service
	db         db.DB
}

func NewMiddlewares(projectSvc project.Service, db db.DB) *Middlewares {
	return &Middlewares{
		projectSvc: projectSvc,
		db:         db,
	}
}
