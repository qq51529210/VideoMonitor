package zlm

import "net/url"

// ResumeRTPCheckReq 是 ResumeRTPCheck 的参数
type ResumeRTPCheckReq struct {
	// 该端口绑定的流id
	StreamID string
}

func (m *ResumeRTPCheckReq) toQuery() url.Values {
	q := make(url.Values)
	if m.StreamID != "" {
		q.Set("stream_id", m.StreamID)
	}
	return q
}

// resumeRTPCheckRes 是 ResumeRTPCheck 的返回值
type resumeRTPCheckRes struct {
	Code int          `json:"code"`
	Data []*MediaInfo `json:"data"`
}

// ResumeRTPCheck 调用 /index/api/getMediaList
// 获取流列表，可选筛选参数
// 返回媒体信息列表
func (s *Server) ResumeRTPCheck(req *ResumeRTPCheckReq) ([]*MediaInfo, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res resumeRTPCheckRes
	err := httpGet(s, s.url("getMediaList"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
