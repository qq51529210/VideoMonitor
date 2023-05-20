package db

import (
	"context"
	"time"

	"github.com/qq51529210/log"
	"gorm.io/gorm/logger"
)

// gormLog 用于接收 gorm 的日志
type gormLog struct {
}

func (lg *gormLog) LogMode(logger.LogLevel) logger.Interface {
	return lg
}

func (lg *gormLog) Info(ctx context.Context, str string, args ...interface{}) {
	log.Info(str)
}

func (lg *gormLog) Warn(ctx context.Context, str string, args ...interface{}) {
	log.Warn(str)
}

func (lg *gormLog) Error(ctx context.Context, str string, args ...interface{}) {
	log.Error(str)
}

func (lg *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	log.Debugf("%s cost %v", sql, time.Since(begin))
	//
	if err != nil {
		log.Error(err)
		return
	}
}
