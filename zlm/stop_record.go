package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// StopRecordReq 是 StopRecord 的参数
type StopRecordReq struct {
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 0为hls，1为mp4
	Type string `query:"type"`
}

// stopRecordRes 是 StopRecord 的返回值
type stopRecordRes struct {
	Code int `json:"code"`
	// 成功与否
	Result bool `json:"result"`
}

// StopRecord 调用 /index/api/stopRecord
// 停止录制流
// 返回是否成功
func (s *Server) StopRecord(req *StopRecordReq) (bool, error) {
	var _res stopRecordRes
	err := util.HTTP[any](http.MethodGet,
		s.url("stopRecord"),
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
