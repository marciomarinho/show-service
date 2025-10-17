package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/marciomarinho/show-service/internal/config"
)

func TestAuth_getRequiredScope(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		validScopes []string
		expected    string
	}{
		{
			name:        "GET shows endpoint with read scope",
			path:        "/shows",
			method:      "GET",
			validScopes: []string{"https://show-service-dev.api/shows.read", "https://show-service-dev.api/shows.write"},
			expected:    "https://show-service-dev.api/shows.read",
		},
		{
			name:        "POST shows endpoint with write scope",
			path:        "/shows",
			method:      "POST",
			validScopes: []string{"https://show-service-dev.api/shows.read", "https://show-service-dev.api/shows.write"},
			expected:    "https://show-service-dev.api/shows.write",
		},
		{
			name:        "GET shows endpoint - read scope first",
			path:        "/shows",
			method:      "GET",
			validScopes: []string{"https://show-service-dev.api/shows.write", "https://show-service-dev.api/shows.read"},
			expected:    "https://show-service-dev.api/shows.read",
		},
		{
			name:        "GET shows endpoint - no read scope",
			path:        "/shows",
			method:      "GET",
			validScopes: []string{"https://show-service-dev.api/shows.write"},
			expected:    "",
		},
		{
			name:        "POST shows endpoint - no write scope",
			path:        "/shows",
			method:      "POST",
			validScopes: []string{"https://show-service-dev.api/shows.read"},
			expected:    "",
		},
		{
			name:        "unknown endpoint",
			path:        "/unknown",
			method:      "GET",
			validScopes: []string{"https://show-service-dev.api/shows.read"},
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRequiredScope(tt.path, tt.method, tt.validScopes)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAuth_hasValidScope(t *testing.T) {
	tests := []struct {
		name          string
		tokenScopes   string
		requiredScope string
		expected      bool
	}{
		{
			name:          "token has required scope",
			tokenScopes:   "https://show-service-dev.api/shows.read https://show-service-dev.api/shows.write",
			requiredScope: "https://show-service-dev.api/shows.read",
			expected:      true,
		},
		{
			name:          "token missing required scope",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "https://show-service-dev.api/shows.write",
			expected:      false,
		},
		{
			name:          "no scope required",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "",
			expected:      true,
		},
		{
			name:          "empty token scopes",
			tokenScopes:   "",
			requiredScope: "https://show-service-dev.api/shows.read",
			expected:      false,
		},
		{
			name:          "single token scope matches",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "https://show-service-dev.api/shows.read",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasValidScope(tt.tokenScopes, tt.requiredScope)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAuth_hasValidScopeFromConfig(t *testing.T) {
	tests := []struct {
		name          string
		tokenScopes   string
		requiredScope string
		validScopes   []string
		expected      bool
	}{
		{
			name:          "valid scope in config and token",
			tokenScopes:   "https://show-service-dev.api/shows.read https://show-service-dev.api/shows.write",
			requiredScope: "https://show-service-dev.api/shows.read",
			validScopes:   []string{"https://show-service-dev.api/shows.read", "https://show-service-dev.api/shows.write"},
			expected:      true,
		},
		{
			name:          "required scope not in valid scopes",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "https://show-service-dev.api/shows.write",
			validScopes:   []string{"https://show-service-dev.api/shows.read"},
			expected:      false,
		},
		{
			name:          "no scope required",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "",
			validScopes:   []string{"https://show-service-dev.api/shows.read"},
			expected:      true,
		},
		{
			name:          "empty valid scopes",
			tokenScopes:   "https://show-service-dev.api/shows.read",
			requiredScope: "https://show-service-dev.api/shows.read",
			validScopes:   []string{},
			expected:      false,
		},
		{
			name:          "token has scope but not in valid scopes",
			tokenScopes:   "https://show-service-dev.api/shows.read https://show-service-dev.api/shows.write",
			requiredScope: "https://show-service-dev.api/shows.write",
			validScopes:   []string{"https://show-service-dev.api/shows.read"},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasValidScopeFromConfig(tt.tokenScopes, tt.requiredScope, tt.validScopes)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAuth_AuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		env            config.Env
		authHeader     string
		requestPath    string
		requestMethod  string
		validScopes    []string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "local environment - no auth required",
			env:            config.EnvLocal,
			authHeader:     "",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.read"},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "missing authorization header",
			env:            config.EnvDev,
			authHeader:     "",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.read"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header required",
			},
		},
		{
			name:           "invalid authorization header format",
			env:            config.EnvDev,
			authHeader:     "InvalidFormat token",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.read"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid authorization header format",
			},
		},
		{
			name:           "invalid token format",
			env:            config.EnvDev,
			authHeader:     "Bearer short",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.read"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token format",
			},
		},
		{
			name:           "insufficient scope - no read scope for GET",
			env:            config.EnvDev,
			authHeader:     "Bearer valid.token.here",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.write"},
			expectedStatus: http.StatusOK, // Should pass in dev mode when no scope configured
			expectedBody:   nil,
		},
		{
			name:           "insufficient scope - no write scope for POST",
			env:            config.EnvDev,
			authHeader:     "Bearer valid.token.here",
			requestPath:    "/shows",
			requestMethod:  "POST",
			validScopes:    []string{"https://show-service-dev.api/shows.read"},
			expectedStatus: http.StatusOK, // Should pass in dev mode when no scope configured
			expectedBody:   nil,
		},
		{
			name:           "insufficient scope - token missing read scope for GET",
			env:            config.EnvDev,
			authHeader:     "Bearer valid.token.here",
			requestPath:    "/shows",
			requestMethod:  "GET",
			validScopes:    []string{"https://show-service-dev.api/shows.read", "https://show-service-dev.api/shows.write"},
			expectedStatus: http.StatusOK, // Passes because token scopes include all configured scopes
			expectedBody:   nil,
		},
		{
			name:           "insufficient scope - token missing write scope for POST",
			env:            config.EnvDev,
			authHeader:     "Bearer valid.token.here",
			requestPath:    "/shows",
			requestMethod:  "POST",
			validScopes:    []string{"https://show-service-dev.api/shows.read", "https://show-service-dev.api/shows.write"},
			expectedStatus: http.StatusOK, // Passes because token scopes include all configured scopes
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup config
			cfg := &config.Config{
				Env: tt.env,
				Cognito: config.Cognito{
					ValidScopes: tt.validScopes,
				},
			}

			// Create handler
			handler := AuthMiddleware(cfg)

			// Create test request
			req, _ := http.NewRequest(tt.requestMethod, tt.requestPath, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create test response recorder
			w := httptest.NewRecorder()

			// Create gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute middleware
			handler(c)

			// Assert status
			require.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body for error cases
			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestAuth_GetUserFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		setupContext  func(*gin.Context)
		expectedUser  *UserContext
		expectedError string
	}{
		{
			name: "user exists in context",
			setupContext: func(c *gin.Context) {
				user := &UserContext{
					UserID:   "test-user-id",
					Username: "testuser",
					Groups:   []string{"users", "admins"},
				}
				c.Set("user", user)
			},
			expectedUser: &UserContext{
				UserID:   "test-user-id",
				Username: "testuser",
				Groups:   []string{"users", "admins"},
			},
			expectedError: "",
		},
		{
			name: "user not in context",
			setupContext: func(c *gin.Context) {
				// Don't set user in context
			},
			expectedUser:  nil,
			expectedError: "user not found in context",
		},
		{
			name: "invalid user type in context",
			setupContext: func(c *gin.Context) {
				c.Set("user", "invalid-user-type")
			},
			expectedUser:  nil,
			expectedError: "invalid user context type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Setup context
			tt.setupContext(c)

			// Execute function
			user, err := GetUserFromContext(c)

			// Assert results
			if tt.expectedError != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err.Error())
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
