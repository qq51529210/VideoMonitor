package zlm

import (
	"context"
	"errors"
	"net/http"

	"github.com/qq51529210/util"
)

var (
	errServerConfigNotFound = errors.New("media server config not found")
)

// GetServerConfigRes 是 GetThreadsLoad 的返回值
type getServerConfigRes struct {
	apiRes
	Data []*Config `json:"data"`
}

const (
	apiPathGetServerConfig = "getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig 获取服务器配置
func (s *Server) GetServerConfig() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.GetServerConfigWithContext(ctx)
}

// GetServerConfigWithContext 调用 /index/api/getServerConfig 获取服务器配置
func (s *Server) GetServerConfigWithContext(ctx context.Context) error {
	var res getServerConfigRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathGetServerConfig),
		s.query(nil),
		nil,
		&res,
		http.StatusOK)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return &Error{
			Code: res.Code,
			Msg:  res.Msg,
			ID:   s.ID,
			API:  apiPathGetServerConfig,
		}
	}
	if len(res.Data) < 1 {
		return errServerConfigNotFound
	}
	for _, cfg := range res.Data {
		if cfg.GeneralMediaServerID == s.ID {
			s.Cfg = cfg
			break
		}
	}
	//
	return nil
}
