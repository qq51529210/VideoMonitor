package sip

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	errHeaderFormat              = errors.New("error header format")
	errHeaderExpiresFormat       = errors.New("error header Expires format")
	errHeaderMaxForwardsFormat   = errors.New("error header Max-Forwards format")
	errHeaderContentLengthFormat = errors.New("error header Content-Length format")
	errMissingHeaderCSeq         = errors.New("missing header CSeq")
	errMissingHeaderCallID       = errors.New("missing header Call-ID")
	errMissingHeaderTo           = errors.New("missing header To")
	errMissingHeaderFrom         = errors.New("missing header From")
	errMissingHeaderVia          = errors.New("missing header Via")
)

// HeaderIntValue 表示 Header 的整型值
type HeaderIntValue[T int | uint | int32 | uint32 | int64 | uint64] struct {
	n T
	s string
}

// Set 设置数值
func (h *HeaderIntValue[T]) Parse(line string) error {
	n, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return err
	}
	h.n = T(n)
	h.s = line
	return nil
}

// Set 设置数值
func (h *HeaderIntValue[T]) Set(n T) {
	h.n = n
	h.s = strconv.FormatInt(int64(n), 10)
}

// Get 返回数值
func (h *HeaderIntValue[T]) Get() T {
	return h.n
}

func (h *HeaderIntValue[T]) String() string {
	return h.s
}

// OK 返回是否有效
func (h *HeaderIntValue[T]) OK() bool {
	return h.s != ""
}

// Header 表示消息的一些必需的头字段
type Header struct {
	Via         []Via
	From        Address
	To          Address
	CallID      string
	CSeq        CSeq
	MaxForwards HeaderIntValue[uint32]
	Contact     URI
	Expires     HeaderIntValue[uint32]
	ContentType string
	UserAgent   string
	// 其他头
	Others []KV
	// 自动
	contentLength int64
}

// CopyTo 将所有字段 copy 到 hh
func (h *Header) CopyTo(hh *Header) {
	hh.Via = hh.Via[:0]
	hh.Via = append(hh.Via, h.Via...)
	h.From.CopyTo(&hh.From)
	h.To.CopyTo(&hh.To)
	hh.CallID = h.CallID
	hh.CSeq.Method = h.CSeq.Method
	hh.CSeq.SN = h.CSeq.SN
	hh.MaxForwards.n = h.MaxForwards.n
	hh.MaxForwards.s = h.MaxForwards.s
	hh.Contact = h.Contact
	hh.ContentType = h.ContentType
	hh.Others = hh.Others[:0]
	hh.Others = append(hh.Others, h.Others...)
	hh.contentLength = h.contentLength
}

// Reset 重置所有字段
func (h *Header) Reset() {
	h.CallID = ""
	h.To.Reset()
	h.From.Reset()
	h.Contact.Reset()
	h.MaxForwards.n = 0
	h.MaxForwards.s = ""
	h.Via = h.Via[:0]
	h.ContentType = ""
	h.UserAgent = ""
	h.Others = h.Others[:0]
	h.contentLength = 0
}

// KeepBasic 重置 contact、contentType、useragent、other
func (h *Header) KeepBasic() {
	h.Contact.Reset()
	h.ContentType = ""
	h.UserAgent = ""
	h.ResetOther()
}

// ResetOther 重置 others
func (h *Header) ResetOther() {
	h.Others = h.Others[:0]
}

// GetOther 返回指定 other，index 表示第几个（某些头有多个，比如 via ），从 0 开始。
func (h *Header) GetOther(key string, index int) string {
	if index < 0 {
		index = 0
	}
	n := 0
	for i := 0; i < len(h.Others); i++ {
		if h.Others[i].Key == key {
			if n == index {
				return h.Others[i].Value
			}
			n++
		}
	}
	return ""
}

// SetOther 设置指定 others，如果没有找到，添加一个
func (h *Header) SetOther(key, value string) {
	for i := 0; i < len(h.Others); i++ {
		if h.Others[i].Key == key {
			h.Others[i].Value = value
			return
		}
	}
	h.Others = append(h.Others, KV{Key: key, Value: value})
}

// ReplaceOther 使用指定 newKey 和 value 替换掉指定的 oldkey ，没有找到就添加
func (h *Header) ReplaceOther(oldkey, newKey, value string) {
	for i := 0; i < len(h.Others); i++ {
		if h.Others[i].Key == oldkey {
			h.Others[i].Key = newKey
			h.Others[i].Value = value
			return
		}
	}
	h.Others = append(h.Others, KV{Key: newKey, Value: value})
}

// RemoveOther 移除指定 key 的 header
func (h *Header) RemoveOther(key string) {
	for i := 0; i < len(h.Others); i++ {
		if h.Others[i].Key == key {
			copy(h.Others[i:], h.Others[i+1:])
			h.Others = h.Others[:len(h.Others)-1]
			return
		}
	}
}

