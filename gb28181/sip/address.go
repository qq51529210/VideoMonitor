package sip

import (
	"errors"
	"strings"
)

var (
	errAddressFormat = errors.New("error address format")
)

// Address 表示 "name" <uri> tag
type Address struct {
	Name           string
	URI            URI
	Tag            string
	OriginalString string `json:"-"`
}

// Reset 重置
func (a *Address) Reset() {
	a.Name = ""
	a.URI.Reset()
	a.Tag = ""
	a.OriginalString = ""
}

// CopyTo copy 数据到 aa
func (a *Address) CopyTo(aa *Address) {
	aa.Name = a.Name
	aa.URI = a.URI
	aa.Tag = a.Tag
	aa.OriginalString = a.OriginalString
}

// FormatTo 格式化到 writer 中。
func (a *Address) FormatTo(writer Writer) error {
	var err error
	// name
	if a.Name != "" {
		_, err = writer.WriteString(a.Name)
		if err != nil {
			return err
		}
		_, err = writer.WriteString(" ")
		if err != nil {
			return err
		}
	}
	// uri
	err = a.URI.FormatTo(writer)
	if err != nil {
		return err
	}
	// tag
	if a.Tag != "" {
		_, err = writer.WriteString(";tag=")
		if err != nil {
			return err
		}
		_, err = writer.WriteString(a.Tag)
		if err != nil {
			return err
		}
	}
	return nil
}

// Parse 从 line 解析数据。
func (a *Address) Parse(line string) error {
	a.OriginalString = line
	// 分段
	part := strings.Split(line, ";")
	// name uri
	fields := strings.Fields(part[0])
	switch len(fields) {
	case 1:
		// 只有 uri
		err := a.URI.Parse(TrimByte(fields[0], '<', '>'))
		if err != nil {
			return err
		}
	case 2:
		// name uri
		a.Name = fields[0]
		err := a.URI.Parse(TrimByte(fields[1], '<', '>'))
		if err != nil {
			return err
		}
	default:
		return errAddressFormat
	}
	// tag
	if len(part) > 1 {
		var kv KV
		err := kv.Parse(part[1])
		if err != nil {
			return errAddressFormat
		}
		if kv.Key == "tag" {
			a.Tag = kv.Value
		}
	}
	// 返回
	return nil
}
