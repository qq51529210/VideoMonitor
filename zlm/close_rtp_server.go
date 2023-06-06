package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// CloseRTPServerReq 是 CloseRTPServer 的参数
type CloseRTPServerReq struct {
	// 调用closeRtpServer接口时提供的流ID
	StreamID string `query:"stream_id"`
}

// closeRTPServerRes 是 CloseRTPServer 的返回值
type closeRTPServerRes struct {
	Code int `json:"code"`
	// 是否找到记录并关闭
	Hit int `json:"hit"`
}

// CloseRTPServer 调用 /index/api/closeRtpServer
// 关闭GB28181 RTP接收端口
// 返回成功的个数
func (s *Server) CloseRTPServer(req *CloseRTPServerReq) (int, error) {
	var _res closeRTPServerRes
	err := util.HTTP[any](http.MethodGet,
		s.url("closeRtpServer"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return 0, err
	}
	if _res.Code != 0 {
		return 0, CodeError(_res.Code)
	}
	return _res.Hit, nil
}
