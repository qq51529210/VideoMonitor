package weekplan

import (
	"net/http"
	"recordplan/db"
	"runtime"
	"sync"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/util"
)

var (
	_checker checker
)

func init() {
	_checker.stream = make(map[string]int)
	_checker.weekplan = make(map[string]*weekplan)
}

// checker 用于检查周计划
type checker struct {
	sync.WaitGroup
	sync.RWMutex
	// 所有的录像计划，key: db.WeekPlan.ID
	weekplan map[string]*weekplan
	// 需要录像流，方便查询
	// key: db.WeekPlanStream.Stream
	stream map[string]int
	// 检查周期
	checkInterval time.Duration
	// 并发数
	concurrency int
}

// init 初始化
func (c *checker) init(checkInterval, concurrency int) error {
	c.checkInterval = time.Duration(checkInterval) * time.Second
	if c.checkInterval < 1 {
		c.checkInterval = time.Second
	}
	c.concurrency = concurrency
	if c.concurrency < 1 {
		c.concurrency = runtime.NumCPU()
	}
	// 数据库
	models, err := db.GetWeekPlanAll()
	if err != nil {
		return err
	}
	// 加载
	c.stream = make(map[string]int)
	c.weekplan = make(map[string]*weekplan)
	for _, model := range models {
		p := new(weekplan)
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
		if weekDay >= len(wp.Peroids) {
			continue
		}
		// 是否需要录像
		needRecord := false
		for _, p := range wp.Peroids[weekDay] {
			// 当前时间在时间段内
			if peroid.After(p.Start) && peroid.Before(p.End) {
				needRecord = true
			}
			// 只要在某一段即可
			if needRecord {
				break
			}
		}
		wp.IsRecording = needRecord
		// 需要录像，调用回调
		if wp.IsRecording {
			c.callback(wp.ID)
		}
	}
}

// callback 调用 id 关联的所有 stream 的 callback
func (c *checker) callback(id string) {
	// 数据库
	models, err := db.GetWeekPlanStreamListByPlanID(id)
	if err != nil {
		log.Errorf("week plan %s record callback get db stream list error: %s", id, err.Error())
		return
	}
	// 回调
	for _, model := range models {
		err := util.HTTP[int, int](http.MethodGet, *model.Callback, nil, nil, nil, http.StatusOK, time.Second)
		if err != nil {
			log.Errorf("week plan %s record stream %s callback error: %s", id, model.Stream, err.Error())
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
		// 不存在
		if weekPlanModel == nil {
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
	// 新加载的数据，比内存的旧
	if p.Version > model.UpdatedAt {
		c.Unlock()
		return
	}
	// 比内存的新，替换
	p = new(weekplan)
	c.weekplan[model.ID] = p
	c.Unlock()
	// 初始化
	p.init(model)
}

// remove 移除数据，id 不存在略过
func (c *checker) remove(id string) {
	c.Lock()
	delete(c.weekplan, id)
	c.Unlock()
}

// batchRemove 批量移除数据，id 不存在略过
func (c *checker) batchRemove(ids []string) {
	c.Lock()
	for _, id := range ids {
		delete(c.weekplan, id)
	}
	c.Unlock()
}
