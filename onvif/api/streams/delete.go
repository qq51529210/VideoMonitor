package streams

import (
	"net/http"
	"recordplan/api/internal"
	"recordplan/db"
	"recordplan/task/week"

	"github.com/gin-gonic/gin"
)

//	@Summary	删除
//	@Tags		媒体流
//	@Param		stream	path	string	true	"Stream"
//	@Accept		json
//	@Success	204
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/streams/{stream} [delete]
func delete(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	rows, err := db.DeleteWeekPlanStreamByStream(id.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if rows > 0 {
		week.DeleteStream(id.ID)
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}

//	@Summary	批量删除
//	@Tags		媒体流
//	@Param		stream	body	[]string	true	"Stream 数组"
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Success	204
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/streams [delete]
func batchDelete(ctx *gin.Context) {
	// 参数
	var req internal.BatchDelete[string]
	err := ctx.ShouldBindJSON(&req.ID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	rows, err := db.BatchDeleteWeekPlanStreamByStream(req.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if rows > 0 {
		for _, id := range req.ID {
			week.DeleteStream(id)
		}
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
