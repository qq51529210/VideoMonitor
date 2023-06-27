package mediaservers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
	"github.com/qq51529210/video-monitor/onvif/api/internal"
	"github.com/qq51529210/video-monitor/zlm"
)

type patchReq struct {
	// 访问密钥
	Secret *string `json:"secret" binding:"omitempty,max=64"`
	// 名称
	Name *string `json:"name" binding:"omitempty,max=32"`
	// 描述
	Describe *string `json:"describe" binding:"omitempty,max=128"`
	// API 地址 (http|https)://ip:port
	APIBaseURL *string `json:"apiBaseURL" binding:"omitempty,http_url"`
	// 外网访问的 ip ，生成播放地址时使用，默认使用 apiBaseURL 中的 host
	PublicIP *string `json:"publicIP" binding:"omitempty,ip_addr"`
	// 内网访问的 ip ，生成播放地址时使用，默认使用 apiBaseURL 中的 host
	PrivateIP *string `json:"privateIP" binding:"omitempty,ip_addr"`
	// 请求服务接口超时时间，单位，毫秒，默认 5000
	APICallTimeout *uint32 `json:"apiCallTimeout" binding:"omitempty,min=1000"`
	// 是否禁用，0/1 ，默认 1
	Enable *int8 `json:"enable" binding:"omitempty,oneof=0 1"`
	// 所属的分组 ，默认 1
	MediaServerGroupID *int64 `json:"mediaServerGroupID" binding:"omitempty,min=1"`
}

// @Summary  修改
// @Tags     流媒体服务
// @Param    id path string true "id"
// @Param    data body patchReq true  "数据"
// @Accept   json
// @Produce  json
// @Success  201 {object} internal.RowResult
// @Failure  400 {object} internal.Error
// @Failure  500 {object} internal.Error
// @Router   /media_servers/{id} [patch]
func patch(ctx *gin.Context) {
	// 参数
	var id internal.IDPath[string]
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
	var model zlm.MediaServer
	util.CopyStruct(&model, &req)
	model.ID = id.ID
	rows, err := zlm.MediaServerDA.Update(&model)
	if err != nil {
		internal.WriteErrorDataBaseAccess(ctx, err)
		return
	}
	// 返回
	ctx.JSON(http.StatusCreated, &internal.RowResult{
		Row: rows,
	})
}
