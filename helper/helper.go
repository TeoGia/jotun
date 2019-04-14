package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pidListOut struct {
	PidLlist []pidOut
}

type pidOut struct {
	Pid          string
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
	if pid != "" && pidList == nil && allPids == false {
		singlePid = true
	} else if pid == "" && (pidList != nil || allPids) {
		singlePid = false
	} else if pid == "" && pidList == nil && !allPids {
		fmt.Println("Neither a pid is specified nor a pid list or --all flag.")
		PrintOptions(false)
	} else if pid != "" && (pidList != nil || allPids) {
		if allPids {
			fmt.Println("You cant have --all flag when checking for a single pid or pid list.")
		} else {
			fmt.Println("You cant have both single pid and pid list.")
		}
		PrintOptions(false)
	}
	if singlePid {
		fmt.Println(string(getSinglePidHeap()))
	}
}

//getFreeRAM returns the system's total RAM in the selected format
func getFreeRAM() string {
	res := exeCmd("cat /proc/meminfo | awk 'FNR==2{print $2}'")
	if strings.Contains(res, "not found") == true {
		fmt.Println("No java process found with pid:", pid, "Exiting..")
		fmt.Println(res)
		os.Exit(1)
	}
	totaLRAM, _ := strconv.ParseFloat(res[:len(res)-1], 64)
	if humanFormat == "MB" {
		totaLRAM = totaLRAM / 1024
	} else if humanFormat == "GB" {
		totaLRAM = totaLRAM / 1024 / 1024
	} else if humanFormat == "B" {
		totaLRAM = totaLRAM * 1024
	}
	return fmt.Sprintf("%.2f", totaLRAM)
}

//getTotalRAM returns the system's total RAM in the selected format
func getTotalRAM() string {
	res := exeCmd("cat /proc/meminfo | awk 'FNR==1{print $2}'")
	if strings.Contains(res, "not found") == true {
		fmt.Println("No java process found with pid:", pid, "Exiting..")
		fmt.Println(res)
		os.Exit(1)
	}
	totaLRAM, _ := strconv.ParseFloat(res[:len(res)-1], 64)
	if humanFormat == "MB" {
		totaLRAM = totaLRAM / 1024
	} else if humanFormat == "GB" {
		totaLRAM = totaLRAM / 1024 / 1024
	} else if humanFormat == "B" {
		totaLRAM = totaLRAM * 1024
	}
	return fmt.Sprintf("%.2f", totaLRAM)
}

//getSinglePidHeap Gets pid's heap usage via jstat
func getSinglePidHeap() []byte {
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
	output := pidOut{
		Pid:          pid,
		Heap:         fmt.Sprintf("%.2f", heap),
		Format:       humanFormat,
		AvailableRAM: getTotalRAM(),
		FreeRAM:      getFreeRAM(),
	}
	jsonRes, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return jsonRes
}

//isValueInList checks if a string is contained in a string list
func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//CheckHumanFormat validates user input of -h flag.
func CheckHumanFormat(format string) {
	var acceptableFormats = []string{"GB", "MB", "kB", "B"}
	if !isValueInList(format, acceptableFormats) {
		fmt.Println("Human format prameter is not one of the accepted values GB, MB, kB, B. Exiting..")
		os.Exit(1)
	}
	humanFormat = format
}

//ValidatePid Checks if there's an existing JAVA process running with the provided pid.
func ValidatePid(pidInput string) {
	res := exeCmd("ps auwx | grep java")
	if strings.Contains(res, pidInput) == false {
		fmt.Println("No java process found with pid:", pidInput, "Exiting..")
		os.Exit(1)
	}
	pid = pidInput
}

//ValidateAll Checks if any JAVA process is running, if not it terminates the execution
func ValidateAll() {
	res := exeCmd("ps auwx | grep java | sed '$d'")
	if strings.Contains(res, "java") == false {
		fmt.Println("No java process found .Exiting..")
		os.Exit(1)
	}
	allPids = true

}

//ParsePidList Parses pid list input and checks is pids exist.
func ParsePidList(pidListInput string) {
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

//PrintOptions print options & help instead of failre or upon request
func PrintOptions(help bool) {
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
	fmt.Println("	-v		Get jotun's varsion")
	if !help {
		os.Exit(1)
	}
	os.Exit(0)
}
