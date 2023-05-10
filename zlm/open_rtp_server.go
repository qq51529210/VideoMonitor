package zlm

import (
	"net/url"
)

// OpenRTPServerReq 是 OpenRTPServer 的参数
type OpenRTPServerReq struct {
	// 接收端口，0则为随机端口
	Port string
	// 创建 udp端口时是否同时监听tcp端口
	EnableTCP string
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	StreamID string
	// 是否重用端口，默认为0，非必选参数，0/1
	ReusePort string
	// 是否指定收流的rtp ssrc, 十进制数字，不指定或指定0时则不过滤rtp，非必选参数
	SSRC string
}

func (m *OpenRTPServerReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Port != "" {
		q.Set("port", m.Port)
	}
	if m.EnableTCP != "" {
		q.Set("enable_tcp", m.EnableTCP)
	}
	if m.StreamID != "" {
		q.Set("stream_id", m.StreamID)
	}
	if m.ReusePort != "" {
		q.Set("re_use_port", m.ReusePort)
	}
	if m.SSRC != "" {
		q.Set("ssrc", m.SSRC)
	}
	return q
}

// openRTPServerRes 是 OpenRTPServer 的返回值
type openRTPServerRes struct {
	Code int `json:"code"`
	// 接收端口，方便获取随机端口号
	Port int `json:"port"`
}

// OpenRTPServer 调用 /index/api/openRtpServer
// 创建GB28181 RTP接收端口，如果该端口接收数据超时，则会自动被回收(不用调用closeRtpServer接口)
// 返回使用的端口
func (s *Server) OpenRTPServer(req *OpenRTPServerReq) (int, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res openRTPServerRes
	err := httpGet(s, s.url("openRtpServer"), query, &res)
	if err != nil {
		return -1, err
	}
	if res.Code != 0 {
		return -1, CodeError(res.Code)
	}
	return res.Port, nil
}
