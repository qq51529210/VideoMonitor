package sip

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qq51529210/log"
)

// udpTransactions 表示 udp 事务表
type udpTransactions struct {
	sync.RWMutex
	sync.Pool
	t map[string]*udpTransaction
}

// init 初始化
func (t *udpTransactions) init() {
	t.t = make(map[string]*udpTransaction)
	t.New = func() any { return new(udpTransaction) }
}

// new 根据 msg 返回 tx ，处理请求消息时使用
func (t *udpTransactions) new(msg *Message) *udpTransaction {
	key := msg.TransactionKey()
	// 查询
	t.Lock()
	tt := t.t[key]
	// 已存在
	if tt != nil {
		t.Unlock()
		//
		return tt
	}
	// 新的
	tt = t.Get().(*udpTransaction)
	t.t[key] = tt
	// 初始化字段
	tt.key = key
	tt.ctx = nil
	tt.quit.Init(0)
	tt.writeData.Reset()
	t.Unlock()
	// 返回
	return tt
}

// get 根据 msg 返回 tx ，处理响应消息时使用
func (t *udpTransactions) get(msg *Message) *udpTransaction {
	key := msg.TransactionKey()
	//
	t.RLock()
	tt := t.t[key]
	t.RUnlock()
	// 返回
	return tt
}

// rm 移除 tt
func (t *udpTransactions) rm(tt *udpTransaction) {
	// 移除
	t.Lock()
	delete(t.t, tt.key)
	t.Unlock()
	// 通知
	tt.quit.Close()
	// 回收
	if atomic.AddInt32(&tt.recovery, -1) == 0 {
		t.Put(tt)
	}
}

// udpTransaction 表示一个 udp 事务
type udpTransaction struct {
	// 事务表的 key
	key string
	// 是否已经启动协程处理请求消息
	handlingReq int32
	// 是否已经启动协程处理响应消息
	handlingRes int32
	// 调用者上下文数据
	ctx context.Context
	// 发送消息数据
	writeData bytes.Buffer
	// 退出信号，用于退出主动发起事务请求消息重发协程，
	// 在收到响应消息的时候，提前退出
	quit safeChan[struct{}]
	// 用于控制回收
	recovery int32
}

func (t *udpTransaction) Key() string {
	return t.key
}

// writeMessage 格式化 msg 到 writeData
func (t *udpTransaction) writeMessage(conn Conn, msg *Message) error {
	t.writeData.Reset()
	msg.FormatTo(&t.writeData)
	log.DebugfTrace(t.key, "write udp %s:%s\n%s", conn.RemoteIP(), conn.RemotePort(), t.writeData.String())
	return nil
}

// handleUDPTransactionRequestRoutine 处理 udp 的事务请求消息
func (s *Server) handleUDPTransactionRequestRoutine(t *udpTransaction, conn Conn, msg *Message) {
	// 计时器
	var rtoTimer *time.Timer
	// 退出清理
	defer func() {
		atomic.StoreInt32(&t.handlingReq, 0)
		// 计时器
		if rtoTimer != nil {
			rtoTimer.Stop()
		}
		// 回收
		s.msgPool.Put(msg)
		// 移除
		s.udptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 回调处理
	s.Handler.HandleRequest(&Request{transaction: t, Conn: conn, Message: msg, s: s, Context: t.ctx})
	// 如果有数据，发送直到超时
	if t.writeData.Len() > 0 {
		startTime := time.Now()
		rtoTimer = time.NewTimer(0)
		for s.isOK() {
			now := <-rtoTimer.C
			// 事务超时
			if now.Sub(startTime) > s.WriteTimeout {
				return
			}
			// 发送消息
			err := conn.write(t.writeData.Bytes())
			if err != nil {
				log.ErrorTrace(t.key, err)
				return
			}
			log.DebugTrace(t.key, "retransmission")
			// 重置计时器
			rtoTimer.Reset(s.rto)
		}
	}
}

// handleUDPTransactionResponseRoutine 处理 udp 的事务响应消息
func (s *Server) handleUDPTransactionResponseRoutine(t *udpTransaction, conn Conn, msg *Message) {
	// 退出清理
	defer func() {
		atomic.StoreInt32(&t.handlingRes, 0)
		// 回收
		s.msgPool.Put(msg)
		// 移除
		if t.ctx != nil {
			<-t.ctx.Done()
		}
		s.udptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 退出超时重发协程
	t.quit.Close()
	// 回调处理
	s.Handler.HandleResponse(&Response{transaction: t, Conn: conn, Message: msg, s: s, Context: t.ctx})
}

// udpTransactionRetransmissionRoutine 用于在协程中发送 udp 消息。
// 主动发起请求的时候使用，ctx 用于控制协程退出
func (s *Server) udpTransactionRetransmissionRoutine(ctx context.Context, conn Conn, t *udpTransaction) {
	// 计时器
	rtoTimer := time.NewTimer(0)
	// 退出清理
	defer func() {
		// 计时器
		rtoTimer.Stop()
		// 移除
		s.udptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 开始发送时间
	startTime := time.Now()
	for s.isOK() {
		select {
		case <-ctx.Done():
			// 调用结束通知
			return
		case <-t.quit.c:
			return
		case now := <-rtoTimer.C:
			// 事务超时
			if now.Sub(startTime) > s.WriteTimeout {
				return
			}
			// 发送消息
			err := conn.write(t.writeData.Bytes())
			if err != nil {
				log.ErrorTrace(t.key, err)
				return
			}
			log.DebugTrace(t.key, "retransmission")
			// 重置计时器
			rtoTimer.Reset(s.rto)
		}
	}
}

// udpTransactionRetransmissionTimeoutRoutine 用于在协程中发送 udp 消息。
// timeout 用于控制整个事务的超时，小于 0 则使用 s.TransactionTimeout 。
func (s *Server) udpTransactionRetransmissionTimeoutRoutine(conn Conn, t *udpTransaction, timeout time.Duration) {
	// 超时
	if timeout < 1 {
		timeout = s.WriteTimeout
	}
	// 超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	t.ctx = ctx
	// 启动
	s.udpTransactionRetransmissionRoutine(ctx, conn, t)
}
