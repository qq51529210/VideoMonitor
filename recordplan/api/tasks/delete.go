package tasks

import (
	"recordplan/api/internal"

	"github.com/gin-gonic/gin"
)

// @Summary	删除
// @Tags		任务
// @Param		id		path	string		true	"Task ID"
// @Success	204
// @Failure	400	{object}	internal.Error
// @Failure	500	{object}	internal.Error
// @Router		/tasks/{id} [delete]
func delete(ctx *gin.Context) {
	// 参数
	var weekplanID internal.IDPath[string]
	err := ctx.ShouldBindUri(&weekplanID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// var req []string
	// err = ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	internal.Handle400(ctx, err)
	// 	return
	// }
	// // 数据
	// keys := make([]db.WeekPlanTaskKey, len(req))
	// for i := 0; i < len(keys); i++ {
	// 	keys[i].WeekPlanID = weekplanID.ID
	// 	keys[i].TaskID = req[i]
	// }
	// // 删除
	// rows, err := db.DeleteWeekPlanTaskDA.Batch(keys)
	// if err != nil {
	// 	internal.HandleDB500(ctx, err)
	// 	return
	// }
	// // 更新
	// if rows > 0 {
	// 	week.Reload(weekplanID.ID)
	// }
	// // 返回
	// ctx.Status(http.StatusNoContent)
}
