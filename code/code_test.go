package code

import "testing"

func TestRegister(t *testing.T) {
	var Test Code = -999999
	Register(Test, `Test`)
	if Get(Test).Text != "Test" {
		t.Fatalf("Unmatch Text")
	}
	if Get(Test).HTTPCode != 200 {
		t.Fatalf("Unmatch HTTPCode")
	}
}
