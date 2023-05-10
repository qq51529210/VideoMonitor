package zlm

// // 对应接口参数
// const (
// 	True  = "1"
// 	False = "0"
// )

// const (
// 	// RTP 表示 rtp 国标的 app
// 	RTP = "rtp"
// 	// PUSH 表示流媒体推流的 app
// 	PUSH = "push"
// 	// PULL 表示流媒体拉流的 app
// 	PULL = "pull"
// )

// // 流的协议
// const (
// 	RTMP = "rtmp"
// 	RTSP = "rtsp"
// 	HLS  = "hls"
// 	TS   = "ts"
// 	FMP4 = "fmp4"
// )

// // 错误
// var (
// 	ErrorServerNotFound     = errors.New("zlm: server not found")
// 	ErrorServerNotInit      = errors.New("zlm: server not init")
// 	ErrorServerNotAvailable = errors.New("not available zlm server")
// )

// var (
// 	lock          sync.RWMutex
// 	wg            sync.WaitGroup
// 	servers       map[string]*Server
// 	checkTimer    *time.Timer
// 	checkQuit     = make(chan struct{})
// 	offlineHandle func(string)
// )

// func init() {
// 	servers = make(map[string]*Server)
// }

// // Init 初始化
// func Init(offline func(string)) error {
// 	offlineHandle = offline
// 	// 加载服务
// 	models, err := db.All[db.ZLM](&db.ZLMListQuery{})
// 	if err != nil {
// 		return err
// 	}
// 	lock.Lock()
// 	defer lock.Unlock()
// 	for _, m := range models {
// 		newServer(m)
// 	}
// 	// 启动检查心跳
// 	wg.Add(1)
// 	go checkOfflineRoutine()
// 	return nil
// }

// // UnInit 释放资源
// func UnInit() {
// 	close(checkQuit)
// 	lock.Lock()
// 	for _, ser := range servers {
// 		// ser.exit()
// 		atomic.StoreInt32(&ser.ok, 2)
// 	}
// 	servers = make(map[string]*Server)
// 	lock.Unlock()
// 	// 等待
// 	wg.Wait()
// }

// // GetMinLoadServer 返回最小负载服务
// func GetMinLoadServer() *Server {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	load := int32(0)
// 	var ser *Server
// 	for _, v := range servers {
// 		if !v.IsOnline() {
// 			continue
// 		}
// 		if load < 1 {
// 			load = atomic.LoadInt32(&v.load)
// 			ser = v
// 		} else {
// 			if load < atomic.LoadInt32(&v.load) {
// 				ser = v
// 			}
// 		}
// 	}
// 	return ser
// }

// // GetMinLoadServerWithAppStream 返回包含 app_stream 的最小负载服务，和最小负载服务
// func GetMinLoadServerWithAppStream(app, stream string) (*Server, *Server) {
// 	key := fmt.Sprintf("%s_%s", app, stream)
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	load := int32(0)
// 	var ser1, ser2 *Server
// 	for _, v := range servers {
// 		if !v.IsOnline() {
// 			continue
// 		}
// 		v.lock.RLock()
// 		playURL := v.stream[key]
// 		v.lock.RUnlock()
// 		if load < 1 {
// 			load = atomic.LoadInt32(&v.load)
// 			if playURL != nil {
// 				ser1 = v
// 			}
// 			ser2 = v
// 		} else {
// 			if load < atomic.LoadInt32(&v.load) {
// 				if playURL != nil {
// 					ser1 = v
// 				}
// 				ser2 = v
// 			}
// 		}
// 	}
// 	return ser1, ser2
// }

// // GetServer 返回指定 id 的服务
// func GetServer(id string) *Server {
// 	lock.RLock()
// 	s := servers[id]
// 	lock.RUnlock()
// 	return s
// }

// // checkOfflineRoutine 是检查服务超时的协程
// func checkOfflineRoutine() {
// 	defer func() {
// 		log.Recover(recover())
// 		checkTimer.Stop()
// 		wg.Done()
// 	}()
// 	checkTimer = time.NewTimer(time.Second)
// 	for {
// 		select {
// 		case <-checkQuit:
// 			return
// 		case <-checkTimer.C:
// 			checkOffline()
// 			checkTimer.Reset(time.Second)
// 		}
// 	}
// }

