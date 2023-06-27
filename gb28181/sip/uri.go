package sip

import (
	"errors"
	"strings"
)

var (
	errURIFormat = errors.New("error uri format")
)

// URI 表示 sip:name@address;transport=tcp
type URI struct {
	Scheme         string
	Name           string
	Address        string
	Transport      string
	OriginalString string `json:"-"`
}

// Reset 重置数据
func (u *URI) Reset() {
	u.Scheme = ""
	u.Scheme = ""
	u.Address = ""
	u.OriginalString = ""
}

// Parse 从 line 解析数据。
func (u *URI) Parse(line string) error {
	u.OriginalString = line
	line = TrimByte(line, '<', '>')
	// 分段
	part := strings.Split(line, ";")
	// scheme:
	i := strings.IndexByte(part[0], ':')
	if i < 0 {
		return errURIFormat
	}
	u.Scheme = part[0][:i]
	if u.Scheme == "" {
		return errURIFormat
	}
	part[0] = part[0][i+1:]
	// name@address
	i = strings.IndexByte(part[0], '@')
	if i > 0 {
		u.Name = part[0][:i]
		if u.Name == "" {
			return errURIFormat
		}
	}
	u.Address = part[0][i+1:]
	if u.Address == "" {
		return errURIFormat
	}
	for i := 1; i < len(part); i++ {
		var kv KV
		err := kv.Parse(part[i])
		if err != nil {
			return errURIFormat
		}
		if kv.Key == "transport" {
			u.Transport = kv.Value
		}
	}
	return nil
}

// FormatTo 格式化到 writer 中。
func (u *URI) FormatTo(writer Writer) error {
	_, err := writer.WriteString("<")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(u.Scheme)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(":")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(u.Name)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("@")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(u.Address)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(">")
	return err
}
