package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Enable  = int8(1)
	Disable = int8(0)
)

var (
	_db *gorm.DB
)

// Init 初始化数据
func Init(uri string) error {
	// 连接
	err := initDB(uri)
	if err != nil {
		return err
	}
	// 表
	err = initTable()
	if err != nil {
		return err
	}
	// 数据
	err = initData()
	if err != nil {
		return err
	}
	//
	return nil
}

// initDB 创建连接
func initDB(uri string) error {
	// 配置
	var config gorm.Config
	config.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	}
	config.Logger = &gormLog{}
	// sqlite
	db, err := gorm.Open(sqlite.Open(uri), &config)
	if err != nil {
		return err
	}
	_db = db
	return err
}

// initTable 创建表
func initTable() error {
	return _db.AutoMigrate(
		&Record{},
	)
}

// initData 写入默认数据
func initData() error {
	return nil
}
