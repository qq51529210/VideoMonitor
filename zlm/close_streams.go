package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// CloseStreamsReq 是 CloseStreams 参数
type CloseStreamsReq struct {
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 是否强制关闭(有人在观看是否还关闭)，0/1
	Force string `query:"force"`
}

// closeStreamsRes 用于解析 close_streams 的返回值
type closeStreamsRes struct {
	Code int `json:"code"`
	// 筛选命中的流个数
	CountHit int `json:"count_hit"`
	// 被关闭的流个数，可能小于count_hit
	CountClosed int `json:"count_closed"`
}

// CloseStreamsRes 是 CloseStreams 的返回值
type CloseStreamsRes struct {
	// 筛选命中的流个数
	CountHit int
	// 被关闭的流个数，可能小于count_hit
	CountClosed int
}

// CloseStreams 调用 /index/api/close_streams
// 关闭流(目前所有类型的流都支持关闭)
func (s *Server) CloseStreams(req *CloseStreamsReq, res *CloseStreamsRes) error {
	var _res closeStreamsRes
	err := util.HTTP[any](http.MethodGet,
		s.url("close_streams"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return err
	}
	if _res.Code != 0 {
		return CodeError(_res.Code)
	}
	res.CountHit = _res.CountHit
	res.CountClosed = _res.CountClosed
	return nil
}
