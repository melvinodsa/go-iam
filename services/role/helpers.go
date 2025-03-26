package role

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func fromSdkToModel(role sdk.Role) models.Role {
	return models.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromSdkResourceListToModel(role.Resources),
		CreatedAt: *role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: *role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelToSdk(role *models.Role) *sdk.Role {
	return &sdk.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromModelResourceListToSdk(role.Resources),
		CreatedAt: &role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: &role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelListToSdk(roles []models.Role) []sdk.Role {
	result := make([]sdk.Role, len(roles))
	for i := range roles {
		result[i] = *fromModelToSdk(&roles[i])
	}
	return result
}

func fromSdkResourceListToModel(resources []sdk.Resources) []models.Resources {
	result := make([]models.Resources, len(resources))
	for i, res := range resources {
		result[i] = models.Resources{
			Key:  res.Key,
			Name: res.Name,
		}
	}
	return result
}

func fromModelResourceListToSdk(resources []models.Resources) []sdk.Resources {
	result := make([]sdk.Resources, len(resources))
	for i, res := range resources {
		result[i] = sdk.Resources{
			Key:  res.Key,
			Name: res.Name,
		}
	}
	return result
}

// Add user to role mapping
func (s *store) addToRoleMap(ctx context.Context, roleId string, userId string) error {
	roleMapMd := models.GetRoleMap()
	_, err := s.db.UpdateOne(
		ctx,
		roleMapMd,
		bson.M{roleMapMd.RoleIdKey: roleId},
		bson.M{"$addToSet": bson.M{roleMapMd.UserIdKey: userId}}, // Prevents duplicates
		options.Update().SetUpsert(true),                         // Creates if not exists
	)
	if err != nil {
		return fmt.Errorf("error updating roleMap: %w", err)
	}
	return nil
}

// Remove user from role mapping
func (s *store) removeFromRoleMap(ctx context.Context, roleId string, userId string) error {
	roleMapMd := models.GetRoleMap()
	_, err := s.db.UpdateOne(
		ctx,
		roleMapMd,
		bson.M{roleMapMd.RoleIdKey: roleId},
		bson.M{"$pull": bson.M{roleMapMd.UserIdKey: userId}}, // Removes the userId
	)
	if err != nil {
		return fmt.Errorf("error updating roleMap: %w", err)
	}
	return nil
}

// Add role to resource mapping
func (s *store) addToResourceMap(ctx context.Context, roleId string, resources []sdk.Resources) error {
	resMd := models.GetResourceMap()
	for _, res := range resources {
		_, err := s.db.UpdateOne(
			ctx,
			resMd,
			bson.M{resMd.ResourceIdKey: res.Key},
			bson.M{"$addToSet": bson.M{resMd.RoleIdKey: roleId}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("error updating resource map: %w", err)
		}
	}
	return nil
}

// Remove role from resource mapping
func (s *store) removeFromResourceMap(ctx context.Context, roleId string, resources []sdk.Resources) error {
	resMd := models.GetResourceMap()
	for _, res := range resources {
		_, err := s.db.UpdateOne(
			ctx,
			resMd,
			bson.M{resMd.ResourceIdKey: res.Key},
			bson.M{"$pull": bson.M{resMd.RoleIdKey: roleId}}, // Removes roleId
		)
		if err != nil {
			return fmt.Errorf("error updating resource map: %w", err)
		}
	}
	return nil
}
