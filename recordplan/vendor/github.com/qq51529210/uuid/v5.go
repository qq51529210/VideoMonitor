package uuid

import "hash"

// V5 初始化 v1 版本的
func (d *UUID) V5(namespace, name []byte) {
	h := v5Pool.Get().(hash.Hash)
	h.Reset()
	h.Write(namespace)
	h.Write(name)
	var sum [20]byte
	h.Sum(sum[:0])
	v5Pool.Put(h)
	copy(d.b[:16], sum[:])
	// version
	d.b[6] &= 0b00001111
	d.b[6] |= 0b01010000
	// variant
	d.b[8] &= 0b00011111
	d.b[8] |= 0b10000000
}

// UpperV5 返回十六进制大写字符串
func UpperV5(namespace, data []byte) string {
	var id UUID
	id.V5(namespace, data)
	return id.UpperString()
}

// LowerV5 返回十六进制小写字符串
func LowerV5(namespace, data []byte) string {
	var id UUID
	id.V5(namespace, data)
	return id.LowerString()
}

// UpperV5WithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func UpperV5WithoutHyphen(namespace, data []byte) string {
	var id UUID
	id.V5(namespace, data)
	return id.UpperStringWithoutHyphen()
}

// LowerV5WithoutHyphen 返回十六进制小写字符串，不包含 '-' 字符
func LowerV5WithoutHyphen(namespace, data []byte) string {
	var id UUID
	id.V5(namespace, data)
	return id.LowerStringWithoutHyphen()
}
