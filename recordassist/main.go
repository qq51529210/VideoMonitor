package main

import (
	"recordassist/api"
	"recordassist/zlm"

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
	// 日志
	err = initLogger()
	if err != nil {
		panic(err)
	}
	// 启动
	zlm.Run(_cfg.ZLM.Dir,
		_cfg.ZLM.CheckInterval,
		_cfg.ZLM.APICallTimeout,
		_cfg.ZLM.MaxDuration)
	// api 服务
	api.Serve(8001)
}
