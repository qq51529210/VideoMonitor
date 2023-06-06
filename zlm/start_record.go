package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// StartRecordReq 是 StartRecord 的参数
type StartRecordReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 0为hls，1为mp4
	Type string `query:"type"`
	// 录像保存目录
	CustomizedPath string `query:"customized_path"`
	// mp4录像切片时间大小,单位秒，置0则采用配置项
	MaxSecond string `query:"max_second"`
}

// startRecordRes 是 StartRecord 的返回值
type startRecordRes struct {
	Code int `json:"code"`
	// 成功与否
	Result bool `json:"result"`
}

// StartRecord 调用 /index/api/startRecord
// 开始录制hls或MP4
// 返回是否成功
func (s *Server) StartRecord(req *StartRecordReq) (bool, error) {
	var _res startRecordRes
	err := util.HTTP[any](http.MethodGet,
		s.url("startRecord"),
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
	return _res.Result, nil
}
