package system

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
)

type accessToCreatedResource struct {
	id      string
	userSvc user.Service
}

func NewAccessToCreatedResource(userSvc user.Service) accessToCreatedResource {
	return accessToCreatedResource{id: "@policy/system/access_to_created_resource", userSvc: userSvc}
}

func (a accessToCreatedResource) ID() string {
	return a.id
}

func (a accessToCreatedResource) HandleEvent(event utils.Event[sdk.Resource]) {
	log.Debugw("received resource event", "event", event.Name())
	userId := event.Metadata().User.Id
	user, err := a.userSvc.GetById(context.Background(), userId)
	if err != nil {
		log.Errorw("error fetching user while handling resource create event", "userId", userId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	if _, exists := user.Policies[a.id]; !exists {
		return
	}
	err = a.userSvc.AddResourceToUser(context.Background(), userId, sdk.AddUserResourceRequest{
		PolicyId: a.id,
		Key:      event.Payload().Key,
		Name:     event.Payload().Name,
	})
	if err != nil {
		log.Errorw("error adding resource to user while handling resource create event", "userId", userId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	log.Infow("successfully added created resource to user", "userId", userId, "resource_id", event.Payload().ID)
}
