package sip

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"sync/atomic"
	"time"
)

// 一些默认的数值
const (
	// 默认的 io 读超时
	DefaultReadTimeout = time.Second * 10
	// 默认的 io 写超时
	DefaultWriteTimeout = time.Second * 10
	// 默认的每个 udp 连接的数据包缓存队列长度
	DefaultUDPDataQueueLen = 16
	// 最小的超时重发
	MinRTO = time.Millisecond * 200
)

const (
	// UDPMaxDataLen 表示一个 sip 消息的最大字节
	UDPMaxDataLen = 65535 - 20 - 8
	// UDPMinDataLen 表示一个 sip 消息的最小字节
	UDPMinDataLen = 576 - 20 - 8
)

var (
	// 服务已关闭
	errServerClosed = errors.New("server closed")
	// 连接已关闭
	errConnClosed = errors.New("conn closed")
	// 地址类型不对
	errAddrType = errors.New("error net.Addr type")
	// 发送请求生成事务发现事务已经存在
	errTransactionExists = errors.New("transaction exists")
)

// Handler 是处理消息的接口
type Handler interface {
	// 返回 true 表示已经处理好当前消息，递增 sn
	HandleRequest(*Request)
	// 返回值暂时没意义，先这样
	HandleResponse(*Response)
}

// Server 表示 sip 服务
type Server struct {
	// 监听端口
	Port int
	// 网关 的地址，有时候不想让别人知道在 nat 后面，隐藏内网地址
	AddrPort string
	// 最小是 UDPMinDataLen ，最大是 UDPMaxDataLen
	MessageLen int
	// 每个 udp 连接的数据包缓存队列长度，默认是 DefaultUDPDataQueueLen
	UDPDataQueueLen int
	// 读取消息的超时，也是事务的超时时间，毫秒。默认是 DefaultReadTimeout
	ReadTimeout time.Duration
	// 发送消息的超时，毫秒。默认是 DefaultWriteTimeout
	WriteTimeout time.Duration
	// sip 事务失效的超时时间，单位毫秒。默认是 DefaultTransactionTimeout
	// TransactionTimeout time.Duration
	// sip 消息重发间隔，单位毫秒，UDP 使用。默认是 DefaultTransactionRTO
	// TransactionRTO time.Duration
	// 回调函数
	Handler Handler
	// udp 超时重发
	rto time.Duration
	// 用于同步等待协程退出
	wg sync.WaitGroup
	// tcp 事务
	tcptx tcpTransactions
	// udp 事务
	udptx udpTransactions
	// tcp 监听
	tcpListener *net.TCPListener
	// tcp 锁
	tcplock sync.RWMutex
	// tcp 客户端列表
	tcpConns map[connKey]*tcpConn
	// udp 连接
	udpConn *net.UDPConn
	// udp 数据缓存
	udpData sync.Pool
	// Message 缓存池
	msgPool sync.Pool
	// buffer 缓存池
	bufPool sync.Pool
	// 状态，1: 正常
	ok int32
	// received 的值，隐藏内网地址
	received string
	// rport 的值，隐藏内网地址
	rport string
}

// Listen 根据 opt 初始化服务
func (s *Server) Listen() error {
	addrPort, err := netip.ParseAddrPort(s.AddrPort)
	if err != nil {
		return err
	}
	s.received = addrPort.Addr().String()
	s.rport = fmt.Sprintf("%d", addrPort.Port())
	// 纠正数据
	if s.MessageLen < UDPMinDataLen {
		s.MessageLen = UDPMinDataLen
	}
	if s.MessageLen > UDPMaxDataLen {
		s.MessageLen = UDPMaxDataLen
	}
	if s.UDPDataQueueLen < 1 {
		s.UDPDataQueueLen = DefaultUDPDataQueueLen
	}
	if s.ReadTimeout < 1 {
		s.ReadTimeout = DefaultReadTimeout
	}
	if s.WriteTimeout < 1 {
		s.WriteTimeout = DefaultWriteTimeout
	}
	s.rto = s.WriteTimeout / 3 * 3
	if s.rto < MinRTO {
		s.rto = MinRTO
	}
	// 事务表
	s.tcptx.init()
	s.udptx.init()
	// tcp 连接池
	s.tcpConns = make(map[connKey]*tcpConn)
	// 缓存池
	s.msgPool.New = func() any { return new(Message) }
	s.bufPool.New = func() any { return bytes.NewBuffer(nil) }
	s.udpData.New = func() any { return &udpData{b: make([]byte, s.MessageLen)} }
	// 开始服务
	atomic.StoreInt32(&s.ok, 1)
	port := fmt.Sprintf(":%d", s.Port)
	err = s.listenUDP(port)
	if err != nil {
		return err
	}
	err = s.listenTCP(port)
	if err != nil {
		s.closeUDP()
		return err
	}
	return nil
}

// Close 停止服务
func (s *Server) Close() error {
	// 修改状态
	if !atomic.CompareAndSwapInt32(&s.ok, 1, 0) {
		return errServerClosed
	}
	// 关闭双服务
	s.closeUDP()
	s.closeTCP()
	// 等待所有协程退出
	s.wg.Wait()
	return nil
}

// isOK 返回状态是否正常
func (s *Server) isOK() bool {
	return s.ok == 1
}

// GetMessage 从缓存池里返回 Message
func (s *Server) GetMessage() *Message {
	return s.msgPool.Get().(*Message)
}

// PutMessage 回收消息到缓存池
func (s *Server) PutMessage(m *Message) {
	s.msgPool.Put(m)
}

