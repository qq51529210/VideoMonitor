package sip

import (
	"net"
	"sync/atomic"
	"time"
)

// tcpConn 表示一个 tcp 连接，实现了 conn 接口
type tcpConn struct {
	key connKey
	// 底层连接
	conn *net.TCPConn
	// 状态
	ok int32
	// io 发送超时时间
	writeTimeout time.Duration
	// 为了方便 via.received
	remoteIP string
	// 为了方便 via.rport
	remotePort string
	// ip:port
	remoteAddr string
}

func (c *tcpConn) Close() error {
	if atomic.CompareAndSwapInt32(&c.ok, 0, 1) {
		return c.conn.Close()
	}
	return errConnClosed
}

func (c *tcpConn) Network() string {
	return "tcp"
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *tcpConn) RemoteIP() string {
	return c.remoteIP
}

func (c *tcpConn) RemotePort() string {
	return c.remotePort
}

func (c *tcpConn) RemoteAddrString() string {
	return c.remoteAddr
}

func (c *tcpConn) write(buf []byte) error {
	if atomic.LoadInt32(&c.ok) != 0 {
		return errConnClosed
	}
	// 是否需要设置发送超时时间
	var err error
	if c.writeTimeout > 0 {
		err = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
		if err != nil {
			return err
		}
	}
	// 发送
	_, err = c.conn.Write(buf)
	return err
}

func (c *tcpConn) isUDP() bool {
	return false
}
