package uuid

import (
	"crypto/md5"
	"crypto/sha1"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

var (
	clockSequence int32
	v1Node        = make([]byte, 6, 6)
	v2GID         = os.Getgid()
	v2UID         = os.Getuid()
	v3Pool        sync.Pool
	v4Rand        = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	v5Pool        sync.Pool
)

func init() {
	// 使用网卡初始化 v1node
	ifs, err := net.Interfaces()
	if nil != err {
		panic(err)
	}
	ok := false
	for i := 0; i < len(ifs); i++ {
		if len(ifs[i].HardwareAddr) >= 6 {
			copy(v1Node[0:], ifs[i].HardwareAddr)
			ok = true
			break
		}
	}
	// 没有就随机
	if !ok {
		v4Rand.Read(v1Node)
	}
	v3Pool.New = func() interface{} {
		return md5.New()
	}
	v5Pool.New = func() interface{} {
		return sha1.New()
	}
}

// SetV1Node 重新设置 version-1 node 的值，默认是找到的第一个网卡的 MAC 地址。
func SetV1Node(node []byte) {
	copy(v1Node, node)
}

// UUID 表示一个 uuid
type UUID struct {
	b [16]byte
}

// String 返回十六进制大写字符串
func (d *UUID) String() string {
	return d.UpperString()
}

// UpperString 返回十六进制大写字符串
func (d *UUID) UpperString() string {
	var str [36]byte
	d.hex(hexUpperTable, str[:])
	return string(str[:])
}

// LowerString 返回十六进制小写写字符串
func (d *UUID) LowerString() string {
	var str [36]byte
	d.hex(hexLowerTable, str[:])
	return string(str[:])
}

// UpperStringWithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func (d *UUID) UpperStringWithoutHyphen() string {
	var str [32]byte
	d.hexWithoutHyphen(hexLowerTable, str[:])
	return string(str[:])
}

// LowerStringWithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func (d *UUID) LowerStringWithoutHyphen() string {
	var str [32]byte
	d.hexWithoutHyphen(hexLowerTable, str[:])
	return string(str[:])
}
