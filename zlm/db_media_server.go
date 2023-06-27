package zlm

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	MediaServerDA *util.GORMCache[string, *MediaServer]
)

func InitMediaServerDA(db *gorm.DB, cache bool) {
	MediaServerDA = util.NewGORMCache(
		db,
		cache,
		func() *MediaServer {
			return new(MediaServer)
		},
		func(m *MediaServer) string {
			return m.ID
		},
		func(db *gorm.DB, k string) *gorm.DB {
			return db.Where("`ID` = ?", k)
		},
		func(db *gorm.DB, ks []string) *gorm.DB {
			return db.Where("`ID` IN ?", ks)
		})
}

// MediaServer 表示一个流媒体服务
type MediaServer struct {
	// id
	ID string `json:"id" gorm:"type:varchar(32);primaryKey"`
	// 访问密钥
	Secret *string `json:"secret" gorm:"type:varchar(64);not null"`
	// 名称，方便记忆
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// API 地址 (http|https)://ip:port
	APIBaseURL *string `json:"apiBaseURL" gorm:"type:varchar(64);not null"`
	// 外网访问的 ip ，生成播放地址时使用
	PublicIP *string `json:"publicIP" gorm:"type:varchar(40);not null"`
	// 内网访问的 ip ，生成播放地址时使用
	PrivateIP *string `json:"privateIP" gorm:"type:varchar(40);not null"`
	// 请求服务接口超时时间，单位，毫秒，默认 5000
	APICallTimeout *uint32 `json:"apiCallTimeout" gorm:"not null;default:5000"`
	// 是否禁用，0/1
	Enable *int8 `json:"enable" gorm:"not null;default:1"`
	// 是否在线，0/1
	Online *int8 `json:"online" gorm:"not null;default:0"`
	// 心跳时间戳
	KeepaliveTime int64 `json:"keepaliveTime" gorm:"-"`
	// 所属的分组
	MediaServerGroupID *int64            `json:"mediaServerGroupID" gorm:""`
	MediaServerGroup   *MediaServerGroup `json:"-" gorm:"foreignKey:MediaServerGroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// MediaServerQuery 是 MediaServer 的查询参数
type MediaServerQuery struct {
	util.GORMPage
	// id，模糊
	ID *string `form:"id" binding:"omitempty,max=32" gq:"like"`
	// 访问密钥，模糊
	Secret *string `form:"secret" binding:"omitempty,max=64" gq:"like"`
	// 名称，模糊
	Name *string `form:"name" binding:"omitempty,max=32" gq:"like"`
	// 描述，模糊
	Describe *string `form:"describe" binding:"omitempty,max=128" gq:"like"`
	// 是否禁用，精确
	Enable *int8 `form:"enable" binding:"omitempty,oneof=0 1" gq:"eq"`
	// 是否在线，精确
	Online *int8 `form:"describe" binding:"omitempty,oneof=0 1" gq:"eq"`
	// 所属的分组，精确
	MediaServerGroupID *int64 `form:"mediaServerGroupID" binding:"omitempty,min=1" gq:"eq"`
}

// Init 实现 Query 接口
func (q *MediaServerQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
