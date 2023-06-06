package db

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	DiscoveryDA *util.GORMCache[string, *Discovery]
)

func initDiscoveryDA(db *gorm.DB, cache bool) {
	DiscoveryDA = util.NewGORMCache(
		db,
		cache,
		func() *Discovery {
			return new(Discovery)
		},
		func(m *Discovery) string {
			return m.IPAddr
		},
		func(db *gorm.DB, k string) *gorm.DB {
			return db.Where("`IPAddr` = ?", k)
		},
		func(db *gorm.DB, ks []string) *gorm.DB {
			return db.Where("`IPAddr` IN ?", ks)
		})
}

// Discovery 表示自动发现的列表
type Discovery struct {
	// 地址
	IPAddr string `json:"ipAddr" gorm:"type:varchar(255);not null;primaryKey"`
	// 数据库的创建时间，时间戳，
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 数据库的更新时间
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// DiscoveryQuery 是 Discovery 的查询参数
type DiscoveryQuery struct {
	// 模糊查询
	IPAddr *string `form:"ipAddr" json:"ipAddr" binding:"omitempty,max=255" gq:"like"`
}

// Init 实现 Query 接口
func (q *DiscoveryQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
