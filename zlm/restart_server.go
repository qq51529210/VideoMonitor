package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// restartServerRes 是 RestartServer 的返回值
type restartServerRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// RestartServer 调用 /index/api/restartServer
// 重启服务器,只有Daemon方式才能重启，否则是直接关闭！
func (s *Server) RestartServer() error {
	var _res restartServerRes
	err := util.HTTP[any](http.MethodGet,
		s.url("restartServer"),
		s.query(nil),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return err
	}
	if _res.Code != 0 {
		return CodeError(_res.Code)
	}
	return nil
}
