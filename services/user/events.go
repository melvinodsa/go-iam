package user

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

func (s *service) HandleEvent(e utils.Event[sdk.Role]) {
	switch e.Name() {
	case utils.EventRoleUpdated:
		s.handleRoleUpdate(e)
	default:
		return
	}

}

func (s *service) handleRoleUpdate(e utils.Event[sdk.Role]) {
	err := s.fetchAndUpdateUsersWithRole(e.Context(), e.Payload())
	if err != nil {
		log.Errorw("error fetching and updating users with role", "error", err)
	}
}

func (s *service) fetchAndUpdateUsersWithRole(ctx context.Context, role sdk.Role) error {
	// fetch all users with this role with page limits
	page := 1
	limit := 10
	for {
		users, err := s.store.GetAll(ctx, sdk.UserQuery{
			RoleId: role.Id,
			Skip:   int64((page - 1) * limit),
			Limit:  int64(limit),
		})
		if err != nil {
			return err
		}
		if len(users.Users) == 0 {
			break
		}
		if err := s.updateUsersWithRole(ctx, role, users.Users); err != nil {
			return err
		}
		page++
	}
	return nil
}

func (s *service) updateUsersWithRole(ctx context.Context, role sdk.Role, users []sdk.User) error {
	log.Debugw("updating users with role", "role_id", role.Id, "role_name", role.Name, "no_of_users", len(users))
	for i := range users {
		err := s.updateUser(ctx, role, &users[i])
		if err != nil {
			log.Errorw("error updating user with role changes. pausing the role update", "user_id", users[i].Id, "user_name", users[i].Name, "role_id", role.Id, "error", err)
			return err
		}
		log.Debugw("successfully updated the role for user", "user_id", users[i].Id, "role_id", role.Id, "index", i)
	}
	return nil
}

func (s *service) updateUser(ctx context.Context, role sdk.Role, user *sdk.User) error {
	// remove the role from the user obj
	removeRoleFromUserObj(user, role)
	// add the role to the user obj
	addRoleToUserObj(user, role)
	// update the user
	err := s.store.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
