package devices

type postReq struct {
	// 名称
	Name *string `json:"name" binding:"required,max=32"`
	// 地址，ip:port
	IPAddress *string `json:"ipAddress" binding:"required,ip_addr"`
	// 账号
	Username *string `json:"username" binding:"omitempty,max=32"`
	// 密码
	Password *string `json:"password" binding:"omitempty,max=64"`
	// MediaServerGroup.ID ，将这个设备绑定到一个流媒体服务组，默认 1
	MediaServerGroupID *int64 `json:"mediaServerGroupID" binding:"omitempty,min=1"`
}

// // @Summary  添加
// // @Tags     ONVIF 设备
// // @Param    data body postReq true "数据"
// // @Accept   json
// // @Produce  json
// // @Success  201 {object} internal.IDResult[int64]
// // @Failure  400 {object} internal.Error
// // @Failure  500 {object} internal.Error
// // @Router   /onvif_devices [post]
// func post(ctx *gin.Context) {
// 	// // 参数
// 	// var req postReq
// 	// err := ctx.ShouldBindJSON(&req)
// 	// if err != nil {
// 	// 	internal.Handle400(ctx, err)
// 	// 	return
// 	// }
// 	// ctx.Set(middleware.ReqDataCtxKey, &req)
// 	// var model db.OnvifDevice
// 	// util.CopyStruct(&model, &req, false)
// 	// // 先查询
// 	// err = ovf.QueryDevice(&model)
// 	// if err != nil {
// 	// 	internal.HandleDB500(ctx, err)
// 	// 	return
// 	// }
// 	// for _, ch := range model.Channel {
// 	// 	ch.KeepStream = req.KeepStream
// 	// }
// 	// // 数据库
// 	// rows, err := db.AddOnvifDevice(&model)
// 	// if err != nil {
// 	// 	internal.HandleDB500(ctx, err)
// 	// 	return
// 	// }
// 	// // 内存
// 	// if rows > 0 {
// 	// 	ovf.LoadDevice(model.ID)
// 	// }
// 	// // 返回
// 	// ctx.JSON(http.StatusCreated, &internal.IDResult[int64]{
// 	// 	ID: model.ID,
// 	// })
// }
