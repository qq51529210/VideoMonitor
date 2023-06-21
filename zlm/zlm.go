package zlm

import (
	"errors"
	"fmt"
)

// 对应接口参数
const (
	True  = "1"
	False = "0"
	Zero  = "0"
)

const (
	// DefaultVHost 默认的
	DefaultVHost = "__defaultVhost__"
)

// 协议
const (
	RTMP = "rtmp"
	RTSP = "rtsp"
	HLS  = "hls"
	TS   = "ts"
	FMP4 = "fmp4"
)

const (
	queryTag    = "query"
	querySecret = "secret"
	queryVHost  = "vhost"
)

var (
	// ErrServerNotFound 表示找不到相应的服务
	ErrServerNotFound = errors.New("media server not found")
	// ErrServerUnavailable 表示当前的服务不可用
	ErrServerUnavailable = errors.New("media server unavailable")
)

type Error struct {
	Code int
	Msg  string
	// 调用的 api
	API string
	// 服务 ID
	ID string
}

func (c *Error) Error() string {
	return fmt.Sprintf("server %s call %s error code: %d, msg: %s", c.ID, c.API, c.Code, c.Msg)
}

type apiRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
