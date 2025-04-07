package policy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	db db.DB
}

func NewStore(db db.DB) Store {
	return store{db: db}
}

// Create adds a new policy to the database
func (s store) Create(ctx context.Context, policy *sdk.Policy) error {
	if policy == nil {
		return errors.New("policy cannot be nil")
	}
	policy.Id = uuid.New().String()
	now := time.Now()
	policy.CreatedAt = &now
	d := fromSdkToModel(policy)
	md := models.GetPolicyModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}
	return nil
}

// Delete removes a policy by ID
func (s store) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("policy ID cannot be empty")
	}

	// Remove the policy key from all users
	userModel := models.GetUserModel()
	update := bson.M{
		"$unset": bson.M{
			fmt.Sprintf("policies.%s", id): "",
		},
	}
	_, err := s.db.UpdateMany(ctx, userModel, bson.M{}, update)
	if err != nil {
		return fmt.Errorf("failed to remove policy from users: %w", err)
	}

	// Delete the policy from the policies collection
	policyModel := models.GetPolicyModel()
	result, err := s.db.DeleteOne(ctx, policyModel, bson.D{{Key: policyModel.IdKey, Value: id}})
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("policy not found")
	}

	return nil
}

// Update modifies an existing policy in the database
func (s store) Update(ctx context.Context, policy *sdk.Policy) error {
	if policy == nil || policy.Id == "" {
		return errors.New("policy ID is required")
	}

	// Get the original policy to compare role changes
	originalPolicy := &sdk.Policy{}
	md := models.GetPolicyModel()
	result := s.db.FindOne(
		ctx,
		md,
		bson.D{{Key: md.IdKey, Value: policy.Id}},
	)
	err := result.Decode(originalPolicy)
	if err != nil {
		return fmt.Errorf("failed to find original policy: %w", err)
	}

	// Step 0: Update the policy in the database
	d := fromSdkToModel(policy)
	_, err = s.db.UpdateOne(
		ctx,
		md,
		bson.D{{Key: md.IdKey, Value: policy.Id}},
		bson.D{{Key: "$set", Value: d}},
	)
	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	// Step 1: Find all users who have any of the roles from this policy (both old and new roles)
	userMd := models.GetUserModel()

	// Collect all role IDs (both from original and updated policy)
	roleIds := map[string]struct{}{}
	for rid := range originalPolicy.Roles {
		roleIds[rid] = struct{}{}
	}
	for rid := range policy.Roles {
		roleIds[rid] = struct{}{}
	}

	roleIdList := make([]string, 0, len(roleIds))
	for id := range roleIds {
		roleIdList = append(roleIdList, id)
	}

	// Find all users with any of these roles
	orQuery := make([]bson.M, 0, len(roleIdList))
	for _, id := range roleIdList {
		orQuery = append(orQuery, bson.M{"roles." + id: bson.M{"$exists": true}})
	}

	var users []models.User
	cursor, err := s.db.Find(ctx, userMd, bson.M{"$or": orQuery})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return err
	}

	// Step 2: Get all roles for these users
	allRoleIds := map[string]struct{}{}
	for _, user := range users {
		for rid := range user.Roles {
			allRoleIds[rid] = struct{}{}
		}
	}

	allRoleIdList := make([]string, 0, len(allRoleIds))
	for id := range allRoleIds {
		allRoleIdList = append(allRoleIdList, id)
	}

	// Step 3: Get all policies associated with these roles
	policyMd := models.GetPolicyModel()
	policyOrQuery := make([]bson.M, 0, len(allRoleIdList))
	for _, id := range allRoleIdList {
		policyOrQuery = append(policyOrQuery, bson.M{"roles." + id: bson.M{"$exists": true}})
	}

	var allPolicies []sdk.Policy
	cursor, err = s.db.Find(ctx, policyMd, bson.M{"$or": policyOrQuery})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &allPolicies); err != nil {
		return err
	}

	// Map role -> policies
	roleToPolicies := map[string][]sdk.Policy{}
	for _, p := range allPolicies {
		for rid := range p.Roles {
			roleToPolicies[rid] = append(roleToPolicies[rid], p)
		}
	}

	// Step 4: Prepare BulkWrite models
	writeModels := make([]mongo.WriteModel, 0, len(users))
	for _, user := range users {
		// User has at least one role from this policy
		affectedUser := false
		for roleID := range user.Roles {
			if _, exists := roleIds[roleID]; exists {
				affectedUser = true
				break
			}
		}

		if !affectedUser {
			continue
		}

		// Rebuild policies based on user's roles
		combinedPolicies := map[string]string{}
		for roleID := range user.Roles {
			for _, pol := range roleToPolicies[roleID] {
				combinedPolicies[pol.Id] = pol.Name
			}
		}

		// Prepare update operation
		filter := bson.M{"id": user.Id}
		update := bson.M{
			"$set": bson.M{
				"policies":   combinedPolicies,
				"updated_at": time.Now(),
			},
		}
		writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	// Step 5: Execute BulkWrite
	if len(writeModels) > 0 {
		_, err = s.db.BulkWrite(ctx, userMd, writeModels)
		if err != nil {
			return fmt.Errorf("failed to bulk update users: %w", err)
		}
	}

	return nil
}

// GetById retrieves a policy by ID
func (s store) Get(ctx context.Context, id string) (*sdk.Policy, error) {
	if id == "" {
		return nil, errors.New("policy ID cannot be empty")
	}

	md := models.GetPolicyModel()
	var policy models.Policy

	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&policy)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("policy with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}
	return fromModelToSdk(&policy), nil
}

// GetAll retrieves all policies from the database
func (s store) GetAll(ctx context.Context) ([]sdk.Policy, error) {
	md := models.GetPolicyModel()
	var policies []models.Policy

	cursor, err := s.db.Find(ctx, md, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch policies: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &policies); err != nil {
		return nil, fmt.Errorf("failed to read policies: %w", err)
	}

	return fromModelListToSdk(policies), nil
}

// GetPoliciesByRoleId retrieves all policies associated with a specific role ID
func (s store) GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.Policy, error) {
	if roleId == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	md := models.GetPolicyModel()
	var policies []models.Policy

	cursor, err := s.db.Find(ctx, md, bson.D{{Key: "roles." + roleId, Value: bson.D{{Key: "$exists", Value: true}}}})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch policies: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &policies); err != nil {
		return nil, fmt.Errorf("failed to read policies: %w", err)
	}

	return fromModelListToSdk(policies), nil
}
