package helper

import (
	"fmt"
	"os"
	"os/exec"
)

//IsValueInList checks if a string is contained in a string list
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//ExeCmd Executes a bash command
func ExeCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return string(out)
}
