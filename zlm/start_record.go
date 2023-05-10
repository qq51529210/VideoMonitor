package zlm

import "net/url"

// StartRecordReq 是 StartRecord 的参数
type StartRecordReq struct {
	// 0为hls，1为mp4
	Type string
	// 筛选虚拟主机
	VHost string
	// 筛选应用名，例如 live
	App string
	// 筛选流id，例如 test
	Stream string
	// 录像保存目录
	CustomizedPath string
	// mp4录像切片时间大小,单位秒，置0则采用配置项
	MaxSecond string
}

func (m *StartRecordReq) toQuery() url.Values {
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
	if m.CustomizedPath != "" {
		q.Set("customized_path", m.CustomizedPath)
	}
	if m.MaxSecond != "" {
		q.Set("max_second", m.MaxSecond)
	}
	return q
}

// startRecordRes 是 StartRecord 的返回值
type startRecordRes struct {
	Code int `json:"code"`
	// 成功与否
	Result bool `json:"result"`
}

// StartRecord 调用 /index/api/startRecord
// 开始录制hls或MP4
// 返回是否成功
func (s *Server) StartRecord(req *StartRecordReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res startRecordRes
	err := httpGet(s, s.url("startRecord"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Result, nil
}
