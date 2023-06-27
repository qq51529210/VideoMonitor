package sip

import (
	"encoding/binary"
	"net"
)

// Conn 主要是为了方便使用时不区分 tcp 和 udp
type Conn interface {
	// 返回网络类型
	Network() string
	// 返回对方的地址
	RemoteAddr() net.Addr
	// 返回对方的 IP
	RemoteIP() string
	// 返回对方的端口
	RemotePort() string
	// 返回对方的 IP:port
	RemoteAddrString() string
	// 写入数据
	write([]byte) error
	// 是否 udp
	isUDP() bool
}

// connKey 表示 udp 虚拟连接的 key
type connKey struct {
	// IPV6地址字符数组前64位
	ip1 uint64
	ip2 uint64
	// 端口
	port uint16
}

// 将128位的ip地址（v4的转成v6）的字节分成两个64位整数，加上端口，作为key
func (k *connKey) Init(ip net.IP, port int) {
	if len(ip) == net.IPv4len {
		k.ip1 = 0
		k.ip2 = uint64(0xff)<<40 | uint64(0xff)<<32 |
			uint64(ip[0])<<24 | uint64(ip[1])<<16 |
			uint64(ip[2])<<8 | uint64(ip[3])
	} else {
		k.ip1 = binary.BigEndian.Uint64(ip[0:])
		k.ip2 = binary.BigEndian.Uint64(ip[8:])
	}
	k.port = uint16(port)
}
