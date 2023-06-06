package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// AddPushStreamerProxyReq 是 AddPushStreamerProxy 参数
type AddPushStreamerProxyReq struct {
	// 流虚拟主机
	VHost string `query:"vhost"`
	// 推流协议，支持rtsp、rtmp，大小写敏感
	Schema string `query:"schema"`
	// 流应用名
	App string `query:"app"`
	// 流ID
	Stream string `query:"stream"`
	// 推流地址，需要与schema字段协议一致
	DstURL string `query:"stream"`
	// rtsp推流时，推流方式，0：tcp，1：udp
	RTPType string `query:"rtp_type"`
	// 拉流超时时间，单位秒，float类型
	TimeoutSec string `query:"timeout_sec"`
	// 推流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string `query:"retry_count"`
}

// addPushStreamerProxyRes 用于解析 addPushStreamerProxy 的返回值
type addPushStreamerProxyRes struct {
	Code int `json:"code"`
	Data struct {
		// 流的唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

// AddPushStreamerProxy 调用 /index/api/addPushStreamerProxy
// 添加rtsp/rtmp主动推流(把本服务器的直播流推送到其他服务器去)
// 返回 key
func (s *Server) AddPushStreamerProxy(req *AddPushStreamerProxyReq) (string, error) {
	var _res addPushStreamerProxyRes
	err := util.HTTP[any](http.MethodGet,
		s.url("addPushStreamerProxy"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return "", err
	}
	if _res.Code != 0 {
		return "", CodeError(_res.Code)
	}
	return _res.Data.Key, nil
}
