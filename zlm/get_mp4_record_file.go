package zlm

import (
	"net/http"

	"github.com/qq51529210/util"
)

// GetMP4RecordFileReq 是 GetMP4RecordFile 的参数
type GetMP4RecordFileReq struct {
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 流的录像日期，格式为2020-02-01,如果不是完整的日期，那么是搜索录像文件夹列表，否则搜索对应日期下的mp4文件列表
	Period string `query:"period"`
	// 自定义搜索路径，与startRecord方法中的customized_path一样，默认为配置文件的路径
	CustomizedPath string `query:"customized_path"`
}

// getMP4RecordFileRes 是 getMp4RecordFile 的返回值
type getMP4RecordFileRes struct {
	Code int                  `json:"code"`
	Data *GetMP4RecordFileRes `json:"data"`
}

// GetMP4RecordFileResData 是 GetMP4RecordFile 的返回值
type GetMP4RecordFileRes struct {
	Path     []string `json:"paths"`
	RootPath string   `json:"rootPath"`
}

// GetMP4RecordFile 调用 /index/api/getMp4RecordFile
// 搜索文件系统，获取流对应的录像文件列表或日期文件夹列表
func (s *Server) GetMP4RecordFile(req *GetMP4RecordFileReq) (*GetMP4RecordFileRes, error) {
	var _res getMP4RecordFileRes
	err := util.HTTP[any](http.MethodGet,
		s.url("getMp4RecordFile"),
		s.query(req),
		nil,
		&_res,
		http.StatusOK,
		s.APICallTimeout)
	if err != nil {
		return nil, err
	}
	if _res.Code != 0 {
		return nil, CodeError(_res.Code)
	}
	//
	return _res.Data, nil
}
