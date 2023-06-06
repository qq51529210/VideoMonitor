package zlm

// OnRecordMP4Req 表示 on_rtsp_auth 提交的数据
type OnRecordMP4Req struct {
	// 服务器id，通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// rtsp或rtsps
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 文件名
	FileName string `json:"file_name"`
	// 文件绝对路径
	FilePath string `json:"file_path"`
	// 文件大小，单位字节
	FileSize int `json:"file_size"`
	// 文件所在目录路径
	Folder string `json:"folder"`
	// 开始录制时间戳
	StartTime int `json:"start_time"`
	// 录制时长，单位秒
	TimeLen int `json:"time_len"`
	// http/rtsp/rtmp点播相对url路径
	URL string `json:"url"`
}

// // OnRecordMP4 处理 zlm 的 on_rtsp_auth 回调
// func OnRecordMP4(data *OnRecordMP4Req) {
// 	lock.Lock()
// 	ser := servers[data.MediaServerID]
// 	lock.Unlock()
// 	if ser == nil {
// 		return
// 	}
// }
