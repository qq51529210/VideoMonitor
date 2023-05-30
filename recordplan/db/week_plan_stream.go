package db

import (
	"fmt"
	"strings"

	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	WeekPlanStreamDA *util.GORMCache[WeekPlanStreamKey, *WeekPlanStream]
)

func initWeekPlanStreamDA(db *gorm.DB, cache bool) {
	WeekPlanStreamDA = util.NewGORMCache(
		db,
		cache,
		func() *WeekPlanStream {
			return new(WeekPlanStream)
		},
		func(m *WeekPlanStream) WeekPlanStreamKey {
			return m.WeekPlanStreamKey
		},
		func(db *gorm.DB, k WeekPlanStreamKey) *gorm.DB {
			return db.Where("`Stream` = ? AND `WeekPlanID` = ?", k.Stream, k.WeekPlanID)
		},
		func(db *gorm.DB, ks []WeekPlanStreamKey) *gorm.DB {
			if len(ks) < 1 {
				return db
			}
			// todo 没有预编译，小心 sql 注入
			var sqlStr strings.Builder
			sqlStr.WriteByte('(')
			fmt.Fprintf(&sqlStr, "(%s,%s)", ks[0].Stream, ks[0].Stream)
			for i := 1; i < len(ks); i++ {
				fmt.Fprintf(&sqlStr, ",(%s,%s)", ks[i].Stream, ks[i].Stream)
			}
			sqlStr.WriteByte(')')
			return db.Where("(`Stream`, `WeekPlanID`) IN %s", sqlStr.String())
		})
}

type WeekPlanStreamKey struct {
	// 用于查询，最大 64 个字符
	Stream string `json:"stream" gorm:"primaryKey;type:varchar(64)"`
	// WeekPlan.ID
	WeekPlanID string `json:"weekPlanID" gorm:"primaryKey;type:varchar(32)"`
}

// WeekPlan 表示周计划媒体流
// Stream 和 WeekPlanID 是多对多的关系
type WeekPlanStream struct {
	WeekPlanStreamKey
	WeekPlan *WeekPlan `json:"-" gorm:"foreignKey:WeekPlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 用于版本控制
	Version int64 `json:"-" gorm:"autoUpdateTime:nano"`
	// 开始录像的回调，Get 方法
	StartCallback *string `json:"startCallback" gorm:"type:varchar(512)"`
	// 停止录像的回调，Get 方法
	StopCallback *string `json:"stopCallback" gorm:"type:varchar(512)"`
}

// WeekPlanStreamQuery 是 WeekPlanStream 的查询参数
type WeekPlanStreamQuery struct {
	// ID，精确匹配
	Stream string `form:"Stream" binding:"omitempty,min=1" gq:"eq=Stream"`
	// WeekPlan.ID，精确匹配
	WeekPlanID string `form:"WeekPlanStream" binding:"omitempty,min=1" gq:"eq=WeekPlanID"`
}

// Init 实现 Query 接口
func (q *WeekPlanStreamQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}

// GetWeekPlanStreamByPlanID 返回指定 planID
func GetWeekPlanStreamByPlanID(planID string) ([]*WeekPlanStream, error) {
	// 缓存
	if WeekPlanStreamDA.IsCache() {
		return WeekPlanStreamDA.Search(func(m *WeekPlanStream) bool {
			return m.WeekPlanID == planID
		})
	}
	// 数据库
	var models []*WeekPlanStream
	err := WeekPlanStreamDA.DB().Where("`WeekPlanID` = ?", planID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// GetWeekPlanStreamByStream 返回指定 stream
func GetWeekPlanStreamByStream(stream string) ([]*WeekPlanStream, error) {
	// 缓存
	if WeekPlanStreamDA.IsCache() {
		return WeekPlanStreamDA.Search(func(m *WeekPlanStream) bool {
			return m.Stream == stream
		})
	}
	// 数据库
	var models []*WeekPlanStream
	err := WeekPlanStreamDA.DB().Where("`Stream` = ?", stream).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// DeleteWeekPlanStreamByStream 删除指定 stream
func DeleteWeekPlanStreamByStream(stream string) (int64, error) {
	// 数据库
	db := WeekPlanStreamDA.Model().Delete(stream, "`Stream` = ?")
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 缓存
	WeekPlanStreamDA.BatchDeleteCache(func(m *WeekPlanStream) bool {
		return m.Stream == stream
	})
	//
	return db.RowsAffected, db.Error
}

// BatchDeleteWeekPlanStreamByStream 批量删除指定 streams
func BatchDeleteWeekPlanStreamByStream(streams []string) (int64, error) {
	// 数据库
	db := WeekPlanStreamDA.Model().Delete(streams, "`Stream` IN ?")
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 缓存
	if WeekPlanStreamDA.IsCache() {
		WeekPlanStreamDA.BatchDeleteCache(func(m *WeekPlanStream) bool {
			for i := 0; i < len(streams); i++ {
				if m.Stream == streams[i] {
					return true
				}
			}
			return false
		})
	}
	//
	return db.RowsAffected, db.Error
}
