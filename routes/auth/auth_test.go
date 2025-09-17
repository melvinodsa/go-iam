package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

func TestLogin(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("redirect to login url successfully", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("GetLoginUrl", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("test-auth", nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/login?client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 307, res.StatusCode, "Expected status code 307")
		assert.Nil(t, err)
	})

	t.Run("fetch login url successfully", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("GetLoginUrl", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("test-auth", nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/login?client_id=10001&postback=true", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		var resp sdk.AuthLoginDataResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)

		assert.Equalf(t, 200, res.StatusCode, "Expected status code 307")
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.LoginUrl)
	})

	t.Run("invalid code challenge", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/login?client_id=10001&code_challenge_method=aaa", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
	})

	t.Run("fetch login url error", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("GetLoginUrl", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("some error")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/login?client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		var resp sdk.AuthLoginDataResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)

		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.NotNil(t, resp)
	})
}

func TestRedirect(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("redirect to login url successfully", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("Redirect", mock.Anything, mock.Anything, mock.Anything).Return(&sdk.AuthRedirectResponse{
			RedirectUrl: "test-auth",
		}, nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/authp-callback?code=1234&state=abcd&client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 307, res.StatusCode, "Expected status code 307")
		assert.Nil(t, err)
	})

	t.Run("fetch redirect url successfully", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("Redirect", mock.Anything, mock.Anything, mock.Anything).Return(&sdk.AuthRedirectResponse{
			RedirectUrl: "test-auth",
		}, nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/authp-callback?code=1234&state=abcd&postback=true&client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		var resp sdk.AuthRedirectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)

		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.RedirectUrl)
	})

	t.Run("redirect url error", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("Redirect", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/authp-callback?code=1234&state=abcd&client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		var resp sdk.AuthRedirectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)

		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.NotNil(t, resp)
	})
}

func TestVerify(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("verify code successfully for backend", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCallback", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sdk.AuthVerifyCodeResponse{
			AccessToken: "test-token",
		}, nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0") // test:test base64 encoded
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.AccessToken)
	})

	t.Run("verify code successfully for frontend", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCallback", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sdk.AuthVerifyCodeResponse{
			AccessToken: "test-token",
		}, nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&code_challenge=234&client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.AccessToken)
	})

	t.Run("verify code error", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCallback", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0") // test:test base64 encoded
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
	})

	t.Run("authorization header missing", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
	})

	t.Run("basic text missing in authorization header", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		req.Header.Set("Authorization", "ttt dGVzdA==") // invalid base64 encoded
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
	})

	t.Run("invalid base64 encoding in authorization header", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		req.Header.Set("Authorization", "Basic test:test") // invalid base64 encoded
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
	})

	t.Run("invalid authorization header", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("GET", "/auth/v1/verify?code=1234&client_id=10001", nil)
		req.Header.Set("Authorization", "Basic dGVzdA==")
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp sdk.AuthVerifyCodeResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
	})
}

func TestClientCredentials(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("client credentials success", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCredentials", mock.Anything, mock.Anything, mock.Anything).Return(&sdk.AuthVerifyCodeResponse{
			AccessToken: "test_access_token",
		}, nil).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`{
			"client_id": "test",
			"client_secret": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("client credentials error", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCredentials", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`{
			"client_id": "test",
			"client_secret": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 401, res.StatusCode, "Expected status code 401")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("client not found", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCredentials", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("invalid client_id")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`{
			"client_id": "test",
			"client_secret": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("client disabled", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth mock

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("ClientCredentials", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("client is disabled")).Once()

		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`{
			"client_id": "test",
			"client_secret": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 403, res.StatusCode, "Expected status code 403")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("empty payload", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("invalid payload", func(t *testing.T) {

		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/auth")

		req, _ := http.NewRequest("POST", "/auth/v1/client", strings.NewReader(`invalid`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)

		var resp sdk.ClientCredentialsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
