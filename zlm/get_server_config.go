package zlm

import "net/url"

// GetServerConfigRes 是 GetThreadsLoad 的返回值
type GetServerConfigRes struct {
	Code int       `json:"code"`
	Data []*Config `json:"data"`
}

// GetServerConfig 调用 /index/api/getServerConfig
// 获取服务器配置
func (s *Server) GetServerConfig() ([]*Config, error) {
	var res GetServerConfigRes
	query := make(url.Values)
	err := httpGet(s, s.url("getServerConfig"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
