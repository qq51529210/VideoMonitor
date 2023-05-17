package week

import (
	"bytes"
	"encoding/json"
	"recordplan/db"
	"time"
)

type peroid struct {
	// 开始时间
	Start time.Time
	// 结束时间
	End time.Time
}

type weekplan struct {
	// db.WeekPlan.ID
	ID string
	// JSON.Decode(db.WeekPlan.Peroids)
	Peroids [][]*peroid
	// 数据库版本
	Version int64
	// 是否在录像时间
	IsRecording bool
}

func (p *weekplan) init(model *db.WeekPlan) {
	p.ID = model.ID
	json.NewDecoder(bytes.NewBufferString(*model.Peroids)).Decode(&p.Peroids)
}
