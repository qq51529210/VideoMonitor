package db

import (
	"fmt"

	"gorm.io/gorm"
)

// BaseModel 基本字段
type BaseModel struct {
	// 数据库ID
	ID string `json:"id" gorm:"primaryKey;type:varchar(40)"`
	TimeModel
}

// TimeModel 基本字段
type TimeModel struct {
	// 数据库的创建时间，时间戳，
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 数据库的更新时间
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// Page 分页查询参数
type Page struct {
	// 偏移，小于 0 不匹配
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `form:"order"`
}

// Init 初始化查询条件
func (q *Page) Init(db *gorm.DB) *gorm.DB {
	// 分页
	if q.Offset != nil {
		db = db.Offset(*q.Offset)
	}
	if q.Count != nil {
		db = db.Limit(*q.Count)
	}
	// 排序
	if q.Order != "" {
		db = db.Order(q.Order)
	}
	return db
}

// Valid 返回参数是否有效
func (q *Page) Valid() bool {
	return q.Order != "" || q.Offset != nil || q.Count != nil
}

// Query 是 All 函数格式化查询参数的接口
type Query interface {
	Init(*gorm.DB) *gorm.DB
}

// All 返回列表查询结果
func All[M any](query Query) ([]*M, error) {
	var t M
	db := _db.Model(&t)
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 列表
	var models []*M
	err := db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	// 返回
	return models, nil
}

// ListData 是 List 的返回值
type ListData[T any] struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Data []T `json:"data"`
}

// List 返回列表查询结果
func List[M any](query Query, page *Page, res *ListData[M]) error {
	var t M
	db := _db.Model(&t)
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 总数
	err := db.Count(&res.Total).Error
	if err != nil {
		return err
	}
	// 列表
	db = page.Init(db)
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}

// TimeRangeQuery 用于时间范围查询
type TimeRangeQuery struct {
	// 创建时间，时间戳，范围匹配，CreatedAt >= CreatedAtAfter
	AfterCreatedAt *int64 `form:"afterCreatedAt" binding:"omitempty,min=0"`
	// 创建时间，时间戳，范围匹配，CreatedAt < CreatedBefore
	BeforeCreatedAt *int64 `form:"beforeCreatedAt" binding:"omitempty,gtfield=AfterCreatedAt"`
	// 更新时间，时间戳，范围匹配，UpdateAt >= UpdatedAtAfter
	AfterUpdatedAt *int64 `form:"afterUpdatedAt" binding:"omitempty,min=0"`
	// 更新时间，时间戳，范围匹配，UpdateAt < UpdatedAtBefore
	BeforeUpdatedAt *int64 `form:"beforeUpdatedAt" binding:"omitempty,gtfield=AfterUpdatedAt"`
}

// Init 实现 Query 接口
func (q *TimeRangeQuery) Init(db *gorm.DB) *gorm.DB {
	if q.AfterCreatedAt != nil {
		db = db.Where("`CreatedAt` >= ?", *q.AfterCreatedAt)
	}
	if q.BeforeCreatedAt != nil {
		db = db.Where("`CreatedAt` < ?", *q.BeforeCreatedAt)
	}
	if q.AfterUpdatedAt != nil {
		db = db.Where("`UpdatedAt` >= ?", *q.AfterUpdatedAt)
	}
	if q.BeforeUpdatedAt != nil {
		db = db.Where("`UpdatedAt` < ?", *q.BeforeUpdatedAt)
	}
	return db
}

// whereLike 生成 column like %name%
func whereLike(db *gorm.DB, column, value string) *gorm.DB {
	if value == "" {
		return db
	}
	return db.Where(fmt.Sprintf("`%s` LIKE '%%%s%%'", column, value))
}

// Save 添加
func Save[M any](m *M) (int64, error) {
	db := _db.Save(m)
	return db.RowsAffected, db.Error
}

// Add 添加
func Add[M any](m *M) (int64, error) {
	db := _db.Create(m)
	return db.RowsAffected, db.Error
}

// Update 根据主键更新
func Update[M any](m *M) (int64, error) {
	db := _db.Updates(m)
	return db.RowsAffected, db.Error
}

// Delete 根据主键删除
func Delete[M any](m *M) (int64, error) {
	db := _db.Delete(m)
	return db.RowsAffected, db.Error
}

// BatchDelete 根据主键批量删除
func BatchDelete[M, ID any](m *M, ids ID) (int64, error) {
	db := _db.Delete(m, ids)
	return db.RowsAffected, db.Error
}

// Get 根据主键查询
func Get[M any](m *M) (bool, error) {
	err := _db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Select 根据主键查询，可选择列
func Select[M any](m *M, c ...string) (bool, error) {
	db := _db
	if len(c) > 0 {
		db = db.Select(c)
	}
	err := db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// In 根据主键查询，where in
func In[M, ID any](ids []ID) ([]*M, error) {
	var m M
	var ms []*M
	err := _db.Model(&m).Find(&ms, ids).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}
