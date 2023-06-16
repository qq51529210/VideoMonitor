package devices

// import (
// 	"gbs/api/internal"
// 	"gbs/api/internal/middleware"
// 	"gbs/db"
// 	"gbs/ovf"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // @Summary  删除
// // @Tags     ONVIF 设备
// // @Param    id path int64 true "数据库 ID"
// // @Security ApiKeyAuth
// // @Accept   json
// // @Produce  json
// // @Success  204 {object} internal.RowResult
// // @Failure  400 {object} internal.Error
// // @Failure  401
// // @Failure  403
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices/{id} [delete]
// func delete(ctx *gin.Context) {
// 	// 参数
// 	var id internal.IDPath[int64]
// 	err := ctx.ShouldBindUri(&id)
// 	if err != nil {
// 		internal.Handle400(ctx, err)
// 		return
// 	}
// 	// 数据库
// 	rows, err := db.DeleteOnvifDevice(id.ID)
// 	if err != nil {
// 		internal.HandleDB500(ctx, err)
// 		return
// 	}
// 	// 内存
// 	if rows > 0 {
// 		ovf.RemoveDevice(id.ID)
// 	}
// 	// 返回
// 	ctx.JSON(http.StatusNoContent, &internal.RowResult{
// 		Row: rows,
// 	})
// }

// // @Summary  批量删除
// // @Tags     ONVIF 设备
// // @Param    data body internal.BatchDelete[int64] true "条件"
// // @Security ApiKeyAuth
// // @Accept   json
// // @Produce  json
// // @Success  204 {object} internal.RowResult
// // @Failure  400 {object} internal.Error
// // @Failure  401
// // @Failure  403
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices [delete]
// func batchDelete(ctx *gin.Context) {
// 	// 参数
// 	var req internal.BatchDelete[int64]
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		internal.Handle400(ctx, err)
// 		return
// 	}
// 	ctx.Set(middleware.ReqDataCtxKey, &req)
// 	// 数据库
// 	rows, err := db.BatchDeleteOnvifDevice(req.ID)
// 	if err != nil {
// 		internal.HandleDB500(ctx, err)
// 		return
// 	}
// 	// 内存
// 	if rows > 0 {
// 		ovf.BatchRemoveDevice(req.ID)
// 	}
// 	// 返回
// 	ctx.JSON(http.StatusNoContent, &internal.RowResult{
// 		Row: rows,
// 	})
// }
