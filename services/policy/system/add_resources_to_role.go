package system

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type addResourcesToRole struct {
	id      string
	userSvc user.Service
	roleSvc role.Service
	pc      PolicyCheck
}

func NewAddResourcesToRole(userSvc user.Service, roleSvc role.Service) addResourcesToRole {
	return addResourcesToRole{id: "@policy/system/add_resources_to_role", userSvc: userSvc, roleSvc: roleSvc, pc: NewPolicyCheck(userSvc)}
}

func (a addResourcesToRole) ID() string {
	return a.id
}

func (a addResourcesToRole) Name() string {
	return "Add resources to role specified in user policy"
}

func (a addResourcesToRole) HandleEvent(event utils.Event[sdk.Resource]) {
	log.Debugw("received resource event", "event", event.Name())
	userId := event.Metadata().User.Id
	user, exists, err := a.pc.RunCheck(event.Context(), a.id, userId)
	if err != nil {
		log.Errorw("error checking user while handling resource create event", "userId", userId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	if !exists {
		return
	}
	targetRoleId, ok := a.getTargetRoleId(user)
	if !ok {
		return
	}
	err = a.roleSvc.AddResource(event.Context(), targetRoleId, sdk.Resources{
		Id:   event.Payload().ID,
		Key:  event.Payload().Key,
		Name: event.Payload().Name,
	})
	if err != nil {
		log.Errorw("error adding resource to role while handling resource create event", "role_id", targetRoleId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	log.Infow("successfully added created resource to role", "role_id", targetRoleId, "resource_id", event.Payload().ID)
}

func (a addResourcesToRole) getTargetRoleId(user *sdk.User) (string, bool) {
	policy, ok := user.Policies[a.id]
	if !ok {
		return "", false
	}
	arg, ok := policy.Mapping.Arguments["@roleId"]
	if !ok {
		return "", false
	}
	if len(arg.Static) == 0 {
		return "", false
	}
	return arg.Static, true
}

func (a addResourcesToRole) PolicyDef() sdk.Policy {
	return sdk.Policy{
		Id:          a.id,
		Name:        a.Name(),
		Description: "This policy adds the created resource to the role specified in the user policy.",
		Definition: sdk.PolicyDefinition{
			Arguments: []sdk.PolicyArgument{
				{
					Name:        "@roleId",
					Description: "The role to which the resource access is granted.",
					DataType:    goiamuniverse.Role,
				},
			},
		},
	}
}
