package zlm

// record 表示一个录像文件
type record struct {
	// 本地路径
	path   string
	app    string
	stream string
	// 文件大小
	size int64
	// 创建时间
	createTime int64
	// 时长，修改时间-创建时间
	duration int64
}
