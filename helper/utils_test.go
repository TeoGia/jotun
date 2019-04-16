package helper

import (
	"strings"
	"testing"
)

func TestIsValueInList(t *testing.T) {
	res1 := IsValueInList("test", []string{"kits", "mits", "test"})
	if res1 {
		t.Log("Positive check on helper.IsValueInList result: True")
	} else {
		t.Error("Positive check on helper.IsValueInList result: False")
	}
	res2 := IsValueInList("test", []string{"kits", "mits", "tsiko"})
	if res2 {
		t.Error("Negative check on helper.IsValueInList result: False")
	} else {
		t.Log("Negative check on helper.IsValueInList result: True")
	}
}

func TestExeCmd(t *testing.T) {
	res := ExeCmd("echo testaki")
	if strings.Contains(res, "testaki") {
		t.Log("ExeCmd test succeded")
	} else {
		t.Error("ExeCmd test failed")
	}
}
