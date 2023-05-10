// package zlm

// import (
// 	"net/url"
// )

// // CloseRTPServerReq 是 CloseRTPServer 的参数
// type CloseRTPServerReq struct {
// 	// 调用closeRtpServer接口时提供的流ID
// 	StreamID string
// }

// func (m *CloseRTPServerReq) toQuery() url.Values {
// 	q := make(url.Values)
// 	if m.StreamID != "" {
// 		q.Set("stream_id", m.StreamID)
// 	}
// 	return q
// }

// // closeRTPServerRes 是 CloseRTPServer 的返回值
// type closeRTPServerRes struct {
// 	Code int `json:"code"`
// 	// 是否找到记录并关闭
// 	Hit int `json:"hit"`
// }

// // CloseRTPServer 调用 /index/api/closeRtpServer
// // 关闭GB28181 RTP接收端口
// // 返回成功的个数
// func (s *Server) CloseRTPServer(req *CloseRTPServerReq) (int, error) {
// 	query := make(url.Values)
// 	if req != nil {
// 		query = req.toQuery()
// 	}
// 	var res closeRTPServerRes
// 	err := httpGet(s, s.url("closeRtpServer"), query, &res)
// 	if err != nil {
// 		return -1, err
// 	}
// 	if res.Code != 0 {
// 		return -1, CodeError(res.Code)
// 	}
// 	return res.Hit, nil
// }
