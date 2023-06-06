package weekplans

import (
	"net/http"

	"recordplan/api/internal"
	"recordplan/db"
	"recordplan/task/week"

	"github.com/gin-gonic/gin"
)

type weekplan struct {
	*db.WeekPlan
	IsRecording bool `json:"isRecording"`
}

//	@Summary	获取周计划
//	@Tags		媒体流
//	@Param		stream	path	string	true	"Stream"
//	@Produce	json
//	@Success	200	{object}	[]weekplan
//	@Failure	404	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/streams/{stream}/week_plans [get]
func list(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	weekPlanTaskModels, err := db.GetWeekPlanStreamByStream(id.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	weekPlanIDs := make([]string, len(weekPlanTaskModels))
	for i := 0; i < len(weekPlanIDs); i++ {
		weekPlanIDs[i] = weekPlanTaskModels[i].WeekPlanID
	}
	weekPlanModels, err := db.GetWeekPlanIn(weekPlanIDs)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	weekplans := make([]weekplan, len(weekPlanModels))
	for i := 0; i < len(weekplans); i++ {
		weekplans[i].WeekPlan = weekPlanModels[i]
		weekplans[i].IsRecording, _ = week.IsRecording(weekPlanModels[i].ID)
	}
	// 返回
	ctx.JSON(http.StatusOK, weekPlanModels)
}
