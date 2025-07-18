package user

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

type store struct {
	db db.DB
}

func NewStore(db db.DB) Store {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, user *sdk.User) error {
	id := uuid.New().String()
	user.Id = id
	t := time.Now()
	user.Enabled = true
	user.CreatedAt = &t
	d := fromSdkToModel(*user)

	md := models.GetUserModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *store) Update(ctx context.Context, user *sdk.User) error {
	now := time.Now()
	user.UpdatedAt = &now
	if user.Id == "" {
		return ErrorUserNotFound
	}
	o, err := s.GetById(ctx, user.Id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	user.CreatedAt = o.CreatedAt
	user.CreatedBy = o.CreatedBy
	d := fromSdkToModel(*user)
	md := models.GetUserModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: user.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *store) GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error) {
	md := models.GetUserModel()
	var usr models.User
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.EmailKey, Value: email}, {Key: md.ProjectIDKey, Value: projectId}}).Decode(&usr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return fromModelToSdk(&usr), nil
}

func (s *store) GetById(ctx context.Context, id string) (*sdk.User, error) {
	md := models.GetUserModel()
	var usr models.User
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&usr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return fromModelToSdk(&usr), nil
}

func (s *store) GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error) {
	md := models.GetUserModel()
	var usr models.User
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.PhoneKey, Value: phone}, {Key: md.ProjectIDKey, Value: projectId}}).Decode(&usr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return fromModelToSdk(&usr), nil
}

func (s *store) GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error) {
	md := models.GetUserModel()
	filter := bson.A{}
	var users []models.User
	if query.SearchQuery != "" {
		//  search by name, email, or phone (case-insensitive)
		filter = append(filter,
			bson.D{{Key: md.NameKey, Value: bson.D{{Key: "$regex", Value: query.SearchQuery}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: md.EmailKey, Value: bson.D{{Key: "$regex", Value: query.SearchQuery}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: md.PhoneKey, Value: bson.D{{Key: "$regex", Value: query.SearchQuery}, {Key: "$options", Value: "i"}}}},
		)
	}
	cond := bson.D{{Key: md.EnabledKey, Value: true}, {Key: md.ProjectIDKey, Value: bson.D{{Key: "$in", Value: query.ProjectIds}}}}

	if len(filter) > 0 {
		cond = append(cond, bson.E{Key: "$or", Value: filter})
	}

	// Get total count
	total, err := s.db.CountDocuments(ctx, md, cond)
	if err != nil {
		return nil, fmt.Errorf("error counting resources: %w", err)
	}

	opts := options.Find().
		SetSkip(query.Skip).
		SetLimit(query.Limit)

	cursor, err := s.db.Find(ctx, md, cond, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding all users: %w", err)
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorw(
				"error closing cursor after reading all users",
				"error", err)
		}
	}()
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, fmt.Errorf("error reading all users: %w", err)
	}
	return &sdk.UserList{
		Users: fromModelListToSdk(users),
		Total: total,
		Skip:  query.Skip,
		Limit: query.Limit,
	}, nil
}
