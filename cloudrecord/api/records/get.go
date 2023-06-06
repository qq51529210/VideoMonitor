package records

import (
	"net/http"

	"github.com/qq51529210/video-monitor/cloudrecord/api/internal"
	"github.com/qq51529210/video-monitor/cloudrecord/db"

	"github.com/gin-gonic/gin"
)

// @Summary	详情
// @Tags		云端录像
// @Param		name	path	string	true	"Record.Name"
// @Produce	json
// @Success	200	{object}	db.Record
// @Failure	404	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/records/{id} [get]
func get(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var model db.Record
	ok, err := db.RecordDA.Get(&model)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	if !ok {
		internal.Handle404(ctx)
		return
	}
	ctx.JSON(http.StatusOK, model)
}
