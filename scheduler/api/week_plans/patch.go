package weekplans

import (
	"net/http"

	"scheduler/api/internal"
	"scheduler/db"
	"scheduler/week"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

type patchReq struct {
	// 名称
	Name *string `json:"name" binding:"omitempty,max=32"`
	// 是否禁用
	// 0: 禁用
	// 1: 启用
	Enable *int8 `json:"enable" binding:"omitempty,oneof=0 1"`
	// 一周的时间数组
	Peroids [][]*db.TimePeroid `json:"peroids" binding:"omitempty,min=1,max=7,dive,required,min=1,dive"`
	// 附加的数据
	Data *string `json:"data" binding:"omitempty"`
}

//	@Summary	修改
//	@Tags		周计划
//	@Param		id		path	string		true	"WeekPlan.ID"
//	@Param		data	body	patchReq	true	"数据"
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	internal.RowResult
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/week_plans/{id} [patch]
func patch(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	var req patchReq
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 参数是空的
	if util.IsNilOrEmpty(&req) {
		// 返回
		ctx.JSON(http.StatusCreated, &internal.RowResult{
			Row: 0,
		})
		return
	}
	// 数据库
	var model db.WeekPlan
	util.CopyStruct(&model, &req)
	model.ID = id.ID
	if len(req.Peroids) > 0 {
		model.Peroids = jsonTimePeroid(req.Peroids)
	}
	rows, err := db.UpdateWeekPlan(&model)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if rows > 0 {
		week.Reload(model.ID)
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.RowResult{
		Row: rows,
	})
}
