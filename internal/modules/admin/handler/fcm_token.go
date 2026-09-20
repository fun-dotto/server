package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/server/gen/admin"
	user_api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// FCMTokenV1List FCMトークン一覧を取得する
func (h *Handler) FCMTokenV1List(c *gin.Context, params api.FCMTokenV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.userClient.FCMTokenV1ListWithResponse(c.Request.Context(), &user_api.FCMTokenV1ListParams{
		UserIds:       params.UserIds,
		Tokens:        params.Tokens,
		UpdatedAtFrom: params.UpdatedAtFrom,
		UpdatedAtTo:   params.UpdatedAtTo,
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

// FCMTokenV1Upsert FCMトークンを作成または更新する
func (h *Handler) FCMTokenV1Upsert(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req user_api.FCMTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userClient.FCMTokenV1UpsertWithResponse(c.Request.Context(), req)
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
