package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/qq51529210/util"
)

var (
	_cfg cfg
)

// cfgLog 日志的配置
type cfgLog struct {
	// 日志保存的根目录
	Dir string `json:"dir" yaml:"dir" validate:"required,path"`
	// 每一份日志文件的最大字节，使用 1.5/K/M/G/T 这样的字符表示。
	MaxFileSize string `json:"maxFileSize" yaml:"maxFileSize"`
	// 保存的最大天数，最小是1天。
	MaxKeepDay float64 `json:"maxKeepDay" yaml:"maxKeepDay"`
	// 同步到磁盘的时间间隔，单位，毫秒。最小是10毫秒。
	SyncInterval int `json:"syncInterval" yaml:"syncInterval" validate:"required,min=1"`
	// 是否输出到控制台，out/err
	Std string `json:"std" yaml:"std" validate:"omitempty,oneof=out err"`
	// 禁用的日志级别
	DisableLevel []string `json:"disableLevel" yaml:"disableLevel" validate:"omitempty,dive,oneof=out err"`
}

// cfgDB 数据库的配置
type cfgDB struct {
	URL         string `json:"url" yaml:"url"`
	EnableCache bool   `json:"enableCache" yaml:"enableCache"`
}

// cfgWeekPlan 周计划的配置
type cfgWeekPlan struct {
	// 检查间隔，最小 1 秒
	CheckInterval int `json:"checkInterval" yaml:"checkInterval" validate:"required,min=1"`
	// 并发检查的个数，0 表示使用 CPU 的个数
	Concurrency int `json:"concurrency" yaml:"concurrency" validate:"required,min=0"`
	// API 调用的超时，单位秒
	APICallTimeout int `json:"apiCallTimeout" yaml:"apiCallTimeout" validate:"required,min=1"`
}

// cfg 用于加载启动配置
type cfg struct {
	Name     string      `json:"name" yaml:"name" validate:"required"`
	Port     int         `json:"port" yaml:"port" validate:"required,min=1"`
	Log      cfgLog      `json:"log" yaml:"log"`
	DB       cfgDB       `json:"db" yaml:"db"`
	WeekPlan cfgWeekPlan `json:"weekPlan" yaml:"weekPlan"`
}

// loadCfg 加载配置
func loadCfg() error {
	// 加载
	path := "cfg.yaml"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	err := util.ReadCfg(path, &_cfg)
	if err != nil {
		return err
	}
	// 检查
	val := validator.New()
	err = val.Struct(_cfg)
	if err != nil {
		return err
	}
	return nil
}
