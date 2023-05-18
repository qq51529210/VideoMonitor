package db

var (
	AppCfg AppConfig
)

// initAppConfig 初始化配置，ID = 1
func initAppConfig() error {
	AppCfg.ID = 1
	// 查询
	ok, err := Get(&AppCfg)
	if err != nil {
		return err
	}
	// 存在
	if ok {
		return nil
	}
	// 添加
	err = _db.Create(&AppCfg).Error
	if err != nil {
		return err
	}
	// 再查询
	_, err = Get(&AppCfg)
	if err != nil {
		return err
	}
	//
	return nil
}

type AppConfig struct {
	ID int64 `json:"name" gorm:"primaryKey"`
	// 用于在线修改验证
	AccessKey string `json:"-" gorm:"not null;default:'record-plan-access-key'"`
	// 是否开启缓存
	EnableCache *int8 `json:"enableCache" gorm:"not null;default:0"`
}
