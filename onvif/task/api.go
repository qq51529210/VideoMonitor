package task

import (
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/video-monitor/onvif/db"
	"github.com/qq51529210/video-monitor/onvif/ovf"
)

// Run 开始自动探测
func Run(cfg *Cfg) error {
	dur := time.Duration(cfg.DiscoveryInterval) * time.Second
	// 不配置不探测
	if dur < 1 || cfg.DiscoveryMulticastAddr == "" {
		return nil
	}
	// 开始探测
	err := ovf.Discovery(cfg.DiscoveryMulticastAddr, dur, handleDiscoveryAddr)
	if err != nil {
		return err
	}
	dur = time.Duration(cfg.UpdateDeviceKeepaliveInterval) * time.Second
	// 检查心跳
	go updateDeviceKeepaliveRoutine(dur)
	//
	return nil
}

// handleDiscoveryAddr 是自动发现的回调，保存数据库即可
func handleDiscoveryAddr(addr string) {
	var model db.Discovery
	model.IPAddr = addr
	model.Time = time.Now().Unix()
	_, err := db.DiscoveryDA.Save(&model)
	if err != nil {
		log.Errorf("[discovery] save db error: %s", err.Error())
	}
}
