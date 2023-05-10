package zlm

import "net/url"

// AddFFMPEGSourceReq 是 AddFFMPEGSource 的参数
type AddFFMPEGSourceReq struct {
	// 拉流地址
	SrcURL string
	// 推流地址
	DstURL string
	// 推流成功超时时间
	TimeoutMS string
	// 是否开启hls录制，0/1
	EnableHLS string
	// 是否开启mp4录制，0/1
	EnableMP4 string
	// 配置文件中 FFmpeg 命令参数模板key(非内容)，
	// 置空则采用默认模板：ffmpeg.cmd
	CmdKey string
}

func (m *AddFFMPEGSourceReq) toQuery(values url.Values) {
	if m.SrcURL != "" {
		values.Set("src_url", m.SrcURL)
	}
	if m.DstURL != "" {
		values.Set("dst_url", m.DstURL)
	}
	if m.TimeoutMS != "" {
		values.Set("timeout_ms", m.TimeoutMS)
	}
	if m.EnableHLS != "" {
		values.Set("enable_hls", m.EnableHLS)
	}
	if m.EnableMP4 != "" {
		values.Set("enable_mp4", m.EnableMP4)
	}
	if m.CmdKey != "" {
		values.Set("ffmpeg_cmd_key", m.CmdKey)
	}
}

// addFFMPEGSourceRes 是 AddFFMPEGSource返回值
type addFFMPEGSourceRes struct {
	Code int                `json:"code"`
	Data AddFFMPEGSourceRes `json:"data"`
}

// AddFFMPEGSourceRes 是 addFFMPEGSourceRes 的 Data 字段
type AddFFMPEGSourceRes struct {
	// 唯一标识
	Key string
}

// AddFFMPEGSource 调用 /index/api/addFFmpegSource
func (s *Server) AddFFMPEGSource(req *AddFFMPEGSourceReq, res *AddFFMPEGSourceRes) error {
	query := make(url.Values)
	if req != nil {
		req.toQuery(query)
	}
	err := httpGet(s, s.url("addFFmpegSource"), query, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", CodeError(res.Code)
	}
	return res.Data.Key, nil
}
