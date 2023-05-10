package zlm

// OnRTSPAuthReq 表示 on_rtsp_auth 提交的数据
type OnRTSPAuthReq struct {
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
	// 播放用户名
	Username string `json:"user_name"`
	// rtsp播放鉴权加密realm
	Realm string `json:"realm"`
	// rtsp url参数
	Params string `json:"params"`
	// 请求的密码是否必须为明文(base64鉴权需要明文密码)
	MustNoEncrypt bool `json:"must_no_encrypt"`
}

// OnRTSPAuth 处理 zlm 的 on_rtsp_auth 回调
func OnRTSPAuth(data *OnRTSPAuthReq) (encrypted bool, passwd string, err error) {
	return false, "", nil
}
