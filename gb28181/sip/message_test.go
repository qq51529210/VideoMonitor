package sip

import (
	"bytes"
	"testing"
)

func Test_Message(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("INVITE sip:bob@biloxi.com SIP/2.0\r\n")
	b.WriteString("Via: SIP/2.0/UDP pc1.atlanta.com;branch=z9hG4bK776asdhds\r\n")
	b.WriteString("Via: SIP/2.0/UDP pc2.atlanta.com\r\n")
	b.WriteString("Max-Forwards: 70\r\n")
	b.WriteString("To: Bob <sip:bob@biloxi.com>\r\n")
	b.WriteString("From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n")
	b.WriteString("Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n")
	b.WriteString("CSeq: 314159 INVITE\r\n")
	b.WriteString("Contact: <sip:alice@pc33.atlanta.com>\r\n")
	b.WriteString("Content-Type: application/sdp\r\n")
	b.WriteString("Content-Length: 5\r\n")
	b.WriteString("\r\n")
	b.WriteString("12345")
	b.WriteString("\r\n")
	// 解析
	var msg Message
	err := msg.ParseFrom(NewReader(&b, -1), 0)
	if err != nil {
		t.Fatal(err)
	}
	// header_test.go 测试过 header
	if msg.Body.String() != "12345" {
		t.FailNow()
	}
}
