package main

import (
	"os"

	"github.com/qq51529210/video-monitor/onvif/task"

	"github.com/go-playground/validator/v10"
	"github.com/qq51529210/util"
)

var (
	_cfg cfg
)

// cfgDB 数据库的配置
type cfgDB struct {
	URL         string `json:"url" yaml:"url"`
	EnableCache bool   `json:"enableCache" yaml:"enableCache"`
}

// cfg 用于加载启动配置
type cfg struct {
	Name      string      `json:"name" yaml:"name" validate:"required"`
	Port      int         `json:"port" yaml:"port" validate:"required,min=1"`
	Log       util.LogCfg `json:"log" yaml:"log"`
	DB        cfgDB       `json:"db" yaml:"db"`
	Discovery task.Cfg    `json:"discovery" yaml:"discovery"`
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
