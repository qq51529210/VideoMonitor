package zlm

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_servers servers
)

func init() {
	_servers.ser = make(map[string]*Server)
}

// Server 表示一个流媒体服务
type Server struct {
	lock   sync.Mutex
	stream map[string]*MediaInfo
	// 配置
	cfg *Config
	// 流媒体夫 ID
	ID string
	// 调用的密钥
	Secret string
	// api 调用超时
	APICallTimeout time.Duration
	// 同步媒体流的间隔
	SyncInterval time.Duration
	// API 地址 (http|https)://ip:port
	APIBaseURL string
	// 外网访问的 ip ，生成播放地址时使用
	PublicIP string
	// 内网访问的 ip ，生成播放地址时使用
	PrivateIP string
	// 心跳时间，后台更新
	keepalive *time.Time
	// 是否在线，后台更新
	online bool
	// 是否可用
	ok int32
}

// IsOK 返回服务是否可用
func (s *Server) IsOK() bool {
	return s != nil && atomic.LoadInt32(&s.ok) == 1
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

// servers 用于管理所有的流媒体服务
type servers struct {
	sync.RWMutex
	// 列表
	ser map[string]*Server
}

// Add 添加一个 ser 到管理
func Add(ser *Server) {
	ser.stream = make(map[string]*MediaInfo)
	_servers.Lock()
	_servers.ser[ser.ID] = ser
	_servers.Unlock()
}

// Get 返回指定 id 的 Server
func Get(id string) *Server {
	// 查找
	_servers.RLock()
	s := _servers.ser[id]
	_servers.RUnlock()
	//
	return s
}

// Remove 移除指定 id 的 Server
func Remove(id string) {
	_servers.Lock()
	delete(_servers.ser, id)
	_servers.Unlock()
}

// BatchRemove 批量移除指定 id 的 Server
func BatchRemove(ids []string) {
	_servers.Lock()
	for _, id := range ids {
		delete(_servers.ser, id)
	}
	_servers.Unlock()
}
