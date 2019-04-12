package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testStruct struct {
	Test string
}

var (
	pid         string
	pidList     []string
	humanFormat = "kB"
	singlePid   bool
	allPids     = false
)

func start() {
	if pid != "" && pidList == nil && allPids == false {
		singlePid = true
	} else if pid == "" && (pidList != nil || allPids) {
		singlePid = false
	} else if pid == "" && pidList == nil && !allPids {
		fmt.Println("Neither a pid is specified nor a pid list or --all flag.")
		printOptions(false)
	} else if pid != "" && (pidList != nil || allPids) {
		if allPids {
			fmt.Print("You cant have --all flag when checking for a single pid or pid list.")
		} else {
			fmt.Println("You cant have both single pid and pid list.")
		}
		printOptions(false)
	}
	if singlePid {
		fmt.Println(getSinglePidHeap())
	}
}

//getSinglePidHeap Gets pid's heap usage via jstat
func getSinglePidHeap() string {
	res := exeCmd("jstat -gc " + pid + " | awk 'FNR==2{print $0}' | awk '{heap=$3+$4+$6+$8+$10+$12; print heap}'")
	if strings.Contains(res, "not found") == true {
		fmt.Println("No java process found with pid:", pid, "Exiting..")
		fmt.Println(res)
		os.Exit(1)
	}
	heap, _ := strconv.ParseFloat(res[:len(res)-1], 64)
	if humanFormat == "MB" {
		heap = heap / 1024
	} else if humanFormat == "GB" {
		heap = heap / 1024 / 1024
	} else if humanFormat == "B" {
		heap = heap * 1024
	}
	output := fmt.Sprintf("%.2f", heap)
	return "Single pid heap: " + output + " " + humanFormat //todo return json not string
}

//checkHumanFormat validates user input of -h flag.
func checkHumanFormat(format string) {
	if format == "GB" || format == "MB" || format == "kB" || format == "B" {
		humanFormat = format
	} else {
		fmt.Println("Human format prameter is not one of the accepted values GB, MB, kB, B. Exiting..")
		os.Exit(1)
	}
}

//validatePid Checks if there's an existing JAVA process running with the provided pid.
func validatePid(pidInput string) {
	res := exeCmd("ps auwx | grep java")
	if strings.Contains(res, pidInput) == false {
		fmt.Println("No java process found with pid:", pidInput, "Exiting..")
		os.Exit(1)
	} else {
		pid = pidInput
	}
}

//validateAll Checks if any JAVA process is running, if not it terminates the execution
func validateAll() {
	res := exeCmd("ps auwx | grep java | sed '$d'")
	if strings.Contains(res, "java") == false {
		fmt.Println("No java process found .Exiting..")
		os.Exit(1)
	} else {
		allPids = true
	}

}

// parsePidList Parses pid list input and checks is pids exist.
func parsePidList(pidListInput string) {
	//todo parse pid list, check if pids exist one by one, else exit 1
	pidList = append(pidList, pidListInput)
}

//exeCmd Executes a bash command
func exeCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return string(out)
}

//printOptions print options & help instead of failre or upon request
func printOptions(help bool) {
	if !help {
		fmt.Println("no valid arguements found. exiting...")
		fmt.Println()
	}
	fmt.Println("You can use the following arguements:")
	fmt.Println()
	fmt.Println("	-p 		Get Heap usage for specific JAVA pid")
	fmt.Println("	--pid-list	Get Heap usage for a list of JAVApids. Provide them like 123,54487,7895. The output will be in JSON format")
	fmt.Println("	--all		Get Heap Usage for all running JAVA processes")
	fmt.Println("	-h 		for human readable format in GB, MB, kB, B")
	fmt.Println("	--help 		to display this help ouput")
	if !help {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func main() {
	if len(os.Args) == 1 {
		printOptions(false)
	} else {
		for i, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-") {
				if arg == "-p" {
					validatePid(os.Args[i+2])
				} else if arg == "--pid-list" {
					parsePidList(os.Args[i+2])
				} else if arg == "--help" {
					fmt.Println("Displaying help")
					printOptions(true)
				} else if arg == "-h" {
					checkHumanFormat(os.Args[i+2])
				} else if arg == "--all" {
					validateAll()
				} else {
					printOptions(false)
				}
			}
		}

		start()
	}
}
