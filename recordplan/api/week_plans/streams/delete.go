package streams

import (
	"net/http"
	"recordplan/api/internal"
	"recordplan/db"

	"github.com/gin-gonic/gin"
)

type putReq struct {
}

//	@Summary	解绑流
//	@Tags		周计划
//	@Param		id		path	string		true	"WeekPlan.ID"
//	@Param		data	body	[]string	true	"流标识数组"
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
	_, err = db.BatchDeleteWeekPlanStream(keys)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
