package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/qq51529210/log"
)

// Error 用于返回错误的响应
type Error struct {
	// 简短信息
	Phrase string `json:"phrase"`
	// 详细信息
	Detail any `json:"detail"`
}

func FormatError(err error) any {
	switch v := err.(type) {
	case validator.ValidationErrors:
		var errs []string
		for _, es := range v {
			switch es.Tag() {
			case "timePeroidEmpty":
				errs = append(errs, fmt.Sprintf("[%s] must not be empty", es.Field()))
			case "timePeroidTimeFormat":
				errs = append(errs, fmt.Sprintf("[%s] error format", es.Field()))
			case "timePeroidError":
				errs = append(errs, "[start] must less than [end]")
			default:
				errs = append(errs, es.Error())
			}
		}
		return errs
	case *validator.InvalidValidationError:
		return v.Error()
	default:
		return err.Error()
	}
}

// Handle400 返回 bind 的错误
func Handle400(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Phrase: "invalid data",
		Detail: FormatError(err),
	})
	ctx.Abort()
}

// Handle404 返回找不到的错误
func Handle404(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, &Error{
		Phrase: "data not found",
	})
	ctx.Abort()
}

// Handle500 返回服务错误
func Handle500(ctx *gin.Context, err error) {
	log.ErrorDepth(1, err)
	ctx.JSON(http.StatusInternalServerError, &Error{
		Phrase: "server error",
		Detail: err.Error(),
	})
	ctx.Abort()
}

// HandleDB500 返回数据库错误
func HandleDB500(ctx *gin.Context, err error) {
	log.ErrorDepth(1, err)
	ctx.JSON(http.StatusInternalServerError, &Error{
		Phrase: "database error",
		Detail: err.Error(),
	})
	ctx.Abort()
}

// HandleAPICall500 返回网络调用错误
func HandleAPICall500(ctx *gin.Context, err error) {
	log.ErrorDepth(1, err)
	ctx.JSON(http.StatusInternalServerError, &Error{
		Phrase: "api call error",
		Detail: err.Error(),
	})
	ctx.Abort()
}
