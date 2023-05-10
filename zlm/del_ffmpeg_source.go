package zlm

import "net/url"

// DelFFmpegSourceReq 是 DelFFmpegSource 参数
type DelFFmpegSourceReq struct {
	// addFFmpegSource接口返回的key
	Key string
}

func (m *DelFFmpegSourceReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Key != "" {
		q.Set("key", m.Key)
	}
	return q
}

// DelFFmpegSourceRes 是 DelFFmpegSource 返回值
type DelFFmpegSourceRes struct {
	Code int                    `json:"code"`
	Data DelFFmpegSourceResData `json:"data"`
}

// DelFFmpegSourceResData 是 DelFFmpegSourceRes 的 Data 字段
type DelFFmpegSourceResData struct {
	// 成功与否
	Flag bool
}

// DelFFmpegSource 调用 /index/api/delFFmpegSource
// 关闭ffmpeg拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelFFmpegSource(req *DelFFmpegSourceReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res DelFFmpegSourceRes
	err := httpGet(s, s.url("delFFmpegSource"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Data.Flag, nil
}
