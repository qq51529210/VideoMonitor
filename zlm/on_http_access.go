package zlm

// OnHTTPAccessReq 表示 on_http_access 提交的数据
type OnHTTPAccessReq struct {
	// TCP链接唯一ID
	ID string `json:"id"`
	// 客户端ip
	IP string `json:"ip"`
	// 客户端口号
	Port int `json:"port"`
	// 访问路径是文件还是目录
	IsDir bool `json:"is_dir"`
	// 推流或播放url参数
	Params string `json:"params"`
	// 请求访问的文件或目录
	Path string `json:"path"`
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
}

// OnHTTPAccess 处理 zlm 的 on_http_access 回调
func OnHTTPAccess(data *OnHTTPAccessReq) (path string, second int, err error) {
	return "", 0, nil
}
