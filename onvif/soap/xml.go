package soap

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Header[Data any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope Header"`
	Data    Data
}

type Body[Data any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope Body"`
	Data    Data
	Fault   *Fault `xml:",omitempty"`
}

type Envelope[header, body any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope Envelope"`
	Header  header
	Body    body
}

// Fault 表示返回的错误
type Fault struct {
	XMLName xml.Name     `xml:"http://www.w3.org/2003/05/soap-envelope Fault"`
	Code    *FaultCode   `xml:"Code"`
	Reason  *FaultReason `xml:"Reason"`
	Detail  *FaultDetail `xml:"Detail"`
}

func (c *Fault) Error() string {
	if c.Detail != nil {
		return fmt.Sprintf("code: %s, %s",
			c.Detail.ErrorCode, c.Detail.Description)
	}
	if c.Reason != nil {
		return c.Reason.Text
	}
	if c.Code != nil {
		var str strings.Builder
		fmt.Fprintf(&str, "code: %s", c.Code.Value)
		if c.Code.Subcode != nil {
			fmt.Fprintf(&str, " sub code: %s", c.Code.Subcode.Value)
		}
		return str.String()
	}
	return "unknown fault"
}

// FaultCode 表示 Fault 的 Code 字段
type FaultCode struct {
	Value   string     `xml:"Value"`
	Subcode *FaultCode `xml:"Subcode,omitempty"`
}

// FaultReason 表示 Fault 的 Reason 字段
type FaultReason struct {
	Text string `xml:"Text"`
}

// FaultDetail 表示 Fault 的 Detail 字段
type FaultDetail struct {
	ErrorCode   string `xml:"ErrorCode"`
	Description string `xml:"Description"`
}
