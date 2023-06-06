package task

// Cfg 配置
type Cfg struct {
	// 多播地址
	MulticastAddr string `json:"multicastAddr" yaml:"multicastAddr" validate:"omitempty,udp_addr"`
	// 探测间隔，最小 1 秒
	Interval int `json:"interval" yaml:"interval" validate:"omitempty,min=1"`
}
