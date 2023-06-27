package sip

import "testing"

func Test_Address(t *testing.T) {
	var addr Address
	// uri_test.go 测试了 uri
	err := addr.Parse(`  aaa  sip:123@456;tag=321`)
	if err != nil {
		t.Fatal(err)
	}
	if addr.Name != "aaa" || addr.Tag != "321" {
		t.FailNow()
	}
	addr.Reset()
	err = addr.Parse(`sip:123@456; tag=321 `)
	if err != nil {
		t.Fatal(err)
	}
	if addr.Name != "" || addr.Tag != "321" {
		t.FailNow()
	}
	addr.Reset()
	err = addr.Parse(`aaa <sip:123@456>`)
	if err != nil {
		t.Fatal(err)
	}
	if addr.Name != "aaa" || addr.Tag != "" {
		t.FailNow()
	}
}

func Test_Address_Error(t *testing.T) {
	// uri_test.go
	// transmission/kv_test.go
}
