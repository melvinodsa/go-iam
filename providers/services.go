package providers

import (
	"github.com/melvinodsa/go-iam/api-server/db"
	"github.com/melvinodsa/go-iam/api-server/services/project"
)

type Service struct {
	Projects project.Service
}

func NewServices(db db.DB, cache *Cache) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	return &Service{Projects: psvc}
}
