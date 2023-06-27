package sip

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qq51529210/log"
)

// tcpTransactions 表示 tcp 事务表
type tcpTransactions struct {
	sync.RWMutex
	sync.Pool
	t map[string]*tcpTransaction
}

// init 初始化
func (t *tcpTransactions) init() {
	t.t = make(map[string]*tcpTransaction)
	t.New = func() any { return new(tcpTransaction) }
}

// new 根据 msg 返回 tx ，处理请求消息时使用
func (t *tcpTransactions) new(msg *Message) *tcpTransaction {
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
	tt = t.Get().(*tcpTransaction)
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
func (t *tcpTransactions) get(msg *Message) *tcpTransaction {
	key := msg.TransactionKey()
	// 获取
	t.RLock()
	tt := t.t[key]
	t.RUnlock()
	// 返回
	return tt
}

// rm 移除 tt
func (t *tcpTransactions) rm(tt *tcpTransaction) {
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

// tcpTransaction 表示一个 tcp 事务
type tcpTransaction struct {
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
	// 退出信号，用于退出主动发起事务请求超时清理协程，
	// 在收到响应消息的时候，提前退出
	quit safeChan[struct{}]
	// 用于控制回收
	recovery int32
}

func (t *tcpTransaction) Key() string {
	return t.key
}

// writeMessage 格式化 msg 到 writeData
func (t *tcpTransaction) writeMessage(conn Conn, msg *Message) error {
	t.writeData.Reset()
	msg.FormatTo(&t.writeData)
	log.DebugfTrace(t.key, "write tcp %s:%s\n%s", conn.RemoteIP(), conn.RemotePort(), t.writeData.String())
	return conn.write(t.writeData.Bytes())
}

// handleTCPTransactionRequestRoutine 处理 tcp 的事务请求消息
func (s *Server) handleTCPTransactionRequestRoutine(t *tcpTransaction, conn Conn, msg *Message) {
	// 退出清理
	defer func() {
		atomic.StoreInt32(&t.handlingReq, 0)
		// 回收
		s.msgPool.Put(msg)
		// 移除
		s.tcptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 回调处理
	s.Handler.HandleRequest(&Request{transaction: t, Conn: conn, Message: msg, s: s, Context: t.ctx})
}

// handleTCPTransactionResponseRoutine 处理 tcp 的事务响应消息
func (s *Server) handleTCPTransactionResponseRoutine(t *tcpTransaction, conn Conn, msg *Message) {
	// 退出清理
	defer func() {
		atomic.StoreInt32(&t.handlingRes, 0)
		// 回收
		s.msgPool.Put(msg)
		// 移除
		s.tcptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 退出清理协程
	t.quit.Close()
	// 回调处理
	s.Handler.HandleResponse(&Response{transaction: t, Conn: conn, Message: msg, s: s, Context: t.ctx})
}

// clearTCPTransactionRoutine 主要用于清理主动发起请求的 tcp 事务
func (s *Server) clearTCPTransactionRoutine(ctx context.Context, t *tcpTransaction) {
	// 退出清理
	defer func() {
		// 移除
		s.tcptx.rm(t)
		// 协程结束
		s.wg.Done()
	}()
	// 等待退出信号
	select {
	case <-ctx.Done():
	case <-t.quit.c:
	}
}

// handleTCPTransactionTimeoutRoutine 主要用于清理主动发起请求的 tcp 事务。
// timeout 用于控制整个事务的超时，小于 0 则使用 s.TransactionTimeout 。
func (s *Server) handleTCPTransactionTimeoutRoutine(t *tcpTransaction, timeout time.Duration) {
	// 超时
	if timeout < 1 {
		timeout = s.WriteTimeout
	}
	// 计时器
	timer := time.NewTimer(timeout)
	// 退出清理
	defer func() {
		// 移除
		s.tcptx.rm(t)
		// 计时器
		timer.Stop()
		// 协程结束
		s.wg.Done()
	}()
	// 等待退出信号
	select {
	case <-timer.C:
	case <-t.quit.c:
	}
}
