package zlm

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/qq51529210/util"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
	// 需要截图的url，可以是本机的，也可以是远程主机的
	URL string `query:"url"`
	// 截图失败超时时间，防止FFmpeg一直等待截图
	TimeoutSec string `query:"timeout_sec"`
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	ExpireSec string `query:"expire_sec"`
}

// GetSnap 调用 /index/api/getSnap
// 获取截图或生成实时截图并返回，jpeg格式的图片，可以在浏览器直接打开
func (s *Server) GetSnap(req *GetSnapReq, writer io.Writer) error {
	return util.HTTPTo[any](http.MethodGet,
		s.url("getSnap"),
		s.query(req),
		nil,
		writer,
		http.StatusOK,
		s.APICallTimeout)
}

// SaveSnap 获取截图并保存到指定路径
func (s *Server) SaveSnap(req *GetSnapReq, dir, file string) error {
	// 请求
	var buf bytes.Buffer
	err := s.GetSnap(req, &buf)
	if err != nil {
		return err
	}
	// 目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// 文件，不使用 writer 模式，避免请求中途失败导致文件损坏，虽然这样也可能损坏
	return ioutil.WriteFile(filepath.Join(dir, file), buf.Bytes(), os.ModePerm)
}
