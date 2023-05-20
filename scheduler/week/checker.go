package week

import (
	"net/http"
	"runtime"
	"scheduler/db"
	"sync"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/util"
)

var (
	_checker checker
)

func init() {
	_checker.weekplan = make(map[string]*weekplan)
}

// checker 用于检查周计划
type checker struct {
	sync.WaitGroup
	sync.RWMutex
	// 所有的录像计划，key: db.WeekPlan.ID
	weekplan map[string]*weekplan
	// 检查周期
	checkInterval time.Duration
	// 并发数
	concurrency int
	// 调用超时
	apiCallTimeout time.Duration
}

// init 初始化
func (c *checker) init(checkInterval, concurrency, apiTimeout int) error {
	c.checkInterval = time.Duration(checkInterval) * time.Second
	if c.checkInterval < 1 {
		c.checkInterval = time.Second
	}
	c.concurrency = concurrency
	if c.concurrency < 1 {
		c.concurrency = runtime.NumCPU()
	}
	c.apiCallTimeout = time.Duration(apiTimeout) * time.Second
	if c.apiCallTimeout < 1 {
		c.apiCallTimeout = time.Second
	}
	// 数据库
	models, err := db.GetWeekPlanAll()
	if err != nil {
		return err
	}
	// 加载
	c.weekplan = make(map[string]*weekplan)
	for _, model := range models {
		p := new(weekplan)
		p.init(model)
		c.weekplan[model.ID] = p
	}
	// 启动
	go c.checkRoutine()
	//
	return nil
}

// all 返回列表形式，解除锁的限制
func (c *checker) all() []*weekplan {
	// 上锁
	c.RLock()
	defer c.RUnlock()
	// 内存
	var models []*weekplan
	for _, v := range c.weekplan {
		models = append(models, v)
	}
	return models
}

// checkRoutine 检查协程
func (c *checker) checkRoutine() {
	timer := time.NewTimer(0)
	for {
		// 时间到
		now := <-timer.C
		// 列表
		wps := c.all()
		// 分块
		n := len(wps) / c.concurrency
		for len(wps) > 0 {
			// 并发检查
			c.Add(1)
			go c.concurrencyCheckRoutine(&now, wps[:n])
			if n > len(wps) {
				n = len(wps)
			}
			wps = wps[n:]
		}
		// 等待所有的结束
		c.Wait()
		// 重置计时器
		timer.Reset(c.checkInterval)
	}
}

// concurrencyCheckRoutine 并发检查协程，检查一部分的 weekplan
func (c *checker) concurrencyCheckRoutine(now *time.Time, wps []*weekplan) {
	defer func() {
		log.Recover(recover())
		// 协程结束
		c.Done()
	}()
	// 今天是周几
	weekDay := int(now.Weekday()) - 1
	if weekDay < 0 {
		weekDay = 6
	}
	peroid := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	for _, wp := range wps {
		// 超过了时间
		if weekDay >= len(wp.peroids) {
			continue
		}
		// 是否需要录像
		needRecord := false
		for _, p := range wp.peroids[weekDay] {
			// 当前时间在时间段内
			if peroid.After(p.Start) && peroid.Before(p.End) {
				needRecord = true
			}
			// 只要在某一段即可
			if needRecord {
				break
			}
		}
		wp.isRecording = needRecord
		// 需要录像，调用回调
		// todo 考虑一下，是否需要判断回调成功后停止
		// 如果这样，对方系统在这段时间出问题，就不能再触发或者停止任务了
		if wp.isRecording {
			c.startCallback(wp)
		} else {
			c.stopCallback(wp)
		}
	}
}

// startCallback 调用 id 关联的所有 task 的 StartCallback
func (c *checker) startCallback(wp *weekplan) {
	for _, task := range wp.allTask() {
		err := util.HTTP[int, int](http.MethodGet, *task.StartCallback, nil, nil, nil, http.StatusOK, c.apiCallTimeout)
		if err != nil {
			if code, ok := err.(util.HTTPStatusError); ok {
				// 没有这个流了，移除
				if code == http.StatusNotFound {
					continue
				}
			}
			log.Errorf("week plan %s start callback task %s error: %s", wp.id, task.TaskID, err.Error())
		}
	}
}

// stopCallback 调用 id 关联的所有 task 的 StartCallback
func (c *checker) stopCallback(wp *weekplan) {
	for _, task := range wp.allTask() {
		err := util.HTTP[int, int](http.MethodGet, *task.StopCallback, nil, nil, nil, http.StatusOK, c.apiCallTimeout)
		if err != nil {
			if code, ok := err.(util.HTTPStatusError); ok {
				// 没有这个流了，移除
				if code == http.StatusNotFound {
					continue
				}
			}
			log.Errorf("week plan %s stop callback task %s error: %s", wp.id, task.TaskID, err.Error())
		}
	}
}

// loadRoutine 在协程中加载数据，
// 保证加载最新的数据，id 不存在不会加载
func (c *checker) loadRoutine(id string) {
	defer func() {
		log.Recover(recover())
		// 协程结束
		c.Done()
	}()
	var weekPlanModel *db.WeekPlan
	var err error
	// 计划
	for {
		weekPlanModel, err = db.GetWeekPlan(id)
		if err != nil {
			log.Errorf("load db week plan %s error: %v", id, err)
			// 1 秒后继续
			time.Sleep(time.Second)
			continue
		}
		// 不存在，或者禁用
		if weekPlanModel == nil || *weekPlanModel.Enable != 1 {
			c.Lock()
			delete(c.weekplan, id)
			c.Unlock()
			return
		}
		// 成功
		c.add(weekPlanModel)
		return
	}
}

// add 添加，或者替换（如果比内存的版本大）数据
func (c *checker) add(model *db.WeekPlan) {
	// 上锁
	c.Lock()
	// 比较数据版本
	p := c.weekplan[model.ID]
	if p != nil {
		// 新加载的数据，比内存的旧
		if p.version > model.Version {
			c.Unlock()
			return
		}
	}
	// 比内存的新，替换
	p = new(weekplan)
	c.weekplan[model.ID] = p
	c.Unlock()
	// 初始化
	p.init(model)
}
