package sip

import (
	"testing"
)

func Test_KV(t *testing.T) {
	var kv KV
	err := kv.Parse("b=2")
	if err != nil {
		t.Fatal(err)
	}
	if kv.Key != "b" || kv.Value != "2" {
		t.FailNow()
	}
	//
	kv.Key = "a"
	kv.Value = "1"
	//
	kv.Key = ""
	kv.Value = ""
	err = kv.Parse("c")
	if err != nil {
		t.Fatal(err)
	}
	if kv.Key != "c" || kv.Value != "" {
		t.FailNow()
	}
}

func Test_KV_Error(t *testing.T) {
	var kv KV
	for _, s := range []string{
		"=c",
		"=",
	} {
		if err := kv.Parse(s); err == nil {
			t.FailNow()
		}
	}
}
