package mediaservergroups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/zlm"
)

// @Summary  删除
// @Tags     流媒体服务分组
// @Param    id path int64 true "id"
// @Success  204
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups/{id} [delete]
func delete(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[int64]
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	_, err = zlm.MediaServerGroupDA.Delete(id.ID)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}

// @Summary  批量删除
// @Tags     流媒体服务分组
// @Param    data body internal.BatchDelete[int64] true "id 数组"
// @Accept   json
// @Success  204
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups [delete]
func batchDelete(ctx *gin.Context) {
	// 参数
	var req []int64
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	_, err = zlm.MediaServerGroupDA.BatchDelete(req)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.Status(http.StatusNoContent)
}
