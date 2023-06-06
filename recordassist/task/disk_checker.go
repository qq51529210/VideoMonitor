package task

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/qq51529210/video-monitor/recordassist/db"

	"github.com/qq51529210/log"
)

var (
	_diskChecker diskChecker
)

// diskChecker 用于检查磁盘的录像
// 防止 api 失败，录像少了
type diskChecker struct {
	sync.WaitGroup
	sync.RWMutex
	// 扫描的目录
	dir string
	// 检查的间隔
	checkInterval time.Duration
	// 录像文件的最大时长
	maxDuration time.Duration
}

// routine 启动检查
// 定时检查指定的目录
// dir/live/obs/2019-09-20/15-53-02.mp4
func (c *diskChecker) routine() {
	defer func() {
		log.Recover(recover())
	}()
	timer := time.NewTimer(0)
	for {
		// 时间到
		<-timer.C
		// 遍历目录
		err := filepath.WalkDir(c.dir, c.walkDirFn)
		if err != nil {
			log.Errorf("read dir %s error: %s", c.dir, err.Error())
			timer.Reset(time.Second)
			continue
		}
		// 重置计时器
		timer.Reset(time.Second)
	}
}

// walkDirFn 处理目录遍历
func (c *diskChecker) walkDirFn(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return nil
	}
	// 路过目录
	if d.IsDir() {
		return nil
	}
	// 根据 zlm 流媒体的录像目录规则，如果是文件，就肯定是目录文件了
	c.handleRecord(path, d)
	//
	return nil
}

// handleRecord 处理文件
func (c *diskChecker) handleRecord(p string, d fs.DirEntry) {
	// 信息
	info, err := d.Info()
	if err != nil {
		log.Errorf("get record file %s info error: %s", p, err.Error())
		return
	}
	// 时长
	cTime, mTime := getTime(info)
	dur := mTime.Sub(cTime)
	if dur <= c.maxDuration {
		return
	}
	// 找出 app 和 stream 目录
	part := filepath.SplitList(p)
	n := len(part)
	if n < 4 {
		log.Warnf("record file %s error path", p)
		// 删除
		if err := os.Remove(p); err != nil {
			log.Errorf("remove file %s error: %s", p, err.Error())
		}
		return
	}
	// date dir
	n--
	// 数据库
	var model db.Record
	model.Path = p
	ok, err := db.Get(&model)
	if err != nil {
		log.Errorf("get db record %s error: %s", p, err.Error())
		return
	}
	if !ok {
		model.Status = db.RecordStatusReady
		model.Time = cTime.Unix()
		model.Duration = float64(dur) / float64(time.Second)
		model.Size = info.Size()
		model.Stream = part[n]
		n--
		model.App = part[n]
		if _, err := db.Add(&model); err != nil {
			log.Errorf("add db record %s error: %s", p, err.Error())
		}
	}
}
