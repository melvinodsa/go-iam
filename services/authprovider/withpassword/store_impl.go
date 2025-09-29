package withpassword

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/hashing"
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

func (s *store) GetUserByUsername(ctx context.Context, email, projectID, password string) (*sdk.WithPasswordUser, error) {
	md := models.GetWithPasswordUserModel()
	var result models.WithPasswordUser
	hashedPassword, err := hashing.HashSecret(password)
	if err != nil {
		return nil, fmt.Errorf("error while hashing the secret. %w", err)
	}
	err = s.db.FindOne(ctx, md, bson.M{md.EmailKey: email, md.PasswordKey: hashedPassword, md.ProjectIDKey: projectID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sdk.ErrUserNotFound
		}
		return nil, err
	}

	return fromModelToSdk(&result), nil
}

func (s *store) CreateUser(ctx context.Context, email string, projectID string, password string) error {
	md := models.GetWithPasswordUserModel()
	hashedPassword, err := hashing.HashSecret(password)
	if err != nil {
		return fmt.Errorf("error while hashing the secret. %w", err)
	}
	now := time.Now()
	user := &models.WithPasswordUser{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: &now,
		CreatedBy: "system",
		UpdatedAt: &now,
		UpdatedBy: "system",
	}
	_, err = s.db.InsertOne(ctx, md, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return sdk.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *store) UpdateUserPassword(ctx context.Context, email, projectID, newPassword string) error {
	md := models.GetWithPasswordUserModel()
	hashedPassword, err := hashing.HashSecret(newPassword)
	if err != nil {
		return fmt.Errorf("error while hashing the secret. %w", err)
	}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			md.PasswordKey:  hashedPassword,
			md.UpdatedAtKey: &now,
			md.UpdatedByKey: "system",
		},
	}
	res, err := s.db.UpdateOne(ctx, md, bson.M{md.EmailKey: email, md.ProjectIDKey: projectID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return sdk.ErrUserNotFound
	}
	return nil
}
