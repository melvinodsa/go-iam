package policy

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchAll(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	t.Run("fetch policies successfully", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		mockPolicySvc := services.MockPolicyService{}
		mockPolicySvc.On("GetAll", mock.Anything, mock.Anything).Return(&sdk.PolicyList{
			Policies: []sdk.Policy{
				{Id: "policy1", Name: "Test Policy 1"},
				{Id: "policy2", Name: "Test Policy 2"},
			},
			Total: 2,
		}, nil).Once()

		svcs.Policy = &mockPolicySvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))

		app.Get("/test", FetchAll)

		req, _ := http.NewRequest("GET", "/test", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp sdk.PoliciesResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Policies fetched successfully", resp.Message)
		assert.Len(t, resp.Data.Policies, 2)
	})

	t.Run("fetch policies with query parameters", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		mockPolicySvc := services.MockPolicyService{}
		mockPolicySvc.On("GetAll", mock.Anything, mock.MatchedBy(func(query sdk.PolicyQuery) bool {
			return query.Query == "test" && query.Skip == 5 && query.Limit == 20
		})).Return(&sdk.PolicyList{
			Policies: []sdk.Policy{{Id: "policy1", Name: "Test Policy"}},
			Total:    1,
		}, nil).Once()

		svcs.Policy = &mockPolicySvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))

		app.Get("/test", FetchAll)

		req, _ := http.NewRequest("GET", "/test?query=test&skip=5&limit=20", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp sdk.PoliciesResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("fetch policies with invalid pagination parameters", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		mockPolicySvc := services.MockPolicyService{}
		mockPolicySvc.On("GetAll", mock.Anything, mock.MatchedBy(func(query sdk.PolicyQuery) bool {
			return query.Skip == 0 && query.Limit == 10 // Should use defaults
		})).Return(&sdk.PolicyList{}, nil).Once()

		svcs.Policy = &mockPolicySvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))

		app.Get("/test", FetchAll)

		req, _ := http.NewRequest("GET", "/test?skip=invalid&limit=invalid", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("fetch policies service error", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		mockPolicySvc := services.MockPolicyService{}
		mockPolicySvc.On("GetAll", mock.Anything, mock.Anything).Return(nil, errors.New("service error")).Once()

		svcs.Policy = &mockPolicySvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))

		app.Get("/test", FetchAll)

		req, _ := http.NewRequest("GET", "/test", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp sdk.PoliciesResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.Message, "failed to get Policy")
	})
}