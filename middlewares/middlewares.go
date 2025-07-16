package middlewares

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils"
)

type Middlewares struct {
	projectSvc  project.Service
	db          db.DB
	AuthEnabled bool
}

func NewMiddlewares(projectSvc project.Service, db db.DB, authEnabled bool) *Middlewares {
	return &Middlewares{
		projectSvc:  projectSvc,
		db:          db,
		AuthEnabled: authEnabled,
	}
}

func (m *Middlewares) Handle(e utils.Event[sdk.Client]) {
	if m.AuthEnabled {
		return
	}
	if e.Name() != sdk.EventClientCreated && e.Name() != sdk.EventClientUpdated {
		return
	}
	if !e.Payload().GoIamClient {
		return
	}
	m.AuthEnabled = true
}
