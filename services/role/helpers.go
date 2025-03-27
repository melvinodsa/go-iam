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
