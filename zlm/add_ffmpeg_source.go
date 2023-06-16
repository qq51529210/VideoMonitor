package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// AddFFMPEGSourceReq 是 AddFFMPEGSource 的参数
type AddFFMPEGSourceReq struct {
	// 拉流地址
	SrcURL string `query:"src_url"`
	// 推流地址
	DstURL string `query:"dst_url"`
	// 推流成功超时时间
	TimeoutMS string `query:"timeout_ms"`
	// 是否开启hls录制，0/1
	EnableHLS string `query:"enable_hls"`
	// 是否开启mp4录制，0/1
	EnableMP4 string `query:"enable_mp4"`
	// 配置文件中 FFmpeg 命令参数模板key(非内容)，
	// 置空则采用默认模板：ffmpeg.cmd
	CmdKey string `query:"ffmpeg_cmd_key"`
}

// addFFMPEGSourceRes 用于解析 addFFmpegSource 的返回值
type addFFMPEGSourceRes struct {
	Error
	Data struct {
		// 唯一标识
		Key string `json:"key"`
	} `json:"data"`
}

// AddFFMPEGSource 调用 /index/api/addFFmpegSource
// 返回 key
func (s *Server) AddFFMPEGSource(req *AddFFMPEGSourceReq) (string, error) {
	var _res addFFMPEGSourceRes
	err := util.HTTP[any](http.MethodGet,
		s.url("addFFmpegSource"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return "", err
	}
	if _res.Code != 0 {
		return "", &_res.Error
	}
	return _res.Data.Key, nil
}
