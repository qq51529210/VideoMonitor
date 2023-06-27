package db

import (
	"time"

	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	DeviceDA *util.GORMCache[int64, *Device]
)

func initDeviceDA(db *gorm.DB, cache bool) {
	DeviceDA = util.NewGORMCache(
		db,
		cache,
		func() *Device {
			return new(Device)
		},
		func(m *Device) int64 {
			return m.ID
		},
		func(db *gorm.DB, k int64) *gorm.DB {
			return db.Where("`ID` = ?", k)
		},
		func(db *gorm.DB, ks []int64) *gorm.DB {
			return db.Where("`ID` IN ?", ks)
		})
}

// Device 表示级联级联通道
type Device struct {
	ID int64 `json:"id" gorm:""`
	// 名称
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// 地址，ip:port
	IPAddr *string `json:"ipAddr" gorm:"type:varchar(40)"`
	// 账号
	Username *string `json:"username" gorm:"type:varchar(32)"`
	// 密码
	Password *string `json:"password" gorm:"type:char(64)"`
	// 厂商
	Manufacturer *string `json:"manufacturer" gorm:"type:char(64)"`
	// 型号
	Model *string `json:"model" gorm:"type:char(64)"`
	// 固件
	FirmwareVersion *string `json:"firmwareVersion" gorm:"type:char(64)"`
	// 序列号
	SerialNumber *string `json:"serialNumber" gorm:"type:char(64)"`
	// 硬件编号
	HardwareID *string `json:"hardwareID" gorm:"type:char(64)"`
	// 是否在线，0/1 ，默认 0
	Online *int8 `json:"online" gorm:"not null;default:0"`
	// 时间
	Time time.Duration `json:"-" gorm:"-"`
}

// DeviceQuery 是 Device 的查询参数
type DeviceQuery struct {
	util.GORMPage
	// 名称，模糊
	Name *string `form:"name" binding:"omitempty,max=32" gq:"like"`
	// 地址，模糊
	IPAddr *string `form:"ipAddr" json:"ipAddr" binding:"omitempty,max=40" gq:"like"`
	// 账号，模糊
	Username *string `form:"username" json:"username" binding:"omitempty,max=32" gq:"like"`
	// 密码，模糊
	Password *string `form:"password" json:"password" binding:"omitempty,max=40" gq:"like"`
	// 厂商，模糊
	Manufacturer *string `json:"manufacturer" binding:"omitempty,max=64" gq:"like"`
	// 型号，模糊
	Model *string `json:"model" binding:"omitempty,max=64" gq:"like"`
	// 固件
	FirmwareVersion *string `json:"firmwareVersion" binding:"omitempty,max=64" gq:"like"`
	// 序列号
	SerialNumber *string `json:"serialNumber" binding:"omitempty,max=64" gq:"like"`
	// 硬件编号
	HardwareID *string `json:"hardwareID" binding:"omitempty,max=64" gq:"like"`
	// 是否在线，0/1 ，精确
	Online *int8 `json:"online" binding:"omitempty,oneof=0 1" gq:"eq"`
}

// Init 实现 Query 接口
func (q *DeviceQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}

// UpdateDeviceOnline 更新 online 字段
func UpdateDeviceOnline(id int64, online int8) error {
	// 数据库
	db := DeviceDA.Model().
		Where("`ID` = ?", id).
		UpdateColumn("Online", online)
	if db.Error != nil {
		return db.Error
	}
	// 缓存
	if DeviceDA.IsCache() {
		DeviceDA.Lock()
		model := DeviceDA.D[id]
		if model != nil {
			model.Online = &online
		}
		DeviceDA.Unlock()
	}
	//
	return nil
}
