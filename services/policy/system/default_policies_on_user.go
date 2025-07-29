package system

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
)

type defaultPoliciesOnUser struct {
	id      string
	userSvc user.Service
}

func NewDefaultPoliciesOnUser(userSvc user.Service) defaultPoliciesOnUser {
	return defaultPoliciesOnUser{id: "@policy/system/default_policies_on_user", userSvc: userSvc}
}

func (a defaultPoliciesOnUser) ID() string {
	return a.id
}

func (a defaultPoliciesOnUser) HandleEvent(event utils.Event[sdk.User]) {
	log.Debugw("received user event", "event", event.Name())
	userId := event.Payload().Id
	user, err := a.userSvc.GetById(context.Background(), userId)
	if err != nil {
		log.Errorw("error fetching user while handling user create event", "userId", userId, "error", err)
		return
	}
	policies := user.Policies
	if len(policies) == 0 {
		policies = map[string]bool{}
	}
	policies[NewAccessToCreatedResource(nil).ID()] = true
	user.Policies = policies
	err = a.userSvc.Update(context.Background(), user)
	if err != nil {
		log.Errorw("error updating user with default policies while handling user create event", "userId", userId, "error", err)
		return
	}
	log.Infow("successfully added default policies to user", "userId", userId)
}