// // checkOffline 检查服务心跳超时
// func checkOffline() {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	for _, s := range servers {
// 		if !s.IsOnline() && atomic.CompareAndSwapInt32(&s.offline, 0, 1) {
// 			serverID := s.model.ServerID
// 			log.Infof("zlm: server %s offline", serverID)
// 			offlineHandle(serverID)
// 		}
// 	}
// }

// type statusError int

// func (e statusError) Error() string {
// 	return fmt.Sprintf("http error status code %d", e)
// }

// // CodeError 表示 zlm code != 0 的响应
// type CodeError int

// func (e CodeError) Error() string {
// 	return fmt.Sprintf("zlm error code %d", e)
// }

// func newServer(model *db.ZLM) (*Server, error) {
// 	s := servers[model.ServerID]
// 	if s != nil {
// 		return s, nil
// 	}
// 	//
// 	s = new(Server)
// 	s.model = model
// 	s.stream = make(map[string]*PlayMediaInfo)
// 	s.keepaliveTime = &time.Time{}
// 	s.rtpPort = make([]uint16, *model.RTPMaxPort-*model.RTPMinPort)
// 	servers[model.ServerID] = s
// 	wg.Add(1)
// 	go s.loadDataRoutine()
// 	return s, nil
// }

// func httpGet[resData any](ser *Server, url string, query url.Values, data *resData) error {
// 	// 初始化请求
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return err
// 	}
// 	// 参数
// 	query.Add("secret", ser.model.Secret)
// 	req.URL.RawQuery = query.Encode()
// 	// 请求
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ser.model.RequestTimeout)*time.Millisecond)
// 	defer cancel()
// 	req = req.WithContext(ctx)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	// 判断状态码
// 	if res.StatusCode != http.StatusOK {
// 		return statusError(res.StatusCode)
// 	}
// 	// 解析数据
// 	return json.NewDecoder(res.Body).Decode(data)
// }

// func httpPost[reqData any, resData any](ser *Server, url string, query url.Values, reqBody *reqData, resBody *resData) error {
// 	// 初始化请求
// 	var body io.Reader = nil
// 	if reqBody != nil {
// 		buf := bytes.NewBuffer(nil)
// 		json.NewEncoder(buf).Encode(reqBody)
// 		body = buf
// 	}
// 	req, err := http.NewRequest(http.MethodPost, url, body)
// 	if err != nil {
// 		return err
// 	}
// 	// secret 参数
// 	query.Add("secret", ser.model.Secret)
// 	req.URL.RawQuery = query.Encode()
// 	// 请求
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ser.model.RequestTimeout)*time.Millisecond)
// 	defer cancel()
// 	req = req.WithContext(ctx)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	// 判断状态码
// 	if res.StatusCode != http.StatusOK {
// 		return statusError(res.StatusCode)
// 	}
// 	// 解析数据
// 	return json.NewDecoder(res.Body).Decode(resBody)
// }

// func httpGet2(ser *Server, url string, query url.Values) error {
// 	// 初始化请求
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return err
// 	}
// 	// secret 参数
// 	query.Add("secret", ser.model.Secret)
// 	req.URL.RawQuery = query.Encode()
// 	// 请求
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ser.model.RequestTimeout)*time.Millisecond)
// 	defer cancel()
// 	req = req.WithContext(ctx)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	// 判断状态码
// 	if res.StatusCode != http.StatusOK {
// 		return statusError(res.StatusCode)
// 	}
// 	return nil
// }

// func httpGet3(ser *Server, url string, query url.Values, data io.Writer) error {
// 	// 初始化请求
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return err
// 	}
// 	// secret 参数
// 	query.Add("secret", ser.model.Secret)
// 	req.URL.RawQuery = query.Encode()
// 	// 请求
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ser.model.RequestTimeout)*time.Millisecond)
// 	defer cancel()
// 	req = req.WithContext(ctx)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	// 判断状态码
// 	if res.StatusCode != http.StatusOK {
// 		return statusError(res.StatusCode)
// 	}
// 	// 拷贝数据
// 	_, err = io.Copy(data, res.Body)
// 	return err
// }

// // Add 添加新的服务
// func Add(model *db.ZLM) {
// 	newServer(model)
// }

// // Remove 移除指定 id 服务
// func Remove(id string) {
// 	lock.Lock()
// 	delete(servers, id)
// 	lock.Unlock()
// }
