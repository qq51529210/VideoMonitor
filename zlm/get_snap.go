package zlm

import (
	"bytes"
	"net/url"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
	// 需要截图的url，可以是本机的，也可以是远程主机的
	URL string
	// 截图失败超时时间，防止FFmpeg一直等待截图
	TimeoutSec string
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	ExpireSec string
}

func (m *GetSnapReq) toQuery() url.Values {
	q := make(url.Values)
	if m.URL != "" {
		q.Set("url", m.URL)
	}
	if m.TimeoutSec != "" {
		q.Set("timeout_sec", m.TimeoutSec)
	}
	if m.ExpireSec != "" {
		q.Set("expire_sec", m.ExpireSec)
	}
	return q
}

// GetSnap 调用 /index/api/getSnap
// 获取截图或生成实时截图并返回，jpeg格式的图片，可以在浏览器直接打开
func (s *Server) GetSnap(req *GetSnapReq) (*bytes.Buffer, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	data := bytes.NewBuffer(nil)
	err := httpGet3(s, s.url("getSnap"), query, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
