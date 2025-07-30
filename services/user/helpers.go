package user

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromSdkToModel(user sdk.User) models.User {
	return models.User{
		Id:         user.Id,
		Email:      user.Email,
		Phone:      user.Phone,
		Name:       user.Name,
		ProjectId:  user.ProjectId,
		Enabled:    user.Enabled,
		Expiry:     user.Expiry,
		ProfilePic: user.ProfilePic,
		Roles:      fromSdkUserRoleMapToModel(user.Roles),
		Resources:  fromSdkUserResourceMapToModel(user.Resources),
		Policies:   fromSdkUserPoliciesToModel(user.Policies),
		CreatedAt:  user.CreatedAt,
		CreatedBy:  user.CreatedBy,
		UpdatedAt:  user.UpdatedAt,
		UpdatedBy:  user.UpdatedBy,
	}
}

func fromModelToSdk(user *models.User) *sdk.User {
	return &sdk.User{
		Id:         user.Id,
		Email:      user.Email,
		Phone:      user.Phone,
		Name:       user.Name,
		ProfilePic: user.ProfilePic,
		ProjectId:  user.ProjectId,
		Expiry:     user.Expiry,
		Enabled:    user.Enabled,
		Roles:      fromModelUserRoleMapToSdk(user.Roles),
		Resources:  fromModelUserResourceMapToSdk(user.Resources),
		Policies:   fromModelUserPoliciesToSdk(user.Policies),
		CreatedAt:  user.CreatedAt,
		CreatedBy:  user.CreatedBy,
		UpdatedAt:  user.UpdatedAt,
		UpdatedBy:  user.UpdatedBy,
	}
}

func fromSdkUserPoliciesToModel(policies map[string]sdk.UserPolicy) map[string]models.UserPolicy {
	result := map[string]models.UserPolicy{}
	for key, policy := range policies {
		result[key] = models.UserPolicy{
			Def:     policy.Def,
			Mapping: fromSdkUserPolicyMappingToModel(policy.Mapping),
		}
	}
	return result
}

func fromSdkUserPolicyMappingToModel(mapping sdk.UserPolicyMapping) models.UserPolicyMapping {
	arguments := map[string]models.UserPolicyMappingValue{}
	for key, argument := range mapping.Arguments {
		arguments[key] = models.UserPolicyMappingValue{
			Static: argument.Static,
		}
	}
	return models.UserPolicyMapping{Arguments: arguments}
}

func fromModelUserPoliciesToSdk(policies map[string]models.UserPolicy) map[string]sdk.UserPolicy {
	result := map[string]sdk.UserPolicy{}
	for key, policy := range policies {
		result[key] = sdk.UserPolicy{
			Def:     policy.Def,
			Mapping: fromModelUserPolicyMappingToSdk(policy.Mapping),
		}
	}
	return result
}

func fromModelUserPolicyMappingToSdk(mapping models.UserPolicyMapping) sdk.UserPolicyMapping {
	arguments := map[string]sdk.UserPolicyMappingValue{}
	for key, argument := range mapping.Arguments {
		arguments[key] = sdk.UserPolicyMappingValue{
			Static: argument.Static,
		}
	}
	return sdk.UserPolicyMapping{Arguments: arguments}
}

// Convert SDK UserRole map to Model UserRoles map (Key: Name)
func fromSdkUserRoleMapToModel(roles map[string]sdk.UserRole) map[string]models.UserRoles {
	userRoles := make(map[string]models.UserRoles)
	for key, role := range roles {
		userRoles[key] = models.UserRoles{
			Name: role.Name,
			Id:   role.Id,
		}
	}
	return userRoles
}

// Convert Model UserRoles map to SDK UserRole map (Key: Name)
func fromModelUserRoleMapToSdk(roles map[string]models.UserRoles) map[string]sdk.UserRole {
	userRoles := make(map[string]sdk.UserRole)
	for key, role := range roles {
		userRoles[key] = sdk.UserRole{
			Name: role.Name,
			Id:   role.Id,
		}
	}
	return userRoles
}

