package devices

// import (
// 	"gbs/api/internal"
// 	"gbs/db"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // @Summary  详情
// // @Tags     ONVIF 设备
// // @Param    id path int64 true "数据库 ID"
// // @Security ApiKeyAuth
// // @Produce  json
// // @Success  200 {object} db.OnvifDevice
// // @Failure  401
// // @Failure  403
// // @Failure  404 {object} internal.Error
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices/{id} [get]
// func get(ctx *gin.Context) {
// 	// 参数
// 	var id internal.IDPath[int64]
// 	err := ctx.ShouldBindUri(&id)
// 	if err != nil {
// 		internal.Handle400(ctx, err)
// 		return
// 	}
// 	// 数据库
// 	model, err := db.GetOnvifDevice(id.ID)
// 	if err != nil {
// 		internal.HandleDB500(ctx, err)
// 		return
// 	}
// 	// 返回
// 	if model == nil {
// 		internal.HandleData404(ctx)
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, model)
// }
