package sip

import "testing"

func Test_CSeq(t *testing.T) {
	var cs CSeq
	s := "123 ACK"
	//
	err := cs.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	if cs.SN != 123 || cs.Method != "ACK" {
		t.FailNow()
	}
}

func Test_CSeq_Error(t *testing.T) {
	var cs CSeq
	s := "123ACK"
	//
	err := cs.Parse(s)
	if err != nil {
		t.FailNow()
	}
	//
	s = "a ACK"
	err = cs.Parse(s)
	if err != nil {
		t.FailNow()
	}
	//
	s = "ACK"
	err = cs.Parse(s)
	if err != nil {
		t.FailNow()
	}
}
