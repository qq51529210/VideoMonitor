package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// StopSendRTPReq 是 StopSendRTP 参数
type StopSendRTPReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 根据 ssrc 关停某路 rtp 推流，不传时关闭所有推流
	SSRC string `query:"ssrc"`
}

const (
	apiPathStopSendRTP = "stopSendRtp"
)

// StopSendRTP 调用 /index/api/stopSendRtp 停止 rtp 推流
func (s *Server) StopSendRTP(req *StopSendRTPReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.StopSendRTPWithContext(ctx, req)
}

// StopSendRTPWithContext 调用 /index/api/stopSendRtp 停止 rtp 推流
func (s *Server) StopSendRTPWithContext(ctx context.Context, req *StopSendRTPReq) error {
	var res apiRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathStopSendRTP),
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
			API:  apiPathStopSendRTP,
		}
	}
	return nil
}
