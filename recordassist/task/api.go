package task

import (
	"time"
)

// Run 开始
func Run(cfg *Cfg) error {
	// 存储
	if err := _minio.init(&cfg.Minio); err != nil {
		return err
	}
	// 磁盘检查
	_diskChecker.dir = cfg.Dir
	_diskChecker.checkInterval = time.Duration(cfg.CheckInterval) * time.Second
	_diskChecker.maxDuration = time.Duration(cfg.MaxDuration) * time.Second
	go _diskChecker.routine()
	// 数据库检查
	_dbChecker.apiCallTimeout = time.Duration(cfg.APICallTimeout) * time.Second
	_dbChecker.apiGetSaveDaysURL = cfg.APIGetSaveDaysURL
	_dbChecker.apiPostRecordURL = cfg.APIPostRecordURL
	go _dbChecker.routine()
	//
	return nil
}
