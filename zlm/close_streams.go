// package zlm

// import "net/url"

// // CloseStreamsReq 是 CloseStreams 参数
// type CloseStreamsReq struct {
// 	// 筛选协议，例如 rtsp或rtmp
// 	Schema string
// 	// 筛选虚拟主机，例如 __defaultVhost__
// 	VHost string
// 	// 筛选应用名，例如 live
// 	App string
// 	// 筛选流id，例如 test
// 	Stream string
// 	// 是否强制关闭(有人在观看是否还关闭)，0/1
// 	Force string
// }

// func (m *CloseStreamsReq) toQuery() url.Values {
// 	q := make(url.Values)
// 	if m.Schema != "" {
// 		q.Set("schema", m.Schema)
// 	}
// 	if m.VHost != "" {
// 		q.Set("vhost", m.VHost)
// 	}
// 	if m.App != "" {
// 		q.Set("app", m.App)
// 	}
// 	if m.Stream != "" {
// 		q.Set("stream", m.Stream)
// 	}
// 	if m.Force != "" {
// 		q.Set("force", m.Force)
// 	}
// 	return q
// }

// // closeStreamsRes 封装 CloseStreamsRes
// type closeStreamsRes struct {
// 	Code int `json:"code"`
// 	CloseStreamsRes
// }

// // CloseStreamsRes 是 closeStreams 返回值
// type CloseStreamsRes struct {
// 	Code int `json:"code"`
// 	// 筛选命中的流个数
// 	CountHit int `json:"count_hit"`
// 	// 被关闭的流个数，可能小于count_hit
// 	CountClosed int `json:"count_closed"`
// }

// // CloseStreams 调用 /index/api/close_streams
// // 关闭流(目前所有类型的流都支持关闭)
// func (s *Server) CloseStreams(req *CloseStreamsReq) (CloseStreamsRes, error) {
// 	query := make(url.Values)
// 	if req != nil {
// 		query = req.toQuery()
// 	}
// 	var res closeStreamsRes
// 	err := httpGet(s, s.url("close_streams"), query, &res)
// 	if err != nil {
// 		return CloseStreamsRes{}, err
// 	}
// 	if res.Code != 0 {
// 		return CloseStreamsRes{}, CodeError(res.Code)
// 	}
// 	return res.CloseStreamsRes, nil
// }
