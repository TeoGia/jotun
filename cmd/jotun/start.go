package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/teogia/jotun/helper"
)

type pidListOut struct {
	PidLlist []pidOut
}

type pidOut struct {
	Pid          string
	PidName      string
	Heap         string
	Format       string
	AvailableRAM string
	FreeRAM      string
}

const (
	//VERSION the application's current version
	VERSION = "1.0.0-alpha"
)

var (
	pid         string
	pidList     []string
	humanFormat = "kB"
	singlePid   bool
	allPids     = false
)

//GetVersion returns current version
func GetVersion() string {
	return VERSION
}

//Start starts jotun after validations are over
func Start() {
	if pid != "" && pidList == nil && !allPids {
		singlePid = true
	} else if pid == "" && (pidList != nil || allPids) {
		singlePid = false
		if pidList != nil && allPids {
			fmt.Println("You cant have --all flag when checking for a single pid or pid list.")
			PrintOptions(false)
		}
		if allPids {
			gatherJavaProcess()
		}
	} else if pid == "" && pidList == nil && !allPids {
		fmt.Println("Neither a pid is specified nor a pid list or --all flag.")
		PrintOptions(false)
	} else if ((pid != "" || pidList != nil) && allPids) || pid != "" && pidList != nil {
		if allPids {
			fmt.Println("You cant have --all flag when checking for a single pid or pid list.")
		} else {
			fmt.Println("You cant have both single pid and pid list.")
		}
		PrintOptions(false)
	}
	if singlePid {
		fmt.Println(string(helper.PrintJSON(getSinglePidHeap(pid))))
	} else {
		fmt.Println(string(helper.PrintJSON(getPidListHeap())))
	}
}

//gatheJavaProcess gathers all running java processes into pid list
func gatherJavaProcess() {
	res := helper.ExeCmd("ps auwx | egrep java")
	if !strings.Contains(res, "java") {
		fmt.Println("No JAVA processes found running. Exiting..")
		os.Exit(1)
	}
	raw := strings.Split(res, "\n")
	for _, row := range raw[:len(raw)-1] {
		if !strings.Contains(row, "grep") {
			pidList = append(pidList, strings.Fields(row)[1])
		}
	}
}

//getFreeRAM returns the system's total RAM in the selected format
func getFreeRAM() string {
	res := helper.ExeCmd("cat /proc/meminfo")
	freeRAM, _ := strconv.ParseFloat(strings.Fields(strings.Split(res, "\n")[1])[1], 64)
	if humanFormat == "MB" {
		freeRAM = freeRAM / 1024
	} else if humanFormat == "GB" {
		freeRAM = freeRAM / 1024 / 1024
	} else if humanFormat == "B" {
		freeRAM = freeRAM * 1024
	}
	return fmt.Sprintf("%.2f", freeRAM)
}

//getTotalRAM returns the system's total RAM in the selected format
func getTotalRAM() string {
	res := helper.ExeCmd("cat /proc/meminfo")
	totalRAM, _ := strconv.ParseFloat(strings.Fields(strings.Split(res, "\n")[0])[1], 64)
	if humanFormat == "MB" {
		totalRAM = totalRAM / 1024
	} else if humanFormat == "GB" {
		totalRAM = totalRAM / 1024 / 1024
	} else if humanFormat == "B" {
		totalRAM = totalRAM * 1024
	}
	return fmt.Sprintf("%.2f", totalRAM)
}

//getPidListHeap Getas pid list's heap usage via jstat
func getPidListHeap() pidListOut {
	res := []pidOut{}
	for _, pid := range pidList {
		res = append(res, getSinglePidHeap(pid))
	}

	output := pidListOut{
		PidLlist: res,
	}
	return output
}

//getSinglePidHeap Gets pid's heap usage via jstat
func getSinglePidHeap(pid string) pidOut {
	res := helper.ExeCmd("jstat -gc " + pid) //| awk 'FNR==2{print $0}' | awk '{heap=$3+$4+$6+$8+$10+$12; print heap}'")
	resFields := strings.Fields(strings.Split(res, "\n")[1])
	heap := 0.0
	for i, mem := range resFields {
		if i == 2 || i == 3 || i == 5 || i == 7 || i == 9 || i == 11 {
			mem, _ := strconv.ParseFloat(mem, 64)
			heap += mem
		}
	}
	pidRaw := helper.ExeCmd("ps awux | grep " + pid)
	pidName := strings.Fields(strings.Split(pidRaw, "\n")[0])[12]
	if strings.Contains(res, "not found") == true {
		fmt.Println("No java process found with pid:", pid, "Exiting..")
		fmt.Println(res)
		os.Exit(1)
	}
	if humanFormat == "MB" {
		heap = heap / 1024
	} else if humanFormat == "GB" {
		heap = heap / 1024 / 1024
	} else if humanFormat == "B" {
		heap = heap * 1024
	}
	output := pidOut{
		Pid:          pid,
		PidName:      pidName,
		Heap:         fmt.Sprintf("%.2f", heap),
		Format:       humanFormat,
		AvailableRAM: getTotalRAM(),
		FreeRAM:      getFreeRAM(),
	}

	return output
}

//CheckHumanFormat validates user input of -h flag.
func CheckHumanFormat(format string) {
	var acceptableFormats = []string{"GB", "MB", "kB", "B"}
	if !helper.IsValueInList(format, acceptableFormats) {
		fmt.Println("Human format prameter is not one of the accepted values GB, MB, kB, B. Exiting..")
		os.Exit(1)
	}
	humanFormat = format
}

//ValidatePid Checks if there's an existing JAVA process running with the provided pid.
func ValidatePid(pidInput string) bool {
	res := helper.ExeCmd("ps auwx | grep java")
	if !strings.Contains(res, pidInput) {
		fmt.Println("No java process found with pid:", pidInput, "Exiting..")
		os.Exit(1)
	}
	return true
}

//ValidateAll Checks if any JAVA process is running, if not it terminates the execution
func ValidateAll() {
	res := helper.ExeCmd("ps auwx | grep java | sed '$d'")
	if !strings.Contains(res, "java") {
		fmt.Println("No JAVA processes found running. Exiting..")
		os.Exit(1)
	}
	allPids = true

}

//ParsePidList Parses pid list input and checks is pids exist.
func ParsePidList(pidListInput string) {
	pidList = strings.Split(pidListInput, ",")
	for _, pid := range pidList {
		ValidatePid(pid)
	}
}

//PrintOptions print options & help instead of failre or upon request
func PrintOptions(help bool) {
	if !help {
		fmt.Println("no valid arguements found. exiting...")
		fmt.Println()
	}
	fmt.Println("You can use the following arguements:")
	fmt.Println()
	fmt.Println("	-p 		Get Heap usage for specific JAVA pid")
	fmt.Println("	--pid-list	Get Heap usage for a list of JAVA pids. Provide them like 123,54487,7895. The output will be in JSON format")
	fmt.Println("	--all		Get Heap Usage for all running JAVA processes")
	fmt.Println("	-h 		for human readable format in GB, MB, kB, B")
	fmt.Println("	--help 		to display this help ouput")
	fmt.Println("	-v		Get jotun's varsion")
	if !help {
		os.Exit(1)
	}
	os.Exit(0)
}
