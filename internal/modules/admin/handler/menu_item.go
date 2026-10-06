package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

func (h *Handler) MenuItemsV1List(c *gin.Context, _ api.MenuItemsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "menu items are not available"})
}

func (h *Handler) MenuItemsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
