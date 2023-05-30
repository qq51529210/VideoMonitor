package records

import (
	"cloudrecord/api/internal"
	"cloudrecord/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

type postReq struct {
	// minio 的标识
	Name string `json:"name" binding:"required,max=40"`
	// stream
	Stream string `json:"stream" binding:"required,max=64"`
	// 大小
	Size int64 `json:"size" binding:"required,min=1"`
	// 时长
	Duration float64 ` json:"duration" binding:"required,min=1"`
	// 创建时间
	Time int64 `json:"time" binding:"required,min=1"`
	// 删除时间
	DeleteTime int64 `json:"deleteTime" binding:"required,min=0"`
}

//	@Summary	添加
//	@Tags		云端录像
//	@Param		data	body	postReq	true	"数据"
//	@Accept		json
//	@Produce	json
//	@Success	201
//	@Failure	400	{object}	internal.Error
//	@Failure	500	{object}	internal.Error
//	@Router		/records [post]
func post(ctx *gin.Context) {
	// 参数
	var req postReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var model db.Record
	util.CopyStruct(&model, &req)
	_, err = db.RecordDA.Add(&model)
	if err != nil {
		internal.HandleDB500(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusCreated)
}
