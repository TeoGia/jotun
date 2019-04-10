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
	Logger.Log("Starting monitoring packets on selected device..\n==================START======================\n")
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	Logger.Check(err)
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		Logger.Log(packet.String())
		getPacketPayload(packet)
	}
}

func start(config conf.Config) {
	print("\033[H\033[2J") //clears terminal screen
	devices, err := pcap.FindAllDevs()
	Logger.Check(err)

	// Print device information
	Logger.Log("Devices found:")
	for _, device := range devices {
		Logger.Log("\nName: " + device.Name)
		Logger.Log("Description: " + device.Description)
		Logger.Log("Devices addresses: " + device.Description)
		for _, address := range device.Addresses {
			Logger.Log("- IP address: " + address.IP.String())
			Logger.Log("- Subnet mask: " + address.Netmask.String())
		}
	}
	Logger.Log("\nType the name of the device you wish to monitor..")
	fmt.Scan(&device)
	print("\033[H\033[2J") //clears terminal screen
	Logger.Log("Selected device: " + device)
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
		Logger.Log("no arguements found. exiting...\n")
	}
	Logger.Log("You can use the following arguements:\n")
	Logger.Log("start - to start tooly (needs to be run with sudo)")
	Logger.Log("--conf - to specify the path to the desirable config file")
	Logger.Log("-h or --help to display this help ouput")
	os.Exit(0)
}

func main() {
	var confile string
	var config conf.Config
	Logger.Init()
	if len(os.Args) == 1 {
		printOptions(false)
	} else {
		for i, arg := range os.Args {
			if arg == "--conf" {
				Logger.Log("found arguement for config file")
				confile = os.Args[i+1]
				config = conf.GetConfig(confile)
			} else if arg == "start" {
				Logger.Log("no arguements found. starting with default config...")
				config = conf.GetConfig(".config")

			} else if arg == "-h" || arg == "--help" {
				Logger.Log("Displaying help")
				printOptions(true)
			}
		}

		start(config)
	}
}
