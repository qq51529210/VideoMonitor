package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

// StartRecordReq 是 StartRecord 的参数
type StartRecordReq struct {
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// 0：hls ，1：mp4
	Type string `query:"type"`
	// 录像文件保存自定义根目录，为空则采用配置文件设置
	CustomizedPath string `query:"customized_path"`
	// MP4 录制的切片时间大小，单位秒，为空则采用配置文件设置
	MaxSecond string `query:"max_second"`
}

// startRecordRes 是 StartRecord 的返回值
type startRecordRes struct {
	apiRes
	// 成功与否
	Result bool `json:"result"`
}

const (
	apiPathStartRecord = "startRecord"
)

// StartRecord 调用 /index/api/startRecord 开始录制，返回是否成功
func (s *Server) StartRecord(req *StartRecordReq) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.StartRecordWithContext(ctx, req)
}

// StartRecordWithContext 调用 /index/api/startRecord 开始录制，返回是否成功
func (s *Server) StartRecordWithContext(ctx context.Context, req *StartRecordReq) (bool, error) {
	var res startRecordRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathStartRecord),
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
			API:  apiPathStartRecord,
		}
	}
	return res.Result, nil
}
