package task

import "time"

// Run 开始
func Run(cfg *Cfg) {
	_diskChecker.dir = cfg.Dir
	_diskChecker.checkInterval = time.Duration(cfg.CheckInterval) * time.Second
	_diskChecker.maxDuration = time.Duration(cfg.MaxDuration) * time.Second
	go _diskChecker.routine()
	//
	_dbChecker.apiCallTimeout = time.Duration(cfg.APICallTimeout) * time.Second
	_dbChecker.apiGetSaveDaysURL = cfg.APIGetSaveDaysURL
	_dbChecker.apiPostRecordURL = cfg.APIPostRecordURL
	go _dbChecker.routine()
}
