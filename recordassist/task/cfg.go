package task

// Cfg 配置
type Cfg struct {
	// 检查的文件根目录
	Dir string `json:"dir" yaml:"dir" validate:"omitempty,filepath"`
	// 检查间隔，最小 1 秒
	CheckInterval int `json:"checkInterval" yaml:"checkInterval" validate:"required,min=1"`
	// 使用的上传处理函数
	// 检查间隔，最小 1 秒
	Uploader string `json:"uploader" yaml:"uploader" validate:"required"`
	// 文件的最大时长，修改时间减去创建时间
	MaxDuration int `json:"maxDuration" yaml:"maxDuration" validate:"required,min=10"`
	// 计划管理，查询保存天数的 API
	APIGetSaveDaysURL string `json:"apiGetSaveDaysURL" yaml:"apiGetSaveDaysURL" validate:"required"`
	// 录像管理，提交数据 API
	APIPostRecordURL string `json:"apiPostRecordURL" yaml:"apiPostRecordURL" validate:"required"`
	// API 调用的超时，单位秒
	APICallTimeout int `json:"apiCallTimeout" yaml:"apiCallTimeout" validate:"required,min=1"`
	// minio
	Minio MinioCfg `json:"minio" yaml:"minio"`
}

// MinioCfg 表示 minio 的配置
type MinioCfg struct {
	// 地址
	Host string `json:"host" yaml:"host" validate:"required"`
	// 用户名
	ID string `json:"id" yaml:"id" validate:"required"`
	// 密码
	Secret string `json:"secret" yaml:"secret" validate:"required"`
	// 操作的桶
	Bucket string `json:"bucket" yaml:"bucket" validate:"required"`
	// 操作超时
	Timeout int `json:"timeout" yaml:"timeout" validate:"required,min=1"`
}
