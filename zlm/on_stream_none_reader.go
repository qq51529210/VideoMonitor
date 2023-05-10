package zlm

// OnStreamNoneReaderReq 表示 on_stream_none_reader 提交的数据
type OnStreamNoneReaderReq struct {
	MediaServerID string `json:"mediaServerId"`
	App           string `json:"app"`
	Schema        string `json:"schema"`
	Stream        string `json:"stream"`
	VHost         string `json:"vhost"`
}

// OnStreamNoneReader 处理 zlm 的 on_stream_none_reader 回调
func OnStreamNoneReader(data *OnStreamNoneReaderReq) bool {
	// 断流
	return true
}
