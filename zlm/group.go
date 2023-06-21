package zlm

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/util"
)

var (
	// 默认的分组
	_group Group
)

// Group 用于管理一组流媒体服务
type Group struct {
	wg   sync.WaitGroup
	lock sync.RWMutex
	// 列表
	ser map[string]*Server
}

// NewGroup 返回一个组
func NewGroup() *Group {
	g := new(Group)
	g.ser = make(map[string]*Server)
	return g
}

// remove 移除列表
func (g *Group) remove(s *Server) {
	// 上锁
	g.lock.Lock()
	defer g.lock.Unlock()
	// 移除
	_s := g.ser[s.ID]
	if _s != nil && _s.version == s.version {
		delete(g.ser, s.ID)
	}
}

// snapshot 截图并保存到本地
func (g *Group) snapshot(s *Server) {
	// 保存目录为空不截图
	if s.SnapDir == "" {
		return
	}
	// copy 一份，不要占用锁
	medias := s.getAllMedia()
	// 截图超时
	timeout := fmt.Sprintf("%d", s.SyncInterval/time.Second)
	// 截图所有
	for _, media := range medias {
		g.wg.Add(1)
		go g.snapshotRoutine(s, media, timeout)
	}
}

// snapshotRoutine 在协程中截图并保存到本地
func (g *Group) snapshotRoutine(s *Server, m *mediaInfo, timeout string) {
	defer func() {
		log.Recover(recover())
		// 协程结束
		g.wg.Done()
	}()
	var req GetSnapReq
	req.TimeoutSec = timeout
	req.ExpireSec = timeout
	// 优先使用 rtsp 然后是 rtmp
	var err error
	if s.Cfg.RTSPPort != "" {
		req.URL = fmt.Sprintf("rtsp://localhost:%s/%s/%s", s.Cfg.RTSPPort, m.App, m.Stream)
		err = s.SaveSnap(&req, filepath.Join(s.SnapDir, m.App), m.Stream)
	} else if s.Cfg.RTMPPort != "" {
		req.URL = fmt.Sprintf("rtmp://localhost:%s/%s/%s", s.Cfg.RTMPPort, m.App, m.Stream)
		err = s.SaveSnap(&req, filepath.Join(s.SnapDir, m.App), m.Stream)
	}
	if err != nil {
		log.ErrorTrace(s.ID, err)
	}
}

// routine 是流媒体服务的协程，主要处理以下几个逻辑
// 确保加载配置和媒体流列表
// 初始化 ssrc 和 rtp 向上推流的端口
// 1.定时检查心跳
// 2.定时更新媒体流列表
// 3.定时保存每一个视频的截图
func (g *Group) routine(s *Server) {
	timer := time.NewTimer(0)
	defer func() {
		log.Recover(recover())
		// 计时器
		timer.Stop()
		// 移除
		g.remove(s)
		// 协程结束
		g.wg.Done()
	}()
	// 读取配置
	if !s.mustLoadConfig(timer) {
		return
	}
	// 读取媒体流列表
	if !s.mustLoadMediaList(timer) {
		return
	}
	// 有了配置，可以用了
	s.ok = true
	// 在线
	s.Online = true
	// 加载媒体流列表的时间
	now := time.Now()
	loadMediaTime := &now
	// 循环
	for {
		select {
		case <-s.quit.C:
			// 全局退出
			return
		case <-s.quit.C:
			// 服务退出
			return
		case now := <-timer.C:
			// 检查离线，心跳超时
			t := s.keepalive
			if now.Sub(*t) > s.keepaliveTimeout {
				s.Online = false
			}
			// 更新媒体流列表
			if s.SyncInterval > 0 && now.Sub(*loadMediaTime) > s.SyncInterval {
				// 记录时间
				loadMediaTime = &now
				// 更新媒体流列表
				err := s.loadMediaList()
				if err != nil {
					log.ErrorTrace(s.ID, err)
					break
				}
				// 截图
				g.snapshot(s)
			}
		}
		// 重置计时器
		timer.Reset(time.Second)
	}
}

// Add 添加一个 ser 到管理
func (g *Group) Add(ser *Server, version int64) {
	// 加入列表
	g.lock.Lock()
	old := g.ser[ser.ID]
	if old != nil {
		old.quit.Close()
	}
	ser.keepalive = &time.Time{}
	ser.version = version
	ser.quit = util.NewSignal()
	g.ser[ser.ID] = ser
	g.lock.Unlock()
	// 启动协程
	g.wg.Add(1)
	go g.routine(ser)
}

// Get 返回指定 id 的 Server
func (g *Group) Get(id string) *Server {
	// 查找
	g.lock.RLock()
	s := g.ser[id]
	g.lock.RUnlock()
	//
	return s
}

// Remove 移除
func (g *Group) Remove(id string) {
	// 上锁
	g.lock.Lock()
	defer g.lock.Unlock()
	// 删除
	ser := g.ser[id]
	if ser != nil {
		ser.quit.Close()
		delete(g.ser, id)
	}
}

// BatchRemove 批量移除
func (g *Group) BatchRemove(ids []string) {
	// 上锁
	g.lock.Lock()
	defer g.lock.Unlock()
	// 循环删除
	for _, id := range ids {
		ser := g.ser[id]
		if ser != nil {
			ser.quit.Close()
			delete(g.ser, id)
		}
	}
}

// Release 移除所有服务，如何等待所有协程结束
func (g *Group) Release() {
	// 移除所有
	g.lock.Lock()
	for _, ser := range g.ser {
		ser.quit.Close()
	}
	g.lock.Unlock()
	g.ser = make(map[string]*Server)
	// 等待
	g.wg.Wait()
}
