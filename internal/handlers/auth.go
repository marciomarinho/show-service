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

func getRequiredScope(path, method string, validScopes []string) string {
	switch {
	case path == "/shows" && method == "GET":
		for _, scope := range validScopes {
			if strings.Contains(scope, "shows.read") {
				return scope
			}
		}
		return ""
	case path == "/shows" && method == "POST":
		for _, scope := range validScopes {
			if strings.Contains(scope, "shows.write") {
				return scope
			}
		}
		return ""
	default:
		return ""
	}
}

func hasValidScope(tokenScopes, requiredScope string) bool {
	if requiredScope == "" {
		return true
	}

	scopes := strings.Fields(tokenScopes)
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

func hasValidScopeFromConfig(tokenScopes, requiredScope string, validScopes []string) bool {
	if requiredScope != "" {
		found := false
		for _, validScope := range validScopes {
			if validScope == requiredScope {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return hasValidScope(tokenScopes, requiredScope)
}

// AuthMiddleware validates JWT tokens for non-local environments
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Env == config.EnvLocal {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// For now, we'll do basic token format validation
		if len(tokenString) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		requiredScope := getRequiredScope(c.Request.URL.Path, c.Request.Method, cfg.Cognito.ValidScopes)

		if requiredScope == "" {
			userCtx := &UserContext{
				UserID:   "placeholder-user-id",
				Username: "placeholder-user",
				Groups:   []string{"users"},
			}
			c.Set("user", userCtx)
			c.Next()
			return
		}

		// For now, we'll use a placeholder scope for development
		// In production, this would come from the validated JWT token claims
		tokenScopes := strings.Join(cfg.Cognito.ValidScopes, " ")

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

		userCtx := &UserContext{
			UserID:   "placeholder-user-id",
			Username: "placeholder-user",
			Groups:   []string{"users"},
		}
		c.Set("user", userCtx)

		c.Next()
	}
}

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
