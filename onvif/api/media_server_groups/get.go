package mediaservergroups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/zlm"
)

// @Summary  详情
// @Tags     流媒体服务分组
// @Param    id path int64 true "id"
// @Produce  json
// @Success  200 {object} zlm.MediaServerGroup
// @Failure  404 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups/{id} [get]
func get(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[int64]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	model, err := zlm.MediaServerGroupDA.Get(id.ID)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	if model == nil {
		internal.WriteErrorDataNotFound(ctx)
		return
	}
	ctx.JSON(http.StatusOK, &model)
}
