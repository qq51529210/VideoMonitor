package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// GetRTPInfoReq 是 GetRTPInfo 的参数
type GetRTPInfoReq struct {
	// RTP的ssrc，16进制字符串或者是流的id(openRtpServer接口指定)
	Stream string `query:"stream_id"`
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
func (s *Server) GetRTPInfo(req *GetRTPInfoReq, res *GetRTPInfoRes) error {
	var _res getRTPInfoRes
	err := util.HTTP[any](http.MethodGet,
		s.url("getRtpInfo"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return err
	}
	if _res.Code != 0 {
		return CodeError(_res.Code)
	}
	res = &_res.GetRTPInfoRes
	//
	return nil
}
