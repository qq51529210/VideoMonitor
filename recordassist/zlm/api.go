package zlm

import "time"

// Run 开始检查
func Run(dir string, interval, apiCallTimeout, maxDuration int) {
	// 参数
	_checker.dir = dir
	_checker.checkInterval = time.Duration(interval) * time.Second
	if _checker.checkInterval < 1 {
		_checker.checkInterval = time.Second * 10
	}
	_checker.apiCallTimeout = time.Duration(apiCallTimeout) * time.Second
	if _checker.apiCallTimeout < 1 {
		_checker.apiCallTimeout = time.Second
	}
	_checker.maxDuration = time.Duration(maxDuration) * time.Second
	if _checker.maxDuration < 1 {
		_checker.maxDuration = time.Second
	}
	_checker.record = make(map[string]int)
	// 启动
	go _checker.routine()
}
