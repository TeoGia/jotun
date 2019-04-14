package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		PrintOptions(false)

	}
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			if arg == "-p" && len(os.Args) > 2 {
				ValidatePid(os.Args[i+2])
			} else if arg == "--pid-list" && len(os.Args) > 2 {
				ParsePidList(os.Args[i+2])
			} else if arg == "--help" {
				fmt.Println("Displaying help")
				PrintOptions(true)
			} else if arg == "-h" && len(os.Args) > 2 {
				CheckHumanFormat(os.Args[i+2])
			} else if arg == "--all" {
				ValidateAll()
			} else if arg == "-v" {
				fmt.Println(GetVersion())
				os.Exit(0)
			} else {
				PrintOptions(false)
			}
		}
	}

	Start()

}
