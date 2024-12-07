package project

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

func (s store) GetAll(ctx context.Context) ([]sdk.Project, error) {
	md := models.GetProjectModel()
	var projects []models.Project
	cursor, err := s.db.Find(ctx, md, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("error finding all projects: %w", err)
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorw(
				"error closing cursor after reading all projects",
				"error", err)
		}
	}()
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, fmt.Errorf("error reading all projects: %w", err)
	}
	return fromModelListToSdk(projects), nil
}
func (s store) Get(ctx context.Context, id string) (*sdk.Project, error) {
	md := models.GetProjectModel()
	var project models.Project
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("error finding project: %w", err)
	}

	return fromModelToSdk(&project), nil
}
func (s store) Create(ctx context.Context, project *sdk.Project) error {
	id := uuid.New().String()
	project.Id = id
	t := time.Now()
	project.CreatedAt = &t
	d := fromSdkToModel(*project)
	md := models.GetProjectModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}
	return nil
}
func (s store) Update(ctx context.Context, project *sdk.Project) error {
	now := time.Now()
	project.UpdatedAt = &now
	if project.Id == "" {
		return ErrProjectNotFound
	}
	o, err := s.Get(ctx, project.Id)
	if err != nil {
		return fmt.Errorf("error finding project: %w", err)
	}
	project.CreatedAt = o.CreatedAt
	project.CreatedBy = o.CreatedBy
	d := fromSdkToModel(*project)
	md := models.GetProjectModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: project.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating project: %w", err)
	}

	return nil
}
