package zlm

import "net/url"

// DelPushStreamerProxyReq 是 DelPushStreamerProxy 参数
type DelPushStreamerProxyReq struct {
	// addPushStreamerProxy接口返回的key
	Key string
}

func (m *DelPushStreamerProxyReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Key != "" {
		q.Set("key", m.Key)
	}
	return q
}

// delPushStreamerProxyRes 是 DelPushStreamerProxy 返回值
type delPushStreamerProxyRes struct {
	Code int                         `json:"code"`
	Data DelPushStreamerProxyResData `json:"data"`
}

// DelPushStreamerProxyResData 是 delPushStreamerProxyRes 的 Data 字段
type DelPushStreamerProxyResData struct {
	// 成功与否
	Flag bool
}

// DelPushStreamerProxy 调用 /index/api/delPushStreamerProxy
// 关闭拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelPushStreamerProxy(req *DelPushStreamerProxyReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res delPushStreamerProxyRes
	err := httpGet(s, s.url("delPushStreamerProxy"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Data.Flag, nil
}
