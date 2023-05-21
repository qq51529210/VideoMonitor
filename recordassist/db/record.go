package db

type Record struct {
	// 状态
	// 0：录像完成
	// 1：已上传，未提交
	// 2：已提交
	State *int8 `json:"state" gorm:"not null;default:0"`
	// 本地路径
	Path *string `json:"path" gorm:"not null;type:varchar(255)"`
}
