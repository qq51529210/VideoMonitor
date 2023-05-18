package db

var (
	AppCfg *AppConfig
)

type AppConfig struct {
	ID int64 `json:"name" gorm:"primaryKey"`
	// 用于在线修改验证
	AccessKey string `json:"-" gorm:"not null;default:'record-plan-access-key'"`
	// 是否开启缓存
	EnableCache *int8 `json:"enableCache" gorm:"not null;default:0"`
}
