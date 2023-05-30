package tasks

import (
	"net/http"
	"recordplan/api/internal"
	"recordplan/db"
	"recordplan/task/week"

	"github.com/gin-gonic/gin"
)

// @Summary	添加任务
// @Tags		周计划
// @Param		id		path	string	true	"WeekPlan.ID"
// @Param		data	body	[]tasks	true	"任务数组"
// @Accept		json
// @Produce	json
// @Success	201	{object}	internal.IDResult[int64]
// @Failure	400	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/week_plans/{id}/tasks [put]
func put(ctx *gin.Context) {
	// 参数
	var weekplanID internal.IDPath[string]
	err := ctx.ShouldBindUri(&weekplanID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	var req []*tasks
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据
	var models []*db.WeekPlanTask
	for _, r := range req {
		model := new(db.WeekPlanTask)
		model.WeekPlanID = weekplanID.ID
		model.TaskID = r.ID
		model.StartCallback = r.StartCallback
		model.StopCallback = r.StopCallback
		models = append(models, model)
	}
	// 添加
	rows, err := db.WeekPlanTaskDA.BatchSave(models)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if rows > 0 {
		week.Reload(weekplanID.ID)
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.RowResult{
		Row: rows,
	})
}
