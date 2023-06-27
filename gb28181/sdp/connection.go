package sdp

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrSDPConnectionFormat 表示 Connection 格式错误
	ErrSDPConnectionFormat = errors.New("error sdp connection format")
)

// Connection 表示连接信息
type Connection struct {
	// 网络类型
	NetType string
	// 地址类型
	AddrType string
	// 连接地址
	Address string
}

// Parse 从 line 中解析
func (m *Connection) Parse(line string) error {
	p := strings.Fields(line)
	if len(p) != 3 {
		return ErrSDPConnectionFormat
	}
	m.NetType = p[0]
	m.AddrType = p[1]
	m.Address = p[2]
	//
	return nil
}

// String 格式化并返回
func (m *Connection) String() string {
	return fmt.Sprintf("%s %s %s",
		m.NetType,
		m.AddrType,
		m.Address)
}
