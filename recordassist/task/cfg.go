package task

// Cfg 配置
type Cfg struct {
	// 检查的文件根目录
	Dir string `json:"dir" yaml:"dir" validate:"omitempty,filepath"`
	// 检查间隔，最小 1 秒
	CheckInterval int `json:"checkInterval" yaml:"checkInterval" validate:"required,min=1"`
	// 文件的最大时长，修改时间减去创建时间
	MaxDuration int `json:"maxDuration" yaml:"maxDuration" validate:"required,min=10"`
	// 计划管理，查询保存天数的 API
	APIGetSaveDaysURL string `json:"apiGetSaveDaysURL" yaml:"apiGetSaveDaysURL" validate:"required"`
	// 录像管理，提交数据 API
	APIPostRecordURL string `json:"apiPostRecordURL" yaml:"apiPostRecordURL" validate:"required"`
	// API 调用的超时，单位秒
	APICallTimeout int `json:"apiCallTimeout" yaml:"apiCallTimeout" validate:"required,min=1"`
}
