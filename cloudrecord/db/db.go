package db

import (
	"github.com/qq51529210/util"

	"gorm.io/gorm"
)

var (
	Enable  = int8(1)
	Disable = int8(0)
)

// Init 初始化数据
func Init(uri string) error {
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
	initDA(db)
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
		&Record{},
	)
}

// initDA 创建各个数据访问
func initDA(db *gorm.DB) {
	RecordDA.Init(db, new(Record))
}
