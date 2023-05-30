package week

import (
	"bytes"
	"encoding/json"
	"recordplan/db"
	"time"

	"github.com/qq51529210/log"
)

type peroid struct {
	// 开始时间
	Start time.Time
	// 结束时间
	End time.Time
}

// weekplan 表示计划
type weekplan struct {
	// db.WeekPlan.ID
	id string
	// 数据库版本
	version int64
	// JSON.Decode(db.WeekPlan.Peroids)
	peroids [][]*peroid
	// 是否在录像时间
	isRecording bool
	// 关联的流
	task []*db.WeekPlanTask
	// task 数据是否有效
	taskOK bool
}

func (p *weekplan) init(model *db.WeekPlan) {
	p.id = model.ID
	p.version = model.Version
	var timePeroids [][]*db.TimePeroid
	json.NewDecoder(bytes.NewBufferString(*model.Peroids)).Decode(&timePeroids)
	for _, timePeroid := range timePeroids {
		var peroids []*peroid
		for _, tp := range timePeroid {
			pp := new(peroid)
			pp.Start, _ = time.ParseInLocation(db.TimePeroidFormat, tp.Start, time.Local)
			pp.End, _ = time.ParseInLocation(db.TimePeroidFormat, tp.End, time.Local)
			peroids = append(peroids, pp)
		}
		p.peroids = append(p.peroids, peroids)
	}
	p.initTask()
}

func (p *weekplan) initTask() {
	// 查询
	models, err := db.GetWeekPlanTaskByPlanID(p.id)
	if err != nil {
		log.ErrorfDepth(1, "load week plan %s task error: %s", p.id, err.Error())
	}
	// 加载
	p.task = models
	p.taskOK = true
}

func (p *weekplan) allTask() []*db.WeekPlanTask {
	// 数据不正常，加载
	if !p.taskOK {
		p.initTask()
	}
	return p.task
}
