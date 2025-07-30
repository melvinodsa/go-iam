package system

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
)

type PolicyCheck struct {
	userSvc user.Service
}

func NewPolicyCheck(userSvc user.Service) PolicyCheck {
	return PolicyCheck{userSvc: userSvc}
}

func (p PolicyCheck) RunCheck(ctx context.Context, id, userId string) (*sdk.User, bool, error) {
	user, err := p.userSvc.GetById(ctx, userId)
	if err != nil {
		return nil, false, err
	}
	if _, exists := user.Policies[id]; !exists {
		return user, false, nil
	}
	return user, true, nil
}
