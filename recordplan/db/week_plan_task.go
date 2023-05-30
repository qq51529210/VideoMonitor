package db

import (
	"fmt"
	"strings"

	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	WeekPlanTaskDA *util.GORMCache[WeekPlanTaskKey, *WeekPlanTask]
)

func initWeekPlanTaskDA(db *gorm.DB, cache bool) {
	WeekPlanTaskDA = util.NewGORMCache(
		db,
		cache,
		func() *WeekPlanTask {
			return new(WeekPlanTask)
		},
		func(m *WeekPlanTask) WeekPlanTaskKey {
			return m.WeekPlanTaskKey
		},
		func(db *gorm.DB, k WeekPlanTaskKey) *gorm.DB {
			return db.Where("`TaskID` = ? AND `WeekPlanID` = ?", k.TaskID, k.WeekPlanID)
		},
		func(db *gorm.DB, ks []WeekPlanTaskKey) *gorm.DB {
			if len(ks) < 1 {
				return db
			}
			// todo 没有预编译，小心 sql 注入
			var sqlStr strings.Builder
			sqlStr.WriteByte('(')
			fmt.Fprintf(&sqlStr, "(%s,%s)", ks[0].TaskID, ks[0].TaskID)
			for i := 1; i < len(ks); i++ {
				fmt.Fprintf(&sqlStr, ",(%s,%s)", ks[i].TaskID, ks[i].TaskID)
			}
			sqlStr.WriteByte(')')
			return db.Where("(`TaskID`, `WeekPlanID`) IN %s", sqlStr.String())
		})
}

type WeekPlanTaskKey struct {
	// 用于查询，最大 128 个字符
	TaskID string `json:"taskID" gorm:"primaryKey;type:varchar(64)"`
	// WeekPlan.ID
	WeekPlanID string `json:"weekPlanID" gorm:"primaryKey;type:varchar(32)"`
}

// WeekPlan 表示周计划任务
// TaskID 和 WeekPlanID 是多对多的关系
type WeekPlanTask struct {
	WeekPlanTaskKey
	WeekPlan *WeekPlan `json:"-" gorm:"foreignKey:WeekPlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 用于版本控制
	Version int64 `json:"-" gorm:"autoUpdateTime:nano"`
	// 开始录像的回调，Get 方法
	StartCallback *string `json:"startCallback" gorm:"type:varchar(512)"`
	// 停止录像的回调，Get 方法
	StopCallback *string `json:"stopCallback" gorm:"type:varchar(512)"`
}

// WeekPlanTaskQuery 是 WeekPlanTask 的查询参数
type WeekPlanTaskQuery struct {
	// ID，精确匹配
	TaskID string `form:"TaskID" binding:"omitempty,min=1" gq:"eq=TaskID"`
	// WeekPlan.ID，精确匹配
	WeekPlanID string `form:"WeekPlanTaskID" binding:"omitempty,min=1" gq:"eq=WeekPlanID"`
}

// Init 实现 Query 接口
func (q *WeekPlanTaskQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}

// GetWeekPlanTaskByPlanID 返回指定 planID
func GetWeekPlanTaskByPlanID(planID string) ([]*WeekPlanTask, error) {
	// 缓存
	if WeekPlanTaskDA.IsCache() {
		return WeekPlanTaskDA.Search(func(m *WeekPlanTask) bool {
			return m.WeekPlanID == planID
		})
	}
	// 数据库
	var models []*WeekPlanTask
	err := WeekPlanTaskDA.DB().Where("`WeekPlanID` = ?", planID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// GetWeekPlanTaskByTaskID 返回指定 taskID
func GetWeekPlanTaskByTaskID(taskID string) ([]*WeekPlanTask, error) {
	// 缓存
	if WeekPlanTaskDA.IsCache() {
		return WeekPlanTaskDA.Search(func(m *WeekPlanTask) bool {
			return m.TaskID == taskID
		})
	}
	// 数据库
	var models []*WeekPlanTask
	err := WeekPlanTaskDA.DB().Where("`TaskID` = ?", taskID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// DeleteWeekPlanTaskByTaskID 删除指定 taskID
func DeleteWeekPlanTaskByTaskID(taskID string) (int64, error) {
	// 数据库
	db := WeekPlanTaskDA.Model().Delete(taskID, "`TaskID` = ?")
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 缓存
	if WeekPlanTaskDA.IsCache() {
		// 上锁
		WeekPlanTaskDA.Lock()
		defer WeekPlanTaskDA.Unlock()
		for k := range WeekPlanTaskDA.D {
			if k.TaskID == taskID {
				delete(WeekPlanTaskDA.D, k)
			}
		}
	}
	//
	return db.RowsAffected, db.Error
}
