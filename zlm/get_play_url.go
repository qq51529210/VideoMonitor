package zlm

import (
	"fmt"
)

// PlayMediaInfoVideoTrack 是 PlayMediaInfoMediaInfo 的 Video 字段
type PlayMediaInfoVideoTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 视频fps
	FPS int `json:"fps"`
	// 视频高
	Height int `json:"height"`
	// 视频宽
	Width int `json:"width"`
}

// PlayMediaInfoAudioTrack 是 PlayMediaInfoMediaInfo 的 Audio 字段
type PlayMediaInfoAudioTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 音频采样位数
	SampleBit int `json:"sampleBit"`
	// 音频采样率
	SampleRate int `json:"sampleRate"`
}

// PlayMediaInfoMediaInfo 是 PlayMediaInfo 的各种字段
type PlayMediaInfoMediaInfo struct {
	URL   []string                   `json:"url"`
	Audio []*PlayMediaInfoAudioTrack `json:"audioTrack"`
	Video []*PlayMediaInfoVideoTrack `json:"videoTrack"`
}

// PlayMediaInfo 表示所有流的信息
// 使用指针先保存到局部，然后判断是否为空
type PlayMediaInfo struct {
	RTSP *PlayMediaInfoMediaInfo `json:"rtsp"`
	RTMP *PlayMediaInfoMediaInfo `json:"rtmp"`
	// FLV  *PlayMediaInfoMediaInfo `json:"flv"`
	HLS  *PlayMediaInfoMediaInfo `json:"hls"`
	TS   *PlayMediaInfoMediaInfo `json:"ts"`
	FMP4 *PlayMediaInfoMediaInfo `json:"fmp4"`
}

// GetPlayMediaInfo 封装了 GetMediaList
func (s *Server) GetPlayMediaInfo(app, stream string) (*PlayMediaInfo, error) {
	var playURL *PlayMediaInfo
	key := fmt.Sprintf("%s_%s", app, stream)
	// 本地缓存
	s.lock.RLock()
	playURL = s.stream[key]
	s.lock.RUnlock()
	if playURL != nil {
		return playURL, nil
	}
	// 从服务查询
	medias, err := s.GetMediaList(&GetMediaListReq{
		App:    app,
		Stream: stream,
	})
	if err != nil {
		return nil, err
	}
	if len(medias) > 0 {
		playURL = new(PlayMediaInfo)
		s.InitPlayMediaInfos(medias, playURL)
		// 加入列表
		s.lock.RLock()
		s.stream[key] = playURL
		s.lock.RUnlock()
	}
	return playURL, nil
}

// InitPlayMediaInfos 根据 mediaInfos 填充 playURL ，mediaInfos 必须是同一个 app/stream 的，不然会混乱
func (s *Server) InitPlayMediaInfos(medias []*MediaInfo, res *PlayMediaInfo) {
	ip := s.model.PublicIP
	cfg := s.cfg
	for _, media := range medias {
		vs, as := media.ParseTracks()
		var vss []*PlayMediaInfoVideoTrack
		var ass []*PlayMediaInfoAudioTrack
		for _, v := range vs {
			vss = append(vss, &PlayMediaInfoVideoTrack{
				CodecName: v.CodecIDName,
				FPS:       v.FPS,
				Height:    v.Height,
				Width:     v.Width,
			})
		}
		for _, a := range as {
			ass = append(ass, &PlayMediaInfoAudioTrack{
				CodecName:  a.CodecIDName,
				SampleBit:  a.SampleBit,
				SampleRate: a.SampleRate,
			})
		}
		switch media.Schema {
		case RTSP:
			res.RTSP = new(PlayMediaInfoMediaInfo)
			res.RTSP.Video, res.RTSP.Audio = vss, ass
			if cfg.RTSPPort != "" {
				res.RTSP.URL = append(res.RTSP.URL, fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, cfg.RTSPPort, media.App, media.Stream))
			}
			if cfg.RTSPSSLPort != "" {
				res.RTSP.URL = append(res.RTSP.URL, fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, cfg.RTSPSSLPort, media.App, media.Stream))
			}
		case RTMP:
			res.RTMP = new(PlayMediaInfoMediaInfo)
			res.RTMP.Video, res.RTMP.Audio = vss, ass
			if cfg.RTMPPort != "" {
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, cfg.RTMPPort, media.App, media.Stream))
			}
			if cfg.RTMPSSLPort != "" {
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, cfg.RTMPSSLPort, media.App, media.Stream))
			}
			if cfg.HTTPPort != "" {
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
			}
			if cfg.HTTPSSLport != "" {
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLport, media.App, media.Stream))
				res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
			}
		case HLS:
			res.HLS = new(PlayMediaInfoMediaInfo)
			res.HLS.Video, res.HLS.Audio = vss, ass
			if cfg.HTTPPort != "" {
				res.HLS.URL = append(res.HLS.URL, fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPPort, media.App, media.Stream))
			}
			if cfg.HTTPSSLport != "" {
				res.HLS.URL = append(res.HLS.URL, fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPSSLport, media.App, media.Stream))
			}
		case TS:
			res.TS = new(PlayMediaInfoMediaInfo)
			res.TS.Video, res.TS.Audio = vss, ass
			if cfg.HTTPPort != "" {
				res.TS.URL = append(res.TS.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, cfg.HTTPPort, media.App, media.Stream))
			}
			if cfg.HTTPSSLport != "" {
				res.TS.URL = append(res.TS.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, cfg.HTTPSSLport, media.App, media.Stream))
			}
		case FMP4:
			res.FMP4 = new(PlayMediaInfoMediaInfo)
			res.FMP4.Video, res.FMP4.Audio = vss, ass
			if cfg.HTTPPort != "" {
				res.FMP4.URL = append(res.FMP4.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPPort, media.App, media.Stream))
			}
			if cfg.HTTPSSLport != "" {
				res.FMP4.URL = append(res.FMP4.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLport, media.App, media.Stream))
			}
		}
	}
}

