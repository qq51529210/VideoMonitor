package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// DelFFmpegSourceReq 是 DelFFmpegSource 参数
type DelFFmpegSourceReq struct {
	// addFFmpegSource接口返回的key
	Key string `query:"key"`
}

// delFFmpegSourceRes 用于解析 delFFmpegSource 的返回值
type delFFmpegSourceRes struct {
	Code int `json:"code"`
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

// DelFFmpegSource 调用 /index/api/delFFmpegSource
// 关闭ffmpeg拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelFFmpegSource(req *DelFFmpegSourceReq) (bool, error) {
	var _res delFFmpegSourceRes
	err := util.HTTP[any](http.MethodGet,
		s.url("delFFmpegSource"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return false, err
	}
	if _res.Code != 0 {
		return false, CodeError(_res.Code)
	}
	return _res.Data.Flag, nil
}
