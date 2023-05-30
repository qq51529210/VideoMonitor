package uuid

// V4 初始化 v4 版本的，随机数版本有可能相同哦
func (d *UUID) V4() {
	v4Rand.Read(d.b[:16])
	// version
	d.b[6] &= 0b00001111
	d.b[6] |= 0b01000000
	// variant
	d.b[8] &= 0b00011111
	d.b[8] |= 0b10000000
}

// UpperV4 返回十六进制大写字符串
func UpperV4() string {
	var id UUID
	id.V4()
	return id.UpperString()
}

// LowerV4 返回十六进制小写字符串
func LowerV4() string {
	var id UUID
	id.V4()
	return id.LowerString()
}

// UpperV4WithoutHyphen 返回十六进制大写字符串，不包含 '-' 字符
func UpperV4WithoutHyphen() string {
	var id UUID
	id.V4()
	return id.UpperStringWithoutHyphen()
}

// LowerV4WithoutHyphen 返回十六进制小写字符串，不包含 '-' 字符
func LowerV4WithoutHyphen() string {
	var id UUID
	id.V4()
	return id.LowerStringWithoutHyphen()
}
