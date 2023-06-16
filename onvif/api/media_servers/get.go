package mediaservers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/onvif/db"
)

// @Summary  详情
// @Tags     流媒体服务
// @Param    id path string true "id"
// @Produce  json
// @Success  200 {object} db.MediaServer
// @Failure  404 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_servers/{id} [get]
func get(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	model, err := db.MediaServerDA.Get(id.ID)
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
