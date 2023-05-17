package weekplans

import (
	"net/http"

	"recordplan/api/internal"
	"recordplan/db"

	"github.com/gin-gonic/gin"
)

// @Summary	详情
// @Tags		周计划
// @Param		id	path	string	true	"WeekPlan.ID"
// @Produce	json
// @Success	200	{object}	db.WeekPlan
// @Failure	404	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/week_plans/{id} [get]
func get(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	model, err := db.GetWeekPlan(id.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	if model == nil {
		internal.Handle404(ctx)
		return
	}
	ctx.JSON(http.StatusOK, model)
}
