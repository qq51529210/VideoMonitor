package zlm

import "net/url"

// StopSendRTPReq 是 StopSendRTP 参数
type StopSendRTPReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string
	// 添加的应用名，例如 live
	App string
	// 添加的流id，例如 test
	Stream string
	// 根据ssrc关停某路rtp推流，不传时关闭所有推流
	SSRC string
}

func (m *StopSendRTPReq) toQuery() url.Values {
	q := make(url.Values)
	if m.VHost != "" {
		q.Set("vhost", m.VHost)
	}
	if m.App != "" {
		q.Set("app", m.App)
	}
	if m.Stream != "" {
		q.Set("stream", m.Stream)
	}
	if m.SSRC != "" {
		q.Set("ssrc", m.SSRC)
	}
	return q
}

// stopSendRTPRes 是 StopSendRTP 返回值
type stopSendRTPRes struct {
	Code int `json:"code"`
}

// StopSendRTP 调用 /index/api/stopSendRtp
// 停止 GB28181 ps-rtp 推流。
func (s *Server) StopSendRTP(req *StopSendRTPReq) error {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res stopSendRTPRes
	err := httpGet(s, s.url("stopSendRtp"), query, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return CodeError(res.Code)
	}
	return nil
}
