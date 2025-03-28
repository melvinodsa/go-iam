package role

import (
	"context"
	"fmt"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Conversion functions with added input validation
func fromSdkToModel(role sdk.Role) models.Role {
	// Use zero time if pointers are nil to prevent nil pointer dereference
	createdAt := time.Time{}
	updatedAt := time.Time{}

	if role.CreatedAt != nil {
		createdAt = *role.CreatedAt
	}
	if role.UpdatedAt != nil {
		updatedAt = *role.UpdatedAt
	}

	return models.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromSdkResourceMapToModel(role.Resources),
		CreatedAt: createdAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelToSdk(role *models.Role) *sdk.Role {
	if role == nil {
		return nil
	}

	return &sdk.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromModelResourceMapToSdk(role.Resources),
		CreatedAt: &role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: &role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelListToSdk(roles []models.Role) []sdk.Role {
	result := make([]sdk.Role, 0, len(roles))
	for i := range roles {
		sdkRole := fromModelToSdk(&roles[i])
		if sdkRole != nil {
			result = append(result, *sdkRole)
		}
	}
	return result
}

// More robust resource map conversion functions
func fromSdkResourceMapToModel(resources map[string]sdk.Resources) map[string]models.Resources {
	result := make(map[string]models.Resources)
	for _, res := range resources {
		// Only add resources with non-empty keys
		if res.Key != "" {
			result[res.Key] = models.Resources{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}
	return result
}

func fromModelResourceMapToSdk(resources map[string]models.Resources) map[string]sdk.Resources {
	result := make(map[string]sdk.Resources)
	for key, res := range resources {
		// Only add resources with non-empty keys
		if key != "" {
			result[key] = sdk.Resources{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}
	return result
}

// Improved role and resource mapping functions with better error handling
func (s *store) addToRoleMap(ctx context.Context, roleId, userId string) error {
	if roleId == "" || userId == "" {
		return fmt.Errorf("role ID and user ID cannot be empty")
	}

	userMd := models.GetUserModel()
	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.M{userMd.IdKey: userId},
		bson.M{"$set": bson.M{fmt.Sprintf("roles.%s", roleId): true}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}
	return nil
}

func (s *store) removeFromRoleMap(ctx context.Context, roleId, userId string) error {
	if roleId == "" || userId == "" {
		return fmt.Errorf("role ID and user ID cannot be empty")
	}

	userMd := models.GetUserModel()
	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.M{userMd.IdKey: userId},
		bson.M{"$unset": bson.M{fmt.Sprintf("roles.%s", roleId): ""}},
	)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}

func (s *store) addToResourceMap(ctx context.Context, userId string, resources map[string]models.Resources) error {
	if userId == "" || len(resources) == 0 {
		return fmt.Errorf("user ID cannot be empty and resources must be non-empty")
	}

	userMd := models.GetUserModel()
	updateData := bson.M{}
	for key, res := range resources {
		updateData[fmt.Sprintf("resource.%s", key)] = res
	}

	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.M{userMd.IdKey: userId},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to add resource to user: %w", err)
	}
	return nil
}

func (s *store) removeFromResourceMap(ctx context.Context, userId string, resources map[string]sdk.Resources) error {
	if userId == "" || len(resources) == 0 {
		return fmt.Errorf("user ID cannot be empty and resources must be non-empty")
	}

	userMd := models.GetUserModel()
	updateData := bson.M{}
	for key := range resources {
		updateData[fmt.Sprintf("resources.%s", key)] = 1
	}

	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.M{userMd.IdKey: userId},
		bson.M{"$unset": updateData},
	)
	if err != nil {
		return fmt.Errorf("failed to remove resource from user: %w", err)
	}
	return nil
}

func (s *service) updateUserRoleAndResources(ctx context.Context, user *models.User) error {
	updatedRoles := make(map[string]models.UserRoles)
	resourceSet := make(map[string]models.UserResource)

	for roleId := range user.Roles {
		updatedRole, err := s.GetById(ctx, roleId)
		if err != nil {
			return fmt.Errorf("failed to fetch role %s: %w", roleId, err)
		}

		updatedRoles[roleId] = models.UserRoles{
			Id:   updatedRole.Id,
			Name: updatedRole.Name,
		}

		for _, res := range updatedRole.Resources {
			resourceSet[res.Key] = models.UserResource{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}

	// Persist changes to the database
	userMd := models.GetUserModel()
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"roles":    updatedRoles,
			"resource": resourceSet,
		}},
	}

	_, err := s.store.(*store).db.UpdateOne(
		ctx,
		userMd,
		bson.D{{Key: userMd.IdKey, Value: user.Id}},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *service) refreshUsersWithRole(ctx context.Context, role *sdk.Role) error {
	userMd := models.GetUserModel()

	// Find users with this role
	users, err := s.store.(*store).db.Find(
		ctx,
		userMd,
		bson.D{{Key: "roles." + role.Id, Value: bson.D{{Key: "$exists", Value: true}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch users with role: %w", err)
	}
	defer users.Close(ctx)

	// Batch update users
	for users.Next(ctx) {
		var user models.User
		if err := users.Decode(&user); err != nil {
			return fmt.Errorf("failed to decode user: %w", err)
		}

		if err := s.updateUserRoleAndResources(ctx, &user); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) calculateResourceUsage(userRoles map[string]models.UserRoles, roleToSkip string) map[string]int {
	resourceUsage := make(map[string]int)

	for roleId := range userRoles {
		if roleId == roleToSkip {
			continue
		}

		role, err := s.GetById(context.Background(), roleId)
		if err != nil {
			// Log the error or handle it appropriately
			continue
		}

		for _, res := range role.Resources {
			resourceUsage[res.Key]++
		}
	}

	return resourceUsage
}

// cleanupUnusedResources removes resources that are no longer used by any role
func (s *service) cleanupUnusedResources(
	currentResources map[string]models.UserResource,
	resourceUsage map[string]int,
) map[string]models.UserResource {
	cleanedResources := make(map[string]models.UserResource)

	for key, resource := range currentResources {
		// Keep the resource if it's used in other roles
		if resourceUsage[key] > 0 {
			cleanedResources[key] = resource
		}
	}

	return cleanedResources
}
