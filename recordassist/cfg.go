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
	Dir string `json:"dir" yaml:"dir" validate:"required,filepath"`
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
	URL string `json:"url" yaml:"url"`
}

// cfgZLM 检查配置
type cfgZLM struct {
	// 文件根目录
	Dir string `json:"dir" yaml:"dir" validate:"omitempty,filepath"`
	// 检查间隔，最小 1 秒
	CheckInterval int `json:"checkInterval" yaml:"checkInterval" validate:"required,min=1"`
	// API 调用的超时，单位秒
	APICallTimeout int `json:"apiCallTimeout" yaml:"apiCallTimeout" validate:"required,min=1"`
	// 文件的最大时长，修改-创建
	MaxDuration int `json:"maxDuration" yaml:"maxDuration" validate:"required,min=1"`
}

// cfg 用于加载启动配置
type cfg struct {
	Name string `json:"name" yaml:"name" validate:"required"`
	Port int    `json:"port" yaml:"port" validate:"required,min=1"`
	Log  cfgLog `json:"log" yaml:"log"`
	DB   cfgDB  `json:"db" yaml:"db"`
	ZLM  cfgZLM `json:"zlm" yaml:"zlm"`
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
