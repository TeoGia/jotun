package helper

import (
	"encoding/json"
	"fmt"
	"log"
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

//PrintJSON coverts struct object into json and returns the output
func PrintJSON(t interface{}) []byte {
	jsonRes, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return jsonRes
}
