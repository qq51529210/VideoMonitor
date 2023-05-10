package zlm

import "net/url"

// StartSendRTPPassiveReq 是 StartSendRTPPassive 参数
type StartSendRTPPassiveReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string
	// 添加的应用名，例如 live
	App string
	// 添加的流id，例如 test
	Stream string
	// 推流的rtp的ssrc,指定不同的ssrc可以同时推流到多个服务器
	SSRC string
	// 使用的本机端口，为0或不传时默认为随机端口
	SrcPort string
	// 发送时，rtp的pt（uint8_t）,不传时默认为96
	PT string
	// 发送时，rtp的负载类型。为1时，负载为ps；为0时，为es；不传时默认为1
	UsePS string
	// rtp es方式打包时，是否只打包音频，该参数非必选参数
	OnlyAudio string
	// 是否推送本地MP4录像，该参数非必选参数，0/1
	FromMP4 string
}

func (m *StartSendRTPPassiveReq) toQuery() url.Values {
	q := make(url.Values)
	if m.VHost != "" {
		q.Set("vhost", m.VHost)
	}
	if m.App != "" {
		q.Set("app", m.App)
	}
	if m.Stream != "" {
		q.Set("stream", m.Stream)
	}
	if m.SSRC != "" {
		q.Set("ssrc", m.SSRC)
	}
	if m.SrcPort != "" {
		q.Set("src_port", m.SrcPort)
	}
	if m.PT != "" {
		q.Set("pt", m.PT)
	}
	if m.UsePS != "" {
		q.Set("use_ps", m.UsePS)
	}
	if m.OnlyAudio != "" {
		q.Set("only_audio", m.OnlyAudio)
	}
	if m.FromMP4 != "" {
		q.Set("from_mp4", m.FromMP4)
	}
	return q
}

// startSendRTPPassiveRes 是 StartSendRTPPassive 返回值
type startSendRTPPassiveRes struct {
	Code int `json:"code"`
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

// StartSendRTPPassive 调用 /index/api/startSendRtpPassive
// 作为GB28181 Passive TCP服务器；该接口支持rtsp/rtmp等协议转ps-rtp被动推流。调用该接口，zlm会启动tcp服务器等待连接请求，连接建立后，zlm会关闭tcp服务器，然后源源不断的往客户端推流。
// 第一次推流失败会直接返回错误，成功一次后，后续失败也将无限重试(不停地建立tcp监听，超时后再关闭)。
// 返回使用的本地端口号
func (s *Server) StartSendRTPPassive(req *StartSendRTPPassiveReq) (int, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res startSendRTPPassiveRes
	err := httpGet(s, s.url("startSendRtpPassive"), query, &res)
	if err != nil {
		return -1, err
	}
	if res.Code != 0 {
		return -1, CodeError(res.Code)
	}
	return res.LocalPort, nil
}
