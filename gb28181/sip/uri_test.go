package sip

import "testing"

func Test_URI(t *testing.T) {
	var uri URI
	err := uri.Parse(`sip:123@456`)
	if err != nil {
		t.Fatal(err)
	}
	if uri.Scheme != "sip" ||
		uri.Name != "123" ||
		uri.Address != "456" {
		t.FailNow()
	}
}

func Test_URI_Error(t *testing.T) {
	var uri URI
	for _, s := range []string{
		":123@456",
		"sip:@456",
		"123@456",
		"sip:123",
	} {
		if err := uri.Parse(s); err == nil {
			t.FailNow()
		}
	}
}
