package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/middleware"
	"github.com/codeready-toolchain/devcluster/pkg/server"
	"github.com/codeready-toolchain/devcluster/test"
	"github.com/codeready-toolchain/toolchain-common/pkg/status"
	authsupport "github.com/codeready-toolchain/toolchain-common/pkg/test/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestAuthMiddlewareSuite struct {
	test.UnitTestSuite
}

func TestRunAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, &TestAuthMiddlewareSuite{test.UnitTestSuite{}})
}

func (s *TestAuthMiddlewareSuite) TestAuthMiddleware() {
	s.Run("create with DefaultTokenParser failing", func() {
		authMiddleware, err := middleware.NewAuthMiddleware()
		require.Nil(s.T(), authMiddleware)
		require.Error(s.T(), err)
		require.Equal(s.T(), "no default TokenParser created, call `InitializeDefaultTokenParser()` first", err.Error())
	})
}

func (s *TestAuthMiddlewareSuite) TestAuthMiddlewareService() {
	// create a TokenGenerator and a key
	tokengenerator := authsupport.NewTokenManager()
	kid0, err := uuid.NewV4()
	require.NoError(s.T(), err)
	_, err = tokengenerator.AddPrivateKey(kid0.String())
	require.NoError(s.T(), err)

	// create some test tokens
	id, err := uuid.NewV4()
	require.NoError(s.T(), err)
	username, err := uuid.NewV4()
	require.NoError(s.T(), err)
	identity0 := authsupport.Identity{
		ID:       id,
		Username: username.String(),
	}
	email, err := uuid.NewV4()
	require.NoError(s.T(), err)
	emailClaim0 := authsupport.WithEmailClaim(email.String() + "@email.tld")
	// valid token
	tokenValid, err := tokengenerator.GenerateSignedToken(identity0, kid0.String(), emailClaim0)
	require.NoError(s.T(), err)
	// invalid token - no email
	tokenInvalidNoEmail, err := tokengenerator.GenerateSignedToken(identity0, kid0.String())
	require.NoError(s.T(), err)
	// invalid token - garbage
	tokenInvalidGarbage, err := uuid.NewV4()
	require.NoError(s.T(), err)
	// invalid token - expired
	expTime := time.Now().Add(-60 * time.Second)
	expClaim := authsupport.WithExpClaim(expTime)
	tokenInvalidExpiredJWT := tokengenerator.GenerateToken(identity0, kid0.String(), emailClaim0, expClaim)
	tokenInvalidExpired, err := tokengenerator.SignToken(tokenInvalidExpiredJWT, kid0.String())
	require.NoError(s.T(), err)

	// start key service
	keysEndpointURL := tokengenerator.NewKeyServer().URL

	// create server
	srv := server.New(s.Config)

	// set the key service url in the config
	err = os.Setenv("DEVCLUSTER_AUTH_CLIENT_PUBLIC_KEYS_URL", keysEndpointURL)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), keysEndpointURL, srv.Config().GetAuthClientPublicKeysURL(), "key url not set correctly")
	err = os.Setenv("DEVCLUSTER_ENVIRONMENT", configuration.UnitTestsEnvironment)
	require.NoError(s.T(), err)

	// Setting up the routes.
	err = srv.SetupRoutes()
	require.NoError(s.T(), err)

	// Check that there are routes registered.
	routes := srv.GetRegisteredRoutes()
	require.NotEmpty(s.T(), routes)

	// Check that Engine() returns the router object.
	require.NotNil(s.T(), srv.Engine())

	s.Run("health check requests", func() {
		health := &status.Health{}
		resp := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
		require.NoError(s.T(), err)

		srv.Engine().ServeHTTP(resp, req)
		// Check the status code is what we expect.
		assert.Equal(s.T(), http.StatusOK, resp.Code, "request returned wrong status code")

		err = json.Unmarshal(resp.Body.Bytes(), health)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), health.Alive, true)
		assert.Equal(s.T(), health.Environment, "unit-tests")
		assert.Equal(s.T(), health.Revision, "0")
		assert.NotEqual(s.T(), health.BuildTime, "")
		assert.NotEqual(s.T(), health.StartTime, "")
	})

	s.Run("static request", func() {
		resp := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/favicon.ico", nil)
		require.NoError(s.T(), err)

		srv.Engine().ServeHTTP(resp, req)
		// Check the status code is what we expect.
		assert.Equal(s.T(), http.StatusOK, resp.Code, "request returned wrong status code")
	})

	s.Run("auth requests", func() {

		// do some requests
		var authtests = []struct {
			name        string
			urlPath     string
			method      string
			tokenHeader string
			status      int
		}{
			{"auth_test, no auth, denied", "/api/v1/auth_test", http.MethodGet, "", http.StatusUnauthorized},
			{"auth_test, valid header auth", "/api/v1/auth_test", http.MethodGet, "Bearer " + tokenValid, http.StatusOK},
			{"auth_test, invalid header auth, no email claim", "/api/v1/auth_test", http.MethodGet, "Bearer " + tokenInvalidNoEmail, http.StatusUnauthorized},
			{"auth_test, invalid header auth, expired", "/api/v1/auth_test", http.MethodGet, "Bearer " + tokenInvalidExpired, http.StatusUnauthorized},
			{"auth_test, invalid header auth, token garbage", "/api/v1/auth_test", http.MethodGet, "Bearer " + tokenInvalidGarbage.String(), http.StatusUnauthorized},
			{"auth_test, invalid header auth, wrong header format", "/api/v1/auth_test", http.MethodGet, tokenValid, http.StatusUnauthorized},
			{"auth_test, invalid header auth, bearer but no token", "/api/v1/auth_test", http.MethodGet, "Bearer ", http.StatusUnauthorized},
		}
		for _, tt := range authtests {
			s.Run(tt.name, func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(tt.method, tt.urlPath, nil)
				require.NoError(s.T(), err)
				if tt.tokenHeader != "" {
					req.Header.Set("Authorization", tt.tokenHeader)
				}
				srv.Engine().ServeHTTP(resp, req)
				// Check the status code is what we expect.
				assert.Equal(s.T(), tt.status, resp.Code, "request returned wrong status code")
			})
		}
	})
}
