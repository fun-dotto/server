package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	user_api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// UsersV1List ユーザー一覧を取得する
func (h *Handler) UsersV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.userClient.UsersV1ListWithResponse(c.Request.Context())
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

// UsersV1Detail ユーザーを取得する
func (h *Handler) UsersV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.userClient.UsersV1DetailWithResponse(c.Request.Context(), id)
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

// UsersV1Upsert ユーザーを作成または更新する
func (h *Handler) UsersV1Upsert(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req user_api.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userClient.UsersV1UpsertWithResponse(c.Request.Context(), id, req)
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
