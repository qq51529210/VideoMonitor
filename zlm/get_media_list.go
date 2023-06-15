package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// // ParseTracks 从 Tracks 解析出 VideoTrack 和 AudioTrack
// func (m *MediaInfo) ParseTracks() ([]*VideoTrack, []*AudioTrack) {
// 	var videos []*VideoTrack
// 	var audios []*AudioTrack
// 	for _, track := range m.Tracks {
// 		codecType, ok := track["codec_type"].(float64)
// 		if !ok {
// 			continue
// 		}
// 		switch codecType {
// 		case 0:
// 			var video VideoTrack
// 			err := video.FromMap(track)
// 			if err != nil {
// 				log.Error(err)
// 				continue
// 			}
// 			videos = append(videos, &video)
// 		case 1:
// 			var audio AudioTrack
// 			err := audio.FromMap(track)
// 			if err != nil {
// 				log.Error(err)
// 				continue
// 			}
// 			audios = append(audios, &audio)
// 		}
// 	}
// 	return videos, audios
// }

// // MediaInfoOriginSock 是 MediaInfo 的 OriginSock 字段
// type MediaInfoOriginSock struct {
// 	Identifier string `json:"identifier"`
// 	LocalIP    string `json:"local_ip"`
// 	LocalPort  int    `json:"local_port"`
// 	PeerIP     string `json:"peer_ip"`
// 	PeerPort   int    `json:"peer_port"`
// }

// // VideoTrack 用于保存 OnStreamChangedReq Tracks 字段中的视频轨道
// type VideoTrack struct {
// 	// H264 = 0, H265 = 1
// 	CodecID int `json:"codec_id"`
// 	// 编码类型名称
// 	CodecIDName string `json:"codec_id_name"`
// 	// 视频fps
// 	FPS int `json:"fps"`
// 	// 视频高
// 	Height int `json:"height"`
// 	// 视频宽
// 	Width int `json:"width"`
// 	// 轨道是否准备就绪
// 	Ready bool `json:"ready"`
// }

// // FromMap 从 data 初始化字段
// func (m *VideoTrack) FromMap(data map[string]any) error {
// 	i, ok := data["codec_id"].(float64)
// 	if !ok {
// 		return errorVideoTrackMissCodecID
// 	}
// 	m.CodecID = int(i)
// 	//
// 	s, ok := data["codec_id_name"].(string)
// 	if !ok {
// 		return errorVideoTrackMissCodecIDName
// 	}
// 	m.CodecIDName = s
// 	//
// 	i, ok = data["fps"].(float64)
// 	if !ok {
// 		return errorVideoTrackMissFPS
// 	}
// 	m.FPS = int(i)
// 	//
// 	i, ok = data["height"].(float64)
// 	if !ok {
// 		return errorVideoTrackMissHeight
// 	}
// 	m.Height = int(i)
// 	//
// 	i, ok = data["width"].(float64)
// 	if !ok {
// 		return errorVideoTrackMissWidth
// 	}
// 	m.Width = int(i)
// 	//
// 	b, ok := data["ready"].(bool)
// 	if !ok {
// 		return errorVideoTrackMissReady
// 	}
// 	m.Ready = b
// 	//
// 	return nil
// }

// // AudioTrack 用于保存 OnStreamChangedReq Tracks 字段中的音频轨道
// type AudioTrack struct {
// 	// 音频通道数
// 	Channels int `json:"channels"`
// 	// AAC = 2, G711A = 3, G711U = 4
// 	CodecID int `json:"codec_id"`
// 	// 编码类型名称
// 	CodecIDName string `json:"codec_id_name"`
// 	// 音频采样位数
// 	SampleBit int `json:"sample_bit"`
// 	// 音频采样率
// 	SampleRate int `json:"sample_rate"`
// 	// 轨道是否准备就绪
// 	Ready bool `json:"ready"`
// }

// // FromMap 从 data 初始化字段
// func (m *AudioTrack) FromMap(data map[string]any) error {
// 	i, ok := data["channels"].(float64)
// 	if !ok {
// 		return errorAudioTrackMissChannels
// 	}
// 	m.Channels = int(i)
// 	//
// 	i, ok = data["codec_id"].(float64)
// 	if !ok {
// 		return errorAudioTrackMissCodecID
// 	}
// 	m.CodecID = int(i)
// 	//
// 	s, ok := data["codec_id_name"].(string)
// 	if !ok {
// 		return errorAudioTrackMissCodecIDName
// 	}
// 	m.CodecIDName = s
// 	//
// 	i, ok = data["sample_bit"].(float64)
// 	if !ok {
// 		return errorAudioTrackMissSampleBit
// 	}
// 	m.SampleBit = int(i)
// 	//
// 	i, ok = data["sample_rate"].(float64)
// 	if !ok {
// 		return errorAudioTrackMissSampleRate
// 	}
// 	m.SampleRate = int(i)
// 	//
// 	b, ok := data["ready"].(bool)
// 	if !ok {
// 		return errorAudioTrackMissReady
// 	}
// 	m.Ready = b
// 	//
// 	return nil
// }

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
}

// getMediaListRes 用于解析 getMediaList 的返回值
type getMediaListRes struct {
	Code int          `json:"code"`
	Data []*MediaList `json:"data"`
}

// MediaList 是 getMediaListRes 的 Data 字段
type MediaList struct {
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 播放或推流的协议，可能是rtsp、rtmp、http
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// VideoTrack/AudioTrack，Video: codec_type= 0, Audio: codec_type= 1
	Tracks []map[string]any `json:"tracks"`
	// 是否正在录像 hls
	IsRecordingHLS bool `json:"isRecordingHLS"`
	// 是否正在录像 mp4
	IsRecordingMP4 bool `json:"isRecordingMP4"`
	// 存活时间，单位秒
	// AliveSecond int `json:"aliveSecond"`
	// 数据产生速度，单位byte
	// BytesSpeed int `json:"bytesSpeed"`
	// unix系统时间戳，单位秒
	// CreateStamp int64 `json:"createStamp"`
	// 本协议观看人数
	// ReaderCount int `json:"readerCount"`
	// 观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv/rtc
	// TotalReaderCount int64 `json:"totalReaderCount"`
	// 产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7,rtc_push=8
	// OriginType    int    `json:"originType"`
	// OriginTypeStr string `json:"originTypeStr"`
	// 产生源的url
	// OriginURL string `json:"originUrl"`
	// OriginSock  *MediaInfoOriginSock `json:"originSock"`
}

// GetMediaList 调用 /index/api/getMediaList
// 获取流列表，可选筛选参数
func (s *Server) GetMediaList(req *GetMediaListReq) ([]*MediaList, error) {
	var _res getMediaListRes
	err := util.HTTP[any](http.MethodGet,
		s.url("getMediaList"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return nil, err
	}
	if _res.Code != 0 {
		return nil, CodeError(_res.Code)
	}
	return _res.Data, nil
}
