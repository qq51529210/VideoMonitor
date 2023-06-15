package zlm

import (
	"errors"
	"fmt"

	"github.com/qq51529210/log"
)

var (
	errVideoTrackMissCodecID     = errors.New("video track missing field codec_id")
	errVideoTrackMissCodecIDName = errors.New("video track missing field codec_id_name")
	errVideoTrackMissFPS         = errors.New("video track missing field fps")
	errVideoTrackMissCodecType   = errors.New("video track missing field codec_type")
	errVideoTrackMissHeight      = errors.New("video track missing field height")
	errVideoTrackMissWidth       = errors.New("video track missing field width")
	errVideoTrackMissReady       = errors.New("video track missing field ready")

	errAudioTrackMissChannels    = errors.New("audio track missing field channels")
	errAudioTrackMissCodecID     = errors.New("audio track missing field codec_id")
	errAudioTrackMissCodecIDName = errors.New("audio track missing field codec_id_name")
	errAudioTrackMissCodecType   = errors.New("audio track missing field codec_type")
	errAudioTrackMissReady       = errors.New("audio track missing field ready")
	errAudioTrackMissSampleBit   = errors.New("audio track missing field sample_bit")
	errAudioTrackMissSampleRate  = errors.New("audio track missing field sample_rate")

	errMediaNotFound = errors.New("media not found")
)

// parseTracks 从 tracks 解析出 MediaInfoVideoTrack 和 MediaInfoAudioTrack
func parseTracks(tracks []map[string]any) ([]*MediaInfoVideoTrack, []*MediaInfoAudioTrack) {
	var videos []*MediaInfoVideoTrack
	var audios []*MediaInfoAudioTrack
	for _, track := range tracks {
		codecType, ok := track["codec_type"].(float64)
		if !ok {
			continue
		}
		switch codecType {
		case 0:
			video := new(MediaInfoVideoTrack)
			err := video.parse(track)
			if err != nil {
				log.Error(err)
				continue
			}
			videos = append(videos, video)
		case 1:
			audio := new(MediaInfoAudioTrack)
			err := audio.parse(track)
			if err != nil {
				log.Error(err)
				continue
			}
			audios = append(audios, audio)
		}
	}
	return videos, audios
}

// MediaInfoVideoTrack 是 MediaInfo 的 Video 字段
type MediaInfoVideoTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 视频fps
	FPS int `json:"fps"`
	// 视频高
	Height int `json:"height"`
	// 视频宽
	Width int `json:"width"`
}

// parse 从 data 初始化字段
func (m *MediaInfoVideoTrack) parse(data map[string]any) error {
	//
	s, ok := data["codec_id_name"].(string)
	if !ok {
		return errVideoTrackMissCodecIDName
	}
	m.CodecName = s
	//
	i, ok := data["fps"].(float64)
	if !ok {
		return errVideoTrackMissFPS
	}
	m.FPS = int(i)
	//
	i, ok = data["height"].(float64)
	if !ok {
		return errVideoTrackMissHeight
	}
	m.Height = int(i)
	//
	i, ok = data["width"].(float64)
	if !ok {
		return errVideoTrackMissWidth
	}
	m.Width = int(i)
	//
	return nil
}

// MediaInfoAudioTrack 是 MediaInfo 的 Audio 字段
type MediaInfoAudioTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 音频采样位数
	SampleBit int `json:"sampleBit"`
	// 音频采样率
	SampleRate int `json:"sampleRate"`
}

// parse 从 data 初始化字段
func (m *MediaInfoAudioTrack) parse(data map[string]any) error {
	s, ok := data["codec_id_name"].(string)
	if !ok {
		return errAudioTrackMissCodecIDName
	}
	m.CodecName = s
	//
	i, ok := data["sample_bit"].(float64)
	if !ok {
		return errAudioTrackMissSampleBit
	}
	m.SampleBit = int(i)
	//
	i, ok = data["sample_rate"].(float64)
	if !ok {
		return errAudioTrackMissSampleRate
	}
	m.SampleRate = int(i)
	//
	return nil
}

// MediaInfo 用于返回
type MediaInfo struct {
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsp://ip:port/app/stream
	RTSP string `json:"rtsp,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsps://ip:port/app/stream
	RTSPs string `json:"rtsps,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmp://ip:port/app/stream
	RTMP string `json:"rtmp,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmps://ip:port/app/stream
	RTMPs string `json:"rtmps,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.flv
	HTTPFLV string `json:"httpFLV,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.flv
	HTTPsFLV string `json:"httpsFLV,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream/hls.m3u8
	HLS string `json:"hls,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream/hls.m3u8
	HTTPsHLS string `json:"httpsHLS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.ts
	HTTPTS string `json:"httpTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.ts
	HTTPsTS string `json:"httpsTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.mp4
	HTTPFMP4 string `json:"httpFMP4,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.mp4
	HTTPsFMP4 string `json:"httpsFMP4,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.flv
	WSFLV string `json:"wsFLV,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.flv
	WSsFLV string `json:"wssFLV,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream/hls.m3u8
	WSHLS string `json:"wsHLS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream/hls.m3u8
	WSsHLS string `json:"wssHLS,omitempty"`
	// ws://ip:port/app/stream.live.ts
	WSTS string `json:"wsTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.ts
	WSsTS string `json:"wssTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.mp4
	WSFMP4 string `json:"wsFMP4,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.mp4
	WSsFMP4 string `json:"wssFMP4,omitempty"`
	// 流的标识
	App string `json:"app"`
	// 流的标识
	Stream string `json:"stream"`
	// 音频轨道信息
	Audio []*MediaInfoAudioTrack `json:"audioTracks"`
	// 视频轨道信息
	Video []*MediaInfoVideoTrack `json:"videoTracks"`
	// 属于的流媒体服务
	Ser *Server `json:"-"`
	// 是否正在录制 hls
	isRecordingHLS bool
	// 是否正在录制 mp4
	isRecordingMP4 bool
}

