package main

import "github.com/qq51529210/log"

var (
	_logger logger
)

// logger 表示日志
type logger struct {
	file *log.File
}

// initLogger 初始化日志
func initLogger() error {
	return nil
}
