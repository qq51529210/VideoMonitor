package mediaservers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/zlm"
)

// @Summary  删除
// @Tags     流媒体服务
// @Param    id path string true "id"
// @Success  204
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_servers/{id} [delete]
func delete(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	_, err = zlm.MediaServerDA.Delete(id.ID)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}

// @Summary  批量删除
// @Tags     流媒体服务
// @Param    data body internal.BatchDelete[string] true "id 数组"
// @Accept   json
// @Success  204
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_servers [delete]
func batchDelete(ctx *gin.Context) {
	// 参数
	var req []string
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	_, err = zlm.MediaServerDA.BatchDelete(req)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
