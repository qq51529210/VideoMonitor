package zlm

import "net/url"

// IsRecordingReq 是 IsRecording 的参数
type IsRecordingReq struct {
	// 0为hls，1为mp4
	Type string
	// 筛选虚拟主机
	VHost string
	// 筛选应用名，例如 live
	App string
	// 筛选流id，例如 test
	Stream string
}

func (m *IsRecordingReq) toQuery() url.Values {
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

// isRecordingRes 是 IsRecording 的返回值
type isRecordingRes struct {
	Code int `json:"code"`
	// false:未录制,true:正在录制
	Status bool `json:"status"`
}

// IsRecording 调用 /index/api/isRecording
// 获取流录制状态
// 返回状态
func (s *Server) IsRecording(req *IsRecordingReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res isRecordingRes
	err := httpGet(s, s.url("isRecording"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Status, nil
}
