package zlm

import "net/url"

// setServerConfigRes 是 SetThreadsLoad 的返回值
type setServerConfigRes struct {
	Code    int `json:"code"`
	Changed int `json:"changed"`
}

// SetServerConfig 调用 /index/api/setServerConfig
// 设置服务器配置
// 返回被修改的个数
func (s *Server) SetServerConfig(cfg *Config) (int, error) {
	query := make(url.Values)
	if cfg != nil {
		query = cfg.toQuery()
	}
	var res setServerConfigRes
	err := httpGet(s, s.url("setServerConfig"), query, &res)
	if err != nil {
		return -1, err
	}
	if res.Code != 0 {
		return -1, CodeError(res.Code)
	}
	return res.Changed, nil
}
