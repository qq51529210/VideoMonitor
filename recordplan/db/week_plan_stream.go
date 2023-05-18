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
		func(m *WeekPlanStream) WeekPlanStreamKey {
			return m.WeekPlanStreamKey
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

// BatchAddWeekPlanStream 批量添加
func BatchAddWeekPlanStream(models []*WeekPlanStream) (int64, error) {
	// 先更新数据库
	rows, err := batchAddWeekPlanStream(models)
	if err != nil {
		return rows, err
	}
	// 缓存
	if *AppCfg.EnableCache == 1 {
		// 上锁
		weekPlanStreamCache.Lock()
		defer weekPlanStreamCache.Unlock()
		// 加载
		for i := 0; i < len(models); i++ {
			// 失败一次就算了，等下次全部加载
			err = weekPlanStreamCache.load(models[i].WeekPlanStreamKey)
			if err != nil {
				break
			}
		}
	}
	return rows, nil
}

// batchAddWeekPlanStream 批量添加
func batchAddWeekPlanStream(models []*WeekPlanStream) (int64, error) {
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

type WeekPlanStreamKey struct {
	// 用于查询，最大 128 个字符
	Stream string `json:"stream" gorm:"primaryKey;type:varchar(64)"`
	// WeekPlan.ID
	WeekPlanID string `json:"WeekPlanStreamID" gorm:"primaryKey;type:varchar(32)"`
}

// WeekPlan 表示周计划
// Stream 和 WeekPlanID 是多对多的关系
type WeekPlanStream struct {
	WeekPlanStreamKey
	WeekPlan *WeekPlan `json:"-" gorm:"foreignKey:WeekPlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 开始录像的回调，Get 方法
	StartCallback *string `json:"startCallback" gorm:"type:varchar(512)"`
	// 停止录像的回调，Get 方法
	StopCallback *string `json:"stopCallback" gorm:"type:varchar(512)"`
	// 用于版本控制
	Version int64 `gorm:"autoUpdateTime:nano"`
}

// WeekPlanStreamQuery 是 WeekPlanStream 的查询参数
type WeekPlanStreamQuery struct {
	// ID，精确匹配
	Stream string `form:"stream" binding:"omitempty,min=1"`
	// WeekPlan.ID，精确匹配
	WeekPlanID string `form:"WeekPlanStreamID" binding:"omitempty,min=1"`
}

// Init 实现 Query 接口
func (q *WeekPlanStreamQuery) Init(db *gorm.DB) *gorm.DB {
	// Stream
	if q.Stream != "" {
		db = db.Where("`Stream` = ?", q.Stream)
	}
	// WeekPlan.ID
	if q.WeekPlanID != "" {
		db = db.Where("`WeekPlanID` = ?", q.WeekPlanID)
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
