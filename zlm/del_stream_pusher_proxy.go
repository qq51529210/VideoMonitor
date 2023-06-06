package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// DelPushStreamerProxyReq 是 DelPushStreamerProxy 参数
type DelPushStreamerProxyReq struct {
	// addPushStreamerProxy接口返回的key
	Key string `query:"key"`
}

// delPushStreamerProxyRes 用于解析 delPushStreamerProxy 返回值
type delPushStreamerProxyRes struct {
	Code int `json:"code"`
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

// DelPushStreamerProxy 调用 /index/api/delPushStreamerProxy
// 关闭拉流代理(流注册成功后，也可以使用close_streams接口替代)
func (s *Server) DelPushStreamerProxy(req *DelPushStreamerProxyReq) (bool, error) {
	var _res delPushStreamerProxyRes
	err := util.HTTP[any](http.MethodGet,
		s.url("delPushStreamerProxy"),
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
