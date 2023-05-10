package zlm

// OnServerStarted 处理 zlm 的 on_server_started 回调
func OnServerStarted(ip string, cfg *Config) {
	// 获取实例
	lock.Lock()
	defer lock.Unlock()
	ser := servers[cfg.GeneralMediaServerID]
	if ser == nil {
		return
	}
}
