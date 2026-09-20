package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// TimetableItemsV1List 時間割を取得する
func (h *Handler) TimetableItemsV1List(c *gin.Context, params api.TimetableItemsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.TimetableItemsV1ListWithResponse(c.Request.Context(), &academic_api.TimetableItemsV1ListParams{
		Year:      params.Year,
		Semesters: convertSlice[api.DottoFoundationV1CourseSemester, academic_api.DottoFoundationV1CourseSemester](params.Semesters),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": response.JSON200.TimetableItems,
	})
}

// TimetableItemsV1Create 時間割に追加する
func (h *Handler) TimetableItemsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.TimetableItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.TimetableItemsV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"item": response.JSON201.TimetableItem,
	})
}

// TimetableItemsV1Delete 時間割を削除する
func (h *Handler) TimetableItemsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.TimetableItemsV1DeleteWithResponse(c.Request.Context(), id)
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

func convertSlice[From, To ~string](src []From) []To {
	result := make([]To, len(src))
	for i, v := range src {
		result[i] = To(v)
	}
	return result
}
