package sip

import (
	"context"

	"github.com/qq51529210/uuid"
)

// Request 表示请求消息
type Request struct {
	transaction
	*Message
	Conn
	context.Context
	s *Server
}

// Response 改造当前的请求消息，添加 to.tag 和 via 的 rport 和 received
// 变为响应消息，然后使用当前的 conn 发送
func (r *Request) Response(status, phrase string) error {
	if phrase == "" {
		phrase = StatusPhrase(status)
	}
	r.InitStartLineOfResponse(status, phrase)
	if r.Message.Header.Via[0].RProt != nil {
		if r.Header.Via[0].RProt != nil {
			ip := r.RemoteIP()
			r.Header.Via[0].Received = &ip
			port := r.RemotePort()
			r.Header.Via[0].RProt = &port
		}
		// if r.s.received != "" {
		// 	r.Header.Via[0].Received = &r.s.received
		// } else {
		// 	ip := r.RemoteIP()
		// 	r.Header.Via[0].Received = &ip
		// }
		// if r.s.rport != "" {
		// 	r.Header.Via[0].RProt = &r.s.rport
		// } else {
		// 	port := r.RemotePort()
		// 	r.Header.Via[0].RProt = &port
		// }
	}
	// 给 To 加个 tag
	if r.Header.To.Tag == "" {
		r.Header.To.Tag = uuid.SnowflakeIDString()
	}
	// 设置 Agent 和 body
	r.Header.UserAgent = ""
	// 发送
	return r.transaction.writeMessage(r.Conn, r.Message)
}

// ResponseWithoutBody 清空 body 和多余的 header ，调用 Response
func (r *Request) ResponseWithoutBody(status, phrase string) error {
	r.KeepBasicHeaders()
	return r.Response(status, phrase)
}
