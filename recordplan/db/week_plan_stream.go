package db

import (
	"gorm.io/gorm"
)

var (
	// 缓存
	weekPlanStreamCache = newMapCache(
		func() *WeekPlanStream {
			return new(WeekPlanStream)
		},
		func(db *gorm.DB, k WeekPlanStreamKey) *gorm.DB {
			return db.Where("`Stream` = ? AND `WeekPlanID` = ?", k.Stream, k.WeekPlanID)
		},
		func(db *gorm.DB, ks []WeekPlanStreamKey) *gorm.DB {
			return db.Where("(`Stream`, `WeekPlanID`) IN ?", ks)
		})
	// GetWeekPlanStream 返回指定 id 的缓存
	GetWeekPlanStream = weekPlanStreamCache.Get
	// GetWeekPlanStreamAll 返回所有缓存
	GetWeekPlanStreamAll = weekPlanStreamCache.All
	// AddWeekPlanStream 添加数据库和缓存
	AddWeekPlanStream = weekPlanStreamCache.Add
	// UpdateWeekPlanStream 修改数据库和缓存
	UpdateWeekPlanStream = weekPlanStreamCache.Update
	// DeleteWeekPlanStream 删除数据库和缓存
	DeleteWeekPlanStream = weekPlanStreamCache.Delete
	// BatchDeleteWeekPlanStream 批量删除数据库和缓存
	BatchDeleteWeekPlanStream = weekPlanStreamCache.BatchDelete
)

type WeekPlanStreamKey struct {
	Stream     string
	WeekPlanID string
}

// WeekPlan 表示周计划
// Stream 和 WeekPlanID 是多对多的关系
type WeekPlanStream struct {
	// 用于查询，最大 128 个字符
	Stream string `json:"stream" gorm:"primary;char(128)"`
	// WeekPlan.ID
	WeekPlanID string    `json:"WeekPlanStreamID" gorm:"primaryKey;char(32)"`
	WeekPlan   *WeekPlan `json:"-" gorm:"foreignKey:WeekPlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 回调
	Callback *string `json:"callback" gorm:"type:text"`
}

// key 实现缓存的接口
func (m *WeekPlanStream) key() WeekPlanStreamKey {
	return WeekPlanStreamKey{
		Stream:     m.Stream,
		WeekPlanID: m.WeekPlanID,
	}
}

// WeekPlanStreamQuery 是 WeekPlanStream 的查询参数
type WeekPlanStreamQuery struct {
	// ID，精确匹配
	Stream string `form:"stream" binding:"omitempty,min=1"`
	// WeekPlan.ID，精确匹配
	WeekPlanID *int64 `form:"WeekPlanStreamID" binding:"omitempty,min=1"`
}

// Init 实现 Query 接口
func (q *WeekPlanStreamQuery) Init(db *gorm.DB) *gorm.DB {
	// Stream
	if q.Stream != "" {
		db = db.Where("`Stream` = ?", q.Stream)
	}
	// WeekPlan.ID
	if q.WeekPlanID != nil {
		db = db.Where("`WeekPlanID` = ?", *q.WeekPlanID)
	}
	//
	return db
}

// GetWeekPlanStreamListByPlanID 返回指定 planID 的缓存
func GetWeekPlanStreamListByPlanID(planID string) ([]*WeekPlanStream, error) {
	// 上锁
	weekPlanStreamCache.Lock()
	defer weekPlanStreamCache.Unlock()
	// 确保数据
	err := weekPlanStreamCache.check()
	if err != nil {
		return nil, err
	}
	// 列表
	var models []*WeekPlanStream
	for _, v := range weekPlanStreamCache.d {
		if v.WeekPlanID == planID {
			models = append(models, v)
		}
	}
	return models, nil
}
