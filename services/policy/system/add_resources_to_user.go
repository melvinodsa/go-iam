package system

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
)

type addResourcesToUser struct {
	id      string
	userSvc user.Service
	pc      PolicyCheck
}

func NewAddResourcesToUser(userSvc user.Service) addResourcesToUser {
	return addResourcesToUser{id: "@policy/system/add_resources_to_user", userSvc: userSvc, pc: NewPolicyCheck(userSvc)}
}

func (a addResourcesToUser) ID() string {
	return a.id
}

func (a addResourcesToUser) HandleEvent(event utils.Event[sdk.Resource]) {
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
	targetUserId, ok := a.getTargetUserId(user)
	if !ok {
		return
	}
	err = a.userSvc.AddResourceToUser(event.Context(), targetUserId, sdk.AddUserResourceRequest{
		PolicyId: a.id,
		Key:      event.Payload().Key,
		Name:     event.Payload().Name,
	})
	if err != nil {
		log.Errorw("error adding resource to user while handling resource create event", "user_id", targetUserId, "resource_id", event.Payload().ID, "error", err)
		return
	}
	log.Infow("successfully added created resource to user", "user_id", targetUserId, "resource_id", event.Payload().ID)
}

func (a addResourcesToUser) getTargetUserId(user *sdk.User) (string, bool) {
	policy, ok := user.Policies[a.id]
	if !ok {
		return "", false
	}
	arg, ok := policy.Mapping.Arguments["@userId"]
	if !ok {
		return "", false
	}
	if len(arg.Static) == 0 {
		return "", false
	}
	return arg.Static, true
}
