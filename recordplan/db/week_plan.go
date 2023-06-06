package db

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

const (
	// TimePeroid 时间段的格式
	TimePeroidFormat = "150405"
)

var (
	WeekPlanDA *util.GORMCache[string, *WeekPlan]
)

func initWeekPlanDA(db *gorm.DB, cache bool) {
	WeekPlanDA = util.NewGORMCache(
		db,
		cache,
		func() *WeekPlan {
			return new(WeekPlan)
		},
		func(m *WeekPlan) string {
			return m.ID
		},
		func(db *gorm.DB, k string) *gorm.DB {
			return db.Where("`ID` = ?", k)
		},
		func(db *gorm.DB, ks []string) *gorm.DB {
			return db.Where("`ID` IN ?", ks)
		})
}

// TimePeroid 表示每一段的时间
type TimePeroid struct {
	// 开始时间戳
	Start string `json:"start"`
	// 结束时间戳
	End string `json:"end"`
}

// WeekPlan 表示周计划
type WeekPlan struct {
	ID string `json:"id" gorm:""`
	// 用于版本控制
	Version int64 `json:"-" gorm:"autoUpdateTime:nano"`
	// 名称
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// 是否禁用
	// 0: 禁用
	// 1: 启用
	Enable *int8 `json:"enable" gorm:"not null;default:0"`
	// 保存的天数
	SaveDay *int16 `json:"saveDay" gorm:"not null;defalt:1"`
	// 是一个 RecordTime 的 JSON 数组
	Peroids *string `json:"peroids" gorm:"not null;type:text"`
}

// WeekPlanQuery 是 WeekPlan 的查询参数
type WeekPlanQuery struct {
	util.GORMPage
	// 名称，模糊匹配
	Name *string `form:"name" binding:"omitempty,max=32" gq:"like"`
	// 是否禁用，精确匹配
	Enable *int8 `form:"enable" binding:"omitempty,oneof=0 1" gq:"eq"`
	// 保存的天数
	SaveDay *int16 `form:"saveDay" binding:"omitempty,min=1" gq:"eq"`
}

// Init 实现 Query 接口
func (q *WeekPlanQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}

// GetWeekPlanIn 返回指定 ids
func GetWeekPlanIn(ids []string) ([]*WeekPlan, error) {
	// 缓存
	if WeekPlanDA.IsCache() {
		return WeekPlanDA.Search(func(m *WeekPlan) bool {
			for _, id := range ids {
				if m.ID == id {
					return true
				}
			}
			return false
		})
	}
	// 数据库
	var models []*WeekPlan
	err := WeekPlanDA.DB().Where("`ID` IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
