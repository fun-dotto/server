package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// FacultyRoomsV1List 教員室割当一覧を取得する
func (h *Handler) FacultyRoomsV1List(c *gin.Context, params api.FacultyRoomsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.FacultyRoomsV1ListWithResponse(c.Request.Context(), &academic_api.FacultyRoomsV1ListParams{
		Year: params.Year,
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

// FacultyRoomsV1Create 教員室割当を作成する
func (h *Handler) FacultyRoomsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.FacultyRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.FacultyRoomsV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, response.JSON201)
}

// FacultyRoomsV1Delete 教員室割当を削除する
func (h *Handler) FacultyRoomsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.FacultyRoomsV1DeleteWithResponse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.StatusCode() != http.StatusNoContent {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.Status(http.StatusNoContent)
}
