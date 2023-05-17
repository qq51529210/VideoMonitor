package main

import (
	"recordplan/db"
	"recordplan/webapi"
	"recordplan/weekplan"

	"github.com/qq51529210/log"
)

//	@Title		接口文档
//	@version	1.0.0
func main() {
	defer func() {
		log.Recover(recover())
	}()
	// 配置
	err := loadCfg()
	if err != nil {
		panic(err)
	}
	// 数据库
	err = db.Init(_cfg.DB.URL)
	if err != nil {
		panic(err)
	}
	// 检查
	err = weekplan.Run(_cfg.WeekPlan.CheckInterval, _cfg.WeekPlan.Concurrency)
	if err != nil {
		panic(err)
	}
	// api 服务
	webapi.Serve(8001)
}
