package sip

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errCSeqFormat = errors.New("error format CSeq")
)

// CSeq 表示 sn method
type CSeq struct {
	SN             uint32
	Method         string
	OriginalString string `json:"-"`
}

// Parse 从 line 解析数据。
func (c *CSeq) Parse(line string) error {
	c.OriginalString = line
	// 分段
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return errCSeqFormat
	}
	// SN
	n, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return errCSeqFormat
	}
	c.SN = uint32(n)
	// Method
	c.Method = strings.ToUpper(fields[1])
	// 返回
	return nil
}

// FormatTo 格式化到 writer 中。
func (c *CSeq) FormatTo(writer Writer) error {
	_, err := writer.WriteString(strconv.FormatInt(int64(c.SN), 10))
	if err != nil {
		return err
	}
	_, err = writer.WriteString(" ")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(c.Method)
	return err
}
