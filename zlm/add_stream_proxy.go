package zlm

import (
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
	// 拉流地址，例如rtmp://live.hkstv.hk.lxdns.com/live/hks2
	URL string `query:"url"`
	// rtsp拉流时，拉流方式，0：tcp，1：udp，2：组播
	RTPType string `query:"rtp_type"`
	// 拉流超时时间，单位秒，float类型
	TimeoutSec string `query:"timeout_sec"`
	// 拉流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string `query:"retry_count"`
	// 是否转换成hls协议，0/1
	EnableHLS string `query:"enable_hls"`
	// 是否转换成mp4协议，0/1
	EnableMP4 string `query:"enable_mp4"`
	// 是否转换成rtsp协议，0/1
	EnableRTSP string `query:"enable_rtsp"`
	// 是否转换成rtmp/flv协议，0/1
	EnableRTMP string `query:"enable_rtmp"`
	// 是否转换成http-ts/ws-ts协议，0/1
	EnableTS string `query:"enable_ts"`
	// 是否转换成http-fmp4/ws-fmp4协议，0/1
	EnableFMP4 string `query:"enable_fmp4"`
	// 转协议时是否开启音频
	EnableAudio string `query:"enable_audio"`
	// 转协议时，无音频是否添加静音aac音频
	AddMuteAudio string `query:"add_mute_audio"`
	// mp4录制文件保存根目录，置空使用默认
	MP4SavePath string `query:"mp4_save_path"`
	// mp4录制切片大小，单位秒
	MP4MaxSecond string `query:"mp4_max_second"`
	// hls文件保存保存根目录，置空使用默认
	HLSSavePath string `query:"hls_save_path"`
}

// addStreamProxyRes 用于解析 addStreamProxy 的返回值
type addStreamProxyRes struct {
	Code int `json:"code"`
	Data struct {
		// 流的唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

// AddStreamProxy 调用 /index/api/addStreamProxy
// 动态添加rtsp/rtmp/hls拉流代理(只支持H264/H265/aac/G711负载)
// 返回 key
func (s *Server) AddStreamProxy(req *AddStreamProxyReq) (string, error) {
	var _res addStreamProxyRes
	err := util.HTTP[any](http.MethodGet,
		s.url("addStreamProxy"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return "", err
	}
	if _res.Code != 0 {
		return "", CodeError(_res.Code)
	}
	return _res.Data.Key, nil
}
