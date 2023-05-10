package zlm

import "net/url"

// DelStreamProxyReq 是 DelStreamProxy 参数
type DelStreamProxyReq struct {
	// addStreamProxy接口返回的key
	Key string
}

func (m *DelStreamProxyReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Key != "" {
		q.Set("key", m.Key)
	}
	return q
}

// delStreamProxyRes 是 DelStreamProxy 返回值
type delStreamProxyRes struct {
	Code int                   `json:"code"`
	Data DelStreamProxyResData `json:"data"`
}

// DelStreamProxyResData 是 delStreamProxyRes 的 Data 字段
type DelStreamProxyResData struct {
	// 成功与否
	Flag bool
}

// DelStreamProxy 调用 /index/api/delStreamProxy
// 关闭拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelStreamProxy(req *DelStreamProxyReq) (bool, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res delStreamProxyRes
	err := httpGet(s, s.url("delStreamProxy"), query, &res)
	if err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, CodeError(res.Code)
	}
	return res.Data.Flag, nil
}
