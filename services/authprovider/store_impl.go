package authprovider

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
	"github.com/melvinodsa/go-iam/services/encrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	enc encrypt.Service
	db  db.DB
}

func NewStore(enc encrypt.Service, db db.DB) Store {
	return store{enc: enc, db: db}
}

func (s store) Get(ctx context.Context, id string) (*sdk.AuthProvider, error) {
	md := models.GetAuthProviderModel()
	var provider models.AuthProvider
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&provider)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrAuthProviderNotFound
		}
		return nil, fmt.Errorf("error finding auth provider: %w", err)
	}

	err = s.decryptSecrets(&provider)
	if err != nil {
		return nil, fmt.Errorf("error decrypting auth provider secrets: %w", err)
	}

	return fromModelToSdk(&provider), nil
}

func (s store) GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error) {
	md := models.GetAuthProviderModel()
	var providers []models.AuthProvider

	// fetching the values from db
	filter := bson.D{}
	if len(params.ProjectIds) > 0 {
		filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: params.ProjectIds}}})
	}
	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding all auth providers: %w", err)
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorw(
				"error closing cursor after reading all auth providers",
				"error", err)
		}
	}()
	err = cursor.All(ctx, &providers)
	if err != nil {
		return nil, fmt.Errorf("error reading all auth providers: %w", err)
	}

	// decrypting the secrets
	for i := range providers {
		err = s.decryptSecrets(&providers[i])
		if err != nil {
			return nil, fmt.Errorf("error decrypting auth provider secrets: %w", err)
		}
	}

	// converting to sdk
	return fromModelListToSdk(providers), nil
}

func (s store) Create(ctx context.Context, provider *sdk.AuthProvider) error {
	id := uuid.New().String()
	provider.Id = id
	t := time.Now()
	provider.Enabled = true
	provider.CreatedAt = &t
	d := fromSdkToModel(*provider)

	// encrypt the secrets before storing
	err := s.encryptSecrets(&d)
	if err != nil {
		return fmt.Errorf("error encrypting auth provider secrets: %w", err)
	}

	md := models.GetAuthProviderModel()
	_, err = s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating auth provider: %w", err)
	}
	return nil
}

func (s store) Update(ctx context.Context, provider *sdk.AuthProvider) error {
	t := time.Now()
	provider.UpdatedAt = &t
	if provider.Id == "" {
		return ErrAuthProviderNotFound
	}
	o, err := s.Get(ctx, provider.Id)
	if err != nil {
		return fmt.Errorf("error finding auth provider: %w", err)
	}

	provider.CreatedAt = o.CreatedAt
	provider.CreatedBy = o.CreatedBy
	d := fromSdkToModel(*provider)

	// encrypt the secrets before storing
	err = s.encryptSecrets(&d)
	if err != nil {
		return fmt.Errorf("error encrypting auth provider secrets: %w", err)
	}

	md := models.GetAuthProviderModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: provider.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating auth provider: %w", err)
	}
	return nil
}

func (s store) decryptSecrets(provider *models.AuthProvider) error {
	for i := range provider.Params {
		if provider.Params[i].IsSecret {
			decrypted, err := s.enc.Decrypt(provider.Params[i].Value)
			if err != nil {
				return fmt.Errorf("error decrypting auth provider secret at %d : %w", i, err)
			}
			provider.Params[i].Value = decrypted
		}
	}
	return nil
}

func (s store) encryptSecrets(provider *models.AuthProvider) error {
	for i := range provider.Params {
		if provider.Params[i].IsSecret {
			encrypted, err := s.enc.Encrypt(provider.Params[i].Value)
			if err != nil {
				return fmt.Errorf("error encrypting auth provider secret at %d : %w", i, err)
			}
			provider.Params[i].Value = encrypted
		}
	}
	return nil
}
