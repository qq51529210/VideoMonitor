package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// StopRecordReq 是 StopRecord 的参数
type StopRecordReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 0 为 hls，1 为 mp4
	Type string `query:"type"`
}

// stopRecordRes 是 StopRecord 的返回值
type stopRecordRes struct {
	apiRes
	// 成功与否
	Result bool `json:"result"`
}

const (
	apiPathStopRecord = "stopRecord"
)

// StopRecord 调用 /index/api/stopRecord 停止录制流，返回是否成功
func (s *Server) StopRecord(req *StopRecordReq) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.StopRecordWithContext(ctx, req)
}

// StopRecordWithContext 调用 /index/api/stopRecord 停止录制流，返回是否成功
func (s *Server) StopRecordWithContext(ctx context.Context, req *StopRecordReq) (bool, error) {
	var res stopRecordRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathStopRecord),
		s.query(req),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathStopRecord,
		}
	}
	return res.Result, nil
}
