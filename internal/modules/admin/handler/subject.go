package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	"github.com/fun-dotto/server/internal/modules/admin/middleware"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context, params api.SubjectsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	clientParams := &academic_api.SubjectsV1ListParams{
		Q:                         params.Q,
		Grades:                    convertSlicePtr[api.DottoFoundationV1Grade, academic_api.DottoFoundationV1Grade](params.Grades),
		Courses:                   convertSlicePtr[api.DottoFoundationV1Course, academic_api.DottoFoundationV1Course](params.Courses),
		Classes:                   convertSlicePtr[api.DottoFoundationV1Class, academic_api.DottoFoundationV1Class](params.Classes),
		Classifications:           convertSlicePtr[api.DottoFoundationV1SubjectClassification, academic_api.DottoFoundationV1SubjectClassification](params.Classifications),
		Year:                      params.Year,
		Semesters:                 convertSlicePtr[api.DottoFoundationV1CourseSemester, academic_api.DottoFoundationV1CourseSemester](params.Semesters),
		RequirementTypes:          convertSlicePtr[api.DottoFoundationV1SubjectRequirementType, academic_api.DottoFoundationV1SubjectRequirementType](params.RequirementTypes),
		CulturalSubjectCategories: convertSlicePtr[api.DottoFoundationV1CulturalSubjectCategory, academic_api.DottoFoundationV1CulturalSubjectCategory](params.CulturalSubjectCategories),
	}

	response, err := h.academicClient.SubjectsV1ListWithResponse(c.Request.Context(), clientParams)
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

// SubjectsV1Detail 科目を詳細取得する
func (h *Handler) SubjectsV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.SubjectsV1DetailWithResponse(c.Request.Context(), id)
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

// SubjectsV1Delete 科目を削除する
func (h *Handler) SubjectsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.SubjectsV1DeleteWithResponse(c.Request.Context(), id)
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

func convertSlicePtr[From, To ~string](src *[]From) *[]To {
	if src == nil {
		return nil
	}
	result := make([]To, len(*src))
	for i, v := range *src {
		result[i] = To(v)
	}
	return &result
}
