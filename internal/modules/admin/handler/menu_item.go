package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/server/gen/admin"
	funch_api "github.com/fun-dotto/server/gen/funch"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// MenuItemsV1List メニュー一覧を取得する
func (h *Handler) MenuItemsV1List(c *gin.Context, params api.MenuItemsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.funchClient.MenuItemsV1ListWithResponse(c.Request.Context(), &funch_api.MenuItemsV1ListParams{
		Date: params.Date,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, response.JSON200)
}
