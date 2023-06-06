package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// StartSendRTPReq 是 StartSendRTP 参数
type StartSendRTPReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 推流的rtp的ssrc,指定不同的ssrc可以同时推流到多个服务器
	SSRC string `query:"ssrc"`
	// 目标ip或域名
	DstURL string `query:"dst_url"`
	// 目标端口
	DstPort string `query:"dst_port"`
	// 是否为udp模式，否则为tcp模式，0/1
	IsUDP string `query:"is_udp"`
	// 使用的本机端口，为0或不传时默认为随机端口
	SrcPort string `query:"src_port"`
	// 发送时，rtp的pt（uint8_t）,不传时默认为96
	PT string `query:"pt"`
	// 发送时，rtp的负载类型。为1时，负载为ps；为0时，为es；不传时默认为1
	UsePS string `query:"use_ps"`
	// rtp es方式打包时，是否只打包音频，该参数非必选参数
	OnlyAudio string `query:"only_audio"`
	// 是否推送本地MP4录像，该参数非必选参数，0/1
	FromMP4 string `query:"from_mp4"`
}

// startSendRTPRes 是 StartSendRTP 返回值
type startSendRTPRes struct {
	Code int `json:"code"`
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

// StartSendRTP 调用 /index/api/startSendRtp
// 作为GB28181客户端，启动ps-rtp推流，支持rtp/udp方式；该接口支持rtsp/rtmp等协议转ps-rtp推流。
// 第一次推流失败会直接返回错误，成功一次后，后续失败也将无限重试。
// 返回使用的本地端口号
func (s *Server) StartSendRTP(req *StartSendRTPReq) (int, error) {
	var _res startSendRTPRes
	err := util.HTTP[any](http.MethodGet,
		s.url("startSendRtp"),
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
	return _res.LocalPort, nil
}
