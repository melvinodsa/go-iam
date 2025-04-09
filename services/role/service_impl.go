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
	userMd := models.GetUserModel()
	s.store.Update(ctx, role)

	// Step 1: Get all users with this role
	var users []models.User
	cursor, err := s.store.(*store).db.Find(ctx, userMd, bson.M{fmt.Sprintf("roles.%s", role.Id): bson.M{"$exists": true}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return err
	}

	// Step 2: Collect all unique role IDs across these users
	roleIds := map[string]struct{}{role.Id: {}}
	for _, user := range users {
		for rid := range user.Roles {
			roleIds[rid] = struct{}{}
		}
	}
	roleIdList := make([]string, 0, len(roleIds))
	for id := range roleIds {
		roleIdList = append(roleIdList, id)
	}

	// Step 3: Fetch all roles
	var allRoles []sdk.Role
	cursor, err = s.store.(*store).db.Find(ctx, models.GetRoleModel(), bson.M{"id": bson.M{"$in": roleIdList}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &allRoles); err != nil {
		return err
	}

	// Build roleResourcesMap
	roleResourcesMap := map[string]map[string]sdk.Resources{}
	for _, r := range allRoles {
		roleResourcesMap[r.Id] = r.Resources
	}

	// Step 4: Fetch all policies
	policyMd := models.GetPolicyModel()
	orQuery := make([]bson.M, 0, len(roleIdList))
	for _, id := range roleIdList {
		orQuery = append(orQuery, bson.M{"roles." + id: bson.M{"$exists": true}})
	}

	var allPolicies []sdk.Policy
	cursor, err = s.store.(*store).db.Find(ctx, policyMd, bson.M{"$or": orQuery})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &allPolicies); err != nil {
		return err
	}

	// Map role -> policies
	roleToPolicies := map[string][]sdk.Policy{}
	for _, policy := range allPolicies {
		for rid := range policy.Roles {
			roleToPolicies[rid] = append(roleToPolicies[rid], policy)
		}
	}

	// Build new policies map
	newPoliciesMap := map[string]string{}
	for _, p := range roleToPolicies[role.Id] {
		newPoliciesMap[p.Id] = p.Name
	}

	// Step 5: Prepare BulkWrite models
	writeModels := make([]mongo.WriteModel, 0, len(users))
	for _, user := range users {
		if _, ok := user.Roles[role.Id]; !ok {
			continue
		}

		// Update role metadata
		user.Roles[role.Id] = models.UserRoles{
			Id:   role.Id,
			Name: role.Name,
		}

		// Rebuild resources
		combinedResources := map[string]models.UserResource{}
		for rid := range user.Roles {
			resMap, ok := roleResourcesMap[rid]
			if !ok {
				continue
			}
			for key, res := range resMap {
				combinedResources[key] = models.UserResource{
					Id:   res.Id,
					Key:  res.Key,
					Name: res.Name,
				}
			}
		}
		user.Resources = combinedResources

		// Rebuild policies
		combinedPolicies := map[string]string{}
		for rid := range user.Roles {
			for _, pol := range roleToPolicies[rid] {
				combinedPolicies[pol.Id] = pol.Name
			}
		}
		user.Policies = combinedPolicies

		// Prepare update operation
		filter := bson.M{"id": user.Id}
		update := bson.M{
			"$set": bson.M{
				"roles":     user.Roles,
				"resources": user.Resources,
				"policies":  user.Policies,
			},
		}
		writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	// Step 6: Execute BulkWrite
	if len(writeModels) > 0 {
		_, err = s.store.(*store).db.BulkWrite(ctx, userMd, writeModels)
		if err != nil {
			return err
		}
	}

	// Step 7: Update role itself
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

	// Update user in the database
	err = s.store.RemoveRoleFromUser(ctx, &user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
