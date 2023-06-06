package streams

import (
	"net/http"
	"recordplan/api/internal"
	"recordplan/db"
	"recordplan/task/week"

	"github.com/gin-gonic/gin"
)

//	@Summary	解绑媒体流
//	@Tags		周计划
//	@Param		id		path	string		true	"WeekPlan.ID"
//	@Param		data	body	[]string	true	"自定义的媒体流 ID 数组"
//	@Accept		json
//	@Success	204
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/week_plans/{id}/streams [delete]
func delete(ctx *gin.Context) {
	// 参数
	var weekplanID internal.IDPath[string]
	err := ctx.ShouldBindUri(&weekplanID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	var req []string
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据
	keys := make([]db.WeekPlanStreamKey, len(req))
	for i := 0; i < len(keys); i++ {
		keys[i].WeekPlanID = weekplanID.ID
		keys[i].Stream = req[i]
	}
	// 删除
	rows, err := db.WeekPlanStreamDA.BatchDelete(keys)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 更新
	if rows > 0 {
		week.Reload(weekplanID.ID)
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
