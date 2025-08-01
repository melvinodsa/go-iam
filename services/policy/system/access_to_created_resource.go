package system

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
)

type accessToCreatedResource struct {
	id      string
	userSvc user.Service
	pc      PolicyCheck
}

func NewAccessToCreatedResource(userSvc user.Service) accessToCreatedResource {
	return accessToCreatedResource{id: "@policy/system/access_to_created_resource", userSvc: userSvc, pc: NewPolicyCheck(userSvc)}
}

func (a accessToCreatedResource) ID() string {
	return a.id
}

func (a accessToCreatedResource) HandleEvent(event utils.Event[sdk.Resource]) {
	log.Debugw("received resource event", "event", event.Name())
	userId := event.Metadata().User.Id
	_, exists, err := a.pc.RunCheck(event.Context(), a.id, userId)
	if err != nil {
		log.Errorw("error checking user while handling resource create event", "userId", userId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	if !exists {
		return
	}
	err = a.userSvc.AddResourceToUser(event.Context(), userId, sdk.AddUserResourceRequest{
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
