package zlm

// OnPlayReq 表示 on_play 提交的数据
type OnPlayReq struct {
	// 流应用名
	App string `json:"app"`
	// TCP链接唯一ID
	ID string `json:"id"`
	// 推流器ip
	IP string `json:"ip"`
	// 推流url参数
	Params string `json:"params"`
	// 推流器端口号
	Port int `json:"port"`
	// 推流的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流ID
	Stream string `json:"stream"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
}

// OnPlay 处理 zlm 的 on_play 回调
func OnPlay(data *OnPlayReq) error {
	return nil
}
