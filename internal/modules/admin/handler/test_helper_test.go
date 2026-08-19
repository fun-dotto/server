package handler

import (
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"
	academic_api "github.com/fun-dotto/server/gen/academic"
	announcement_api "github.com/fun-dotto/server/gen/announcement"
	funch_api "github.com/fun-dotto/server/gen/funch"
	user_api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
	"github.com/gin-gonic/gin"
)

func newTestHandler(t *testing.T, baseURL string) *Handler {
	t.Helper()

	academicClient, err := academic_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new academic client: %v", err)
	}
	announcementClient, err := announcement_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new announcement client: %v", err)
	}
	funchClient, err := funch_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new funch client: %v", err)
	}
	userClient, err := user_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new user client: %v", err)
	}

	return NewHandler(academicClient, announcementClient, funchClient, userClient)
}

func setAdminClaim(c *gin.Context) {
	c.Set(middleware.FirebaseTokenContextKey, &firebaseauth.Token{
		Claims: map[string]interface{}{
			"admin": true,
		},
	})
}
