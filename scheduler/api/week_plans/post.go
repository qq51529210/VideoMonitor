package weekplans

import (
	"net/http"
	"scheduler/api/internal"
	"scheduler/db"
	"scheduler/week"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
	"github.com/qq51529210/uuid"
)

type postReq struct {
	// 名称
	Name *string `json:"name" binding:"required,max=32"`
	// 是否禁用，默认是 1
	// 0: 禁用
	// 1: 启用
	Enable *int8 `json:"enable" binding:"omitempty,oneof=0 1"`
	// 一周的时间数组
	Peroids [][]*db.TimePeroid `json:"peroids" binding:"required,min=1,max=7,dive,required,min=1,dive"`
	// 附加的数据
	Data *string `json:"data" binding:"omitempty"`
}

// @Summary	添加
// @Tags		周计划
// @Param		data	body	postReq	true	"数据"
// @Accept		json
// @Produce	json
// @Success	201	{object}	internal.IDResult[int64]
// @Failure	400	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/week_plans [post]
func post(ctx *gin.Context) {
	// 参数
	var req postReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var model db.WeekPlan
	util.CopyStruct(&model, &req)
	model.ID = uuid.LowerV1WithoutHyphen()
	model.Peroids = jsonTimePeroid(req.Peroids)
	row, err := db.AddWeekPlan(&model)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if row > 0 {
		week.Reload(model.ID)
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.IDResult[string]{
		ID: model.ID,
	})
}
