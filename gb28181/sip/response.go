package sip

import "context"

// Response 表示响应消息
type Response struct {
	transaction
	*Message
	Conn
	context.Context
	s *Server
}
