package db

var (
	recordModel = new(Record)
)

const (
	RecordStatusReady     = "ready"
	RecordStatusUploaded  = "uploaded"
	RecordStatusSubmitted = "submitted"
)

// Record 表示一个录像
type Record struct {
	// 本地路径
	// ../live/obs/2019-09-20/15-53-02.mp4
	Path string `gorm:"primaryKey;type:varchar(255)"`
	// 创建时间
	Time int64 `gorm:""`
	// 时长
	Duration float64 ` gorm:""`
	// 大小
	Size int64 `gorm:""`
	// 存储的 key
	Name string `gorm:"type:varchar(40)"`
	// app dir
	App string `gorm:"type:varchar(64)"`
	// stream dir
	Stream string `gorm:"type:varchar(64)"`
	// 状态
	// ready 录像完成
	// uploaded 已上传
	// submitted 已提交
	Status string `gorm:"not null"`
}

// GetRecordList 获取指定数量的数据
func GetRecordList(count int) ([]*Record, error) {
	var models []*Record
	err := _db.Limit(count).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// UpdateRecordState 更新 State 字段
func UpdateRecordState(m *Record) error {
	return _db.Model(recordModel).
		Where("`Path` = ?", m.Path).
		UpdateColumn("Status", m.Status).Error
}

// UpdateRecord 更新
func UpdateRecord(path string, m map[string]any) error {
	return _db.Model(recordModel).
		Where("`Path` = ?", path).
		UpdateColumns(m).Error
}
