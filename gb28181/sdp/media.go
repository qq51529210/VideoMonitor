package sdp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 一些常量
const (
	ProtoUDP = "RTP/AVP"
	ProtoTCP = "TCP/RTP/AVP"
)

var (
	// ErrSDPMediaFormat 表示 Media 格式错误
	ErrSDPMediaFormat = errors.New("error sdp media format")
)

// Media 表示会话时间
type Media struct {
	// 媒体类型，video / audio
	Type string
	// 端口号
	Port int
	// 传输协议
	Proto string
	// 格式列表
	FMT []string
}

// Parse 从 line 中解析
func (m *Media) Parse(line string) error {
	p := strings.Fields(line)
	if len(p) < 3 {
		return ErrSDPMediaFormat
	}
	m.Type = p[0]
	n, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		return ErrSDPMediaFormat
	}
	m.Port = int(n)
	m.Proto = p[2]
	m.FMT = p[3:]
	//
	return nil
}

// String 格式化并返回
func (m *Media) String() string {
	return fmt.Sprintf("%s %d %s %s",
		m.Type,
		m.Port,
		m.Proto,
		strings.Join(m.FMT, " "))
}
