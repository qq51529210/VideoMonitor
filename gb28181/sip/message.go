package sip

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// // PhraseError 表示解析的错误的短语
// type PhraseError string

// func (e PhraseError) Error() string {
// 	return string(e)
// }

var (
	errStartLineFormat = errors.New("error start line format")
	errReadingBody     = errors.New("error reading body")
)

var (
	// ErrLargeMessage 表示读取消息的数据大小超过了设置的值
	ErrLargeMessage = errors.New("large message")
)

// Message 表示 start line + header + body 的结构。
type Message struct {
	StartLine [3]string
	// 顺序的
	Header Header
	// 由于 sip 协议的数据都特别小，这里就直接使用内存。
	Body bytes.Buffer
	// 是否请求消息
	isRequest bool
	// 事务的key
	tKey strings.Builder
}

// FormatTo 格式化 header 和 body（如果 body 不为空）到 writer 中。
// Content-Length 字段是根据 body 的大小自动添加的。
func (m *Message) FormatTo(writer Writer) error {
	// start line ，这样写是为了减少内存逃逸
	part := make([]string, 6)
	part[0] = m.StartLine[0]
	part[1] = " "
	part[2] = m.StartLine[1]
	part[3] = " "
	part[4] = m.StartLine[2]
	part[5] = "\r\n"
	for i := 0; i < len(part); i++ {
		_, err := writer.WriteString(part[i])
		if err != nil {
			return err
		}
	}
	// Content-Length
	m.Header.contentLength = int64(m.Body.Len())
	// header
	err := m.Header.FormatTo(writer)
	if err != nil {
		return err
	}
	// header 和 body 的空行
	_, err = writer.WriteString("\r\n")
	if err != nil {
		return err
	}
	// body
	if m.Header.contentLength > 0 {
		_, err = writer.Write(m.Body.Bytes())
	}
	return err
}

// ParseFrom 从 reader 中读取并解析一个完整的 Message ，
// max 表示消息的最大字节，小于 1 表示不限制的读取。
// start line 的 [0][2] 转换为大写
// CSeq 的 method 转换为大写
func (m *Message) ParseFrom(reader Reader, max int) (err error) {
	// start line
	max, err = m.parseStartLine(reader, max)
	if err != nil {
		return err
	}
	// header
	max, err = m.Header.ParseFrom(reader, max)
	if err != nil {
		return err
	}
	// body
	if m.Header.contentLength > 0 {
		if m.Header.contentLength > int64(max) {
			return ErrLargeMessage
		}
		_, err = io.CopyN(&m.Body, reader, m.Header.contentLength)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseStartLine 解析 start line
func (m *Message) parseStartLine(reader Reader, max int) (int, error) {
	line, err := reader.ReadLine()
	if err != nil {
		return max, err
	}
	max -= len(line)
	if max < 0 {
		return max, ErrLargeMessage
	}
	// 0
	line = strings.TrimSpace(line)
	i := strings.Index(line, " ")
	if i < 0 {
		return max, errStartLineFormat
	}
	m.StartLine[0] = strings.ToUpper(line[:i])
	// 1
	line = strings.TrimSpace(line[i+1:])
	i = strings.Index(line, " ")
	if i < 0 {
		return max, errStartLineFormat
	}
	m.StartLine[1] = strings.ToUpper(line[:i])
	// 2
	m.StartLine[2] = strings.TrimSpace(line[i+1:])
	// 检查
	if m.StartLine[2] == SIPVersion {
		m.isRequest = true
	} else if m.StartLine[0] != SIPVersion {
		return max, errStartLineFormat
	}
	return max, nil
}

// TransactionKey 返回事务的 key
func (m *Message) TransactionKey() string {
	if m.tKey.Len() < 1 {
		m.tKey.WriteString(m.Header.CSeq.Method)
		m.tKey.WriteString(m.Header.CallID)
		m.tKey.WriteString(m.Header.Via[0].Branch)
	}
	return m.tKey.String()
}

// String 返回格式化后的字符串。
func (m *Message) String() string {
	var str strings.Builder
	m.FormatTo(&str)
	return str.String()
}

// Reset 重置所有字段
func (m *Message) Reset() {
	m.isRequest = false
	m.tKey.Reset()
	for i := 0; i < 3; i++ {
		m.StartLine[i] = ""
	}
	m.Header.Reset()
	m.Body.Reset()
}

// KeepBasicHeaders 重置其他的数据，只保留基本的头字段
func (m *Message) KeepBasicHeaders() {
	m.tKey.Reset()
	m.Header.KeepBasic()
	m.Body.Reset()
}

// CopyTo 拷贝字段到 mm 。
func (m *Message) CopyTo(mm *Message) {
	mm.StartLine = m.StartLine
	m.Header.CopyTo(&mm.Header)
	mm.Body.Write(m.Body.Bytes())
}

// IsRequest 返回是否请求消息
func (m *Message) IsRequest() bool {
	return m.isRequest
}

// IsStatus 返回是否与 code 相等，因为 StartLine[1] 不好记忆
func (m *Message) IsStatus(code string) bool {
	return m.StartLine[1] == code
}

// InitStartLineOfRequest 初始化请求的 start line
func (m *Message) InitStartLineOfRequest(method string, uri string) {
	m.StartLine[0] = method
	m.StartLine[1] = uri
	m.StartLine[2] = SIPVersion
}

// InitStartLineOfResponse 初始化响应的 start line
func (m *Message) InitStartLineOfResponse(code, reason string) {
	m.StartLine[0] = SIPVersion
	m.StartLine[1] = string(code)
	m.StartLine[2] = reason
}

// RequestMethod 返回 StartLine[0]
func (m *Message) RequestMethod() string {
	return m.StartLine[0]
}

// RequestURI 返回 StartLine[1]
func (m *Message) RequestURI() string {
	return m.StartLine[1]
}

// RequestVersion 返回 StartLine[2]
func (m *Message) RequestVersion() string {
	return m.StartLine[2]
}

// ResponseVersion 返回 StartLine[0]
func (m *Message) ResponseVersion() string {
	return m.StartLine[0]
}

// ResponseStatus 返回 StartLine[1]
func (m *Message) ResponseStatus() string {
	return m.StartLine[1]
}

// ResponsePhrase 返回 StartLine[2]
func (m *Message) ResponsePhrase() string {
	return m.StartLine[2]
}
