package db

import (
	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	record util.GORMDB[string, *Record]
	//
	AddRecord         = record.Add
	DeleteRecord      = record.Delete
	BatchDeleteRecord = record.BatchDelete
	GetRecord         = record.Get
	GetRecordList     = record.List
)

// Record 表示录像信息
type Record struct {
	// 名称
	ID string `json:"id" gorm:"type:varchar(40);primaryKey"`
	// stream
	Stream string `json:"stream" gorm:"type:varchar(64)"`
	// 大小，字节
	Size int64 `json:"size" gorm:""`
	// 时长，单位秒
	Duration float64 ` json:"duration" gorm:""`
	// 时间戳
	Time int64 `json:"createTime" gorm:""`
	// 删除时间戳
	DeleteTime int64 `json:"deleteTime" gorm:""`
	// 是否在录像时间内的文件
	IsRecording *int8 `json:"isRecording" gorm:""`
	// 软删除
	IsDeleted *int8 `json:"isDeleted" gorm:"not null;default:0"`
	// app
	// App string `json:"app" gorm:"type:varchar(64)"`
}

// RecordQuery 实现接口 util.GORMQuery
type RecordQuery struct {
	util.GORMPage
	// stream ，精确
	Stream string `json:"stream" form:"stream" binding:"omitempty,max=64" gq:"eq"`
	// 大于创建时间戳，比较
	AfterTime *int64 `json:"afterTime" form:"afterTime" binding:"omitempty" gq:"gte=Time"`
	// 小于删除时间戳 ，比较
	BeforeTime *int64 `json:"beforeTime" form:"beforeTime" binding:"omitempty,gtfield=AfterTime" gq:"lte=Time"`
	// 创建时间戳 ，精确
	// Time *int64 `json:"createTime" form:"createTime" binding:"omitempty,min=0" gq:"eq=Time"`
	// app ，精确
	// App string `json:"app" form:"app" binding:"omitempty,max=64" gq:"eq"`
	// 大小，字节 ，精确
	// Size *int64 `json:"size" form:"size" binding:"omitempty,min=0" gq:"eq"`
	// 时长，单位秒 ，精确
	// Duration *float64 ` json:"duration" form:"duration" binding:"omitempty,min=0" gq:"eq"`
	// 删除时间戳 ，精确
	// DeleteTime *int64 `json:"deleteTime" form:"deleteTime" binding:"omitempty,min=0" gq:"eq"`
	// 是否在录像时间内的文件 ，精确
	// IsRecording *int8 `json:"isRecording" form:"isRecording" binding:"omitempty,oneof=0 1" gq:"eq"`
	// 软删除 ，精确
	// IsDeleted *int8 `json:"isDeleted" form:"isDeleted" binding:"omitempty,oneof=0 1" gq:"eq"`
}

// Init 实现接口 util.GORMQuery
func (q *RecordQuery) Init(db *gorm.DB) *gorm.DB {
	return util.GORMInitQuery(db, q)
}
