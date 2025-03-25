package middlewares

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/services/project"
	"go.mongodb.org/mongo-driver/bson"
)

type Middlewares struct {
	projectSvc project.Service
	db         db.DB
}

func NewMiddlewares(projectSvc project.Service, db db.DB) *Middlewares {
	return &Middlewares{
		projectSvc: projectSvc,
		db:         db,
	}
}

var ResourceMapContextKey = "resourceMap"

func (m *Middlewares) ResourceMapper() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("ResourceMapper middleware")
		model := models.GetRoleMap()
		resourceMap := make(map[string][]string)

		// Fetch all resources from MongoDB
		cursor, err := m.db.Find(c.Context(), model, bson.D{})
		if err != nil {
			log.Errorw("failed to fetch resources", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load resources")
		}
		defer cursor.Close(context.Background())

		// Iterate through results and populate resource map
		for cursor.Next(context.Background()) {
			var res models.ResourceMap
			if err := cursor.Decode(&res); err != nil {
				log.Errorw("error decoding resource map", "error", err)
				continue
			}
			resourceMap[res.ResourcecId] = res.RoleId
		}

		c.Context().SetUserValue(ResourceMapContextKey, resourceMap)

		return c.Next()
	}
}
