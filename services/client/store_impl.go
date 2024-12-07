package client

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
	return store{db: db}
}

func (s store) GetAll(ctx context.Context) ([]sdk.Client, error) {
	md := models.GetClientModel()
	var clients []models.Client
	cursor, err := s.db.Find(ctx, md, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("error finding all clients: %w", err)
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorw(
				"error closing cursor after reading all clients",
				"error", err)
		}
	}()
	err = cursor.All(ctx, &clients)
	if err != nil {
		return nil, fmt.Errorf("error reading all clients: %w", err)
	}
	return fromModelListToSdk(clients), nil
}
func (s store) Get(ctx context.Context, id string) (*sdk.Client, error) {
	md := models.GetClientModel()
	var client models.Client
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&client)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("error finding client: %w", err)
	}

	return fromModelToSdk(&client), nil
}
func (s store) Create(ctx context.Context, client *sdk.Client) error {
	id := uuid.New().String()
	client.Id = id
	t := time.Now()
	client.Enabled = true
	client.CreatedAt = &t
	d := fromSdkToModel(*client)

	// hash the client secret before storing it
	secret, err := hashSecret(d.Secret)
	if err != nil {
		return fmt.Errorf("error hashing client secret: %w", err)
	}
	d.Secret = secret
	md := models.GetClientModel()
	_, err = s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	return nil
}
func (s store) Update(ctx context.Context, client *sdk.Client) error {
	now := time.Now()
	client.UpdatedAt = &now
	if client.Id == "" {
		return ErrClientNotFound
	}
	o, err := s.Get(ctx, client.Id)
	if err != nil {
		return fmt.Errorf("error finding client: %w", err)
	}
	client.CreatedAt = o.CreatedAt
	client.CreatedBy = o.CreatedBy
	d := fromSdkToModel(*client)
	md := models.GetClientModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: client.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating client: %w", err)
	}

	return nil
}