package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
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
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, role *sdk.Role) error {
	id := uuid.New().String()
	role.Id = id
	t := time.Now()
	role.CreatedAt = &t
	d := fromSdkToModel(*role)

	md := models.GetRoleModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating role: %w", err)
	}
	return nil
}

func (s *store) Update(ctx context.Context, role *sdk.Role) error {
	now := time.Now()
	role.UpdatedAt = &now
	if role.Id == "" {
		return errors.New("role not found")
	}
	o, err := s.GetById(ctx, role.Id)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}
	role.CreatedAt = o.CreatedAt
	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: role.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating role: %w", err)
	}
	return nil
}

func (s *store) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	md := models.GetRoleModel()
	var role models.Role
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("error finding role: %w", err)
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
		return nil, fmt.Errorf("error finding roles: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Errorw("error closing cursor after reading roles", "error", err)
		}
	}()

	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, fmt.Errorf("error reading roles: %w", err)
	}
	return fromModelListToSdk(roles), nil
}
