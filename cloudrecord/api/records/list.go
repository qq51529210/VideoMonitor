package records

import (
	"cloudrecord/api/internal"
	"cloudrecord/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

//	@Summary	列表
//	@Tags		云端录像
//	@Param		query	query	db.RecordQuery	false	"条件"
//	@Produce	json
//	@Success	200	{object}	util.GORMList[db.Record]
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/records [get]
func list(ctx *gin.Context) {
	// 参数
	var req db.RecordQuery
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var res util.GORMList[*db.Record]
	err = db.RecordDA.List(&req.GORMPage, &req, &res)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusOK, &res)
}
