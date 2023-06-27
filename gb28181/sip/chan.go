package sip

import "sync/atomic"

// safeChan 用于安全的关闭 ch
type safeChan[T any] struct {
	c  chan T
	ok int32
}

// close 安全的关闭 ch
func (c *safeChan[T]) Close() bool {
	if atomic.CompareAndSwapInt32(&c.ok, 1, 0) {
		close(c.c)
		return true
	}
	return false
}

// init 初始化 ch
func (c *safeChan[T]) Init(n int) {
	c.ok = 1
	c.c = make(chan T, n)
}

// isOK 返回状态
func (c *safeChan[T]) IsOK() bool {
	return atomic.LoadInt32(&c.ok) == 1
}
