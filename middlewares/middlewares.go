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
var RoleMapContextKey = "roleMap"

func (m *Middlewares) ResourceMapper() fiber.Handler {
	return func(c *fiber.Ctx) error {
		model := models.GetResourceMap()
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
			resourceMap[res.ResourceId] = res.RoleId
		}
		c.Context().SetUserValue(ResourceMapContextKey, resourceMap)
		return c.Next()
	}
}

func (m *Middlewares) RoleMapper() fiber.Handler {
	return func(c *fiber.Ctx) error {
		model := models.GetRoleMap()
		roleMap := make(map[string][]string)
		// Fetch all roles from MongoDB
		cursor, err := m.db.Find(c.Context(), model, bson.D{})
		if err != nil {
			log.Errorw("failed to fetch roles", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load roles")
		}
		defer cursor.Close(context.Background())
		// Iterate through results and populate role map
		for cursor.Next(context.Background()) {
			var role models.RoleMap
			if err := cursor.Decode(&role); err != nil {
				log.Errorw("error decoding role map", "error", err)
				continue
			}
			roleMap[role.RoleId] = role.UserId
		}
		c.Context().SetUserValue(RoleMapContextKey, roleMap)
		fmt.Println(roleMap)
		return c.Next()
	}
}
