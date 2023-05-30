package streams

import (
	"net/http"
	"recordplan/api/internal"
	"recordplan/db"
	"recordplan/task/week"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

//	@Summary	绑定媒体流
//	@Tags		周计划
//	@Param		id		path	string		true	"WeekPlan.ID"
//	@Param		data	body	[]stream	true	"媒体流数组"
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	internal.IDResult[int64]
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/week_plans/{id}/streams [put]
func put(ctx *gin.Context) {
	// 参数
	var weekplanID internal.IDPath[string]
	err := ctx.ShouldBindUri(&weekplanID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	var req []*stream
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据
	var models []*db.WeekPlanStream
	for _, r := range req {
		model := new(db.WeekPlanStream)
		util.CopyStruct(model, r)
		model.WeekPlanID = weekplanID.ID
		model.Stream = r.Stream
		model.StartCallback = r.StartCallback
		model.StopCallback = r.StopCallback
		models = append(models, model)
	}
	// 添加
	rows, err := db.WeekPlanStreamDA.BatchSave(models)
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