// ParseFrom 从 msg 解析出字段，注意解析的 other 的 key 是大写
func (h *Header) ParseFrom(reader Reader, max int) (int, error) {
	h.Reset()
	for {
		// 读取一行数据
		line, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return max, err
		}
		// 最后一个空行
		if line == "" {
			break
		}
		max = max - len(line) - 2
		if max < 0 {
			return max, ErrLargeMessage
		}
		// 第一个':'
		line = strings.TrimSpace(line)
		i := strings.IndexByte(line, ':')
		if i < 0 {
			// 不支持多行 values 。
			return max, errHeaderFormat
		}
		key := strings.TrimSpace(line[:i])
		value := strings.TrimSpace(line[i+1:])
		uKey := strings.ToUpper(key)
		// 挑选出必要的头
		switch uKey {
		case "CALL-ID":
			h.CallID = value
		case "CSEQ":
			err = h.CSeq.Parse(value)
		case "TO":
			err = h.To.Parse(value)
		case "FROM":
			err = h.From.Parse(value)
		case "MAX-FORWARDS":
			err = h.MaxForwards.Parse(value)
			if err != nil {
				err = errHeaderMaxForwardsFormat
			}
		case "VIA":
			var via Via
			err = via.Parse(value)
			if err == nil {
				h.Via = append(h.Via, via)
			}
		case "EXPIRES":
			err = h.Expires.Parse(value)
			if err != nil {
				err = errHeaderExpiresFormat
			}
		case "CONTENT-TYPE":
			h.ContentType = value
		case "CONTACT":
			value = TrimByte(value, '<', '>')
			// 有这种格式的
			if value == "*" {
				h.Contact.Address = value
				h.Contact.OriginalString = value
			} else {
				err = h.Contact.Parse(value)
			}
		case "CONTENT-LENGTH":
			n, _err := strconv.ParseInt(value, 10, 64)
			if _err != nil || n < 0 {
				err = errHeaderContentLengthFormat
			} else {
				h.contentLength = n
			}
		default:
			h.Others = append(h.Others, KV{Key: key, Value: value})
		}
		if err != nil {
			return max, err
		}
	}
	// Via
	if len(h.Via) < 1 {
		return max, errMissingHeaderVia
	}
	// From
	if h.From.OriginalString == "" {
		return max, errMissingHeaderFrom
	}
	// To
	if h.To.OriginalString == "" {
		return max, errMissingHeaderTo
	}
	// CSeq
	if h.CSeq.OriginalString == "" {
		return max, errMissingHeaderCSeq
	}
	// Call-ID
	if h.CallID == "" {
		return max, errMissingHeaderCallID
	}
	return max, nil
}

// FormatTo 将数据写入到 writer ，不包含 ContentLength
func (h *Header) FormatTo(writer Writer) error {
	var err error
	// Via
	for i := 0; i < len(h.Via); i++ {
		_, err = writer.WriteString("Via: ")
		if err != nil {
			return err
		}
		err = h.Via[i].FormatTo(writer)
		if err != nil {
			return err
		}
		_, err = writer.WriteString("\r\n")
		if err != nil {
			return err
		}
	}
	// From
	err = formatHeaderTo2(writer, "From", &h.From)
	if err != nil {
		return err
	}
	// To
	err = formatHeaderTo2(writer, "To", &h.To)
	if err != nil {
		return err
	}
	// Call-ID
	err = formatHeaderTo(writer, "Call-ID", h.CallID)
	if err != nil {
		return err
	}
	// CSeq
	err = formatHeaderTo2(writer, "CSeq", &h.CSeq)
	if err != nil {
		return err
	}
	// Contact
	if h.Contact.Scheme != "" && h.Contact.Name != "" && h.Contact.Address != "" {
		err = formatHeaderTo2(writer, "Contact", &h.Contact)
		if err != nil {
			return err
		}
	}
	// Expires
	if h.Expires.s != "" {
		err = formatHeaderTo(writer, "Expires", h.Expires.s)
		if err != nil {
			return err
		}
	}
	// Max-Forwards
	if h.MaxForwards.s != "" {
		err = formatHeaderTo(writer, "Max-Forwards", h.MaxForwards.String())
		if err != nil {
			return err
		}
	}
	// Content-Type
	if h.ContentType != "" {
		err = formatHeaderTo(writer, "Content-Type", h.ContentType)
		if err != nil {
			return err
		}
	}
	// Others
	for i := 0; i < len(h.Others); i++ {
		err = formatHeaderTo(writer, h.Others[i].Key, h.Others[i].Value)
		if err != nil {
			return err
		}
	}
	// User-Agent
	if h.UserAgent == "" {
		h.UserAgent = "gbs"
	}
	err = formatHeaderTo(writer, "User-Agent", h.UserAgent)
	if err != nil {
		return err
	}
	// Content-Length
	err = formatHeaderTo(writer, "Content-Length", strconv.FormatInt(h.contentLength, 10))
	if err != nil {
		return err
	}
	return nil
}
