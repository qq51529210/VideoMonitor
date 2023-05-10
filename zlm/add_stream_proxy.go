package zlm

import "net/url"

// AddStreamProxyReq 是 AddStreamProxy 参数
type AddStreamProxyReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string
	// 添加的应用名，例如 live
	App string
	// 添加的流id，例如 test
	Stream string
	// 拉流地址，例如rtmp://live.hkstv.hk.lxdns.com/live/hks2
	URL string
	// rtsp拉流时，拉流方式，0：tcp，1：udp，2：组播
	RTPType string
	// 拉流超时时间，单位秒，float类型
	TimeoutSec string
	// 拉流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string
	// 是否转换成hls协议，0/1
	EnableHLS string
	// 是否转换成mp4协议，0/1
	EnableMP4 string
	// 是否转换成rtsp协议，0/1
	EnableRTSP string
	// 是否转换成rtmp/flv协议，0/1
	EnableRTMP string
	// 是否转换成http-ts/ws-ts协议，0/1
	EnableTS string
	// 是否转换成http-fmp4/ws-fmp4协议，0/1
	EnableFMP4 string
	// 转协议时是否开启音频
	EnableAudio string
	// 转协议时，无音频是否添加静音aac音频
	AddMuteAudio string
	// mp4录制文件保存根目录，置空使用默认
	MP4SavePath string
	// mp4录制切片大小，单位秒
	MP4MaxSecond string
	// hls文件保存保存根目录，置空使用默认
	HLSSavePath string
}

func (m *AddStreamProxyReq) toQuery() url.Values {
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
	if m.URL != "" {
		q.Set("url", m.URL)
	}
	if m.RetryCount != "" {
		q.Set("retry_count", m.RetryCount)
	}
	if m.RTPType != "" {
		q.Set("rtp_type", m.RTPType)
	}
	if m.TimeoutSec != "" {
		q.Set("timeout_sec	", m.TimeoutSec)
	}
	if m.EnableHLS != "" {
		q.Set("enable_hls", m.EnableHLS)
	}
	if m.EnableMP4 != "" {
		q.Set("enable_mp4", m.EnableMP4)
	}
	if m.EnableRTSP != "" {
		q.Set("enable_rtsp", m.EnableRTSP)
	}
	if m.EnableRTMP != "" {
		q.Set("enable_rtmp", m.EnableRTMP)
	}
	if m.EnableTS != "" {
		q.Set("enable_ts", m.EnableTS)
	}
	if m.EnableFMP4 != "" {
		q.Set("enable_fmp4", m.EnableFMP4)
	}
	if m.EnableAudio != "" {
		q.Set("enable_audio", m.EnableAudio)
	}
	if m.AddMuteAudio != "" {
		q.Set("add_mute_audio", m.AddMuteAudio)
	}
	if m.MP4SavePath != "" {
		q.Set("mp4_save_path", m.MP4SavePath)
	}
	if m.MP4MaxSecond != "" {
		q.Set("mp4_max_second", m.MP4MaxSecond)
	}
	if m.HLSSavePath != "" {
		q.Set("hls_save_path", m.HLSSavePath)
	}
	return q
}

// addStreamProxyRes 是 AddStreamProxy 返回值
type addStreamProxyRes struct {
	Code int                   `json:"code"`
	Data AddStreamProxyResData `json:"data"`
}

// AddStreamProxyResData 是 addStreamProxyRes 的 Data 字段
type AddStreamProxyResData struct {
	// 流的唯一标识
	Key string
}

// AddStreamProxy 调用 /index/api/addStreamProxy
// 动态添加rtsp/rtmp/hls拉流代理(只支持H264/H265/aac/G711负载)
// 返回 key
func (s *Server) AddStreamProxy(req *AddStreamProxyReq) (string, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res addStreamProxyRes
	err := httpGet(s, s.url("addStreamProxy"), query, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", CodeError(res.Code)
	}
	return res.Data.Key, nil
}
