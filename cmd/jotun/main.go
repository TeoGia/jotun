package main

import (
	"fmt"
	"os"

	"github.com/teogia/tooly/Logger"
)

type testStruct struct {
	Test string
}

var (
	pid         string
	pidList     []string
	humanFormat string
)

func start() {

}

func checkHumanFormat(format string) {
	//TODO validate user input
}

func validatePid(pidInput string) {
	//todo check if pid exists . if not exit 1
}

func parsePidList(pidListInput string) {
	//todo parse pid list, check if pids exist one by one, else exit 1
}

//printOptions print options & help instead of failre or upon request
func printOptions(help bool) {
	fmt.Println(help)
	if !help {
		fmt.Println("no valid arguements found. exiting...")
		fmt.Println()
	}
	fmt.Println("You can use the following arguements:")
	fmt.Println()
	fmt.Println("-p Get Heap usage for specific pid")
	fmt.Println("--pid-list Get Heap usage for a list of pids. Provide them like 123,54487,7895. The output will be in JSON format")
	fmt.Println("-h for human readable format in GB, MB, kB")
	fmt.Println("--help to display this help ouput")
	if !help {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func main() {
	Logger.Init()
	if len(os.Args) == 1 {
		printOptions(false)
	} else {
		for i, arg := range os.Args {
			if arg == "-p" {
				validatePid(os.Args[i+1])
			} else if arg == "--pid-list" {
				parsePidList(os.Args[i+1])
			} else if arg == "--help" {
				fmt.Println("Displaying help")
				printOptions(true)
			} else if arg == "-h" {
				checkHumanFormat(os.Args[i+1])
			}
		}

		start()
	}
}
