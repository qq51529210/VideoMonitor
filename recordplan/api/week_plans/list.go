package weekplans

import (
	"net/http"

	"recordplan/api/internal"
	"recordplan/db"

	"github.com/gin-gonic/gin"
)

// @Summary	列表
// @Tags		周计划
// @Param		query	query	db.WeekPlanQuery	false	"条件"
// @Produce	json
// @Success	200	{object}	db.ListData[db.WeekPlan]
// @Failure	400	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/week_plans [get]
func list(ctx *gin.Context) {
	// 参数
	var req db.WeekPlanQuery
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var res db.ListData[db.WeekPlan]
	err = db.List(&req, &req.Page, &res)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusOK, &res)
}