// HasTCPConn 是否有到 remoteAddr 的连接
func (s *Server) HasTCPConn(ip net.IP, port int) bool {
	var key connKey
	key.Init(ip, port)
	s.tcplock.RLock()
	_, ok := s.tcpConns[key]
	s.tcplock.RUnlock()
	return ok
}

// sendUDP 发送一个新的 udp 事务请求，调用要注意 data 的资源释放。
func (s *Server) sendUDP(ctx context.Context, conn Conn, msg *Message) error {
	// 新的事务
	t := s.udptx.new(msg)
	t.ctx = ctx
	// 数据
	t.writeMessage(conn, msg)
	// 事务回收计数
	atomic.AddInt32(&t.recovery, 1)
	// 启动发送协程
	s.wg.Add(1)
	go s.udpTransactionRetransmissionRoutine(ctx, conn, t)
	// 返回
	return nil
}

// sendUDPTimeout 发送一个新的 udp 事务请求。
// timeout 用于控制整个事务的超时，小于 0 则使用 s.TransactionTimeout ，调用要注意 data 的资源释放。
func (s *Server) sendUDPTimeout(conn Conn, msg *Message, timeout time.Duration) error {
	// 新的事务
	t := s.udptx.new(msg)
	// 数据
	t.writeMessage(conn, msg)
	// 事务回收计数
	atomic.AddInt32(&t.recovery, 1)
	// 启动发送协程
	s.wg.Add(1)
	go s.udpTransactionRetransmissionTimeoutRoutine(conn, t, timeout)
	// 返回
	return nil
}

// sendTCP 发送一个新的 tcp 事务请求，调用要注意 data 的资源释放。
func (s *Server) sendTCP(ctx context.Context, conn Conn, msg *Message) error {
	// 新的事务
	t := s.tcptx.new(msg)
	t.ctx = ctx
	// 发送
	err := t.writeMessage(conn, msg)
	if err != nil {
		return err
	}
	// 事务回收计数
	atomic.AddInt32(&t.recovery, 1)
	// 启动超时清理协程
	s.wg.Add(1)
	go s.clearTCPTransactionRoutine(ctx, t)
	// 返回
	return nil
}

// sendTCPTimeout 发送一个新的 tcp 事务请求。
// timeout 用于控制整个事务的超时，小于 0 则使用 s.TransactionTimeout 。
func (s *Server) sendTCPTimeout(conn Conn, msg *Message, timeout time.Duration) error {
	// 新的事务
	t := s.tcptx.new(msg)
	// 发送
	err := t.writeMessage(conn, msg)
	if err != nil {
		return err
	}
	// 事务回收计数
	atomic.AddInt32(&t.recovery, 1)
	// 启动超时清理协程
	s.wg.Add(1)
	go s.handleTCPTransactionTimeoutRoutine(t, timeout)
	// 返回
	return nil
}

// SendRequest 发送一个新的事务请求。如果 addr 是 tcp 且没有相应的连接则主动发起连接，
// ctx 在 tcp 发起连接时使用，在事务成功后，ctx 控制事务的销毁。
func (s *Server) SendRequest(ctx context.Context, addr net.Addr, msg *Message) error {
	if !s.isOK() {
		return errServerClosed
	}
	// tcp 地址
	if a, ok := addr.(*net.TCPAddr); ok {
		// 拿到/建立连接
		conn, err := s.getTCPConn(ctx, a)
		if err != nil {
			return err
		}
		// 发送
		return s.sendTCP(ctx, conn, msg)
	}
	// udp 地址
	if a, ok := addr.(*net.UDPAddr); ok {
		// 连接
		conn := new(udpConn)
		s.initUDPConn(conn, a)
		// 发送
		return s.sendUDP(ctx, conn, msg)
	}
	return errAddrType
}

// SendRequestTimeout 发送一个新的事务请求。
// timeout 用于控制整个事务的超时，小于 0 则使用 s.TransactionTimeout 。
func (s *Server) SendRequestTimeout(addr net.Addr, msg *Message, timeout time.Duration) error {
	if !s.isOK() {
		return errServerClosed
	}
	// tcp 地址
	if a, ok := addr.(*net.TCPAddr); ok {
		// 超时
		if timeout < 1 {
			// timeout = s.TransactionTimeout
			timeout = s.WriteTimeout
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// 拿到/建立连接
		conn, err := s.getTCPConn(ctx, a)
		if err != nil {
			return err
		}
		// 发送
		return s.sendTCPTimeout(conn, msg, timeout)
	}
	// udp 地址
	if a, ok := addr.(*net.UDPAddr); ok {
		// 连接
		conn := new(udpConn)
		s.initUDPConn(conn, a)
		// 发送
		return s.sendUDPTimeout(conn, msg, timeout)
	}
	return errAddrType
}

// SendRequestWithConn 使用当前的 conn 来发送新的事务请求，就不需要到连接表里查找了。
// ctx.Done 用于控制事务的销毁。
func (s *Server) SendRequestWithConn(ctx context.Context, conn Conn, msg *Message) error {
	if !s.isOK() {
		return errServerClosed
	}
	// 发送
	if conn.isUDP() {
		return s.sendUDP(ctx, conn, msg)
	}
	return s.sendTCP(ctx, conn, msg)
}

// SendRequestWithConnTimeout 使用当前的 conn 来发送新的事务请求，就不需要到连接表里查找了。
// timeout 用于控制整个事务的超时销毁，小于 0 则使用 s.TransactionTimeout 。
func (s *Server) SendRequestWithConnTimeout(conn Conn, msg *Message, timeout time.Duration) error {
	if !s.isOK() {
		return errServerClosed
	}
	// 发送
	if conn.isUDP() {
		return s.sendUDPTimeout(conn, msg, timeout)
	}
	return s.sendTCPTimeout(conn, msg, timeout)
}
