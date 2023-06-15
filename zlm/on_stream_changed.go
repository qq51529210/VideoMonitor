package zlm

// OnStreamChangedReq 表示 on_stream_changed 提交的数据，保存注册和注销的所有字段
type OnStreamChangedReq struct {
	// 服务器id，通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流注册或注销
	Regist bool `json:"regist"`
	// 流的媒体信息
	MediaInfo
}

// // OnStreamChanged 处理 zlm 的 on_stream_changed 回调
// func (s *Server) OnStreamChanged(data *OnStreamChangedReq) *PlayMediaInfo {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	//
// 	key := fmt.Sprintf("%s_%s", data.App, data.Stream)
// 	playURL, ok := s.stream[key]
// 	// 流注销
// 	if !data.Regist {
// 		if !ok {
// 			return playURL
// 		}
// 		// 移除
// 		switch data.Schema {
// 		case RTMP:
// 			playURL.RTMP = nil
// 		case RTSP:
// 			playURL.RTSP = nil
// 		case HLS:
// 			playURL.HLS = nil
// 		case TS:
// 			playURL.TS = nil
// 		case FMP4:
// 			playURL.FMP4 = nil
// 		}
// 		if playURL.RTMP == nil ||
// 			playURL.RTSP == nil ||
// 			playURL.HLS == nil ||
// 			playURL.TS == nil ||
// 			playURL.FMP4 == nil {
// 			delete(s.stream, key)
// 		}
// 		return nil
// 	}
// 	// 流注册
// 	if !ok {
// 		playURL = new(PlayMediaInfo)
// 		s.InitPlayMediaInfo(&data.MediaInfo, playURL)
// 		s.stream[key] = playURL
// 	}
// 	return playURL
// }
