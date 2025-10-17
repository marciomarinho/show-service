package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marciomarinho/show-service/internal/config"
)

// CognitoClaims represents the claims in a Cognito JWT token
type CognitoClaims struct {
	Sub             string   `json:"sub"`
	CognitoGroups   []string `json:"cognito:groups"`
	TokenUse        string   `json:"token_use"`
	Scope           string   `json:"scope"`
	AuthTime        int64    `json:"auth_time"`
	Iss             string   `json:"iss"`
	CognitoUsername string   `json:"cognito:username"`
	Exp             int64    `json:"exp"`
	Iat             int64    `json:"iat"`
	ClientID        string   `json:"client_id"`
	Username        string   `json:"username"`
}

// UserContext represents authenticated user information
type UserContext struct {
	UserID   string
	Username string
	Groups   []string
}

// getRequiredScope determines the required scope for an endpoint from configured scopes
func getRequiredScope(path, method string, validScopes []string) string {
	switch {
	case path == "/shows" && method == "GET":
		// Look for read scope for shows endpoint
		for _, scope := range validScopes {
			if strings.Contains(scope, "shows.read") {
				return scope
			}
		}
		return "" // No read scope found
	case path == "/shows" && method == "POST":
		// Look for write scope for shows endpoint
		for _, scope := range validScopes {
			if strings.Contains(scope, "shows.write") {
				return scope
			}
		}
		return "" // No write scope found
	default:
		return ""
	}
}

// hasValidScope checks if the token has the required scope
func hasValidScope(tokenScopes, requiredScope string) bool {
	if requiredScope == "" {
		return true // No scope required
	}

	// Split token scopes (space-separated)
	scopes := strings.Fields(tokenScopes)
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// hasValidScopeFromConfig checks if the token has the required scope and validates against configured valid scopes
func hasValidScopeFromConfig(tokenScopes, requiredScope string, validScopes []string) bool {
	// First check if the required scope is in the configured valid scopes
	if requiredScope != "" {
		found := false
		for _, validScope := range validScopes {
			if validScope == requiredScope {
				found = true
				break
			}
		}
		if !found {
			return false // Required scope not in valid scopes list
		}
	}

	// Then check if token has the required scope
	return hasValidScope(tokenScopes, requiredScope)
}

// AuthMiddleware validates JWT tokens for non-local environments
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for local environment
		if cfg.Env == config.EnvLocal {
			c.Next()
			return
		}

		// Extract bearer token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// TODO: Implement proper JWT validation with AWS Cognito
		// For now, we'll do basic token format validation
		if len(tokenString) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// TODO: Add proper JWT validation using github.com/golang-jwt/jwt/v5
		// This would involve:
		// 1. Fetching Cognito JWKS from the configured URL
		// 2. Validating token signature against public keys
		// 3. Checking token expiration and claims
		// 4. Verifying issuer, audience, and token use

		// Get required scope for this endpoint
		requiredScope := getRequiredScope(c.Request.URL.Path, c.Request.Method, cfg.Cognito.ValidScopes)

		// If no scope is required for this endpoint, skip scope validation
		if requiredScope == "" {
			// Placeholder: Add basic user context for development
			userCtx := &UserContext{
				UserID:   "placeholder-user-id",
				Username: "placeholder-user",
				Groups:   []string{"users"},
			}
			c.Set("user", userCtx)
			c.Next()
			return
		}

		// TODO: In full implementation, extract scope from validated JWT token
		// For now, we'll use a placeholder scope for development
		// In production, this would come from the validated JWT token claims
		tokenScopes := strings.Join(cfg.Cognito.ValidScopes, " ")

		// Validate scope against configured valid scopes
		if !hasValidScopeFromConfig(tokenScopes, requiredScope, cfg.Cognito.ValidScopes) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Insufficient scope",
				"required_scope": requiredScope,
				"token_scopes":   tokenScopes,
				"valid_scopes":   cfg.Cognito.ValidScopes,
			})
			c.Abort()
			return
		}

		// Placeholder: Add basic user context for development
		userCtx := &UserContext{
			UserID:   "placeholder-user-id",
			Username: "placeholder-user",
			Groups:   []string{"users"},
		}
		c.Set("user", userCtx)

		c.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from gin context
func GetUserFromContext(c *gin.Context) (*UserContext, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	userCtx, ok := user.(*UserContext)
	if !ok {
		return nil, fmt.Errorf("invalid user context type")
	}

	return userCtx, nil
}
