package sip

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/qq51529210/log"
)

// closeTCP 停止 tcp 监听，关闭所有的 tcp 连接。
func (s *Server) closeTCP() {
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}
	// 停止监听
	s.tcpListener.Close()
	// 关闭所有conn
	for _, c := range s.tcpConns {
		c.Close()
	}
}

// listenTCP 初始化 tcp 监听，启动 1 个监听协程接入客户端连接。
func (s *Server) listenTCP(port string) error {
	// 初始化
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return err
	}
	s.tcpListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	// 监听协程
	s.wg.Add(1)
	go s.listenTCPRoutine()
	// 返回
	return nil
}

// listenTCPRoutine 循环监听客户端连接，然后启动 1 个 tcp 读协程
func (s *Server) listenTCPRoutine() {
	log.Debug("tcp listen routine start")
	defer func() {
		// log.Recover(recover())
		// 日志
		log.Debug("tcp listen routine end")
		// 协程结束
		s.wg.Done()
	}()
	for s.isOK() {
		// 监听
		conn, err := s.tcpListener.AcceptTCP()
		if err != nil {
			log.Error(err)
			continue
		}
		// 加入到列表
		c := s.newTCPConn(conn)
		s.tcplock.Lock()
		s.tcpConns[c.key] = c
		s.tcplock.Unlock()
		// 处理协程
		s.wg.Add(1)
		go s.readTCPRoutine(c)
	}
}

// readTCPRoutine 读取数据并处理的协程
func (s *Server) readTCPRoutine(c *tcpConn) {
	log.Debugf("tcp %s read routine start", c.RemoteAddrString())
	defer func() {
		// log.Recover(recover())
		// 日志
		log.Debugf("tcp %s read routine end", c.RemoteAddrString())
		// 关闭
		s.closeTCPConn(c)
		// 协程结束
		s.wg.Done()
	}()
	// 开始
	var err error
	reader := NewReader(c.conn, s.MessageLen).(*reader)
	for {
		// 读超时，避免空闲连接
		err = c.conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		if err != nil {
			log.Error(err)
			return
		}
		// 读取并解析消息
		msg := s.msgPool.Get().(*Message)
		msg.Reset()
		err = msg.ParseFrom(reader, s.MessageLen)
		if err != nil {
			s.msgPool.Put(msg)
			log.Errorf("read tcp %v %v\n%s", c.RemoteAddrString(), err, string(reader.buf[reader.begin:reader.end]))
			return
		}
		// 处理
		s.handleTCPMessage(c, msg)
	}
}

// handleTCPMessage 处理 tcp 消息
func (s *Server) handleTCPMessage(conn *tcpConn, msg *Message) {
	// 请求消息
	if msg.isRequest {
		t := s.tcptx.new(msg)
		if atomic.CompareAndSwapInt32(&t.handlingReq, 0, 1) {
			// 回收计数
			atomic.AddInt32(&t.recovery, 1)
			// 启动协程处理
			s.wg.Add(1)
			go s.handleTCPTransactionRequestRoutine(t, conn, msg)
			//
			return
		}
	} else {
		// 响应消息
		t := s.tcptx.get(msg)
		if t != nil {
			// 预处理 1xx
			if msg.StartLine[1] != "" && msg.StartLine[1][0] == '1' {
				s.msgPool.Put(msg)
				return
			}
			// 第一个消息
			if atomic.CompareAndSwapInt32(&t.handlingRes, 0, 1) {
				// 回收计数
				atomic.AddInt32(&t.recovery, 1)
				// 启动协程处理
				s.wg.Add(1)
				go s.handleTCPTransactionResponseRoutine(t, conn, msg)
				//
				return
			}
		}
	}
	// 已经在处理
	s.msgPool.Put(msg)
}

// closeTCPConn 关闭并移除 c
func (s *Server) closeTCPConn(c *tcpConn) {
	s.tcplock.Lock()
	delete(s.tcpConns, c.key)
	s.tcplock.Unlock()
	// 关闭底层连接
	c.conn.Close()
}

// getTCPConn 返回 rAddr 对应的客户端连接，如果没有，就创建新的连接(tcp)
func (s *Server) getTCPConn(ctx context.Context, rAddr *net.TCPAddr) (*tcpConn, error) {
	var key connKey
	key.Init(rAddr.IP, rAddr.Port)
	// 获取存在的连接
	var c *tcpConn
	s.tcplock.RLock()
	c = s.tcpConns[key]
	s.tcplock.RUnlock()
	if c != nil {
		return c, nil
	}
	// 没有，创建连接
	dialer := new(net.Dialer)
	conn, err := dialer.DialContext(ctx, rAddr.Network(), rAddr.String())
	if err != nil {
		return nil, err
	}
	cc := s.newTCPConn(conn.(*net.TCPConn))
	// 再次看看有没有并发创建了
	s.tcplock.Lock()
	c = s.tcpConns[key]
	if c == nil {
		s.tcpConns[key] = cc
		s.tcplock.Unlock()
		//
		s.wg.Add(1)
		go s.readTCPRoutine(cc)
		return cc, nil
	}
	// 已经有其他创建了
	s.tcplock.Unlock()
	// 关闭这个
	conn.Close()
	// 返回并发创建的那个
	return c, nil
}

// newTCPConn 根据 conn 创建并返回新的 tcpConn
func (s *Server) newTCPConn(conn *net.TCPConn) *tcpConn {
	rAddr := conn.RemoteAddr().(*net.TCPAddr)
	// 初始化
	c := &tcpConn{
		conn:         conn,
		writeTimeout: s.WriteTimeout,
		remoteIP:     rAddr.IP.String(),
		remotePort:   strconv.Itoa(rAddr.Port),
	}
	c.remoteAddr = fmt.Sprintf("%s:%s", c.remoteIP, c.remotePort)
	c.key.Init(rAddr.IP, rAddr.Port)
	// 返回
	return c
}
