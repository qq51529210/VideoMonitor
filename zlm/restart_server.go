package zlm

import "net/url"

// restartServerRes 是 RestartServer 的返回值
type restartServerRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// RestartServer 调用 /index/api/restartServer
// 重启服务器,只有Daemon方式才能重启，否则是直接关闭！
func (s *Server) RestartServer() error {
	var res restartServerRes
	query := make(url.Values)
	err := httpGet(s, s.url("restartServer"), query, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return CodeError(res.Code)
	}
	return nil
}
