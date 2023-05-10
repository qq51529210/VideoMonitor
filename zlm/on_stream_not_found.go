package zlm

// OnStreamNotFoundReq 表示 on_stream_not_found 提交的数据
type OnStreamNotFoundReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流应用名
	App string `json:"app"`
	// TCP链接唯一ID
	ID string `json:"id"`
	// 播放器ip
	IP string `json:"ip"`
	// 播放url参数
	Params string `json:"params"`
	// 播放器端口号
	Port int `json:"port"`
	// 播放的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流ID
	Stream string `json:"stream"`
	// 流虚拟主机
	VHost string `json:"vhost"`
}

// OnStreamNotFound 处理 zlm 的 on_stream_not_found 回调
func OnStreamNotFound(data *OnStreamNotFoundReq) {
	// todo 去拉流
}