// Convert SDK UserResource map to Model UserResource map (Key: Key)
func fromSdkUserResourceMapToModel(resources map[string]sdk.UserResource) map[string]models.UserResource {
	userResources := make(map[string]models.UserResource)
	for key, res := range resources {
		userResources[key] = models.UserResource{
			PolicyIds: res.PolicyIds,
			RoleIds:   res.RoleIds,
			Key:       res.Key,
			Name:      res.Name,
		}
	}
	return userResources
}

// Convert Model UserResource map to SDK UserResource map (Key: Key)
func fromModelUserResourceMapToSdk(resources map[string]models.UserResource) map[string]sdk.UserResource {
	userResources := make(map[string]sdk.UserResource)
	for key, res := range resources {
		userResources[key] = sdk.UserResource{
			PolicyIds: res.PolicyIds,
			RoleIds:   res.RoleIds,
			Key:       res.Key,
			Name:      res.Name,
		}
	}
	return userResources
}

// Convert list of Model Users to list of SDK Users
func fromModelListToSdk(users []models.User) []sdk.User {
	result := []sdk.User{}
	for i := range users {
		result = append(result, *fromModelToSdk(&users[i]))
	}
	return result
}

func removeRoleFromUserObj(user *sdk.User, role sdk.Role) {
	// Ensure user's fields are initialized
	if user.Roles == nil {
		user.Roles = make(map[string]sdk.UserRole)
	}
	if user.Resources == nil {
		user.Resources = make(map[string]sdk.UserResource)
	}

	// update user roles
	delete(user.Roles, role.Id)

	// Remove resources only if no other roles require them
	for _, res := range role.Resources {

		vl, exists := user.Resources[res.Key]
		if !exists {
			continue
		}
		delete(vl.RoleIds, role.Id)

		// there are no requirement of the resource as no one needs it
		if len(vl.RoleIds) == 0 && len(vl.PolicyIds) == 0 {
			delete(user.Resources, res.Key)
		} else {
			user.Resources[res.Key] = vl
		}
	}
}

func addRoleToUserObj(user *sdk.User, role sdk.Role) {
	// Initialize user's fields if nil
	if user.Roles == nil {
		user.Roles = make(map[string]sdk.UserRole)
	}
	if user.Resources == nil {
		user.Resources = make(map[string]sdk.UserResource)
	}

	// Add new role
	user.Roles[role.Id] = sdk.UserRole{
		Id:   role.Id,
		Name: role.Name,
	}

	// Add unique resources from role
	for _, res := range role.Resources {
		// other ran roleids policy ids cuold also exist that is why special treatment for resources
		existingResource, exists := user.Resources[res.Key]
		if !exists {
			existingResource = sdk.UserResource{
				RoleIds: map[string]bool{role.Id: true},
				Key:     res.Key,
				Name:    res.Name,
			}
		} else {
			if len(existingResource.RoleIds) == 0 {
				existingResource.RoleIds = map[string]bool{}
			}
			existingResource.RoleIds[role.Id] = true
		}
		user.Resources[res.Key] = existingResource
	}
}

func addResourceToUserObj(user *sdk.User, res sdk.AddUserResourceRequest) {
	// Initialize user's fields if nil
	if user.Resources == nil {
		user.Resources = make(map[string]sdk.UserResource)
	}

	// other ran roleids policy ids cuold also exist that is why special treatment for resources
	existingResource, exists := user.Resources[res.Key]
	if !exists {
		existingResource = sdk.UserResource{
			PolicyIds: map[string]bool{res.PolicyId: true},
			Key:       res.Key,
			Name:      res.Name,
		}
	} else {
		if len(existingResource.PolicyIds) == 0 {
			existingResource.PolicyIds = map[string]bool{}
		}
		existingResource.PolicyIds[res.PolicyId] = true
	}
	user.Resources[res.Key] = existingResource
}

func addPoliciesToUserObj(user *sdk.User, policies map[string]sdk.UserPolicy) {
	// Initialize user's fields if nil
	if user.Policies == nil {
		user.Policies = make(map[string]sdk.UserPolicy)
	}

	for key, policy := range policies {
		user.Policies[key] = policy
	}
}
