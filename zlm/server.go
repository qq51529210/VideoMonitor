package zlm

import "sync"

// Server 表示一个流媒体服务
type Server struct {
	lock sync.RWMutex
	// // 是否有效
	// ok int32
	// // 数据
	// model *db.ZLM
	// // on_server_keepalive 回调时间
	// keepaliveTime *time.Time
	// // 负载，由心跳的 tcpsession 确定
	// load int32
	// // 流列表
	// stream map[string]*PlayMediaInfo
	// // 心跳间隔
	// keepaliveTimeout time.Duration
	// // 是否调用过离线回调
	// offline int32
	// // 用于上级 invite 的推流端口
	// rtpPort []uint16
}

// // IsOnline 返回是否在线
// func (s *Server) IsOnline() bool {
// 	if atomic.LoadInt32(&s.ok) == 1 {
// 		t := s.keepaliveTime
// 		return time.Since(*t) <= s.keepaliveTimeout
// 	}
// 	return false
// }

// // loadDataRoutine 后台读取配置和媒体流列表
// func (s *Server) loadDataRoutine() {
// 	defer func() {
// 		log.Recover(recover())
// 		atomic.CompareAndSwapInt32(&s.ok, 0, 1)
// 		wg.Done()
// 	}()
// 	// 先获取配置
// 	for atomic.LoadInt32(&s.ok) == 0 {
// 		cfgs, err := s.GetServerConfig()
// 		if err != nil {
// 			log.Error(err)
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		if len(cfgs) < 1 {
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		s.cfg = cfgs[0]
// 		n, err := strconv.ParseInt(s.cfg.HTTPKeepaliveSecond, 10, 64)
// 		if err != nil {
// 			log.Error(err)
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		// 加 3 秒作为网络缓冲，已经足够了
// 		s.keepaliveTimeout = time.Duration(n+1) * time.Second
// 		break
// 	}
// 	for atomic.LoadInt32(&s.ok) == 0 {
// 		medias, err := s.GetMediaList(&GetMediaListReq{
// 			Schema: RTSP,
// 		})
// 		if err != nil {
// 			log.Error(err)
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		// 加载到内存
// 		for _, m := range medias {
// 			key := fmt.Sprintf("%s_%s", m.App, m.Stream)
// 			s.lock.Lock()
// 			playURL, ok := s.stream[key]
// 			if !ok {
// 				playURL = new(PlayMediaInfo)
// 				s.stream[key] = playURL
// 			}
// 			s.InitPlayMediaInfo(m, playURL)
// 			s.lock.Unlock()
// 		}
// 		break
// 	}
// }

// // ServerID 返回流媒体的 ServerID
// func (s *Server) ServerID() string {
// 	return s.model.ServerID
// }

// // Config 返回配置
// func (s *Server) Config() *Config {
// 	return s.cfg
// }

// // PublicIP 返回服务的公网 ip
// func (s *Server) PublicIP() string {
// 	return s.model.PublicIP
// }

// // GetRTPPort 返回一个可用的端口，-1 表示没有端口。
// func (s *Server) GetRTPPort() int {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	for i := range s.rtpPort {
// 		if i == 0 {
// 			s.rtpPort[i] = 1
// 			return int(*s.model.RTPMinPort) + i
// 		}
// 	}
// 	return -1
// }

// // PutRTPPort 释放端口。
// func (s *Server) PutRTPPort(port int) {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	i := port - int(*s.model.RTPMinPort)
// 	if i >= 0 {
// 		s.rtpPort[i] = 0
// 	}
// }

// func (s *Server) url(path string) string {
// 	return fmt.Sprintf("%s/index/api/%s", s.model.BaseURL, path)
// }
