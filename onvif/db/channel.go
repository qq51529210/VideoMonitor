package db

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	ChannelDA *util.GORMCache[int64, *Channel]
)

func initChannelDA(db *gorm.DB, cache bool) {
	ChannelDA = util.NewGORMCache(
		db,
		cache,
		func() *Channel {
			return new(Channel)
		},
		func(m *Channel) int64 {
			return m.ID
		},
		func(db *gorm.DB, k int64) *gorm.DB {
			return db.Where("`ID` = ?", k)
		},
		func(db *gorm.DB, ks []int64) *gorm.DB {
			return db.Where("`ID` IN ?", ks)
		})
}

// Channel 表示级联级联通道
type Channel struct {
	ID int64 `json:"id" gorm:""`
	// Device.ID
	DeviceID int64   `json:"deviceID" gorm:"uniqueIndex:ChannelUNI"`
	Device   *Device `json:"-" gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 名称
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// 做唯一键
	ProfileName string `json:"-" gorm:"type:varchar(64);not null;uniqueIndex:ChannelUNI"`
	// 访问的 token ，好像都是不变的
	ProfileToken string `json:"-" gorm:"type:varchar(64);not null"`
	// 媒体流的地址，不必每次都去走一遍流程，加快获取速度
	StreamURI string `json:"streamURI" gorm:"type:varchar(64);not null"`
	// 分辨率
	Resolution *string `json:"resolution" gorm:"type:varchar(32)"`
	// 视频编码
	VideoEncoding *string `json:"videoEncoding" gorm:"type:varchar(16)"`
	// 音频编码
	AudioEncoding *string `json:"audioEncoding" gorm:"type:varchar(16)"`
	// 是否支持云台控制，0/1
	PTZEnable *int8 `json:"ptzEnable" gorm:"not null;default:0"`
	// 无人观看是否保留媒体流，0/1 ，默认 0
	KeepStream *int8 `json:"keepStream" gorm:"not null;default:0"`
}

// ChannelQuery 是 Channel 的查询参数
type ChannelQuery struct {
	util.GORMPage
	// Channel.ID ，精确匹配
	DeviceID *int64 `form:"deviceID" binding:"omitempty,min=1" gq:"eq"`
	// 名称，模糊匹配
	Name string `form:"name" binding:"omitempty,max=32" gq:"like"`
	// 无人观看是否保留媒体流，0/1 ，精确
	KeepStream *int8 `json:"keepStream" binding:"omitempty,oneof=0 1" gq:"eq"`
}

// Init 实现 Query 接口
func (q *ChannelQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
