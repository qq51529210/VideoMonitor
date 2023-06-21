package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// AddStreamProxyReq 是 AddStreamProxy 参数
type AddStreamProxyReq struct {
	// 流虚拟主机
	VHost string `query:"vhost"`
	// 流应用名
	App string `query:"app"`
	// 流ID
	Stream string `query:"stream"`
	// 拉流地址，例如 rtmp://live.hkstv.hk.lxdns.com/live/hks2
	URL string `query:"url"`
	// rtsp 拉流方式，0：tcp ，1：udp ，2：组播
	RTPType string `query:"rtp_type"`
	// 超时时间，单位秒，float 类型
	TimeoutSec string `query:"timeout_sec"`
	// 重试次数,不传或 0 无限重试
	RetryCount string `query:"retry_count"`
	// 是否转换成 hls 协议，0 / 1
	EnableHLS string `query:"enable_hls"`
	// 是否 mp4 录制，0 / 1
	EnableMP4 string `query:"enable_mp4"`
	// 是否转换成 rtsp/webrtc 协议，0 / 1
	EnableRTSP string `query:"enable_rtsp"`
	// 是否转换成 rtmp/flv 协议，0 / 1
	EnableRTMP string `query:"enable_rtmp"`
	// 是否转换成 http-ts/ws-ts 协议，0 / 1
	EnableTS string `query:"enable_ts"`
	// 是否转换成 http-fmp4/ws-fmp4 协议，0 / 1
	EnableFMP4 string `query:"enable_fmp4"`
	// 转协议是否开启音频，0 / 1
	EnableAudio string `query:"enable_audio"`
	// 转协议无音频是否添加静音 aac 音频，0 / 1
	AddMuteAudio string `query:"add_mute_audio"`
	// mp4 录制文件保存根目录，置空使用默认
	MP4SavePath string `query:"mp4_save_path"`
	// mp4 录制切片大小，单位秒
	MP4MaxSecond string `query:"mp4_max_second"`
	// hls 文件保存保存根目录，置空使用默认
	HLSSavePath string `query:"hls_save_path"`
	// 是否修改原始时间戳
	// 0：不改变
	// 1：采用接收数据时的系统时间戳
	// 2：采用源视频流时间戳相对时间戳
	ModifyStamp string `query:"modify_stamp"`
}

// addStreamProxyRes 用于解析 addStreamProxy 的返回值
type addStreamProxyRes struct {
	apiRes
	Data struct {
		// 流的唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

const (
	apiPathAddStreamProxy = "addStreamProxy"
)

// AddStreamProxy 调用 /index/api/addStreamProxy 返回 key
func (s *Server) AddStreamProxy(req *AddStreamProxyReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.AddStreamProxyWithContext(ctx, req)
}

// AddStreamProxyWithContext 调用 /index/api/addStreamProxy 返回 key
func (s *Server) AddStreamProxyWithContext(ctx context.Context, req *AddStreamProxyReq) (string, error) {
	var res addStreamProxyRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathAddStreamProxy),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathAddStreamProxy,
		}
	}
	return res.Data.Key, nil
}
