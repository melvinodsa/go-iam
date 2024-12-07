package providers

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/project"
)

type Service struct {
	Projects project.Service
	Clients  client.Service
}

func NewServices(db db.DB, cache *Cache) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)
	csvc := client.NewService(cstr)
	return &Service{Projects: psvc, Clients: csvc}
}
