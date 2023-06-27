package sdp

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrSDPOriginFormat 表示 Origin 格式错误
	ErrSDPOriginFormat = errors.New("error sdp origin format")
)

// Origin 表示会话创建者信息
type Origin struct {
	// 用户名
	Username string
	// 会话 ID
	SessionID string
	// 会话版本
	SessionVersion string
	// 网络类型，一般为 IN ，表示 internet
	NetType string
	// 地址类型，IP4 / IP6
	AddrType string
	// 地址
	Address string
}

// Parse 从 line 中解析
func (m *Origin) Parse(line string) error {
	p := strings.Fields(line)
	n := len(p)
	if n < 6 {
		return ErrSDPOriginFormat
	}
	n--
	m.Address = p[n]
	n--
	m.AddrType = p[n]
	n--
	m.NetType = p[n]
	n--
	m.SessionVersion = p[n]
	n--
	m.SessionID = p[n]
	p = p[:n]
	if len(p) > 1 {
		m.Username = strings.Join(p[:n], " ")
	} else {
		m.Username = p[0]
	}
	//
	return nil
}

// String 格式化并返回
func (m *Origin) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		m.Username,
		m.SessionID,
		m.SessionVersion,
		m.NetType,
		m.AddrType,
		m.Address)
}
