package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// OpenRTPServerReq 是 OpenRTPServer 的参数
type OpenRTPServerReq struct {
	// 接收端口，0则为随机端口
	Port string `query:"port"`
	// 创建 udp端口时是否同时监听tcp端口
	EnableTCP string `query:"enable_tcp"`
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	StreamID string `query:"stream_id"`
	// 是否重用端口，默认为0，非必选参数，0/1
	ReusePort string `query:"reuse_port"`
	// 是否指定收流的rtp ssrc, 十进制数字，不指定或指定0时则不过滤rtp，非必选参数
	SSRC string `query:"ssrc"`
}

// openRTPServerRes 是 OpenRTPServer 的返回值
type openRTPServerRes struct {
	Code int `json:"code"`
	// 接收端口，方便获取随机端口号
	Port int `json:"port"`
}

// OpenRTPServer 调用 /index/api/openRtpServer
// 创建GB28181 RTP接收端口，如果该端口接收数据超时，则会自动被回收(不用调用closeRtpServer接口)
// 返回使用的端口
func (s *Server) OpenRTPServer(req *OpenRTPServerReq) (int, error) {
	var _res openRTPServerRes
	err := util.HTTP[any](http.MethodGet,
		s.url("openRtpServer"),
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
	return _res.Port, nil
}
