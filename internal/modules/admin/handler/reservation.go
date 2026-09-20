package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// ReservationsV1List 教室の予約一覧を取得する
func (h *Handler) ReservationsV1List(c *gin.Context, params api.ReservationsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.ReservationsV1ListWithResponse(c.Request.Context(), &academic_api.ReservationsV1ListParams{
		RoomIds: params.RoomIds,
		From:    params.From,
		Until:   params.Until,
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

// ReservationsV1Create 教室を予約する
func (h *Handler) ReservationsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.ReservationsV1CreateWithResponse(c.Request.Context(), req)
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

// ReservationsV1Delete 予約を削除する
func (h *Handler) ReservationsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.ReservationsV1DeleteWithResponse(c.Request.Context(), id)
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
