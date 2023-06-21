package zlm

// OnPlayReq 表示 on_play 提交的数据
type OnPlayReq struct {
	// 服务器id
	MediaServerID string `json:"mediaServerId"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 协议
	Schema string `json:"schema"`
	// 应用名，例如 live
	App string `query:"app"`
	// 流 id，例如 obs
	Stream string `query:"stream"`
	// url 参数
	Params string `json:"params"`
}

// OnPlay 处理 zlm 的 on_play 回调
func OnPlay(data *OnPlayReq) error {
	return nil
}
