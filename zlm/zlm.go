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

type CodeError int

func (c CodeError) Error() string {
	return fmt.Sprintf("error code %d", c)
}
