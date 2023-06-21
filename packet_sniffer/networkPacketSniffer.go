package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	// "golang.org/x/crypto/openpgp/packet"
)

/**
 * May need to install
 * ```sh
 * sudo apt-get install libpcap-dev
 * ```
 *
 * To get the "Network interface"
 * ```sh
 * $ ifconfig
 * $ # or
 * $ ifconfig | grep mtu | awk  '{print $1;}' |  sed 's/\:$//g'
 * ```
 *
 *
 * Suggestion to run
 * ```sh
 * $ go build -o packCap
 * $ sudo ./packCap
 * ```
**/

var (
	DevName = "wlp3s0" // Network interface
	Found   = false
)

func main() {
}

func snifferHttps() {
	sniffer(443)
}

func sniffer(port int) {
	// Find network interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln("Unable to fetch network interfaces")
	}
	for _, ifDev := range devices {
		if ifDev.Name == DevName {
			Found = true
		}
	}
	if !Found {
		log.Panicln("Desired device not found")
	}

	// Open live capture
	handle, err := pcap.OpenLive(DevName, 1600, false, pcap.BlockForever)
	if err != nil {
		fmt.Print(err)
		log.Panicln("Unable to open handle on the device")
	}
	defer handle.Close()

	// BPF -> Berkely Packet Filter
	if err := handle.SetBPFFilter(fmt.Sprintf("tcp and port %d", port)); err != nil {
		log.Panicln(err)
	}
	// Display filtered packets
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		isFtp := port == 443
		if isFtp {
			appLayer := packet.ApplicationLayer()
			if appLayer == nil {
				continue
			}
			data := appLayer.Payload()
			if bytes.Contains(data, []byte("USER")) || bytes.Contains(data, []byte("PASS")) {
				fmt.Println(string(data))
			}
		}

		fmt.Println(packet)
	}
}
