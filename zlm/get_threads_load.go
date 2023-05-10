package zlm

import "net/url"

// getThreadsLoadRes 是 GetThreadsLoad 的返回值
type getThreadsLoadRes struct {
	Code int
	Data []*GetThreadsLoadResData
}

// GetThreadsLoadResData 是 getThreadsLoadRes 的 Data 字段
type GetThreadsLoadResData struct {
	Delay int `json:"delay"`
	Load  int `json:"load"`
}

// GetThreadsLoad 调用 /index/api/getThreadsLoad
// 获取各epoll(或select)线程负载以及延时
// 返回[]*GetThreadsLoadResData
func (s *Server) GetThreadsLoad() ([]*GetThreadsLoadResData, error) {
	var res getThreadsLoadRes
	query := make(url.Values)
	err := httpGet(s, s.url("getThreadsLoad"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}

// GetWorkThreadsLoad 调用 /index/api/getWorkThreadsLoad
// 获取各后台epoll(或select)线程负载以及延时
// 返回[]*GetThreadsLoadResData
func (s *Server) GetWorkThreadsLoad() ([]*GetThreadsLoadResData, error) {
	var res getThreadsLoadRes
	query := make(url.Values)
	err := httpGet(s, s.url("getWorkThreadsLoad"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
