package zlm

import (
	"time"
)

// OnServerKeepaliveReq 表示 on_server_keepalive 提交的数据
type OnServerKeepaliveReq struct {
	Data          *OnServerKeepaliveDataModel `json:"data"`
	MediaServerID string                      `json:"mediaServerId"`
}

// OnServerKeepaliveDataModel 是 OnServerKeepaliveReq 的 data 字段
type OnServerKeepaliveDataModel struct {
	Buffer                int `json:"Buffer"`
	BufferLikeString      int `json:"BufferLikeString"`
	BufferList            int `json:"BufferList"`
	BufferRaw             int `json:"BufferRaw"`
	Frame                 int `json:"Frame"`
	FrameImp              int `json:"FrameImp"`
	MediaSource           int `json:"MediaSource"`
	MultiMediaSourceMuxer int `json:"MultiMediaSourceMuxer"`
	RTMPPacket            int `json:"RtmpPacket"`
	RTPPacket             int `json:"RtpPacket"`
	Socket                int `json:"Socket"`
	TCPClient             int `json:"TcpClient"`
	TCPServer             int `json:"TcpServer"`
	TCPSession            int `json:"TcpSession"`
	UDPServer             int `json:"UdpServer"`
	UDPSession            int `json:"UdpSession"`
}

// OnServerKeepalive 处理 zlm 的 on_server_keepalive 回调
func (g *Group) OnServerKeepalive(req *OnServerKeepaliveReq) {
	// 服务
	ser := g.Get(req.MediaServerID)
	if !ser.IsOK() {
		return
	}
	now := time.Now()
	ser.keepalive = &now
	ser.Online = true
}
