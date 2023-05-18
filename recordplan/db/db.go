package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	gormmysql "gorm.io/driver/mysql"
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
	// mysql
	err := createMysqlSchema(uri)
	if err != nil {
		return err
	}
	_db, err = gorm.Open(gormmysql.Open(uri), &config)
	return err
}

func createMysqlSchema(uri string) error {
	// 解析出 schema
	cfg, err := mysql.ParseDSN(uri)
	if err != nil {
		return err
	}
	_uri := strings.Replace(uri, cfg.DBName, "", 1)
	db, err := sql.Open("mysql", _uri)
	if err != nil {
		return err
	}
	defer db.Close()
	// 创建
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;", cfg.DBName))
	return err
}

// initTable 创建表
func initTable() error {
	return _db.AutoMigrate(
		&WeekPlan{},
		&WeekPlanStream{},
	)
}
