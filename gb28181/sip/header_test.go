package sip

import (
	"bytes"
	"testing"
)

func Test_Header(t *testing.T) {
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
	b.WriteString("Content-Length: 0\r\n")

	var h Header
	_, err := h.ParseFrom(NewReader(&b, -1), 0)
	if err != nil {
		t.Fatal(err)
	}

	var b2 bytes.Buffer
	h.FormatTo(&b2)
	// os.Stderr.Write(b2.Bytes())
}
