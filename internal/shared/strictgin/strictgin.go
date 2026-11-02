// Package strictgin は oapi-codegen の strict handler (StrictGinServerOptions)
// に渡すエラーハンドラ群を提供する。既定の実装と同じレスポンスを返しつつ、
// エラーを logging.RecordError に記録してアクセスログに出力させる。
package strictgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fun-dotto/server/internal/shared/logging"
)

// RequestErrorHandler はリクエストのパースに失敗したときに 400 を返す。
func RequestErrorHandler(c *gin.Context, err error) {
	logging.RecordError(c.Request.Context(), err)
	c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
}

// HandlerErrorHandler はハンドラがエラーを返したときに 500 を返す。
func HandlerErrorHandler(c *gin.Context, err error) {
	logging.RecordError(c.Request.Context(), err)
	c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
}

// ResponseErrorHandler はレスポンスの書き込みに失敗したときに 500 を返す。
func ResponseErrorHandler(c *gin.Context, err error) {
	logging.RecordError(c.Request.Context(), err)
	c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
}
