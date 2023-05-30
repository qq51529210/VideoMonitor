package db

import (
	"errors"

	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	Enable  = int8(1)
	Disable = int8(0)
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
		&WeekPlan{},
		&WeekPlanStream{},
	)
}

// initDA 创建各个数据访问
func initDA(db *gorm.DB, cache bool) {
	initWeekPlanDA(db, cache)
	initWeekPlanStreamDA(db, cache)
}
