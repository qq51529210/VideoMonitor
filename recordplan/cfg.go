package main

import (
	"github.com/go-playground/validator/v10"
)

var (
	_cfg cfg
)

// cfgLog 日志的配置
type cfgLog struct {
}

// cfgDB 数据库的配置
type cfgDB struct {
	URL string `json:"url" yaml:"url"`
}

// cfgWeekPlan 周计划的配置
type cfgWeekPlan struct {
	// 检查间隔，最小 1 秒
	CheckInterval int `json:"checkInterval" yaml:"checkInterval" validate:"required,min=1"`
	// 并发检查的个数，0 表示使用 CPU 的个数
	Concurrency int `json:"concurrency" yaml:"concurrency" validate:"required,min=0"`
}

// cfg 用于加载启动配置
type cfg struct {
	Name     string      `json:"name" yaml:"name" validate:"required"`
	Port     int         `json:"port" yaml:"port" validate:"required,min=1"`
	Log      cfgLog      `json:"log" yaml:"log"`
	DB       cfgDB       `json:"db" yaml:"db"`
	WeekPlan cfgWeekPlan `json:"weekPlan" yaml:"weekPlan"`
}

func loadCfg() {
	// 加载
	// 检查
	val := validator.New()
	err := val.Struct(_cfg)
	if err != nil {
		panic(err)
	}
}
