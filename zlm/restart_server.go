package zlm

import (
	"context"
	"net/http"

	"github.com/qq51529210/util"
)

const (
	apiPathRestartServer = "restartServer"
)

// RestartServer 调用 /index/api/restartServer 重启服务器
func (s *Server) RestartServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.APICallTimeout)
	defer cancel()
	return s.RestartServerWithContext(ctx)
}

// RestartServerWithContext 调用 /index/api/restartServer 重启服务器
func (s *Server) RestartServerWithContext(ctx context.Context) error {
	var res apiRes
	err := util.HTTPWithContext[any](ctx, http.MethodGet,
		s.url(apiPathRestartServer),
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
			API:  apiPathRestartServer,
		}
	}
	return nil
}
