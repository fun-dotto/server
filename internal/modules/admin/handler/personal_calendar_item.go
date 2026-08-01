package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// PersonalCalendarItemsV1List 個人カレンダーアイテム一覧を取得する
func (h *Handler) PersonalCalendarItemsV1List(c *gin.Context, params api.PersonalCalendarItemsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.PersonalCalendarItemsV1ListWithResponse(c.Request.Context(), &academic_api.PersonalCalendarItemsV1ListParams{
		UserId: params.UserId,
		Dates:  params.Dates,
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
