package sip

import (
	"errors"
	"strings"
)

var (
	errHeaderViaFormat = errors.New("error header Via format")
)

// NewVia 返回 sip.Via
func NewVia(proto, addr string, rport, received *string) Via {
	return Via{
		Version:  SIPVersion,
		Proto:    proto,
		Address:  addr,
		Branch:   NewBranch(),
		RProt:    rport,
		Received: received,
	}
}

// Via 表示 version address;rport=x;branch=x
type Via struct {
	Version        string
	Proto          string
	Address        string
	Branch         string
	RProt          *string
	Received       *string
	OriginalString string `json:"-"`
}

// Reset 重置数据
func (v *Via) Reset() {
	v.Version = ""
	v.Address = ""
	v.RProt = nil
	v.Received = nil
	v.Branch = ""
	v.OriginalString = ""
}

// Parse 从 line 解析数据。
func (v *Via) Parse(line string) error {
	v.OriginalString = line
	// 分段
	parts := strings.Split(line, ";")
	// version/proto address
	fields := strings.Fields(parts[0])
	if len(fields) != 2 {
		return errHeaderViaFormat
	}
	// version/proto
	i := strings.LastIndexByte(fields[0], '/')
	if i < 0 {
		return errHeaderViaFormat
	}
	v.Version = fields[0][:i]
	v.Proto = fields[0][i+1:]
	// address
	v.Address = fields[1]
	// rport=x;branch=x
	var kv KV
	for _, part := range parts[1:] {
		err := kv.Parse(part)
		if err != nil {
			return errHeaderViaFormat
		}
		switch kv.Key {
		case "rport":
			v.RProt = &kv.Value
		case "branch":
			v.Branch = kv.Value
		case "received":
			v.Received = &kv.Value
		}
	}
	return nil
}

// FormatTo 格式化到 writer 中。
func (v *Via) FormatTo(writer Writer) error {
	var err error
	// version address
	_, err = writer.WriteString(v.Version)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("/")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(v.Proto)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(" ")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(v.Address)
	if err != nil {
		return err
	}
	// rport
	if v.RProt != nil {
		_, err = writer.WriteString(";rport")
		if err != nil {
			return err
		}
		if *v.RProt != "" {
			err = writer.WriteByte('=')
			if err != nil {
				return err
			}
			_, err = writer.WriteString(*v.RProt)
			if err != nil {
				return err
			}
		}
	}
	// branch
	if v.Branch != "" {
		_, err = writer.WriteString(";branch=")
		if err != nil {
			return err
		}
		_, err = writer.WriteString(v.Branch)
		if err != nil {
			return err
		}
	}
	// received
	if v.Received != nil {
		_, err = writer.WriteString(";received")
		if err != nil {
			return err
		}
		if *v.Received != "" {
			err = writer.WriteByte('=')
			if err != nil {
				return err
			}
			_, err = writer.WriteString(*v.Received)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
