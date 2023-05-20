package db

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	// 缓存
	weekPlanTaskCache = newMapCache(
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
	// GetWeekPlanTask 返回指定 id 的缓存
	GetWeekPlanTask = weekPlanTaskCache.Get
	// GetWeekPlanTaskAll 返回所有缓存
	GetWeekPlanTaskAll = weekPlanTaskCache.All
	// AddWeekPlanTask 添加数据库和缓存
	AddWeekPlanTask = weekPlanTaskCache.Add
	// UpdateWeekPlanTask 修改数据库和缓存
	UpdateWeekPlanTask = weekPlanTaskCache.Update
	// DeleteWeekPlanTask 删除数据库和缓存
	DeleteWeekPlanTask = weekPlanTaskCache.Delete
	// BatchDeleteWeekPlanTask 批量删除数据库和缓存
	BatchDeleteWeekPlanTask = weekPlanTaskCache.BatchDelete
)

// BatchAddWeekPlanTask 批量添加
func BatchAddWeekPlanTask(models []*WeekPlanTask) (int64, error) {
	// 先更新数据库
	rows, err := batchAddWeekPlanTask(models)
	if err != nil {
		return rows, err
	}
	// 缓存
	if enableCache {
		// 上锁
		weekPlanTaskCache.Lock()
		defer weekPlanTaskCache.Unlock()
		// 加载
		for i := 0; i < len(models); i++ {
			// 失败一次就算了，等下次全部加载
			err = weekPlanTaskCache.load(models[i].WeekPlanTaskKey)
			if err != nil {
				break
			}
		}
	}
	return rows, nil
}

// batchAddWeekPlanTask 批量添加
func batchAddWeekPlanTask(models []*WeekPlanTask) (int64, error) {
	rows := int64(0)
	err := _db.Transaction(func(tx *gorm.DB) error {
		db := tx.Save(models)
		if db.Error != nil {
			return db.Error
		}
		rows = db.RowsAffected
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rows, nil
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
	TaskID string `form:"TaskID" binding:"omitempty,min=1"`
	// WeekPlan.ID，精确匹配
	WeekPlanID string `form:"WeekPlanTaskID" binding:"omitempty,min=1"`
}

// Init 实现 Query 接口
func (q *WeekPlanTaskQuery) Init(db *gorm.DB) *gorm.DB {
	// TaskID
	if q.TaskID != "" {
		db = db.Where("`TaskID` = ?", q.TaskID)
	}
	// WeekPlan.ID
	if q.WeekPlanID != "" {
		db = db.Where("`WeekPlanID` = ?", q.WeekPlanID)
	}
	//
	return db
}

// GetWeekPlanTaskListByPlanID 返回指定 planID 的缓存
func GetWeekPlanTaskListByPlanID(planID string) ([]*WeekPlanTask, error) {
	// 上锁
	weekPlanTaskCache.Lock()
	defer weekPlanTaskCache.Unlock()
	// 确保数据
	err := weekPlanTaskCache.check()
	if err != nil {
		return nil, err
	}
	// 列表
	var models []*WeekPlanTask
	for _, v := range weekPlanTaskCache.d {
		if v.WeekPlanID == planID {
			models = append(models, v)
		}
	}
	return models, nil
}
