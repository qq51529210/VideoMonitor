package zlm

import (
	"errors"
	"net/http"

	"github.com/qq51529210/util"
)

var (
	errServerConfigNotFound = errors.New("media server config not found")
)

// GetServerConfigRes 是 GetThreadsLoad 的返回值
type getServerConfigRes struct {
	Code int       `json:"code"`
	Data []*Config `json:"data"`
}

// GetServerConfig 调用 /index/api/getServerConfig
// 获取服务器配置
func (s *Server) GetServerConfig() error {
	var _res getServerConfigRes
	err := util.HTTP[any](http.MethodGet,
		s.url("getServerConfig"),
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
	if len(_res.Data) < 1 {
		return errServerConfigNotFound
	}
	for _, cfg := range _res.Data {
		if cfg.GeneralMediaServerID == s.ID {
			s.Cfg = cfg
			break
		}
	}
	//
	return nil
}
