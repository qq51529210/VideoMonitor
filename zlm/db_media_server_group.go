package zlm

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

// DefaultMediaServerGroupID 默认的服务分组
const DefaultMediaServerGroupID = 1

// InitDefaultMediaServerGroup 添加默认的流媒体服务分组，ID = 1
func InitDefaultMediaServerGroup() error {
	model, err := MediaServerGroupDA.Get(DefaultMediaServerGroupID)
	if err != nil {
		return err
	}
	if model != nil {
		return nil
	}
	model = new(MediaServerGroup)
	name := "默认分组"
	model.Name = &name
	describe := "默认生成的流媒体分组"
	model.Describe = &describe
	_, err = MediaServerGroupDA.Add(model)
	return err
}

var (
	MediaServerGroupDA *util.GORMCache[int64, *MediaServerGroup]
)

func InitMediaServerGroupDA(db *gorm.DB, cache bool) {
	MediaServerGroupDA = util.NewGORMCache(
		db,
		cache,
		func() *MediaServerGroup {
			return new(MediaServerGroup)
		},
		func(m *MediaServerGroup) int64 {
			return m.ID
		},
		func(db *gorm.DB, k int64) *gorm.DB {
			return db.Where("`ID` = ?", k)
		},
		func(db *gorm.DB, ks []int64) *gorm.DB {
			return db.Where("`ID` IN ?", ks)
		})
}

// MediaServerGroup 表示流媒体服务分组
type MediaServerGroup struct {
	ID int64 `json:"id" gorm:""`
	// 名称
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// 描述
	Describe *string `json:"describe" gorm:"type:varchar(128)"`
}

// MediaServerGroupQuery 是 MediaServerGroup 的查询参数
type MediaServerGroupQuery struct {
	util.GORMPage
	// 名称，模糊
	Name *string `form:"name" binding:"omitempty,max=32" gq:"like"`
	// 描述，模糊
	Describe *string `form:"describe" binding:"omitempty,max=128" gq:"like"`
}

// Init 实现 Query 接口
func (q *MediaServerGroupQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
