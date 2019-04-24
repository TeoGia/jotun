package main

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	res := GetVersion()
	if len(res) >= 5 {
		t.Log("Get version ok")
	} else {
		t.Error("Wrong version")
	}
}

//TODO add more tests
