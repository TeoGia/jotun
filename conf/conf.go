package conf

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/teogia/servant/errorHandler"
)

//Config the struct to hold the config data
type Config struct {
	Port   int
	SrcDir string
	Secure bool
}

var conf Config

func readConfFile(fileName string) {
	bs, err := ioutil.ReadFile(fileName)
	errorhandler.Check(err)

	data := strings.Split(string(bs), "\n")

	confMap := make(map[string]string)

	for _, entry := range data {
		if entry != "" && strings.Index(entry, "#") < 0 {
			kvPair := strings.Split(entry, " ")
			confMap[kvPair[0]] = kvPair[1]
		}

	}

	assignToStruct(confMap)
}

func assignToStruct(m map[string]string) {

	var err error

	if m["port"] != "" {
		conf.Port, err = strconv.Atoi(m["port"])
		errorhandler.Check(err)
	}

	if m["secure"] != "" {
		if m["secure"] == "true" {
			conf.Secure = true
		} else {
			conf.Secure = false
		}

	}
	if m["path"] != "" {
		conf.SrcDir = m["path"]
	} else {
		conf.SrcDir = "./"
	}

}

//GetConfig getter for the config struct of the conf file
func GetConfig(file string) Config {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("config file not found, create new y/n?")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		fmt.Println("answer:", answer)
		if answer == "y\n" {
			fmt.Println("creating", file, "Feel free to change it at anytime!")
			createConf(file)
		} else {
			fmt.Println("exiting...")
			os.Exit(1)
		}

	}
	fmt.Println("reading file ", file)
	readConfFile(file)
	return conf
}

func createConf(file string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	errorhandler.Check(err)
	dir = strings.Trim(dir, "bin")
	errorhandler.Log("creating: " + dir + file)
	data := []byte("port 442\npath #please fill in\nsecure false\n")
	erro := ioutil.WriteFile(dir+file, data, 0644)
	errorhandler.Check(erro)
}
