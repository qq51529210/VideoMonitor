package devices

// import (
// 	"net/http"

// 	"gbs/api/internal"
// 	"gbs/db"

// 	"github.com/gin-gonic/gin"
// )

// // @Summary  列表
// // @Tags     ONVIF 设备
// // @Param    query query db.OnvifDeviceQuery false "条件"
// // @Security ApiKeyAuth
// // @Produce  json
// // @Success  200 {object} db.ListData[db.OnvifDevice]
// // @Failure  400 {object} internal.Error
// // @Failure  401
// // @Failure  403
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices [get]
// func list(ctx *gin.Context) {
// 	// 参数
// 	var req db.OnvifDeviceQuery
// 	err := ctx.ShouldBindQuery(&req)
// 	if err != nil {
// 		internal.Handle400(ctx, err)
// 		return
// 	}
// 	// 数据库
// 	var res db.ListData[db.OnvifDevice]
// 	err = db.List(&req, &req.Page, &res)
// 	if err != nil {
// 		internal.HandleDB500(ctx, err)
// 		return
// 	}
// 	// 返回
// 	ctx.JSON(http.StatusOK, &res)
// }
