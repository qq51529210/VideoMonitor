package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// OpenRTPServerReq 是 OpenRTPServer 的参数
type OpenRTPServerReq struct {
	// 监听的端口，0 为随机端口
	Port string `query:"port"`
	// 0：udp ，1：tcp ，2：tcp 主动连接模式
	TCPMode string `query:"tcp_mode"`
	// 绑定的流 id
	StreamID string `query:"stream_id"`
	// 是否重用端口，0 / 1
	ReusePort string `query:"reuse_port"`
	// rtp ssrc
	SSRC string `query:"ssrc"`
	// 是否为单音频 0 / 1 ，用于语音对讲
	OnlyAudio string `query:"only_audio"`
}

// openRTPServerRes 是 OpenRTPServer 的返回值
type openRTPServerRes struct {
	apiRes
	// 本机端口
	Port int `json:"port"`
}

const (
	apiPathOpenRTPServer = "openRtpServer"
)

// OpenRTPServer 调用 /index/api/openRtpServer 返回使用的端口
func (s *Server) OpenRTPServer(req *OpenRTPServerReq) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.OpenRTPServerWithContext(ctx, req)
}

// OpenRTPServerWithContext 调用 /index/api/openRtpServer 返回使用的端口
func (s *Server) OpenRTPServerWithContext(ctx context.Context, req *OpenRTPServerReq) (int, error) {
	var res openRTPServerRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathOpenRTPServer),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return 0, err
	}
	if res.Code != 0 {
		return 0, &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathOpenRTPServer,
		}
	}
	return res.Port, nil
}
