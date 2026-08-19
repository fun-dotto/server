package middleware

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/appcheck"
	api "github.com/fun-dotto/server/gen/app"
	"github.com/gin-gonic/gin"
)

// AppCheckMiddleware returns a Gin middleware that verifies Firebase App Check tokens.
func AppCheckMiddleware(appCheckClient *appcheck.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requiresFirebaseAppCheck(c) {
			c.Next()
			return
		}

		token := c.GetHeader("X-Firebase-AppCheck")
		if token == "" {
			message := "X-Firebase-AppCheck header is required"
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
			return
		}

		_, err := appCheckClient.VerifyToken(strings.Replace(token, "Bearer ", "", 1))
		if err != nil {
			message := "X-Firebase-AppCheck header value is invalid"
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
			return
		}

		c.Next()
	}
}

func requiresFirebaseAppCheck(c *gin.Context) bool {
	_, ok := c.Get(api.FirebaseAppCheckAuthScopes)
	return ok
}
