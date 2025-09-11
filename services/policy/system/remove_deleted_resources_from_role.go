package system

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/utils"
)

type removeDeletedResourceFromRole struct {
	id      string
	roleSvc role.Service
}

func NewRemoveDeletedResourceFromRole(roleSvc role.Service) removeDeletedResourceFromRole {
	return removeDeletedResourceFromRole{id: "@policy/system/remove_deleted_resources_from_role", roleSvc: roleSvc}
}

func (a removeDeletedResourceFromRole) ID() string {
	return a.id
}

func (a removeDeletedResourceFromRole) Name() string {
	return "Remove deleted resources from role specified in user policy"
}

func (a removeDeletedResourceFromRole) HandleEvent(event utils.Event[sdk.Resource]) {
	log.Debugw("received resource event", "event", event.Name())
	err := a.roleSvc.RemoveResourceFromAll(event.Context(), event.Payload().Key)
	if err != nil {
		log.Errorw("error removing resource from all roles while handling resource delete event", "resource_id", event.Payload().ID, "error", err)
		return
	}
	log.Infow("successfully removed deleted resource from all roles", "resource_id", event.Payload().ID)
}

func (a removeDeletedResourceFromRole) PolicyDef() sdk.Policy {
	return sdk.Policy{
		Id:          a.id,
		Name:        a.Name(),
		Description: "This policy removes the deleted resource from all roles.",
		Definition: sdk.PolicyDefinition{
			Arguments: []sdk.PolicyArgument{},
		},
	}
}
