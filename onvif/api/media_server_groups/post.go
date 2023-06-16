package mediaservergroups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/onvif/db"
)

type postReq struct {
	// 名称
	Name *string `json:"name" binding:"required,max=32"`
	// 描述
	Describe *string `json:"describe" binding:"omitempty,max=128"`
}

// @Summary  添加
// @Tags     流媒体服务分组
// @Param    data body postReq true "数据"
// @Accept   json
// @Produce  json
// @Success  201 {object} internal.IDResult[int64]
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups [post]
func post(ctx *gin.Context) {
	// 参数
	var req postReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	var model db.MediaServerGroup
	util.CopyStruct(&model, &req)
	_, err = db.MediaServerGroupDA.Add(&model)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.IDResult[int64]{
		ID: model.ID,
	})
}
