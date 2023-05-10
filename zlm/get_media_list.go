package zlm

import "net/url"

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
	// 筛选协议，例如 rtsp或rtmp
	Schema string
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string
	// 筛选应用名，例如 live
	App string
	// 筛选流id，例如 test
	Stream string
}

func (m *GetMediaListReq) toQuery() url.Values {
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
	return q
}

// getMediaListRes 是 GetMediaList 的返回值
type getMediaListRes struct {
	Code int          `json:"code"`
	Data []*MediaInfo `json:"data"`
}

// GetMediaList 调用 /index/api/getMediaList
// 获取流列表，可选筛选参数
func (s *Server) GetMediaList(req *GetMediaListReq) ([]*MediaInfo, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res getMediaListRes
	err := httpGet(s, s.url("getMediaList"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}

	return res.Data, nil
}
