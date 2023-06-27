package sip

import "testing"

func Test_Via(t *testing.T) {
	var via Via
	//
	err := via.Parse(`v1 a2 ; rport;branch=3`)
	if err != nil {
		t.Fatal(err)
	}
	if via.Version != "v1" || via.Address != "a2" || via.RProt == nil || *via.RProt != "" || via.Branch != "3" {
		t.FailNow()
	}
	via.Reset()
	err = via.Parse(`v1 a2 ;branch=3`)
	if err != nil {
		t.Fatal(err)
	}
	if via.Version != "v1" || via.Address != "a2" || via.RProt != nil || via.Branch != "3" {
		t.FailNow()
	}
	via.Reset()
	err = via.Parse(`v1 a2 ;rport=3`)
	if err != nil {
		t.Fatal(err)
	}
	if via.Version != "v1" || via.Address != "a2" || via.RProt == nil || *via.RProt != "3" || via.Branch != "" {
		t.FailNow()
	}
	via.Reset()
	err = via.Parse(`v1 a2`)
	if err != nil {
		t.Fatal(err)
	}
	if via.Version != "v1" || via.Address != "a2" || via.RProt != nil || via.Branch != "" {
		t.FailNow()
	}
}

func Test_Via_Error(t *testing.T) {
	var via Via
	for _, s := range []string{
		`v1`,
		`v1;rport;branch=3`,
		`v1 a2;rport=abc;branch=3`,
	} {
		if err := via.Parse(s); err == nil {
			t.FailNow()
		}
	}
}
