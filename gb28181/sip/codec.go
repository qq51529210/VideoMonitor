package sip

import (
	"errors"
	"io"
)

const (
	// 默认的 reader 缓存大小
	defaultReaderBuffer = 1024 * 4
)

var (
	// ErrLargeLine 表示读取一行数据大小超过了 reader 设置的缓存
	ErrLargeLine = errors.New("large line")
)

// Reader 主要用于解析以 crlf 结尾的一行数据
type Reader interface {
	// 返回一行字符串，或者错误
	ReadLine() (string, error)
	// 如果缓存有数据，读取缓存的；没有则从底层 reader 读取
	io.Reader
}

// NewReader 返回 n 字节缓存的 Reader ，它从 r 读取数据。
func NewReader(r io.Reader, n int) Reader {
	if n < 1 {
		n = defaultReaderBuffer
	}
	return &reader{r: r, buf: make([]byte, n)}
}

// reader 实现了 Reader 接口
type reader struct {
	r io.Reader
	// 缓存
	buf []byte
	// 有效数据起始下标
	begin int
	// 有效数据终止下标
	end int
	// 已经解析的下标
	parsed int
}

func (r *reader) Reset(reader io.Reader) {
	r.r = reader
	r.begin = 0
	r.end = 0
	r.parsed = 0
}

// func (r *reader) ReadLine() (string, error) {
// 	i := 0
// 	for {
// 		// 解析
// 		for r.parsed < r.end {
// 			if r.buf[r.parsed] == '\r' {
// 				i = r.parsed + 1
// 				if i == r.end {
// 					break
// 				}
// 				if r.buf[i] == '\n' {
// 					line := string(r.buf[r.begin:r.parsed])
// 					r.parsed = i + 1
// 					r.begin = r.parsed
// 					r.checkEmpty()
// 					return line, nil
// 				}
// 			}
// 			r.parsed++
// 		}
// 		// 是否需要读取数据
// 		if r.parsed == r.end {
// 			// 这一行数据太大了
// 			if r.end == len(r.buf) {
// 				if r.begin == 0 {
// 					return "", ErrLargeLine
// 				}
// 				// 缓存向前移
// 				copy(r.buf, r.buf[r.begin:r.end])
// 				r.end -= r.begin
// 				r.parsed -= r.begin
// 				r.begin = 0
// 			}
// 			// 继续读
// 			n, err := r.r.Read(r.buf[r.end:])
// 			if err != nil {
// 				return "", err
// 			}
// 			r.end += n
// 		}
// 	}
// }

func (r *reader) ReadLine() (string, error) {
	for {
		// 尝试查找
		for r.parsed < r.end {
			// 找到回车
			if r.buf[r.parsed] == '\n' {
				i := r.parsed - 1
				// 找到换行
				if r.buf[i] == '\r' {
					line := string(r.buf[r.begin:i])
					r.parsed++
					r.begin = r.parsed
					r.checkEmpty()
					return line, nil
				}
			}
			r.parsed++
		}
		// 缓存中没有，那么要读取数据了
		if r.end == len(r.buf) {
			// 这一行数据太大了，指定的缓存装不下了
			if r.begin == 0 {
				return "", ErrLargeLine
			}
			// 缓存向前移
			if r.begin > 0 {
				copy(r.buf, r.buf[r.begin:r.end])
				r.end -= r.begin
				r.parsed -= r.begin
				r.begin = 0
			}
		}
		// 继续读
		n, err := r.r.Read(r.buf[r.end:])
		if err != nil {
			return "", err
		}
		r.end += n
	}
}

func (r *reader) Read(b []byte) (int, error) {
	if r.begin == r.end {
		return r.r.Read(b)
	}
	n := copy(b, r.buf[r.begin:r.end])
	r.begin += n
	r.parsed += n
	r.checkEmpty()
	return n, nil
}

func (r *reader) checkEmpty() {
	if r.begin == r.end {
		r.begin = 0
		r.parsed = 0
		r.end = 0
	}
}

// Writer 主要用来格式化 Message
type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
	io.Writer
}

// Formatter 格式化的接口
type Formatter interface {
	FormatTo(writer Writer) error
}

// formatHeaderTo 为了减少内存逃逸
func formatHeaderTo(writer Writer, key, value string) error {
	_, err := writer.WriteString(key)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(": ")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(value)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\r\n")
	return err
}

// formatHeaderTo2 为了减少内存逃逸
func formatHeaderTo2(writer Writer, key string, formatter Formatter) error {
	_, err := writer.WriteString(key)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(": ")
	if err != nil {
		return err
	}
	err = formatter.FormatTo(writer)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\r\n")
	return err
}
