package db

import (
	"errors"

	"github.com/qq51529210/util"
	"github.com/qq51529210/video-monitor/zlm"
	"gorm.io/gorm"
)

var (
	True  = int8(1)
	False = int8(0)
)

var (
	// ErrNotFound 表示没有相关的数据
	ErrNotFound = errors.New("data not found")
)

// Init 初始化数据
func Init(uri string, cache bool) error {
	// 连接
	db, err := initDB(uri)
	if err != nil {
		return err
	}
	// 表
	err = initTable(db)
	if err != nil {
		return err
	}
	// 数据访问
	initDA(db, cache)
	// 数据
	err = initData()
	if err != nil {
		return err
	}
	//
	return nil
}

// initDB 创建连接
func initDB(uri string) (*gorm.DB, error) {
	cfg := util.NewGORMConfig()
	cfg.Logger = &util.GORMLog{}
	return util.InitGORM(uri, cfg)
}

// initTable 创建表
func initTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&zlm.MediaServerGroup{},
		&zlm.MediaServer{},
		&Discovery{},
		// &WeekPlanStream{},
	)
}

// initDA 创建各个数据访问
func initDA(db *gorm.DB, cache bool) {
	zlm.InitMediaServerGroupDA(db, false)
	zlm.InitMediaServerDA(db, false)
	initDiscoveryDA(db, false)
	initDeviceDA(db, cache)
	initChannelDA(db, cache)
}

// initData 初始化默认数据
func initData() error {
	err := zlm.InitDefaultMediaServerGroup()
	if err != nil {
		return err
	}
	//
	return nil
}
