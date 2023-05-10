package zlm

import "net/url"

// GetAllSessionReq 是 GetAllSession 参数
type GetAllSessionReq struct {
	// 筛选本机端口，例如筛选rtsp链接：554
	LocalPort string
	// 筛选客户端ip
	PeerIP string
}

func (m *GetAllSessionReq) toQuery() url.Values {
	q := make(url.Values)
	if m.LocalPort != "" {
		q.Set("local_port", m.LocalPort)
	}
	if m.PeerIP != "" {
		q.Set("peer_ip", m.PeerIP)
	}
	return q
}

// getAllSessionRes 是 GetAllSession 返回值
type getAllSessionRes struct {
	Code int                     `json:"code"`
	Data []*GetAllSessionResData `json:"data"`
}

// GetAllSessionResData 是 getAllSessionRes 的 Data 字段
type GetAllSessionResData struct {
	// 该tcp链接唯一id
	ID string `json:"id"`
	// 本机网卡ip
	LocalIP string `json:"local_ip"`
	// 本机端口号(这是个rtmp播放器或推流器)
	LocalPort int `json:"local_port"`
	// 客户端ip
	PeerIP string `json:"peer_ip"`
	// 客户端端口号
	PeerPort int `json:"peer_port"`
	// 客户端TCPSession typeid
	TypeID string `json:"typeid"`
}

// GetAllSession 调用 /index/api/getAllSession
// 获取所有TcpSession列表(获取所有tcp客户端相关信息)
func (s *Server) GetAllSession(req *GetAllSessionReq) ([]*GetAllSessionResData, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res getAllSessionRes
	err := httpGet(s, s.url("getAllSession"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
