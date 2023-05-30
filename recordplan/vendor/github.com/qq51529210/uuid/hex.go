package uuid

var (
	hexLowerTable = []byte("0123456789abcdef")
	hexUpperTable = []byte("0123456789ABCDEF")
)

// hex
func (d *UUID) hex(table, buff []byte) {
	i, j := 0, 0
	for i < 4 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
	}
	buff[j] = '-'
	j++
	for i < 10 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
		buff[j] = '-'
		j++
	}
	for i < 16 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
	}
}

func (d *UUID) hexWithoutHyphen(table, buff []byte) {
	i, j := 0, 0
	for i < 4 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
	}
	for i < 10 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
	}
	for i < 16 {
		buff[j] = table[d.b[i]>>4]
		j++
		buff[j] = table[d.b[i]&0x0f]
		j++
		i++
	}
}
