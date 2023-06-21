package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
}

// getMediaListRes 用于解析 getMediaList 的返回值
type getMediaListRes struct {
	apiRes
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
	// 流 ID
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

const (
	apiPathGetMediaList = "getMediaList"
)

// GetMediaList 调用 /index/api/getMediaList 获取流列表
func (s *Server) GetMediaList(req *GetMediaListReq) ([]*MediaList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.GetMediaListWithContext(ctx, req)
}

// GetMediaListWithContext 调用 /index/api/getMediaList 获取流列表
func (s *Server) GetMediaListWithContext(ctx context.Context, req *GetMediaListReq) ([]*MediaList, error) {
	var res getMediaListRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathGetMediaList),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathGetMediaList,
		}
	}
	return res.Data, nil
}
