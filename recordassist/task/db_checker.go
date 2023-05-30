package task

import (
	"fmt"
	"net/http"
	"os"
	"recordassist/db"
	"sync"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/util"
)

const (
	dayDuration = 24 * 3600
)

var (
	_dbChecker dbChecker
)

// dbChecker 用于检查数据库
type dbChecker struct {
	sync.WaitGroup
	// 最大并发数
	concurrency int
	// 调用超时
	apiCallTimeout time.Duration
	// 查询保存天数的 API
	apiGetSaveDaysURL string
	// 提交录像数据 API
	apiPostRecordURL string
	// 上传函数
	upload func(path string) (url string, err error)
}

// routine 检查数据库中的任务
// 1. 上传
// 2. 查询天数，然后提交
// 3. 删除本地文件
// 为了避免同时处理相同的任务，这里不用 chan
func (c *dbChecker) routine() {
	defer func() {
		log.Recover(recover())
	}()
	for {
		// 查询数据
		models, err := db.GetRecordList(c.concurrency)
		if err != nil {
			log.Errorf("get db record list error: %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		// 没有数据
		if len(models) < 1 {
			time.Sleep(time.Second)
			continue
		}
		// 并发处理
		for _, model := range models {
			switch model.Status {
			case db.RecordStatusReady:
				c.Add(1)
				go c.handleStep0Routine(model)
			case db.RecordStatusUploaded:
				c.Add(1)
				go c.handleStep1Routine(model)
			case db.RecordStatusSubmitted:
				c.Add(1)
				go c.handleStep2Routine(model)
			}
		}
		// 等待所有并发结束
		c.Wait()
	}
}

// handleStep0Routine 从步骤 0 开始
func (c *dbChecker) handleStep0Routine(model *db.Record) {
	defer func() {
		// 尝试更新状态
		if model.Status == db.RecordStatusUploaded {
			if err := db.UpdateRecord(model.Path, map[string]any{
				"Status": model.Status,
				"Name":   model.Name,
			}); err != nil {
				log.ErrorfDepth(1, "update db record error: %s", err.Error())
			}
		} else if model.Status == db.RecordStatusSubmitted {
			if err := db.UpdateRecordState(model); err != nil {
				log.ErrorfDepth(1, "update db record error: %s", err.Error())
			}
		}
		// 捕获异常
		log.Recover(recover())
		// 结束
		c.Done()
	}()
	if !c.handleStep0(model) {
		return
	}
	if !c.handleStep1(model) {
		return
	}
	if !c.handleStep2(model) {
		return
	}
}

// handleStep0 上传
func (c *dbChecker) handleStep0(model *db.Record) bool {
	// 上传
	name, err := c.upload(model.Path)
	if err != nil {
		log.Errorf("upload file %s error: %s", model.Path, err.Error())
		return false
	}
	//
	model.Name = name
	model.Status = db.RecordStatusUploaded
	return true
}

// handleStep1Routine 从步骤 1 开始
func (c *dbChecker) handleStep1Routine(model *db.Record) {
	defer func() {
		// 尝试更新状态
		if model.Status == db.RecordStatusSubmitted {
			if err := db.UpdateRecordState(model); err != nil {
				log.ErrorfDepth(1, "update db record error: %s", err.Error())
			}
		}
		// 捕获异常
		log.Recover(recover())
		// 结束
		c.Done()
	}()
	if !c.handleStep1(model) {
		return
	}
	if !c.handleStep2(model) {
		return
	}
}

// handleStep1 到计划管理查询保存天数，然后提交数据到录像管理
func (c *dbChecker) handleStep1(model *db.Record) bool {
	// 查询
	stream := fmt.Sprintf("%s_%s", model.App, model.Stream)
	saveDay, err := c.getRecordSaveDay(stream)
	if err != nil {
		log.Errorf("api call: query save days of stream %s error: %s", stream, err.Error())
		return false
	}
	// 提交录像，如果保存天数小于 1 直接下一步删除
	if saveDay > 0 {
		if err := c.postRecord(model, stream, saveDay); err != nil {
			log.Errorf("api call: submit data of file %s error: %s", model.Path, err.Error())
			return false
		}
	}
	//
	model.Status = db.RecordStatusSubmitted
	return true
}

// handleStep2Routine 从步骤 2 开始
func (c *dbChecker) handleStep2Routine(model *db.Record) {
	defer func() {
		// 捕获异常
		log.Recover(recover())
		// 结束
		c.Done()
	}()
	if !c.handleStep2(model) {
		return
	}
}

// handleStep2 移除本地文件
func (c *dbChecker) handleStep2(model *db.Record) bool {
	// 移除
	if err := os.Remove(model.Path); err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("remove file %s error: %s", model.Path, err.Error())
			return false
		}
	}
	// 清理数据库
	if _, err := db.Delete(model); err != nil {
		log.Errorf("delete db record %s error: %s", model.Path, err.Error())
		return false
	}
	//
	return true
}

type getRecordSaveDaysRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SaveDay     int64  `json:"saveDay"`
	IsRecording bool   `json:"isRecording"`
}

func (c *dbChecker) getRecordSaveDay(stream string) (int64, error) {
	// 查询请求
	var res []*getRecordSaveDaysRes
	err := util.HTTP[any](http.MethodGet, c.apiGetSaveDaysURL, nil, nil, &res, http.StatusOK, c.apiCallTimeout)
	if err != nil {
		return 0, err
	}
	// 没有关联计划
	if len(res) < 1 {
		return 0, nil
	}
	// 查找最大天数
	maxDay := res[0].SaveDay
	for i := 1; i < len(res); i++ {
		if res[i].SaveDay > maxDay {
			maxDay = res[i].SaveDay
		}
	}
	return maxDay, nil
}

type postRecordReq struct {
	// minio 的标识
	Name string `json:"name" binding:"required,max=40"`
	// stream
	Stream string `json:"stream" binding:"required,max=64"`
	// 大小
	Size int64 `json:"size" binding:"required,min=1"`
	// 时长
	Duration float64 ` json:"duration" binding:"required,min=1"`
	// 创建时间
	Time int64 `json:"time" binding:"required,min=1"`
	// 删除时间
	DeleteTime int64 `json:"deleteTime" binding:"required,min=0"`
}

func (c *dbChecker) postRecord(model *db.Record, stream string, saveDay int64) error {
	var postReq postRecordReq
	postReq.Name = model.Name
	postReq.Stream = stream
	postReq.Size = model.Size
	postReq.Duration = model.Duration
	postReq.Time = model.Time
	postReq.DeleteTime = model.Time + saveDay*dayDuration
	return util.HTTP[postRecordReq, any](http.MethodPost, c.apiPostRecordURL, nil, &postReq, nil, http.StatusCreated, c.apiCallTimeout)
}
