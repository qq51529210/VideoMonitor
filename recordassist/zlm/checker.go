package zlm

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/qq51529210/log"
)

var (
	_checker checker
)

// checker 用于检查 zlm 的录像
type checker struct {
	sync.WaitGroup
	sync.RWMutex
	// 扫描的目录
	dir string
	// 检查的间隔
	checkInterval time.Duration
	// 调用超时
	apiCallTimeout time.Duration
	// 录像文件的最大时长
	maxDuration time.Duration
	// 录像文件
	record map[string]int
}

// routine 启动检查
// 定时检查指定的目录
// dir/live/obs/2019-09-20/15-53-02.mp4
func (c *checker) routine() {
	defer func() {
		log.Recover(recover())
	}()
	timer := time.NewTimer(0)
	for {
		// 时间到
		<-timer.C
		// 目录
		appDirs, err := os.ReadDir(c.dir)
		if err != nil {
			log.Errorf("read dir %s error: %s", c.dir, err.Error())
			timer.Reset(time.Second)
			continue
		}
		// app 目录
		for _, appDir := range appDirs {
			// 只要目录
			if !appDir.IsDir() {
				continue
			}
			// 一个目录一个协程处理
			c.Add(1)
			go c.handleAppDirRoutine(filepath.Join(c.dir, appDir.Name()))
		}
		c.Wait()
		// 重置计时器
		timer.Reset(c.checkInterval)
	}
}

// handleAppDirRoutine 循环 appDir 下的目录
func (c *checker) handleAppDirRoutine(appDir string) {
	// stream 目录
	streamDirs, err := os.ReadDir(appDir)
	if err != nil {
		log.Errorf("read app dir %s error: %s", appDir, err.Error())
		return
	}
	for _, streamDir := range streamDirs {
		if streamDir.IsDir() {
			c.handleStreamDir(appDir, streamDir.Name())
		}
	}
}

// handleStreamDir 循环 streamDir 下的目录
func (c *checker) handleStreamDir(appDir, streamDir string) {
	// date 目录
	dateDirs, err := os.ReadDir(streamDir)
	if err != nil {
		log.Errorf("read stream dir %s error: %s", streamDir, err.Error())
		return
	}
	for _, dateDir := range dateDirs {
		if dateDir.IsDir() {
			c.handleDateDir(appDir, streamDir, dateDir.Name())
		}
	}
}

// handleDateDir 循环 dateDir 下的目录
func (c *checker) handleDateDir(appDir, streamDir, dateDir string) {
	// 录像文件
	recordFiles, err := os.ReadDir(dateDir)
	if err != nil {
		log.Errorf("read date dir %s error: %s", dateDir, err.Error())
		return
	}
	for _, recordFile := range recordFiles {
		if !recordFile.IsDir() {
			c.handleRecordFile(appDir, streamDir, dateDir, recordFile)
		}
	}
}

// handleRecordFile 处理录像文件
func (c *checker) handleRecordFile(appDir, streamDir, dateDir string, fs fs.DirEntry) {
	recordPath := filepath.Join(appDir, streamDir, dateDir, fs.Name())
	// 看看时长
	info, err := fs.Info()
	if err != nil {
		log.Errorf("get record file %s info error: %s", recordPath, err.Error())
		return
	}
	cTime, mTime := getTime(info)
	if mTime.Sub(cTime) <= c.maxDuration {
		return
	}
	// 看看有没有处理过
	c.Lock()
	state, ok := c.record[recordPath]
	if !ok {
		c.record[recordPath] = 0
	}
	c.Unlock()
	// 看看状态
	switch state {
	case 0:
		// 第一次
		c.handleRecordFileStep1()
	case 1:
		// 上传成功
		c.handleRecordFileStep2()
	case 2:
		// 查询保存天数
		c.handleRecordFileStep3()
	case 3:
		// 提交数据管理
		c.handleRecordFileStep4()
	}
}

func (c *checker) handleRecordFileStep1() {
}

func (c *checker) handleRecordFileStep2() {
}

func (c *checker) handleRecordFileStep3() {
}

func (c *checker) handleRecordFileStep4() {
}
