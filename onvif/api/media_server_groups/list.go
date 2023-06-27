package mediaservergroups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/zlm"
)

// @Summary  列表
// @Tags     流媒体服务分组
// @Param    query query zlm.MediaServerGroupQuery false "条件"
// @Produce  json
// @Success  200 {object} util.GORMList[zlm.MediaServerGroup]
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_server_groups [get]
func list(ctx *gin.Context) {
	// 参数
	var req zlm.MediaServerGroupQuery
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		internal.WriteErrorSubmittedData(ctx, err)
		return
	}
	// 数据库
	var res util.GORMList[*zlm.MediaServerGroup]
	err = zlm.MediaServerGroupDA.List(&req.GORMPage, &req, &res)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusOK, &res)
}
