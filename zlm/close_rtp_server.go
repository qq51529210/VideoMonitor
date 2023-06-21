package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// CloseRTPServerReq 是 CloseRTPServer 的参数
type CloseRTPServerReq struct {
	// 绑定的流 id
	StreamID string `query:"stream_id"`
}

// // closeRTPServerRes 是 CloseRTPServer 的返回值
// type closeRTPServerRes struct {
// 	Code int `json:"code"`
// 	// 是否找到记录并关闭
// 	Hit int `json:"hit"`
// }

const (
	apiPathCloseRTPServer = "closeRtpServer"
)

// CloseRTPServer 调用 /index/api/closeRtpServer 关闭 RTP 接收端口
func (s *Server) CloseRTPServer(req *CloseRTPServerReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.CloseRTPServerWithContext(ctx, req)
}

// CloseRTPServerWithContext 调用 /index/api/closeRtpServer 关闭 RTP 接收端口
func (s *Server) CloseRTPServerWithContext(ctx context.Context, req *CloseRTPServerReq) error {
	var res apiRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathCloseRTPServer),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathCloseRTPServer,
		}
	}
	return nil
}
