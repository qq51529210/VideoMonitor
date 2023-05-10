package zlm

import "net/url"

// GetRTPInfoReq 是 GetRTPInfo 的参数
type GetRTPInfoReq struct {
	// RTP的ssrc，16进制字符串或者是流的id(openRtpServer接口指定)
	Stream string
}

func (m *GetRTPInfoReq) toQuery() url.Values {
	q := make(url.Values)
	if m.Stream != "" {
		q.Set("stream_id", m.Stream)
	}
	return q
}

// GetRTPInfoRes 是 GetRTPInfo 的返回值
type GetRTPInfoRes struct {
	// 是否存在
	Exist bool `json:"exist"`
	// 本地监听的网卡ip
	LocalIP string `json:"local_ip"`
	// 本机端口号
	LocalPort int `json:"local_port"`
	// 推流客户端ip
	PeerIP string `json:"peer_ip"`
	// 客户端端口号
	PeerPort int `json:"peer_port"`
}

// getRTPInfoRes 包装 GetRTPInfoRes
type getRTPInfoRes struct {
	Code int `json:"code"`
	GetRTPInfoRes
}

// GetRTPInfo 调用 /index/api/getRtpInfo
// 获取rtp代理时的某路ssrc rtp信息
func (s *Server) GetRTPInfo(req *GetRTPInfoReq) (*GetRTPInfoRes, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res getRTPInfoRes
	err := httpGet(s, s.url("getRtpInfo"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return &res.GetRTPInfoRes, nil
}
