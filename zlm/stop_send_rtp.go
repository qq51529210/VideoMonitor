package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// StopSendRTPReq 是 StopSendRTP 参数
type StopSendRTPReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 根据ssrc关停某路rtp推流，不传时关闭所有推流
	SSRC string `query:"ssrc"`
}

// stopSendRTPRes 是 StopSendRTP 返回值
type stopSendRTPRes struct {
	Code int `json:"code"`
}

// StopSendRTP 调用 /index/api/stopSendRtp
// 停止 GB28181 ps-rtp 推流。
func (s *Server) StopSendRTP(req *StopSendRTPReq) error {
	var _res stopSendRTPRes
	err := util.HTTP[any](http.MethodGet,
		s.url("stopSendRtp"),
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
	return nil
}
