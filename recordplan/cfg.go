package main

// cfgLog 日志的配置
type cfgLog struct {
}

// cfg 用于加载启动配置
type cfg struct {
	Name string `json:"name" yaml:"name" validate:"required"`
	Port int    `json:"port" yaml:"port" validate:"required,min=1"`
	Log  cfgLog `json:"log" yaml:"log" validate:"required,min=1"`
}
