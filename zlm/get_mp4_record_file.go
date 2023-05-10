package zlm

import "net/url"

// GetMP4RecordFileReq 是 GetMP4RecordFile 的参数
type GetMP4RecordFileReq struct {
	// 筛选虚拟主机
	VHost string
	// 筛选应用名，例如 live
	App string
	// 筛选流id，例如 test
	Stream string
	// 流的录像日期，格式为2020-02-01,如果不是完整的日期，那么是搜索录像文件夹列表，否则搜索对应日期下的mp4文件列表
	Period string
	// 自定义搜索路径，与startRecord方法中的customized_path一样，默认为配置文件的路径
	CustomizedPath string
}

func (m *GetMP4RecordFileReq) toQuery() url.Values {
	q := make(url.Values)
	if m.VHost != "" {
		q.Set("vhost", m.VHost)
	}
	if m.App != "" {
		q.Set("app", m.App)
	}
	if m.Stream != "" {
		q.Set("stream", m.Stream)
	}
	if m.Period != "" {
		q.Set("period", m.Period)
	}
	if m.CustomizedPath != "" {
		q.Set("customized_path", m.CustomizedPath)
	}
	return q
}

// getMP4RecordFileRes 是 GetMP4RecordFile 的返回值
type getMP4RecordFileRes struct {
	Code int `json:"code"`
	// 是否存在
	Data *GetMP4RecordFileResData `json:"data"`
}

// GetMP4RecordFileResData 是 getMP4RecordFileRes 的 Data 字段
type GetMP4RecordFileResData struct {
	Path     []string `json:"paths"`
	RootPath string   `json:"rootPath"`
}

// GetMP4RecordFile 调用 /index/api/getMp4RecordFile
// 搜索文件系统，获取流对应的录像文件列表或日期文件夹列表
func (s *Server) GetMP4RecordFile(req *GetMP4RecordFileReq) (*GetMP4RecordFileResData, error) {
	query := make(url.Values)
	if req != nil {
		query = req.toQuery()
	}
	var res getMP4RecordFileRes
	err := httpGet(s, s.url("getMp4RecordFile"), query, &res)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, CodeError(res.Code)
	}
	return res.Data, nil
}