// Init 根据 App ，Stream ，Ser 字段初始化各个 url
// public 决定使用内网还是外网 ip
// token 加入 token
func (m *MediaInfo) Init(token string, public bool) bool {
	// 上锁
	m.Ser.lock.RLock()
	app := m.Ser.media[m.App]
	if app == nil {
		m.Ser.lock.RUnlock()
		return false
	}
	stream := app[m.Stream]
	if stream == nil {
		m.Ser.lock.RUnlock()
		return false
	}
	m.Ser.lock.RUnlock()
	// 轨道
	m.Video = stream.Video
	m.Audio = stream.Audio
	m.isRecordingHLS = stream.IsRecordingMP4
	m.isRecordingMP4 = stream.IsRecordingMP4
	// url
	cfg := m.Ser.Cfg
	var ip string
	if public {
		ip = m.Ser.PublicIP
	} else {
		ip = m.Ser.PrivateIP
	}
	// RTMP
	if cfg.RTMPPort != "" && cfg.RTMPPort != "0" {
		if token != "" {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s?token=%s", ip, cfg.RTMPPort, m.App, m.Stream, token)
		} else {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, cfg.RTMPPort, m.App, m.Stream)
		}
	}
	if cfg.RTMPSSLPort != "" && cfg.RTMPSSLPort != "0" {
		if token != "" {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s?token=%s", ip, cfg.RTMPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, cfg.RTMPSSLPort, m.App, m.Stream)
		}
	}
	// RTSP
	if cfg.RTSPPort != "" && cfg.RTSPPort != "0" {
		if token != "" {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s?token=%s", ip, cfg.RTSPPort, m.App, m.Stream, token)
		} else {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, cfg.RTSPPort, m.App, m.Stream)
		}
	}
	if cfg.RTSPSSLPort != "" && cfg.RTSPSSLPort != "0" {
		if token != "" {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s?token=%s", ip, cfg.RTSPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, cfg.RTSPSSLPort, m.App, m.Stream)
		}
	}
	if cfg.HTTPPort != "" && cfg.HTTPPort != "0" {
		if token != "" {
			// FLV
			m.HTTPFLV = fmt.Sprintf("http://%s:%s/%s/%s.live.flv?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			m.WSFLV = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			// HLS
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			m.WSHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			// TS
			m.HTTPTS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			m.WSTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			// MP4
			m.HTTPFMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
			m.WSFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4?token=%s", ip, cfg.HTTPPort, m.App, m.Stream, token)
		} else {
			// FLV
			m.HTTPFLV = fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, m.App, m.Stream)
			m.WSFLV = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, m.App, m.Stream)
			// HLS
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPPort, m.App, m.Stream)
			m.WSHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPPort, m.App, m.Stream)
			// TS
			m.HTTPTS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, cfg.HTTPPort, m.App, m.Stream)
			m.WSTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts", ip, cfg.HTTPPort, m.App, m.Stream)
			// MP4
			m.HTTPFMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			m.WSFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLPort, m.App, m.Stream)
		}
	}
	if cfg.HTTPSSLPort != "" && cfg.HTTPSSLPort != "0" {
		if token != "" {
			// FLV
			m.HTTPsFLV = fmt.Sprintf("https://%s:%s/%s/%s.live.flv?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			m.WSsFLV = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			// HLS
			m.HLS = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			m.WSsHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			// TS
			m.HTTPsTS = fmt.Sprintf("https://%s:%s/%s/%s.live.ts?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			m.WSsTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			// MP4
			m.HTTPsFMP4 = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
			m.WSsFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4?token=%s", ip, cfg.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			// FLV
			m.HTTPsFLV = fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			m.WSsFLV = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			// HLS
			m.HLS = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			m.WSsHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			// TS
			m.HTTPsTS = fmt.Sprintf("https://%s:%s/%s/%s.live.ts", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			m.WSsTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			// MP4
			m.HTTPsFMP4 = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLPort, m.App, m.Stream)
			m.WSsFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLPort, m.App, m.Stream)
		}
	}
	//
	return true
}

// mediaInfo 用于简单的保存媒体流信息
type mediaInfo struct {
	App    string
	Stream string
	// 音频轨道信息
	Audio []*MediaInfoAudioTrack
	// 视频轨道信息
	Video []*MediaInfoVideoTrack
	// 是否正在录像 hls
	IsRecordingHLS bool
	// 是否正在录像 mp4
	IsRecordingMP4 bool
}
