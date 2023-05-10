package main

import (
	"flag"
	"os"
	"recordplan/webapi"

	"github.com/qq51529210/log"
)

// @Title		接口文档
// @version	1.0.0
func main() {
	// 启动参数
	port := flag.Int("port", 8010, "http listen port")
	flag.Parse()
	// 后台服务
	go webapi.Ser.Serve(*port)
	// 等待退出
	quit := make(chan os.Signal, 1)
	s := <-quit
	log.Warn(s.String())
}
