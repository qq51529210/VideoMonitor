package uuid

import (
	"sync/atomic"
	"time"
)

// V1 初始化 v1 版本的
func (d *UUID) V1() {
	timestamp := time.Now().UTC().UnixNano()
	clock := uint16(atomic.AddInt32(&clockSequence, 1))
	// time low
	d.b[0] = byte(timestamp >> 24)
	d.b[1] = byte(timestamp >> 16)
	d.b[2] = byte(timestamp >> 8)
	d.b[3] = byte(timestamp)
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
	d.b[9] = byte(clock)
	// node
	copy(d.b[10:15], v1Node)
}

// UpperV1 返回十六进制大写字符串
func UpperV1() string {
	var id UUID
	id.V1()
	return id.UpperString()
}

// LowerV1 返回十六进制小写字符串
func LowerV1() string {
	var id UUID
	id.V1()
	return id.LowerString()
}

// UpperV1WithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func UpperV1WithoutHyphen() string {
	var id UUID
	id.V1()
	return id.UpperStringWithoutHyphen()
}

// LowerV1WithoutHyphen 返回十六进制小写字符串，不包含 '-' 字符
func LowerV1WithoutHyphen() string {
	var id UUID
	id.V1()
	return id.LowerStringWithoutHyphen()
}
