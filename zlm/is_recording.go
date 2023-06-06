package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// IsRecordingReq 是 IsRecording 的参数
type IsRecordingReq struct {
	// 筛选虚拟主机
	VHost string `query:"vhost"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 0为hls，1为mp4
	Type string `query:"type"`
}

// isRecordingRes 是 IsRecording 的返回值
type isRecordingRes struct {
	Code int `json:"code"`
	// false:未录制,true:正在录制
	Status bool `json:"status"`
}

// IsRecording 调用 /index/api/isRecording
// 获取流录制状态
// 返回状态
func (s *Server) IsRecording(req *IsRecordingReq) (bool, error) {
	var _res isRecordingRes
	err := util.HTTP[any](http.MethodGet,
		s.url("isRecording"),
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
	return _res.Status, nil
}
