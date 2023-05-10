package zlm

// OnRTSPRealmReq 表示 on_rtsp_realm 提交的数据
type OnRTSPRealmReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// TCP链接唯一ID
	ID string `json:"id"`
	// 流应用名
	App string `json:"app"`
	// rtsp或rtsps
	Schema string `json:"schema"`
	// 流ID
	Stream string `json:"stream"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// rtsp播放器ip
	IP string `json:"ip"`
	// rtsp播放器端口号
	Port int `json:"port"`
	// 播放rtsp url参数
	Params string `json:"params"`
}

// OnRTSPRealm 处理 zlm 的 on_rtsp_realm 回调，返回 realm
func OnRTSPRealm(data *OnRTSPRealmReq) string {

	return ""
}
