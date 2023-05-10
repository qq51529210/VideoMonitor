package zlm

import (
	"sync/atomic"
	"time"
)

// OnServerKeepaliveReq 表示 on_server_keepalive 提交的数据
type OnServerKeepaliveReq struct {
	Data          *OnServerKeepaliveDataModel `json:"data"`
	MediaServerID string                      `json:"mediaServerId"`
}

// OnServerKeepaliveDataModel 是 OnServerKeepaliveReq 的 data 字段
type OnServerKeepaliveDataModel struct {
	Buffer                int
	BufferLikeString      int
	BufferList            int
	BufferRaw             int
	Frame                 int
	FrameImp              int
	MediaSource           int
	MultiMediaSourceMuxer int
	RtmpPacket            int
	RtpPacket             int
	Socket                int
	TCPClient             int `json:"TcpClient"`
	TCPServer             int `json:"TcpServer"`
	TCPSession            int `json:"TcpSession"`
	UDPServer             int `json:"UdpServer"`
	UDPSession            int `json:"UdpSession"`
}

// OnServerKeepalive 处理 zlm 的 on_server_keepalive 回调
func OnServerKeepalive(ip string, data *OnServerKeepaliveReq) {
	// 获取实例
	lock.RLock()
	defer lock.RUnlock()
	ser := servers[data.MediaServerID]
	if ser == nil {
		return
	}
	now := time.Now()
	// 如果是超时的，那么重新加载一下数据
	if atomic.LoadInt32(&ser.ok) == 0 && now.Sub(*ser.keepaliveTime) > ser.keepaliveTimeout {
		wg.Add(1)
		go ser.loadDataRoutine()
		return
	}
	ser.keepaliveTime = &now
	if data.Data != nil {
		atomic.StoreInt32(&ser.load, int32(data.Data.TCPSession))
	}
	atomic.StoreInt32(&ser.offline, 0)
}
