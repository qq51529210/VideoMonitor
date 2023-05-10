package zlm

import "net/url"

// KickSessionReq 是 KickSession 参数
type KickSessionReq struct {
	// 客户端唯一id，可以通过getAllSession接口获取
	ID string
}

func (m *KickSessionReq) toQuery() url.Values {
	q := make(url.Values)
	if m.ID != "" {
		q.Set("id", m.ID)
	}
	return q
}

// kickSessionRes 是 KickSession 返回值
type kickSessionRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// 筛选命中的流个数
	CountHit int `json:"count_hit"`
}

// KickSession 调用 /index/api/kick_session
// 断开tcp连接，比如说可以断开rtsp、rtmp播放器等
// 返回成功的个数
func (s *Server) KickSession(req *KickSessionReq) (int, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res kickSessionRes
	err := httpGet(s, s.url("kick_session"), query, &res)
	if err != nil {
		return -1, err
	}
	if res.Code != 0 {
		return -1, CodeError(res.Code)
	}
	return res.CountHit, nil
}
