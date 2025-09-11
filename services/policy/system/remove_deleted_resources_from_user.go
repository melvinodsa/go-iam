package system

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
)

type removeDeletedResourceFromUser struct {
	id      string
	userSvc user.Service
}

func NewRemoveDeletedResourceFromUser(userSvc user.Service) removeDeletedResourceFromUser {
	return removeDeletedResourceFromUser{id: "@policy/system/remove_deleted_resources_from_user", userSvc: userSvc}
}

func (a removeDeletedResourceFromUser) ID() string {
	return a.id
}

func (a removeDeletedResourceFromUser) Name() string {
	return "Remove deleted resources from user specified in user policy"
}

func (a removeDeletedResourceFromUser) HandleEvent(event utils.Event[sdk.Resource]) {
	log.Debugw("received resource event", "event", event.Name())
	err := a.userSvc.RemoveResourceFromAll(event.Context(), event.Payload().Key)
	if err != nil {
		log.Errorw("error removing resource from all users while handling resource delete event", "resource_id", event.Payload().ID, "error", err)
		return
	}
	log.Infow("successfully removed deleted resource from all users", "resource_id", event.Payload().ID)
}

func (a removeDeletedResourceFromUser) PolicyDef() sdk.Policy {
	return sdk.Policy{
		Id:          a.id,
		Name:        a.Name(),
		Description: "This policy removes the deleted resource from all users.",
		Definition: sdk.PolicyDefinition{
			Arguments: []sdk.PolicyArgument{},
		},
	}
}
