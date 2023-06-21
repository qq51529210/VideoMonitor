package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// StartSendRTPPassiveReq 是 StartSendRTPPassive 参数
type StartSendRTPPassiveReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 推流的 ssrc
	SSRC string `query:"ssrc"`
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
	// 发送同时接收， 0 / 1
	// 一般用于双向语言对讲, 不为空开启接收，
	// 值为接收流的 id
	RecvStreamID string `query:"recv_stream_id"`
	// 等待 tcp 连接超时时间，单位毫秒，默认 5000
	CloseDelayMS string `query:"close_delay_ms"`
}

// startSendRTPPassiveRes 是 StartSendRTPPassive 返回值
type startSendRTPPassiveRes struct {
	apiRes
	// 本地端口号
	LocalPort int `json:"local_port"`
}

const (
	apiPathStartSendRTPPassive = "startSendRtpPassive"
)

// StartSendRTPPassive 调用 /index/api/startSendRtpPassive 启动 rtp 被动推流，返回使用的本地端口号
func (s *Server) StartSendRTPPassive(req *StartSendRTPPassiveReq) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.StartSendRTPPassiveWithContext(ctx, req)
}

// StartSendRTPPassiveWithContext 调用 /index/api/startSendRtpPassive 启动 rtp 被动推流，返回使用的本地端口号
func (s *Server) StartSendRTPPassiveWithContext(ctx context.Context, req *StartSendRTPPassiveReq) (int, error) {
	var res startSendRTPPassiveRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathStartSendRTPPassive),
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
			API:  apiPathStartSendRTPPassive,
		}
	}
	return res.LocalPort, nil
}
