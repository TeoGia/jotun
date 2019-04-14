package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/teogia/jotun/helper"
)

func main() {
	if len(os.Args) == 1 {
		helper.PrintOptions(false)

	}
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			if arg == "-p" && len(os.Args) > 2 {
				helper.ValidatePid(os.Args[i+2])
			} else if arg == "--pid-list" && len(os.Args) > 2 {
				helper.ParsePidList(os.Args[i+2])
			} else if arg == "--help" {
				fmt.Println("Displaying help")
				helper.PrintOptions(true)
			} else if arg == "-h" && len(os.Args) > 2 {
				helper.CheckHumanFormat(os.Args[i+2])
			} else if arg == "--all" {
				helper.ValidateAll()
			} else if arg == "-v" {
				fmt.Println(helper.GetVersion())
				os.Exit(0)
			} else {
				helper.PrintOptions(false)
			}
		}
	}

	helper.Start()

}
