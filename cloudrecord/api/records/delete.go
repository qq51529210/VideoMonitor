package records

import (
	"cloudrecord/api/internal"
	"cloudrecord/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Summary	删除
//	@Tags		云端录像
//	@Param		id	path	string	true	"主键"
//	@Accept		json
//	@Success	204
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/records/{id} [delete]
func delete(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	_, err = db.RecordDA.Delete(id.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}

//	@Summary	批量删除
//	@Tags		云端录像
//	@Param		ids	body	[]string	true	"条件"
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Success	204
//	@Failure	400	{object}	internal.Error
//	@Failure	401
//	@Failure	403
//	@Failure	500	{object}	internal.Error
//	@Router		/records [delete]
func batchDelete(ctx *gin.Context) {
	// 参数
	var req internal.BatchDelete[string]
	err := ctx.ShouldBindJSON(&req.ID)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	_, err = db.RecordDA.BatchDelete(req.ID)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
