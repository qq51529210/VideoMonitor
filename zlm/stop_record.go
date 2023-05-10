package zlm

import "net/url"

// StopRecordReq 是 StopRecord 的参数
type StopRecordReq struct {
	// 0为hls，1为mp4
	Type string
	// 筛选虚拟主机
	VHost string
	// 筛选应用名，例如 live
	App string
	// 筛选流id，例如 test
	Stream string
}

func (m *StopRecordReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Type != "" {
		q.Set("type", m.Type)
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
	return q
}

// stopRecordRes 是 StopRecord 的返回值
type stopRecordRes struct {
	Code int `json:"code"`
	// 成功与否
	Result bool `json:"result"`
}

// StopRecord 调用 /index/api/stopRecord
// 停止录制流
// 返回是否成功
func (s *Server) StopRecord(req *StopRecordReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res stopRecordRes
	err := httpGet(s, s.url("stopRecord"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Result, nil
}
