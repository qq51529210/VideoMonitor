package sip

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/qq51529210/log"
)

// closeUDP 关闭 udp 连接。
func (s *Server) closeUDP() {
	// 关闭 udp 连接
	if s.udpConn != nil {
		s.udpConn.Close()
	}
}

// listenUDP 初始化 udp 连接，启动 cup*2 个协程用于读取 udp 原始数据包，在 Serve 函数中调用
func (s *Server) listenUDP(port string) error {
	// 初始化地址
	address, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return err
	}
	// 初始化底层连接
	s.udpConn, err = net.ListenUDP("udp", address)
	if err != nil {
		return err
	}
	// 读取协程
	for i := 0; i < runtime.NumCPU()*2; i++ {
		s.wg.Add(1)
		go s.readUDPRoutine(i)
	}
	// 返回
	return nil
}

// readUDPRoutine 读取 udp 数据并处理
func (s *Server) readUDPRoutine(i int) {
	log.Debugf("udp read routine %d start", i)
	defer func() {
		// log.Recover(recover())
		// 日志
		log.Debugf("udp read routine %d end", i)
		// 协程结束
		s.wg.Done()
	}()
	// 开始
	var err error
	reader := NewReader(nil, s.MessageLen).(*reader)
	for s.isOK() {
		data := s.udpData.Get().(*udpData)
		// 读取 udp 数据
		data.n, data.a, err = s.udpConn.ReadFromUDP(data.b)
		if err != nil {
			log.Error(err)
			s.udpData.Put(data)
			continue
		}
		data.i = 0
		reader.Reset(data)
		// 解析并处理
		s.handleUDPData(reader, data)
		// 回收
		s.udpData.Put(data)
	}
}

// handleUDPData 处理 udp 数据
func (s *Server) handleUDPData(reader *reader, data *udpData) {
	// 连接
	var c udpConn
	s.initUDPConn(&c, data.a)
	// 一个 udp 数据包可能有多个消息
	for s.isOK() {
		// 解析
		msg := s.msgPool.Get().(*Message)
		msg.Reset()
		err := msg.ParseFrom(reader, s.MessageLen)
		if err != nil {
			s.msgPool.Put(msg)
			if err != io.EOF {
				log.Errorf("read udp %v %v\n%s", data.a, err, string(data.b[:data.n]))
			}
			return
		}
		// 处理
		s.handleUDPMessage(&c, msg)
	}
}

// handleUDPMessage 处理 udp 消息
func (s *Server) handleUDPMessage(conn *udpConn, msg *Message) {
	if msg.isRequest {
		// 请求消息
		t := s.udptx.new(msg)
		if atomic.CompareAndSwapInt32(&t.handlingReq, 0, 1) {
			// 事务回收计数
			atomic.AddInt32(&t.recovery, 1)
			s.wg.Add(1)
			go s.handleUDPTransactionRequestRoutine(t, conn, msg)
			return
		}
	} else {
		// 响应消息
		t := s.udptx.get(msg)
		if t != nil {
			// 预处理 1xx
			if msg.StartLine[1] != "" && msg.StartLine[1][0] == '1' {
				s.msgPool.Put(msg)
				return
			}
			// 第一个消息
			if atomic.CompareAndSwapInt32(&t.handlingRes, 0, 1) {
				// 事务回收计数
				atomic.AddInt32(&t.recovery, 1)
				s.wg.Add(1)
				go s.handleUDPTransactionResponseRoutine(t, conn, msg)
				return
			}
		}
	}
	// 已经在处理
	s.msgPool.Put(msg)
}

// initUDPConn 初始化 udpConn
func (s *Server) initUDPConn(c *udpConn, addr *net.UDPAddr) {
	c.conn = s.udpConn
	c.remote = addr
	c.remoteIP = addr.IP.String()
	c.remotePort = strconv.Itoa(addr.Port)
	c.remoteAddr = fmt.Sprintf("%s:%s", c.remoteIP, c.remotePort)
}
