package db

import (
	"gorm.io/gorm"
)

// WeekPlan 表示周计划
// Stream 和 WeekPlanID 是多对多的关系
type WeekPlanStream struct {
	// 用于查询
	Stream string `json:"stream" gorm:"primary"`
	// WeekPlan.ID
	WeekPlanID int64     `json:"weekPlanID" gorm:"primary"`
	WeekPlan   *WeekPlan `json:"-" gorm:"foreignKey:WeekPlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 回调
	Callback *string `json:"callback" gorm:"type:text"`
}

// WeekPlanStreamQuery 是 WeekPlanStream 的查询参数
type WeekPlanStreamQuery struct {
	// ID，精确匹配
	Stream string `form:"stream" binding:"omitempty,min=1"`
	// WeekPlan.ID，精确匹配
	WeekPlanID *int64 `form:"weekPlanID" binding:"omitempty,min=1"`
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
