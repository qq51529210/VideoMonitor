package zlm

// OnShellLoginReq 表示 on_shell_login 提交的数据
type OnShellLoginReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// TCP链接唯一ID
	ID string `json:"id"`
	// telnet 终端ip
	IP string `json:"ip"`
	// telnet 终端端口号
	Port int `json:"port"`
	// telnet 终端登录用户名
	Username string `json:"user_name"`
	// telnet 终端登录用户密码
	Passwd string `json:"passwd"`
}

// OnShellLogin 处理 zlm 的 on_shell_login 回调
func OnShellLogin(data *OnShellLoginReq) error {
	return nil
}
