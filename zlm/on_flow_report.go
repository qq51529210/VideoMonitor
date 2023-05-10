package zlm

// OnFlowReportReq 表示 on_flow_report 提交的数据
type OnFlowReportReq struct {
	// 流应用名
	App string `json:"app"`
	// tcp链接维持时间，单位秒
	Duration int `json:"duration"`
	// 推流或播放url参数
	Params string `json:"params"`
	// true为播放器，false为推流器
	Player bool `json:"player"`
	// 播放或推流的协议，可能是rtsp、rtmp、http
	Schema string `json:"schema"`
	// 流ID
	Stream string `json:"stream"`
	// 耗费上下行流量总和，单位字节
	TotalBytes int `json:"totalBytes"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 客户端ip
	IP string `json:"ip"`
	// 客户端口号
	Port int `json:"port"`
	// TCP链接唯一ID
	ID string `json:"id"`
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
}

// OnFlowReport 处理 zlm 的 on_flow_report 回调
func OnFlowReport(data *OnFlowReportReq) {

}
