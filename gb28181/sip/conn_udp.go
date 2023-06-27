package sip

import (
	"io"
	"net"
)

// udpData 实现 io.Reader ，用于读取 udp 数据包
type udpData struct {
	// udp 数据
	b []byte
	// 数据的大小
	n int
	// 用于保存 read 的下标
	i int
	// 地址
	a *net.UDPAddr
}

// Len 返回剩余的数据
func (p *udpData) Len() int {
	return p.n - p.i
}

// Read 实现 io.Reader
func (p *udpData) Read(buf []byte) (int, error) {
	// 没有数据
	if p.i == p.n {
		return 0, io.EOF
	}
	// 还有数据，copy
	n := copy(buf, p.b[p.i:p.n])
	// 增加下标
	p.i += n
	// 返回
	return n, nil
}

// udpConn 表示一个虚拟的 udp 连接，实现了 conn 接口
type udpConn struct {
	// 底层连接
	conn *net.UDPConn
	// 对方地址
	remote *net.UDPAddr
	// 为了方便 via.received
	remoteIP string
	// 为了方便 via.rport
	remotePort string
	// ip:port
	remoteAddr string
}

func (c *udpConn) Close() error {
	return nil
}

func (c *udpConn) Network() string {
	return "udp"
}

func (c *udpConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *udpConn) RemoteIP() string {
	return c.remoteIP
}

func (c *udpConn) RemotePort() string {
	return c.remotePort
}

func (c *udpConn) RemoteAddrString() string {
	return c.remoteAddr
}

func (c *udpConn) write(buf []byte) error {
	_, err := c.conn.WriteTo(buf, c.remote)
	return err
}

func (c *udpConn) isUDP() bool {
	return true
}
