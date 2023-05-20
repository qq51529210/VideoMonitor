package db

import "gorm.io/gorm"

const (
	// TimePeroid 时间段的格式
	TimePeroidFormat = "150405"
)

var (
	// 缓存
	weekPlanCache = newMapCache(
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
	// GetWeekPlan 返回指定 id 的缓存
	GetWeekPlan = weekPlanCache.Get
	// GetWeekPlanAll 返回所有缓存
	GetWeekPlanAll = weekPlanCache.All
	// AddWeekPlan 添加数据库和缓存
	AddWeekPlan = weekPlanCache.Add
	// UpdateWeekPlan 修改数据库和缓存
	UpdateWeekPlan = weekPlanCache.Update
	// DeleteWeekPlan 删除数据库和缓存
	DeleteWeekPlan = weekPlanCache.Delete
	// BatchDeleteWeekPlan 批量删除数据库和缓存
	BatchDeleteWeekPlan = weekPlanCache.BatchDelete
)

// TimePeroid 表示每一段的时间
type TimePeroid struct {
	// 开始时间戳
	Start string `json:"start"`
	// 结束时间戳
	End string `json:"end"`
}

// WeekPlan 表示周计划
type WeekPlan struct {
	BaseModel
	// 用于版本控制
	Version int64 `json:"-" gorm:"autoUpdateTime:nano"`
	// 名称
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// 是否禁用
	// 0: 禁用
	// 1: 启用
	Enable *int8 `json:"enable" gorm:"not null;default:0"`
	// 保存的天数
	SaveDays *uint32 `json:"saveDays" gorm:"not null;default:1"`
	// 是一个 RecordTime 的 JSON 数组
	Peroids *string `json:"peroids" gorm:"not null;type:text"`
}

// WeekPlanQuery 是 WeekPlan 的查询参数
type WeekPlanQuery struct {
	Page
	TimeRangeQuery
	// 名称，模糊匹配
	Name string `form:"name" binding:"omitempty,max=32"`
	// 是否禁用，精确匹配
	Enable *int8 `form:"enable" binding:"omitempty,oneof=0 1"`
	// 保存的天数
	SaveDays *int32 `form:"saveDays" binding:"omitempty,min=1"`
}

// Init 实现 Query 接口
func (q *WeekPlanQuery) Init(db *gorm.DB) *gorm.DB {
	db = q.TimeRangeQuery.Init(db)
	// 名称
	if q.Name != "" {
		db = whereLike(db, "Name", q.Name)
	}
	// 是否禁用
	if q.Enable != nil {
		db = db.Where("`Enable` = ?", *q.Enable)
	}
	// 保存的天数
	if q.SaveDays != nil {
		db = db.Where("`SaveDays` = ?", *q.SaveDays)
	}
	//
	return db
}
