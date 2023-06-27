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
	IPAddr string `json:"ipAddr" gorm:"type:varchar(40);not null;primaryKey"`
	// 发现的时间戳
	Time int64 `json:"time" gorm:""`
}

// DiscoveryQuery 是 Discovery 的查询参数
type DiscoveryQuery struct {
	// 地址，模糊查询
	IPAddr *string `form:"ipAddr" json:"ipAddr" binding:"omitempty,max=40" gq:"like"`
}

// Init 实现 Query 接口
func (q *DiscoveryQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
