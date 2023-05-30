package util

import (
	"sync"
)

// SafeChan 用于安全的并发写入和关闭
// 不至于 panic
type SafeChan[T any] struct {
	c  chan T
	m  sync.Mutex
	ok bool
}

// NewSafeChan 返回新的 SafeChan
func NewSafeChan[T any](len int) *SafeChan[T] {
	c := new(SafeChan[T])
	c.c = make(chan T, len)
	c.ok = true
	return c
}

// Close 关闭
func (s *SafeChan[T]) Close() {
	s.m.Lock()
	if !s.ok {
		s.m.Unlock()
		return
	}
	close(s.c)
	s.ok = false
	s.m.Unlock()
}

// Send 写入
func (s *SafeChan[T]) Send(v T) bool {
	s.m.Lock()
	// 已经关闭
	if !s.ok {
		s.m.Unlock()
		return false
	}
	// 写入
	select {
	case s.c <- v:
		s.m.Unlock()
		return true
	default:
		s.m.Unlock()
		return false
	}
}
