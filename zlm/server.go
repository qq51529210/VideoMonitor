package zlm

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/util"
)

// Server 表示一个流媒体服务
type Server struct {
	// 配置
	Cfg *Config
	// 流媒体夫 ID
	ID string
	// 调用的密钥
	Secret string
	// api 调用超时，单位秒
	APICallTimeout time.Duration
	// 同步媒体流的间隔，单位秒
	SyncInterval time.Duration
	// API 地址 (http|https)://ip:port
	APIBaseURL string
	// 外网访问的 ip ，生成播放地址时使用
	PublicIP string
	// 内网访问的 ip ，生成播放地址时使用
	PrivateIP string
	// 截图目录，不为空在更新媒体流列表时会自动截图并保存
	SnapDir string
	// 负载
	Load int32
	// 心跳时间，后台更新
	keepalive *time.Time
	// 心跳超时时间，从流媒体服务配置中得到
	keepaliveTimeout time.Duration
	// 是否在线，后台更新
	Online bool
	// 数据版本，用于移除列表的时候判断
	version int64
	// 是否可用
	ok bool
	// 退出信号
	quit *util.Signal
	// 同步锁
	lock sync.RWMutex
	// 媒体流列表，app:stream:
	media map[string]map[string]*mediaInfo
}

// IsOK 返回服务是否可用
func (s *Server) IsOK() bool {
	return s != nil && s.Online && s.ok
}

func (s *Server) url(path string) string {
	return fmt.Sprintf("%s/index/api/%s", s.APIBaseURL, path)
}

// query 使用反射来初始化 v 并返回 url.Values
func (s *Server) query(v any) url.Values {
	q := make(url.Values)
	//
	if v != nil {
		vv := reflect.ValueOf(v)
		vk := vv.Kind()
		if vk == reflect.Pointer {
			vv = vv.Elem()
			vk = vv.Kind()
		}
		if vk != reflect.Struct {
			panic("v must be struct")
		}
		//
		vt := vv.Type()
		for i := 0; i < vt.NumField(); i++ {
			fv := vv.Field(i)
			if !fv.IsValid() {
				continue
			}
			ft := vt.Field(i)
			name := ft.Tag.Get(queryTag)
			q.Set(name, fmt.Sprintf("%v", fv.Interface()))
		}
	}
	q.Set(querySecret, s.Secret)
	if q.Get(queryVHost) == "" {
		q.Set(queryVHost, DefaultVHost)
	}
	//
	return q
}

// mustLoadMediaList 从 zlm 读取配置，循环读取直到成功
func (s *Server) mustLoadMediaList(timer *time.Timer) bool {
	timer.Reset(0)
	for {
		select {
		case <-s.quit.C:
			// 服务退出
			return false
		case <-timer.C:
			// 加载
			err := s.loadMediaList()
			if err != nil {
				log.ErrorTrace(s.ID, err)
				break
			}
			// 成功返回
			return true
		}
		// 请求失败，一会儿重试
		timer.Reset(time.Second)
	}
}

// loadConfig 读取配置
func (s *Server) loadConfig() error {
	// 请求
	err := s.GetServerConfig()
	if err != nil {
		return err
	}
	// 心跳间隔
	n, err := strconv.ParseFloat(s.Cfg.HookAliveInterval, 64)
	if err != nil {
		return err
	}
	// 加 3 秒作为网络缓冲
	s.keepaliveTimeout = time.Duration(n+3) * time.Second
	// 返回
	return nil
}

// mustLoadConfig 读取配置直到成功
func (s *Server) mustLoadConfig(timer *time.Timer) bool {
	timer.Reset(0)
	for {
		select {
		case <-s.quit.C:
			// 服务退出
			return false
		case <-timer.C:
			// 加载
			err := s.loadConfig()
			if err != nil {
				log.ErrorTrace(s.ID, err)
				break
			}
			// 成功返回
			return true
		}
		// 请求失败，一会儿重试
		timer.Reset(time.Second)
	}
}

// loadConfig 从 zlm 读取流媒体列表，然后检查推流和拉流
func (s *Server) loadMediaList() error {
	// 请求
	data, err := s.GetMediaList(&GetMediaListReq{Schema: RTMP})
	if err != nil {
		return err
	}
	s.UpdateMediaList(data)
	// 返回
	return nil
}

// getAllMedia 返回所有的媒体流
func (s *Server) getAllMedia() []*mediaInfo {
	var ms []*mediaInfo
	s.lock.RLock()
	for _, app := range s.media {
		for _, stream := range app {
			ms = append(ms, stream)
		}
	}
	s.lock.RUnlock()
	//
	return ms
}

// UpdateMediaList 更新保存的媒体流列表
func (s *Server) UpdateMediaList(m []*MediaList) {
	// 解析
	medias := make(map[string]map[string]*mediaInfo)
	for _, d := range m {
		// 轨道
		info := new(mediaInfo)
		info.Video, info.Audio = parseTracks(d.Tracks)
		info.App, info.Stream = d.App, d.Stream
		info.IsRecordingMP4, info.IsRecordingHLS = d.IsRecordingMP4, d.IsRecordingHLS
		// 列表
		app := medias[d.App]
		if app == nil {
			app = make(map[string]*mediaInfo)
			medias[d.App] = app
		}
		app[d.Schema] = info
	}
	// 替换
	s.lock.Lock()
	s.media = medias
	s.lock.Unlock()
}
