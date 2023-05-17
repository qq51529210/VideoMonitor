package main

import (
	"recordplan/db"
	"recordplan/webapi"

	"github.com/qq51529210/log"
)

// @Title		接口文档
// @version	1.0.0
func main() {
	defer func() {
		log.Recover(recover())
	}()
	// 配置
	loadCfg()
	// 数据库
	db.Init(_cfg.DB.URL)
	// 服务
	webapi.Serve(8001)
}
