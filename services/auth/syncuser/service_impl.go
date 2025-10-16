package syncuser

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/auth"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type service struct {
	authSvc auth.Service
}

func NewService(authSvc auth.Service) Service {
	return &service{
		authSvc: authSvc,
	}
}

func (s *service) HandleEvent(event utils.Event[sdk.User]) {
	if event == nil {
		return
	}
	if event.Name() == goiamuniverse.EventUserUpdated {
		user := event.Payload()
		if user.Id == "" {
			return
		}
		err := s.authSvc.SynchronizeIdentity(event.Context(), user.Id)
		if err != nil {
			log.Errorw("failed to synchronize identity", "error", err, "userId", user.Id)
		}
	}
}
