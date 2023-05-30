package uuid

import "hash"

// V3 初始化 v3 版本的
func (d *UUID) V3(namespace, name []byte) {
	h := v3Pool.Get().(hash.Hash)
	h.Reset()
	h.Write(namespace)
	h.Write(name)
	h.Sum(d.b[:0])
	v3Pool.Put(h)
	// version
	d.b[6] &= 0b00001111
	d.b[6] |= 0b00110000
	// variant
	d.b[8] &= 0b00011111
	d.b[8] |= 0b10000000
}

// UpperV3 返回十六进制大写字符串
func UpperV3(namespace, data []byte) string {
	var id UUID
	id.V3(namespace, data)
	return id.UpperString()
}

// LowerV3 返回十六进制小写字符串
func LowerV3(namespace, data []byte) string {
	var id UUID
	id.V3(namespace, data)
	return id.LowerString()
}

// UpperV3WithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func UpperV3WithoutHyphen(namespace, data []byte) string {
	var id UUID
	id.V3(namespace, data)
	return id.UpperStringWithoutHyphen()
}

// LowerV3WithoutHyphen 返回十六进制小写字符串，不包含 '-' 字符
func LowerV3WithoutHyphen(namespace, data []byte) string {
	var id UUID
	id.V3(namespace, data)
	return id.LowerStringWithoutHyphen()
}
