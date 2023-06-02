package soap

import (
	"encoding/xml"
	"os"
	"testing"
)

type testXML1 struct {
	XMLName xml.Name `xml:"testXML1"`
	A       string
}

type testXML2 struct {
	XMLName xml.Name `xml:"testXML2"`
	B       string
}

type testXML3 struct {
	XMLName xml.Name `xml:"C"`
}

func Test_XML(t *testing.T) {
	m1 := new(testXML1)
	m1.A = "1"
	m2 := new(testXML2)
	m2.B = "2"
	m3 := new(testXML3)
	m3.XMLName.Local = "3"

	header := make([]any, 0)
	header = append(header, m1)
	header = append(header, m2)
	body := make([]any, 0)
	body = append(body, m1)
	body = append(body, m3)

	var msg Envelope[*Header[[]any], *Body[[]any]]
	msg.Header = &Header[[]any]{
		Data: header,
	}
	msg.Body = &Body[[]any]{
		Data: body,
	}

	encoder := xml.NewEncoder(os.Stderr)
	encoder.Indent("", "  ")
	err := encoder.Encode(&msg)
	if err != nil {
		t.Fatal(err)
	}
}
