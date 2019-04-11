package Logger

import (
	"fmt"
	"os"
	"time"
)

var (
	file    *os.File
	logFile = false
	path    = "logs/tooly-" + time.Now().Format(time.RFC850) + ".log"
)

// LogEntry - err :true for error, false otherwise
type LogEntry struct {
	message string
	err     bool
}

var errorLog []LogEntry

//Init create log file
func Init() {
	if logFile {
		_, err := os.Create(path)
		Check(err)
	}
}

// GetLogFileStatus returns if a local log file is to be kept.
func GetLogFileStatus() bool {
	return logFile
}

func writeFile(entry string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString(entry)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		fmt.Println(err.Error())
	}
}

//Check for error, print it and push it in the log
func Check(e error) {
	if e != nil {
		entry := LogEntry{message: e.Error(), err: true}
		errorLog = append(errorLog, entry)
		fmt.Println(e)
		if logFile {
			writeFile(e.Error())
		}
	}
}

//Log something that is not of type error
func Log(msg string) {
	entry := LogEntry{msg, false}
	errorLog = append(errorLog, entry)
	fmt.Println(msg)
	if logFile {
		writeFile(msg)
	}
}

//GetErrorLog : Getter method for error log slice in HTML format
func GetErrorLog() []string {
	var log []string

	for _, entry := range errorLog {
		if entry.err {
			log = append(log, "<p style="+"'color:#FF0000'>"+entry.message+"</p>")
		} else {
			log = append(log, "<p style="+"'color:#009933'>"+entry.message+"</p>")
		}
	}
	return log
}
