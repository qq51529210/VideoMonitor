package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// StartSendRTPReq 是 StartSendRTP 参数
type StartSendRTPReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 推流的 ssrc
	SSRC string `query:"ssrc"`
	// 目标 ip 或域名
	DstURL string `query:"dst_url"`
	// 目标端口
	DstPort string `query:"dst_port"`
	// 0：udp ，1：tcp
	IsUDP string `query:"is_udp"`
	// 本机端口，0 或不传为随机端口
	SrcPort string `query:"src_port"`
	// 是否推送本地MP4录像，0 / 1
	FromMP4 string `query:"from_mp4"`
	// rtp 打包， 0：es ，1：ps ，默认 1
	UsePS string `query:"use_ps"`
	// rtp 的 pt ，默认 96
	PT string `query:"pt"`
	// es 打包时，只打包音频 0 / 1
	OnlyAudio string `query:"only_audio"`
	// udp 方式推流时，开启 rtcp 发送和接收超时判断，默认关闭
	// 如果超时，将导致主动停止 rtp 发送
	UDPRTCPTimeout string `query:"udp_rtcp_timeout"`
	// 发送同时接收， 0 / 1
	// 一般用于双向语言对讲, 不为空开启接收，
	// 值为接收流的 id
	RecvStreamID string `query:"recv_stream_id"`
}

// startSendRTPRes 是 StartSendRTP 返回值
type startSendRTPRes struct {
	apiRes
	// 本地端口号
	LocalPort int `json:"local_port"`
}

const (
	apiPathStartSendRTP = "startSendRtp"
)

// StartSendRTP 调用 /index/api/startSendRtp 启动 rtp 推流，返回使用的本地端口号
func (s *Server) StartSendRTP(req *StartSendRTPReq) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.StartSendRTPWithContext(ctx, req)
}

// StartSendRTPWithContext 调用 /index/api/startSendRtp 启动 rtp 推流，返回使用的本地端口号
func (s *Server) StartSendRTPWithContext(ctx context.Context, req *StartSendRTPReq) (int, error) {
	var res startSendRTPRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathStartSendRTP),
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
			API:  apiPathStartSendRTP,
		}
	}
	return res.LocalPort, nil
}
