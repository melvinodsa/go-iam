package role

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	db db.DB
}

func NewStore(db db.DB) Store {
	return &store{
		db: db,
	}
}

// Create adds a new role to the database
func (s *store) Create(ctx context.Context, role *sdk.Role) error {
	if role == nil {
		return errors.New("role cannot be nil")
	}
	role.Id = uuid.New().String()
	now := time.Now()
	role.CreatedAt = &now
	role.UpdatedAt = &now
	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// Update only handles database update, removes complex logic
func (s *store) Update(ctx context.Context, role *sdk.Role) error {
	if role == nil || role.Id == "" {
		return errors.New("role ID is required")
	}

	now := time.Now()
	role.UpdatedAt = &now

	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	result, err := s.db.UpdateOne(
		ctx,
		md,
		bson.D{{Key: md.IdKey, Value: role.Id}},
		bson.D{{Key: "$set", Value: d}},
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("role not found")
	}

	return nil
}

func (s *store) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	if id == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	md := models.GetRoleModel()
	var role models.Role

	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("role with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return fromModelToSdk(&role), nil
}

func (s *store) GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error) {
	md := models.GetRoleModel()
	var roles []models.Role

	filter := bson.D{}
	if query.ProjectId != "" {
		filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: query.ProjectId})
	}
	if query.SearchQuery != "" {
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{Key: md.NameKey, Value: bson.D{{Key: "$regex", Value: query.SearchQuery}, {Key: "$options", Value: "i"}}}},
			},
		})
	}

	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &roles); err != nil {
		return nil, fmt.Errorf("failed to read roles: %w", err)
	}

	return fromModelListToSdk(roles), nil
}

// Simplified database operations for AddRoleToUser and RemoveRoleFromUser
func (s *store) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	userMd := models.GetUserModel()
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("roles.%s", roleId): true,
		},
	}

	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.D{{Key: userMd.IdKey, Value: userId}},
		update,
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}
	return nil
}

func (s *store) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	userMd := models.GetUserModel()
	update := bson.M{
		"$unset": bson.M{
			fmt.Sprintf("roles.%s", roleId): "",
		},
	}

	_, err := s.db.UpdateOne(
		ctx,
		userMd,
		bson.D{{Key: userMd.IdKey, Value: userId}},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}
