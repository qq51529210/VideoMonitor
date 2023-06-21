package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// AddFFMPEGSourceReq 是 AddFFMPEGSource 的参数
type AddFFMPEGSourceReq struct {
	// 拉流地址
	SrcURL string `query:"src_url"`
	// 推流地址
	DstURL string `query:"dst_url"`
	// 超时时间，单位毫秒
	TimeoutMS string `query:"timeout_ms"`
	// 是否开启 hls 录制，0/1
	EnableHLS string `query:"enable_hls"`
	// 是否开启 mp4 录制，0/1
	EnableMP4 string `query:"enable_mp4"`
	// 配置文件中 FFmpeg 命令参数模板key(非内容)，
	// 置空则采用默认模板：ffmpeg.cmd
	CmdKey string `query:"ffmpeg_cmd_key"`
}

// addFFMPEGSourceRes 用于解析 addFFmpegSource 的返回值
type addFFMPEGSourceRes struct {
	apiRes
	Data struct {
		// 唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

const (
	apiPathAddFFMPEGSource = "addFFmpegSource"
)

// AddFFMPEGSource 调用 /index/api/addFFmpegSource 返回 key
func (s *Server) AddFFMPEGSource(req *AddFFMPEGSourceReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.AddFFMPEGSourceWithContext(ctx, req)
}

// AddFFMPEGSourceWithContext 调用 /index/api/addFFmpegSource 返回 key
func (s *Server) AddFFMPEGSourceWithContext(ctx context.Context, req *AddFFMPEGSourceReq) (string, error) {
	var res addFFMPEGSourceRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathAddFFMPEGSource),
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
			API:  apiPathAddFFMPEGSource,
		}
	}
	return res.Data.Key, nil
}
