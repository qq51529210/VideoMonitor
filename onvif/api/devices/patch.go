package devices

// import (
// 	"net/http"

// 	"gbs/api/internal"
// 	"gbs/api/internal/middleware"
// 	"gbs/db"
// 	"gbs/ovf"
// 	"gbs/util"

// 	"github.com/gin-gonic/gin"
// )

// type patchReq struct {
// 	// 数据库 ID
// 	ID int64 `json:"id" binding:"required,min=1"`
// 	// 名称
// 	Name *string `json:"name" binding:"omitempty,max=32"`
// 	// 地址，ip:port
// 	IPAddress *string `json:"ipAddress" binding:"omitempty,ipAddr"`
// 	// 账号
// 	Username *string `json:"username" binding:"omitempty,max=32"`
// 	// 密码
// 	Password *string `json:"password" binding:"omitempty,max=64"`
// }

// // @Summary  修改
// // @Tags     ONVIF 设备
// // @Param    data body patchReq true "数据"
// // @Security ApiKeyAuth
// // @Accept   json
// // @Produce  json
// // @Success  201 {object} internal.RowResult
// // @Failure  400 {object} internal.Error
// // @Failure  401
// // @Failure  403
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices [patch]
// func patch(ctx *gin.Context) {
// 	// 参数
// 	var req patchReq
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		internal.Handle400(ctx, err)
// 		return
// 	}
// 	ctx.Set(middleware.ReqDataCtxKey, &req)
// 	// 数据库
// 	var model db.OnvifDevice
// 	util.CopyStruct(&model, &req, false)
// 	rows, err := db.UpdateOnvifDevice(&model)
// 	if err != nil {
// 		internal.HandleDB500(ctx, err)
// 		return
// 	}
// 	// 内存
// 	if rows > 0 {
// 		ovf.LoadDevice(model.ID)
// 	}
// 	// 返回
// 	ctx.JSON(http.StatusCreated, &internal.RowResult{
// 		Row: rows,
// 	})
// }
