package zlm

import (
	"io"
	"net/http"

	"github.com/qq51529210/util"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
	// 需要截图的url，可以是本机的，也可以是远程主机的
	URL string `query:"url"`
	// 截图失败超时时间，防止FFmpeg一直等待截图
	TimeoutSec string `query:"timeout_sec"`
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	ExpireSec string `query:"expire_sec"`
}

// GetSnap 调用 /index/api/getSnap
// 获取截图或生成实时截图并返回，jpeg格式的图片，可以在浏览器直接打开
func (s *Server) GetSnap(req *GetSnapReq, writer io.Writer) error {
	return util.HTTPTo[any](http.MethodGet,
		s.url("getSnap"),
		s.query(req),
		nil,
		writer,
		http.StatusOK,
		s.APICallTimeout)
}
