package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// listRTPServerRes 是 ListRTPServer 返回值
type listRTPServerRes struct {
	Code int                 `json:"code"`
	Data []*ListRTPServerRes `json:"data"`
}

// ListRTPServerRes 是 ListRTPServerRes 的 Data 字段
type ListRTPServerRes struct {
	// 绑定的端口号
	Port string `json:"port"`
	// 绑定的流ID
	StreamID int `json:"stream_id"`
}

// ListRTPServer 调用 /index/api/listRtpServer
// 获取openRtpServer接口创建的所有RTP服务器
// 返回 []*ListRTPServerResData
func (s *Server) ListRTPServer() ([]*ListRTPServerRes, error) {
	var _res listRTPServerRes
	err := util.HTTP[any](http.MethodGet,
		s.url("listRtpServer"),
		s.query(nil),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return nil, err
	}
	if _res.Code != 0 {
		return nil, CodeError(_res.Code)
	}
	return _res.Data, nil
}
