package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetById(t *testing.T) {

	t.Run("fetch user successfully", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()

		// client cursor
		clientCursor, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		clientCursor1, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor, nil).Once()
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor1, nil).Once()
		d.On("FindOne", mock.Anything, models.GetProjectModel(), mock.Anything, mock.Anything).Return(mongo.NewSingleResultFromDocument(models.Project{Id: "0001"}, nil, nil), nil)

		md := models.GetUserModel()
		mockResult := mongo.NewSingleResultFromDocument(models.User{Id: "0001"}, nil, nil)
		d.On("FindOne", mock.Anything, md, mock.Anything, mock.Anything).Return(mockResult, nil).Once()

		server.SetupTestServer(app, d)

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "0001", resp.Data.Id)
	})

	t.Run("user not found", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()

		// client cursor
		clientCursor, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		clientCursor1, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor, nil).Once()
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor1, nil).Once()
		d.On("FindOne", mock.Anything, models.GetProjectModel(), mock.Anything, mock.Anything).Return(mongo.NewSingleResultFromDocument(models.Project{Id: "0001"}, nil, nil), nil)

		md := models.GetUserModel()
		mockResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		d.On("FindOne", mock.Anything, md, mock.Anything, mock.Anything).Return(mockResult, mongo.ErrNoDocuments).Once()

		server.SetupTestServer(app, d)

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("internal error", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()

		// client cursor
		clientCursor, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		clientCursor1, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor, nil).Once()
		d.On("Find", mock.Anything, models.GetClientModel(), mock.Anything, mock.Anything).Return(clientCursor1, nil).Once()
		d.On("FindOne", mock.Anything, models.GetProjectModel(), mock.Anything, mock.Anything).Return(mongo.NewSingleResultFromDocument(models.Project{Id: "0001"}, nil, nil), nil)

		md := models.GetUserModel()
		mockResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("internal error"), nil)
		d.On("FindOne", mock.Anything, md, mock.Anything, mock.Anything).Return(mockResult, errors.New("internal error")).Once()

		server.SetupTestServer(app, d)

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

}
