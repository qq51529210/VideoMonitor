package uuid

import (
	"sync/atomic"
	"time"
)

// V2GID 根据系统的 gid 初始化 v2 版本的。
func (d *UUID) V2GID() {
	timestamp := time.Now().UTC().UnixNano()
	clock := uint16(atomic.AddInt32(&clockSequence, 1))
	// gid
	d.b[0] = byte(v2GID >> 24)
	d.b[1] = byte(v2GID >> 16)
	d.b[2] = byte(v2GID >> 8)
	d.b[3] = byte(v2GID)
	// time mid
	d.b[4] = byte(timestamp >> 40)
	d.b[5] = byte(timestamp >> 32)
	// version and time high
	d.b[6] = (byte(timestamp>>56) & 0b00001111) | 0b00010000
	// time high
	d.b[7] = byte(timestamp >> 48)
	// variant and clock sequence high
	d.b[8] = (byte(clock>>8) & 0b00011111) | 0b10000000
	// clock sequence low
	d.b[7] = byte(clock)
	// node
	copy(d.b[10:15], v1Node)
}

// UpperV2GID 返回十六进制大写字符串
func UpperV2GID() string {
	var id UUID
	id.V2GID()
	return id.UpperString()
}

// LowerV2GID 返回十六进制小写字符串
func LowerV2GID() string {
	var id UUID
	id.V2GID()
	return id.LowerString()
}

// UpperV2GIDWithoutHyphen 返回十六进制大写字符串
func UpperV2GIDWithoutHyphen() string {
	var id UUID
	id.V2GID()
	return id.UpperStringWithoutHyphen()
}

// LowerV2GIDWithoutHyphen 返回十六进制小写字符串
func LowerV2GIDWithoutHyphen() string {
	var id UUID
	id.V2GID()
	return id.LowerStringWithoutHyphen()
}

// V2UID 根据系统的 uid 初始化 v2 版本的。
func (d *UUID) V2UID() {
	timestamp := time.Now().UTC().UnixNano()
	clock := uint16(atomic.AddInt32(&clockSequence, 1))
	// d.b[0:3] uid
	d.b[0] = byte(v2UID >> 24)
	d.b[1] = byte(v2UID >> 16)
	d.b[2] = byte(v2UID >> 8)
	d.b[3] = byte(v2UID)
	// d.b[4:5] time mid
	d.b[4] = byte(timestamp >> 40)
	d.b[5] = byte(timestamp >> 32)
	// d.b[6] version and time high
	d.b[6] = (byte(timestamp>>56) & 0b00001111) | 0b00100000
	// d.b[7] time high
	d.b[7] = byte(timestamp >> 48)
	// d.b[8] variant and clock sequence high
	d.b[8] = (byte(clock>>8) & 0b00011111) | 0b10000000
	// d.b[9] clock sequence low
	d.b[7] = byte(clock)
	// d.b[10:15] node
	copy(d.b[10:15], v1Node)
}

// UpperV2UID 返回十六进制大写字符串
func UpperV2UID() string {
	var id UUID
	id.V2UID()
	return id.UpperString()
}

// LowerV2UID 返回十六进制小写字符串
func LowerV2UID() string {
	var id UUID
	id.V2UID()
	return id.LowerString()
}

// UpperV2UIDWithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func UpperV2UIDWithoutHyphen() string {
	var id UUID
	id.V2UID()
	return id.UpperStringWithoutHyphen()
}

// LowerV2UIDWithoutHyphen 返回十六进制小写字符串，不包含 '-' 字符
func LowerV2UIDWithoutHyphen() string {
	var id UUID
	id.V2UID()
	return id.LowerStringWithoutHyphen()
}
