package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// CloseStreamsReq 是 CloseStreams 参数
type CloseStreamsReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 协议，例如 rtsp 或 rtmp
	Schema string `query:"schema"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 是否强制关闭(有人在观看是否还关闭)，0/1
	Force string `query:"force"`
}

// // closeStreamsRes 用于解析 close_streams 的返回值
// type closeStreamsRes struct {
// 	apiRes
// 	// 筛选命中的流个数
// 	CountHit int `json:"count_hit"`
// 	// 被关闭的流个数，可能小于count_hit
// 	CountClosed int `json:"count_closed"`
// }

// // CloseStreamsRes 是 CloseStreams 的返回值
// type CloseStreamsRes struct {
// 	// 筛选命中的流个数
// 	CountHit int
// 	// 被关闭的流个数，可能小于count_hit
// 	CountClosed int
// }

const (
	apiPathCloseStreams = "close_streams"
)

// CloseStreams 调用 /index/api/close_streams 关闭流
func (s *Server) CloseStreams(req *CloseStreamsReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.CloseStreamsWithContext(ctx, req)
}

// CloseStreamsWithContext 调用 /index/api/close_streams 关闭流
func (s *Server) CloseStreamsWithContext(ctx context.Context, req *CloseStreamsReq) error {
	var res apiRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathCloseStreams),
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
			API:  apiPathCloseStreams,
		}
	}
	return nil
}