// InitPlayMediaInfo 根据 mediaInfos 填充 playURL ，mediaInfos 必须是同一个 app/stream 的，不然会混乱
func (s *Server) InitPlayMediaInfo(media *MediaInfo, res *PlayMediaInfo) {
	ip := s.model.PublicIP
	cfg := s.cfg
	vs, as := media.ParseTracks()
	var vss []*PlayMediaInfoVideoTrack
	var ass []*PlayMediaInfoAudioTrack
	for _, v := range vs {
		vss = append(vss, &PlayMediaInfoVideoTrack{
			CodecName: v.CodecIDName,
			FPS:       v.FPS,
			Height:    v.Height,
			Width:     v.Width,
		})
	}
	for _, a := range as {
		ass = append(ass, &PlayMediaInfoAudioTrack{
			CodecName:  a.CodecIDName,
			SampleBit:  a.SampleBit,
			SampleRate: a.SampleRate,
		})
	}
	switch media.Schema {
	case RTSP:
		res.RTSP = new(PlayMediaInfoMediaInfo)
		res.RTSP.Video, res.RTSP.Audio = vss, ass
		if cfg.RTSPPort != "" {
			res.RTSP.URL = append(res.RTSP.URL, fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, cfg.RTSPPort, media.App, media.Stream))
		}
		if cfg.RTSPSSLPort != "" {
			res.RTSP.URL = append(res.RTSP.URL, fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, cfg.RTSPSSLPort, media.App, media.Stream))
		}
	case RTMP:
		res.RTMP = new(PlayMediaInfoMediaInfo)
		res.RTMP.Video, res.RTMP.Audio = vss, ass
		if cfg.RTMPPort != "" {
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, cfg.RTMPPort, media.App, media.Stream))
		}
		if cfg.RTMPSSLPort != "" {
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, cfg.RTMPSSLPort, media.App, media.Stream))
		}
		if cfg.HTTPPort != "" {
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
		}
		if cfg.HTTPSSLport != "" {
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLport, media.App, media.Stream))
			res.RTMP.URL = append(res.RTMP.URL, fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, media.App, media.Stream))
		}
	case HLS:
		res.HLS = new(PlayMediaInfoMediaInfo)
		res.HLS.Video, res.HLS.Audio = vss, ass
		if cfg.HTTPPort != "" {
			res.HLS.URL = append(res.HLS.URL, fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPPort, media.App, media.Stream))
		}
		if cfg.HTTPSSLport != "" {
			res.HLS.URL = append(res.HLS.URL, fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPSSLport, media.App, media.Stream))
		}
	case TS:
		res.TS = new(PlayMediaInfoMediaInfo)
		res.TS.Video, res.TS.Audio = vss, ass
		if cfg.HTTPPort != "" {
			res.TS.URL = append(res.TS.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, cfg.HTTPPort, media.App, media.Stream))
		}
		if cfg.HTTPSSLport != "" {
			res.TS.URL = append(res.TS.URL, fmt.Sprintf("https://%s:%s/%s/%s.live.ts", ip, cfg.HTTPSSLport, media.App, media.Stream))
		}
	case FMP4:
		res.FMP4 = new(PlayMediaInfoMediaInfo)
		res.FMP4.Video, res.FMP4.Audio = vss, ass
		if cfg.HTTPPort != "" {
			res.FMP4.URL = append(res.FMP4.URL, fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPPort, media.App, media.Stream))
		}
		if cfg.HTTPSSLport != "" {
			res.FMP4.URL = append(res.FMP4.URL, fmt.Sprintf("https://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLport, media.App, media.Stream))
		}
	}
}

// GetPlayURL 返回拼接的各种协议的 URL 列表，不管服务在不在
func (s *Server) GetPlayURL(app, stream string) []string {
	ip := s.model.PublicIP
	cfg := s.cfg
	urls := make([]string, 0, 14)
	if cfg.RTSPPort != "" {
		urls = append(urls, fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, cfg.RTSPPort, app, stream))
	}
	if cfg.RTSPSSLPort != "" {
		urls = append(urls, fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, cfg.RTSPSSLPort, app, stream))
	}
	if cfg.RTMPPort != "" {
		urls = append(urls, fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, cfg.RTMPPort, app, stream))
	}
	if cfg.RTMPSSLPort != "" {
		urls = append(urls, fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, cfg.RTMPSSLPort, app, stream))
	}
	if cfg.HTTPPort != "" {
		urls = append(urls, fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, app, stream))
		urls = append(urls, fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, cfg.HTTPPort, app, stream))
		urls = append(urls, fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPPort, app, stream))
		urls = append(urls, fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, cfg.HTTPPort, app, stream))
		urls = append(urls, fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPPort, app, stream))
	}
	if cfg.HTTPSSLport != "" {
		urls = append(urls, fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLport, app, stream))
		urls = append(urls, fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, cfg.HTTPSSLport, app, stream))
		urls = append(urls, fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, cfg.HTTPSSLport, app, stream))
		urls = append(urls, fmt.Sprintf("https://%s:%s/%s/%s.live.ts", ip, cfg.HTTPSSLport, app, stream))
		urls = append(urls, fmt.Sprintf("https://%s:%s/%s/%s.live.mp4", ip, cfg.HTTPSSLport, app, stream))
	}
	return urls
}
