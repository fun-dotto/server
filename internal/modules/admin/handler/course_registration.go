package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// CourseRegistrationsV1List 履修情報を取得する
func (h *Handler) CourseRegistrationsV1List(c *gin.Context, params api.CourseRegistrationsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1ListWithResponse(c.Request.Context(), &academic_api.CourseRegistrationsV1ListParams{
		UserId:    params.UserId,
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
		"registrations": response.JSON200.CourseRegistrations,
	})
}

// CourseRegistrationsV1Create 履修情報を作成する
func (h *Handler) CourseRegistrationsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.CourseRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"registration": response.JSON201.CourseRegistration,
	})
}

// CourseRegistrationsV1Delete 履修情報を削除する
func (h *Handler) CourseRegistrationsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1DeleteWithResponse(c.Request.Context(), id)
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
