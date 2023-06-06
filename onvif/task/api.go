package task

import (
	"onvif/discovery"
	"time"
)

// Run 开始自动探测
func Run(cfg *Cfg) error {
	dur := time.Duration(cfg.Interval) * time.Second
	// 不配置不探测
	if dur < 1 || cfg.MulticastAddr == "" {
		return nil
	}
	return discovery.Discovery(cfg.MulticastAddr, dur, handleDiscoveryAddr)
}

func handleDiscoveryAddr(addr string) {

}
