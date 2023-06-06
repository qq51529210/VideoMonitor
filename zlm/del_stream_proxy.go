package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// DelStreamProxyReq 是 DelStreamProxy 参数
type DelStreamProxyReq struct {
	// addStreamProxy接口返回的key
	Key string `query:"key"`
}

// delStreamProxyRes 用于解析 delStreamProxy 返回值
type delStreamProxyRes struct {
	Code int `json:"code"`
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

// DelStreamProxy 调用 /index/api/delStreamProxy
// 关闭拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelStreamProxy(req *DelStreamProxyReq) (bool, error) {
	var _res delStreamProxyRes
	err := util.HTTP[any](http.MethodGet,
		s.url("delStreamProxy"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return false, err
	}
	if _res.Code != 0 {
		return false, CodeError(_res.Code)
	}
	return _res.Data.Flag, nil
}
