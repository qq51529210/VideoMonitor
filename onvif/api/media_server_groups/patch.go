package mediaservergroups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/onvif/db"
)

type patchReq struct {
	// 名称
	Name *string `json:"name" binding:"omitempty,max=32"`
	// 描述
	Describe *string `json:"describe" binding:"omitempty,max=128"`
}

// @Summary  修改
// @Tags     流媒体服务分组
// @Param    id path int64 true "id"
// @Param    data body patchReq true  "数据"
// @Accept   json
// @Produce  json
// @Success  201 {object} internal.RowResult
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups/{id} [patch]
func patch(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[int64]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	var req patchReq
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 是否为空
	if util.IsNilOrEmpty(&req) {
		internal.WriteErrorEmptySubmittedData(ctx)
		return
	}
	// 数据库
	var model db.MediaServerGroup
	util.CopyStruct(&model, &req)
	rows, err := db.MediaServerGroupDA.Update(&model)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.RowResult{
		Row: rows,
	})
}
