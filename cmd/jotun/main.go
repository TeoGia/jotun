package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/teogia/tooly/Logger"
	"github.com/teogia/tooly/conf"
)

type testStruct struct {
	Test string
}

var (
	device            = "eth0"
	snapshotLen int32 = 1024
	promiscuous       = false
	err         error
	timeout     = 30 * time.Second
	handle      *pcap.Handle
)

func monitor() {
	fmt.Println("Starting monitoring packets on selected device..\n==================START======================\n")
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	Logger.Check(err)
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		fmt.Println(packet.String())
		getPacketPayload(packet)
	}
}

func start(config conf.Config) {
	print("\033[H\033[2J") //clears terminal screen
	devices, err := pcap.FindAllDevs()
	Logger.Check(err)

	// Print device information
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: " + device.Name)
		fmt.Println("Description: " + device.Description)
		fmt.Println("Devices addresses: " + device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: " + address.IP.String())
			fmt.Println("- Subnet mask: " + address.Netmask.String())
		}
	}
	fmt.Println("\nType the name of the device you wish to monitor..")
	fmt.Scan(&device)
	print("\033[H\033[2J") //clears terminal screen
	fmt.Println("Selected device: " + device)
	monitor()
}

func getPacketPayload(packet gopacket.Packet) {
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())

		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}

}

//printOptions print options & help instead of failre or upon request
func printOptions(help bool) {
	if !help {
		fmt.Println("no arguements found. exiting...\n")
	} //TODO update the args usage
	fmt.Println("You can use the following arguements:\n")
	fmt.Println("start - to start tooly (needs to be run with sudo)")
	fmt.Println("--conf - to specify the path to the desirable config file")
	fmt.Println("-h or --help to display this help ouput")
	os.Exit(0)
}

func main() { // TODO install sweagle task runner for testimg purposes or create my own app.
	var config conf.Config
	Logger.Init()
	if len(os.Args) == 1 {
		printOptions(false)
	} else {
		for _, arg := range os.Args {
			// if arg == "--conf" { //TODO update the accepted args
			// fmt.Println("found arguement for config file")
			// confile = os.Args[i+1]
			// config = conf.GetConfig(confile) //TODO find a way to load it at the beginning without --conf. Al this has to be removed
			// } else if arg == "start" { //TODO start should initiate the interactive live mode
			// fmt.Println("no arguements found. starting with default config...")
			// config = conf.GetConfig(".config")

			// } else if arg == "-h" || arg == "--help" {
			if arg == "-h" || arg == "--help" {
				fmt.Println("Displaying help")
				printOptions(true)
			}
		}

		start(config)
	}
}
