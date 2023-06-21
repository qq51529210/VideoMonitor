package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// AddPushStreamerProxyReq 是 AddPushStreamerProxy 参数
type AddPushStreamerProxyReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 协议，例如 rtsp 或 rtmp
	Schema string `query:"schema"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 推流地址
	DstURL string `query:"dst_url"`
	// rtsp 推流方式，0：tcp ，1：udp
	RTPType string `query:"rtp_type"`
	// 推流超时时间，单位秒，float 类型
	TimeoutSec string `query:"timeout_sec"`
	// 推流重试次数，不传或 0 ，则无限重试
	RetryCount string `query:"retry_count"`
}

// addPushStreamerProxyRes 用于解析 addPushStreamerProxy 的返回值
type addPushStreamerProxyRes struct {
	apiRes
	Data struct {
		// 流的唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

const (
	apiPathAddPushStreamerProxy = "addPushStreamerProxy"
)

// AddPushStreamerProxy 调用 /index/api/addPushStreamerProxy 主动推流 rtsp / rtmp ，返回 key
func (s *Server) AddPushStreamerProxy(req *AddPushStreamerProxyReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.AddPushStreamerProxyWithContext(ctx, req)
}

// AddPushStreamerProxyWithContext 调用 /index/api/addStreamPusherProxy 主动推流 rtsp / rtmp ，返回 key
func (s *Server) AddPushStreamerProxyWithContext(ctx context.Context, req *AddPushStreamerProxyReq) (string, error) {
	var res addPushStreamerProxyRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathAddPushStreamerProxy),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathAddPushStreamerProxy,
		}
	}
	return res.Data.Key, nil
}
