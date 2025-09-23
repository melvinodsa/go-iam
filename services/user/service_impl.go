package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type service struct {
	store   Store
	e       utils.Emitter[utils.Event[sdk.User], sdk.User]
	roleSvc role.Service
}

func NewService(store Store, roleSvc role.Service) Service {
	return &service{
		store:   store,
		roleSvc: roleSvc,
		e:       utils.NewEmitter[utils.Event[sdk.User]](),
	}
}

func (s *service) Create(ctx context.Context, user *sdk.User) error {
	err := s.store.Create(ctx, user)
	if err != nil {
		return err
	}
	s.Emit(newEvent(ctx, goiamuniverse.EventUserCreated, *user, middlewares.GetMetadata(ctx)))
	return nil
}

func (s *service) Update(ctx context.Context, user *sdk.User) error {
	return s.store.Update(ctx, user)
}

func (s *service) GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error) {
	return s.store.GetByEmail(ctx, email, projectId)
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.User, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error) {
	return s.store.GetByPhone(ctx, phone, projectId)
}

func (s *service) GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error) {
	return s.store.GetAll(ctx, query)
}

func (s *service) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return err
	}

	// Fetch Role
	role, err := s.roleSvc.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	// Skip if role already exists
	if _, exists := user.Roles[role.Id]; exists {
		return nil
	}

	addRoleToUserObj(user, *role)

	err = s.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}

	return nil
}

// removing a role from user, handled all scenarios in it [hopefully T-T]
func (s *service) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return err
	}

	// Skip if role does not exist
	if _, exists := user.Roles[roleId]; !exists {
		return nil
	}

	// Fetch Role
	role, err := s.roleSvc.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	removeRoleFromUserObj(user, *role)

	// Update user in the database
	err = s.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *service) AddResourceToUser(ctx context.Context, userId string, request sdk.AddUserResourceRequest) error {
	usr, err := s.store.GetById(ctx, userId)
	if err != nil {
		return err
	}

	addResourceToUserObj(usr, request)

	// Update user in the database
	err = s.store.Update(ctx, usr)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *service) AddPolicyToUser(ctx context.Context, userId string, policies map[string]sdk.UserPolicy) error {
	usr, err := s.store.GetById(ctx, userId)
	if err != nil {
		return err
	}

	addPoliciesToUserObj(usr, policies)

	// Update user in the database
	err = s.store.Update(ctx, usr)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *service) RemovePolicyFromUser(ctx context.Context, userId string, policyIds []string) error {
	usr, err := s.store.GetById(ctx, userId)
	if err != nil {
		return err
	}

	removePoliciesFromUserObj(usr, policyIds)

	// Update user in the database
	err = s.store.Update(ctx, usr)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *service) TransferOwnership(ctx context.Context, userId, newOwnerId string) error {
	if userId == "" || newOwnerId == "" {
		return errors.New("user ID and new owner ID are required")
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return err
	}

	newOwner, err := s.GetById(ctx, newOwnerId)
	if err != nil {
		return err
	}

	// transfer roles
	for roleId := range user.Roles {
		// Skip if role already exists
		if _, exists := newOwner.Roles[roleId]; exists {
			continue
		}
		role, err := s.roleSvc.GetById(ctx, roleId)
		if err != nil {
			log.Warnf("failed to fetch role %s: %v", roleId, err)
			continue
		}
		addRoleToUserObj(newOwner, *role)
	}

	// transfer resources
	for resourceKey, resource := range user.Resources {
		if _, exists := newOwner.Resources[resourceKey]; exists {
			// merge policies if resource already exists
			for policyId, policy := range resource.PolicyIds {
				newOwner.Resources[resourceKey].PolicyIds[policyId] = policy
			}
		} else {
			newOwner.Resources[resourceKey] = resource
		}
	}

	// transfer policies
	for policyId, policy := range user.Policies {
		newOwner.Policies[policyId] = policy
	}

	// Update new owner in the database
	err = s.store.Update(ctx, newOwner)
	if err != nil {
		return fmt.Errorf("failed to update new owner: %w", err)
	}

	// remove all roles, resources, policies from old user
	user.Roles = make(map[string]sdk.UserRole)
	user.Resources = make(map[string]sdk.UserResource)
	user.Policies = make(map[string]sdk.UserPolicy)

	// Update old user in the database
	err = s.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update old user: %w", err)
	}

	return nil
}

func (s *service) RemoveResourceFromAll(ctx context.Context, resourceKey string) error {
	return s.store.RemoveResourceFromAll(ctx, resourceKey)
}

func (s service) Emit(event utils.Event[sdk.User]) {
	if event == nil {
		return
	}
	s.e.Emit(event)
}

func (s service) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.User], sdk.User]) {
	s.e.Subscribe(eventName, subscriber)
}

type event struct {
	name     goiamuniverse.Event
	payload  sdk.User
	metadata sdk.Metadata
	ctx      context.Context
}

func (e event) Name() goiamuniverse.Event {
	return e.name
}

func (e event) Payload() sdk.User {
	return e.payload
}

func (e event) Metadata() sdk.Metadata {
	return e.metadata
}

func (e event) Context() context.Context {
	return e.ctx
}

func newEvent(ctx context.Context, name goiamuniverse.Event, payload sdk.User, metadata sdk.Metadata) utils.Event[sdk.User] {
	return event{ctx: ctx, name: name, payload: payload, metadata: metadata}
}
