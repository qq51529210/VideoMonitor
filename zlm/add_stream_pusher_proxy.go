package zlm

import "net/url"

// AddPushStreamerProxyReq 是 AddPushStreamerProxy 参数
type AddPushStreamerProxyReq struct {
	// 推流协议，支持rtsp、rtmp，大小写敏感
	Schema string
	// 已注册流的虚拟主机，一般为__defaultVhost__
	VHost string
	// 已注册流的应用名，例如live
	App string
	// 已注册流的id名，例如test
	Stream string
	// 推流地址，需要与schema字段协议一致
	DstURL string
	// rtsp推流时，推流方式，0：tcp，1：udp
	RTPType string
	// 拉流超时时间，单位秒，float类型
	TimeoutSec string
	// 推流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string
}

func (m *AddPushStreamerProxyReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Schema != "" {
		q.Set("schema", m.Schema)
	}
	if m.VHost != "" {
		q.Set("vhost", m.VHost)
	}
	if m.App != "" {
		q.Set("app", m.App)
	}
	if m.Stream != "" {
		q.Set("stream", m.Stream)
	}
	if m.DstURL != "" {
		q.Set("dst_url", m.DstURL)
	}
	if m.RTPType != "" {
		q.Set("rtp_type", m.RTPType)
	}
	if m.TimeoutSec != "" {
		q.Set("timeout_sec	", m.TimeoutSec)
	}
	if m.RetryCount != "" {
		q.Set("retry_count", m.RetryCount)
	}
	return q
}

// addPushStreamerProxyRes 是 AddPushStreamerProxy 返回值
type addPushStreamerProxyRes struct {
	Code int                         `json:"code"`
	Data AddPushStreamerProxyResData `json:"data"`
}

// AddPushStreamerProxyResData 是 addPushStreamerProxyRes 的 Data 字段
type AddPushStreamerProxyResData struct {
	// 流的唯一标识
	Key string
}

// AddPushStreamerProxy 调用 /index/api/addPushStreamerProxy
// 添加rtsp/rtmp主动推流(把本服务器的直播流推送到其他服务器去)
// 返回 key
func (s *Server) AddPushStreamerProxy(req *AddPushStreamerProxyReq) (string, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res addPushStreamerProxyRes
	err := httpGet(s, s.url("addPushStreamerProxy"), query, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", CodeError(res.Code)
	}
	return res.Data.Key, nil
}
