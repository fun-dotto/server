package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	api "github.com/fun-dotto/server/gen/app"
	"github.com/gin-gonic/gin"
)

const (
	userIDKey    = "userID"
	userEmailKey = "userEmail"
)

// AuthMiddleware returns a Gin middleware that verifies Firebase ID tokens
// and stores the user ID in the context.
func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requiresBearerAuth(c) {
			c.Next()
			return
		}

		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.TrimPrefix(authorization, "Bearer ")

		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid ID token"})
			return
		}

		c.Set(userIDKey, token.UID)

		if email, ok := token.Claims["email"].(string); ok {
			c.Set(userEmailKey, email)
		}

		c.Next()
	}
}

func requiresBearerAuth(c *gin.Context) bool {
	_, ok := c.Get(api.BearerAuthScopes)
	return ok
}

// UserIDFromContext extracts the user ID from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// UserEmailFromContext extracts the user email from context.
func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}
