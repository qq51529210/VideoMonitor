package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error 用于返回错误的响应
type Error struct {
	// 简短信息
	Phrase string `json:"phrase,omitempty"`
	// 详细信息
	Detail string `json:"detail,omitempty"`
}

const (
	ErrorSubmittedData      = "error submitted data"
	ErrorServerException    = "server exception"
	ErrorDataBaseAccess     = "database access error"
	ErrorEmptySubmittedData = "submitted data is empty"
	ErrorDataNotFound       = "data not found"
)

// WriteErrorSubmittedData 返回 bind 的错误
func WriteErrorSubmittedData(ctx *gin.Context, detail error) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Phrase: ErrorSubmittedData,
		Detail: detail.Error(),
	})
	ctx.Abort()
}

// WriteErrorEmptySubmittedData 返回提交数据为空的错误
func WriteErrorEmptySubmittedData(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Phrase: ErrorSubmittedData,
		Detail: ErrorEmptySubmittedData,
	})
	ctx.Abort()
}

// WriteErrorDataNotFound 返回查询不到数据的错误
func WriteErrorDataNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Phrase: ErrorDataNotFound,
	})
	ctx.Abort()
}

// WriteErrorDataBaseAccess 返回数据库访问的错误
func WriteErrorDataBaseAccess(ctx *gin.Context, detail error) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Phrase: ErrorDataBaseAccess,
		Detail: detail.Error(),
	})
	ctx.Abort()
}
