package providers

import (
	"github.com/melvinodsa/go-iam/api-server/db"
	"github.com/melvinodsa/go-iam/api-server/services/projects"
)

type Service struct {
	Projects projects.Service
}

func NewServices(db db.DB, cache *Cache) *Service {
	pstr := projects.NewStore(db)
	psvc := projects.NewService(pstr)
	return &Service{Projects: psvc}
}
