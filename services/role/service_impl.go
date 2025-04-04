package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/policy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type service struct {
	store         Store
	policyService policy.Service
}

func NewService(store Store, policySvc policy.Service) Service {
	return &service{
		store:         store,
		policyService: policySvc,
	}
}
func (s *service) Create(ctx context.Context, role *sdk.Role) error {
	return s.store.Create(ctx, role)
}

func (s *service) Update(ctx context.Context, role *sdk.Role) error {
	if err := s.store.Update(ctx, role); err != nil {
		return err
	}
	return nil
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error) {
	return s.store.GetAll(ctx, query)
}

func (s *service) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	userMd := models.GetUserModel()
	var user models.User

	// Fetch User
	err := s.store.(*store).db.FindOne(ctx, userMd, bson.M{userMd.IdKey: userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user with ID %s not found", userId)
		} else {
			return fmt.Errorf("failed to fetch user: %w", err)
		}
	}

	// Fetch Role
	role, err := s.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	// Fetch Policies
	policies, err := s.policyService.GetPoliciesByRoleId(ctx, roleId)
	if err != nil {
		return err
	}

	// Initialize user's fields if nil
	if user.Roles == nil {
		user.Roles = make(map[string]models.UserRoles)
	}
	if user.Resources == nil {
		user.Resources = make(map[string]models.UserResource)
	}
	if user.Policies == nil {
		user.Policies = make(map[string]string)
	}

	// Skip if role already exists
	if _, exists := user.Roles[roleId]; exists {
		return nil
	}

	// Add unique policies
	for _, policy := range policies {
		if _, exists := user.Policies[policy.Id]; !exists {
			user.Policies[policy.Id] = policy.Name
		}
	}

	// Add new role
	user.Roles[roleId] = models.UserRoles{
		Id:   role.Id,
		Name: role.Name,
	}

	// Add unique resources from role
	for _, res := range role.Resources {
		if _, exists := user.Resources[res.Key]; !exists {
			user.Resources[res.Key] = models.UserResource{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}

	s.store.AddRoleToUser(ctx, &user)

	return nil
}

// removing a role from user, handled all scenarios in it [hopefully T-T]
func (s *service) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	fmt.Println("Removing role from user:", userId, roleId)

	userMd := models.GetUserModel()
	var user models.User

	// Fetch User
	err := s.store.(*store).db.FindOne(ctx, userMd, bson.M{userMd.IdKey: userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user with ID %s not found", userId)
		} else {
			return fmt.Errorf("failed to fetch user: %w", err)
		}
	}

	// Skip if role does not exist
	if _, exists := user.Roles[roleId]; !exists {
		return nil
	}

	// Fetch Role
	role, err := s.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	// Fetch Policies
	policies, err := s.policyService.GetPoliciesByRoleId(ctx, roleId)
	if err != nil {
		return err
	}

	// Ensure user's fields are initialized
	if user.Roles == nil {
		user.Roles = make(map[string]models.UserRoles)
	}
	if user.Resources == nil {
		user.Resources = make(map[string]models.UserResource)
	}
	if user.Policies == nil {
		user.Policies = make(map[string]string)
	}

	fmt.Print("Removing role from user: ", userId, " ", roleId, " ", user.Roles[roleId], "\n")

	// update user roles
	delete(user.Roles, roleId)

	fmt.Println("Removing role from user:", user)

	// Remove policies only if no other roles require them
	for _, policy := range policies {
		policyStillNeeded := false
		for rId := range user.Roles {
			otherPolicies, _ := s.policyService.GetPoliciesByRoleId(ctx, rId)
			for _, otherPolicy := range otherPolicies {
				if otherPolicy.Id == policy.Id {
					policyStillNeeded = true
					break
				}
			}
			if policyStillNeeded {
				break
			}
		}
		if !policyStillNeeded {
			delete(user.Policies, policy.Id)
		}
	}

	// Remove resources only if no other roles require them
	for _, res := range role.Resources {
		resourceStillNeeded := false
		for rId := range user.Roles {
			otherRole, _ := s.GetById(ctx, rId)
			for _, otherRes := range otherRole.Resources {
				if otherRes.Key == res.Key {
					resourceStillNeeded = true
					break
				}
			}
			if resourceStillNeeded {
				break
			}
		}
		if !resourceStillNeeded {
			delete(user.Resources, res.Key)
		}
	}
	fmt.Println("Removing role from user:", user)

	// Update user in the database
	err = s.store.RemoveRoleFromUser(ctx, &user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
