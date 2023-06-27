package task

// Cfg 配置
type Cfg struct {
	// 多播地址
	DiscoveryMulticastAddr string `json:"discoveryMulticastAddr" yaml:"discoveryMulticastAddr" validate:"omitempty,udp_addr"`
	// 探测的间隔，最小 1 秒
	DiscoveryInterval int `json:"discoveryInterval" yaml:"discoveryInterval" validate:"omitempty,min=1"`
	// 获取设备时间的间隔，最小 1 秒
	UpdateDeviceKeepaliveInterval int `json:"updateDeviceKeepaliveInterval" yaml:"updateDeviceKeepaliveInterval" validate:"omitempty,min=1"`
}
