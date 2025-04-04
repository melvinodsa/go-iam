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
	md := models.GetPolicyModel()
	result, err := s.db.DeleteOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}})
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
	d := fromSdkToModel(policy)
	md := models.GetPolicyModel()
	_, err := s.db.UpdateOne(
		ctx,
		md,
		bson.D{{Key: md.IdKey, Value: policy.Id}},
		bson.D{{Key: "$set", Value: d}},
	)
	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
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
