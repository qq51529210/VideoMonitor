package streams

import (
	"github.com/gin-gonic/gin"
)

type getReq struct {
	// 是否需要录像
	Stream string `form:"stream" binding:"required"`
}

type getRes struct {
	// 是否需要录像
	Recording bool `json:"recording"`
}

//	@Summary		录像状态
//	@Description	查询指定的流的录像状态
//	@Tags			周计划
//	@Param			stream	path	string	true	"stream"
//	@Produce		json
//	@Success		200	{object}	getRes
//	@Failure		404	{object}	internal.Error
//	@Failure		500	{object}	internal.Error
//	@Router			/streams/{stream}/status [get]
func get(ctx *gin.Context) {
	// // 参数
	// var req getReq
	// err := ctx.ShouldBindUri(&req)
	// if err != nil {
	// 	internal.Handle400(ctx, err)
	// 	return
	// }
	// // 数据库
	// var model db.WeekPlan
	// model.ID = id.ID
	// ok, err := db.Get(&model)
	// if err != nil {
	// 	internal.HandleDB500(ctx, err)
	// 	return
	// }
	// // 返回
	// if !ok {
	// 	internal.Handle404(ctx)
	// 	return
	// }
	// ctx.JSON(http.StatusOK, &model)
}
