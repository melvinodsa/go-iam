package resource

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type store struct {
	db db.DB
}

func NewStore(db db.DB) Store {
	return store{db: db}
}

func (s store) Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error) {
	md := models.GetResourceModel()
	filter := bson.D{}

	if query.Name != "" {
		filter = append(filter, bson.E{Key: md.NameKey, Value: query.Name})
	}
	if query.Description != "" {
		filter = append(filter, bson.E{Key: md.DescriptionKey, Value: query.Description})
	}

	// Get total count
	total, err := s.db.CountDocuments(ctx, md, filter)
	if err != nil {
		return nil, fmt.Errorf("error counting resources: %w", err)
	}

	// Set up options for pagination
	opts := options.Find().
		SetSkip(query.Skip).
		SetLimit(query.Limit)

	var resources []models.Resource
	cursor, err := s.db.Find(ctx, md, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding resources: %w", err)
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorw(
				"error closing cursor after reading resources",
				"error", err)
		}
	}()

	err = cursor.All(ctx, &resources)
	if err != nil {
		return nil, fmt.Errorf("error reading resources: %w", err)
	}

	return &sdk.ResourceList{
		Resources: fromModelListToSdk(resources),
		Total:     total,
		Skip:      query.Skip,
		Limit:     query.Limit,
	}, nil
}

func (s store) Get(ctx context.Context, id string) (*sdk.Resource, error) {
	md := models.GetResourceModel()
	var resource models.Resource
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&resource)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrResourceNotFound
		}
		return nil, fmt.Errorf("error finding resource: %w", err)
	}

	return fromModelToSdk(&resource), nil
}

func (s store) Create(ctx context.Context, resource *sdk.Resource) (string, error) {
	id := uuid.New().String()
	resource.ID = id
	t := time.Now()
	resource.CreatedAt = &t
	d := fromSdkToModel(*resource)
	md := models.GetResourceModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return "", fmt.Errorf("error creating resource: %w", err)
	}
	return id, nil
}

func (s store) Update(ctx context.Context, resource *sdk.Resource) error {
	now := time.Now()
	resource.UpdatedAt = &now
	if resource.ID == "" {
		return ErrResourceNotFound
	}
	o, err := s.Get(ctx, resource.ID)
	if err != nil {
		return fmt.Errorf("error finding resource: %w", err)
	}
	resource.CreatedAt = o.CreatedAt
	resource.CreatedBy = o.CreatedBy
	d := fromSdkToModel(*resource)
	md := models.GetResourceModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: resource.ID}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating resource: %w", err)
	}

	return nil
}

func (s store) Delete(ctx context.Context, id string) error {
	md := models.GetResourceModel()
	//mark it isenabled false
	_, err := s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}, bson.D{{Key: "$set", Value: bson.D{{Key: md.EnabledKey, Value: false}}}})
	if err != nil {
		return fmt.Errorf("error deleting resource: %w", err)
	}
	return nil
}
