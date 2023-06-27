package device

import (
	"encoding/xml"
	"fmt"
	"wsdl/soap"
)

const (
	Namespace = "http://www.onvif.org/ver10/device/wsdl"
)

// ActionURL 返回 http://{host}/ver10/device/wsdl/{action}
func ActionURL(host, action string) string {
	return fmt.Sprintf("http://%s/ver10/device/wsdl/%s", host, action)
}

type envelope struct {
	soap.Envelope
}

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:tt",
		},
		Value: Namespace,
	}
}

func (m *envelope) Init() {
	m.Envelope.Attr = append(m.Envelope.Attr, soap.NewNamespaceAttr())
	m.Envelope.Attr = append(m.Envelope.Attr, NewNamespaceAttr())
}
