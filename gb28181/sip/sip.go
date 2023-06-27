package sip

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/qq51529210/uuid"
)

const (
	// MethodRegister 表示 REGISTER 消息
	MethodRegister string = "REGISTER"
	// MethodInvite 表示 INVITE 消息
	MethodInvite string = "INVITE"
	// MethodACK 表示 ACK 消息
	MethodACK string = "ACK"
	// MethodBye 表示 BYE 消息
	MethodBye string = "BYE"
	// MethodMessage 表示 MESSAGE 消息
	MethodMessage string = "MESSAGE"
	// MethodNotify 表示 NOTIFY 消息
	MethodNotify string = "NOTIFY"
	// MethodSubscribe 表示 SUBSCRIBE 消息
	MethodSubscribe string = "SUBSCRIBE"
	// MethodInfo 表示 INFO 消息
	MethodInfo string = "INFO"
)

const (
	// BranchPrefix 事务必须以z9hG4bK开头
	BranchPrefix = "z9hG4bK"
	// SIPVersion 表示支持的sip协议版本
	SIPVersion = "SIP/2.0"
)

var (
	// buffer 缓存池
	bufPool sync.Pool
	// CSeq 的递增 SN
	// sn32 = uint32(time.Now().Unix())
	sn32   = uint32(0)
	cseq32 = uint32(0)
)

func init() {
	bufPool.New = func() any {
		return bytes.NewBuffer(nil)
	}
}

// GetSN 返回全局递增的 sn
func GetSN() uint32 {
	return atomic.AddUint32(&sn32, 1)
}

// GetCSeq 返回全局递增的 sn
func GetCSeq() uint32 {
	return atomic.AddUint32(&cseq32, 1)
}

// GetSNString 返回字符串形式的全局递增的 sn
func GetSNString() string {
	return fmt.Sprintf("%d", atomic.AddUint32(&sn32, 1))
}

// NewBranch 使用雪花算法生成 z9hG4bK-xxxxxx
func NewBranch() string {
	return fmt.Sprintf("%s-%s", BranchPrefix, uuid.SnowflakeIDString())
}

// TrimByte 去掉两端的字符。
func TrimByte(str string, left, right byte) string {
	if str == "" {
		return str
	}
	for str[0] == left {
		str = str[1:]
		if str == "" {
			return str
		}
	}
	i := len(str) - 1
	for str[i] == right {
		str = str[:i]
		if str == "" {
			return str
		}
		i = len(str) - 1
	}
	return str
}
