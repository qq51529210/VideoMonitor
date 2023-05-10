package zlm

import "net/url"

// ListRTPServerRes 是 ListRTPServer 返回值
type ListRTPServerRes struct {
	Code int                     `json:"code"`
	Data []*ListRTPServerResData `json:"data"`
}

// ListRTPServerResData 是 ListRTPServerRes 的 Data 字段
type ListRTPServerResData struct {
	// 绑定的端口号
	Port string `json:"port"`
	// 绑定的流ID
	StreamID int `json:"stream_id"`
}

// ListRTPServer 调用 /index/api/listRtpServer
// 获取openRtpServer接口创建的所有RTP服务器
// 返回 []*ListRTPServerResData
func (s *Server) ListRTPServer() ([]*ListRTPServerResData, error) {
	var res ListRTPServerRes
	query := make(url.Values)
	err := httpGet(s, s.url("listRtpServer"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
