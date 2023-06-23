package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	externalip "github.com/glendc/go-external-ip"
)

// configureWireguardSubnet asks the user to input a Wireguard IPv4 subnet through the console and
// then parses the input into IP network format. It displays some recommendations about subnet choice
// and allows the user to either input a custom subnet or accept the default one.
//
// This function first prints a couple of messages to guide the user in choosing a suitable subnet.
// Then it reads the user's input from the console. If the user just presses enter without typing
// anything, the function uses a default subnet.
//
// The function then tries to parse the user's input (or the default subnet) into the net.IP and
// net.IPNet types that can be used with the rest of the net package's IP networking functions.
//
// Returns:
//     net.IP: The IP address part of the inputted subnet.
//     net.IPNet: The network and mask part of the inputted subnet.
//     error: An error object indicating any errors that occurred during parsing.
//
// Usage:
//     ip, subnet, err := configureWireguardSubnet()
func configureWireguardSubnet() (net.IP, *net.IPNet, error) {
	fmt.Println("\nConfigure the Wireguard IPv4 subnet:")
	fmt.Println("\t1. You can use any IPv4 subnet if it does not conflict with local addresses.")
	fmt.Println("\t2. It is recommended to use private IPv4 subnet, e.g 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16.")
	fmt.Printf("Enter the Wireguard IPv4 subnet or press Enter to use the suggested one [%s]:",
		defaultWireguardSubnet)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return net.ParseCIDR(defaultWireguardSubnet)
	}

	return net.ParseCIDR(input)
}

// configureWireguardEndpoint asks the user to input a Wireguard server endpoint through the console and
// then configures the endpoint with an auto-detected external IP address and available UDP port. It also provides
// guidance about endpoint configuration and allows the user to either input a custom endpoint or accept the
// suggested one.
//
// This function first gets the external IP using the externalip package and finds an unused UDP port.
// Then it constructs the endpoint string in the format IP:Port and reads user's input from the console.
// If the user types something, it parses the input to extract the hostname and port and uses them to update
// the endpoint and serverPort values.
//
// Returns:
//     string: The final endpoint, in the format of "IP:Port" or "Hostname:Port".
//     int: The final server port.
//
// Usage:
//     endpoint, serverPort := configureWireguardEndpoint()
func configureWireguardEndpoint() (string, int) {
	// Create the default consensus,
	// using the default configuration and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)
	// Get your IP,
	// which is never <nil> when err is <nil>.
	externalIP, err := consensus.ExternalIP()
	if err != nil {
		fmt.Println(externalIP.String()) // print IPv4/IPv6 in string format
	}

	serverPort, err := GetUnusedUdpPort()
	if err != nil {
		log.Fatalf("Failed to obtain available UDP port")
	}

	endpoint := fmt.Sprintf("%s:%d", externalIP.String(), serverPort)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nConfigure the Wireguard Server endpoint:")
	fmt.Println("\t1. You can enter DNS or dynamic DNS host name if you have one configured.")
	fmt.Println("\t2. Don't forget to map the chosen UDP port on your router or VPS provider.")
	fmt.Println("\t3. Enter the Wireguard Server endpoint below or just press Enter to use the suggested one.")
	fmt.Printf("Auto-detected external IP address and UDP port [%s]:", endpoint)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		hostString, portString, err := net.SplitHostPort(input)
		if err == nil {
			port, err := strconv.Atoi(portString)
			if err == nil {
				endpoint = fmt.Sprintf("%s:%s", hostString, portString)
				serverPort = port
			}
		}
	}
	return endpoint, serverPort
}
