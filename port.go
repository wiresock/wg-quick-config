package main

import (
	"math/big"
	"net"
	"strconv"
)

// GetUnusedUdpPort attempts to listen for UDP connections on an automatically
// chosen available port and returns that port number. If an error occurs while
// listening for a connection or parsing the port number, the function will return
// an error alongside the value 0.
//
// Returns:
//     int: The number of the unused UDP port.
//     error: An error object indicating any errors that occurred during the process.
//
// Usage:
//     port, err := GetUnusedUdpPort()
func GetUnusedUdpPort() (int, error) {
	conn, err := net.ListenUDP("udp", nil)

	if err != nil {
		return 0, err
	}
	defer conn.Close()
	hostString := conn.LocalAddr().String()
	_, portString, err := net.SplitHostPort(hostString)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portString)
}

// CheckUdpPort checks if a given UDP port is available by attempting to listen
// for UDP connections on that port. If the port is available, the function returns
// the port number; otherwise, it returns an error.
//
// Parameters:
//     Port (int): The number of the UDP port to check.
//
// Returns:
//     int: The number of the checked UDP port if it is available.
//     error: An error object indicating any errors that occurred during the process, e.g., if the port is not available.
//
// Usage:
//     port, err := CheckUdpPort(12345)
func CheckUdpPort(Port int) (int, error) {
	address := net.UDPAddr{
		Port: Port,
	}

	conn, err := net.ListenUDP("udp", &address)

	if err != nil {
		return 0, err
	}
	defer conn.Close()
	hostString := conn.LocalAddr().String()
	_, portString, err := net.SplitHostPort(hostString)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portString)
}

// NextIP calculates and returns the next sequential IP address from the given IP.
// The function treats the IP address as a big integer, increments it, and returns
// the resulting IP address. If the provided IP is the highest possible IP (255.255.255.255),
// this function will return an invalid IP address (0.0.0.0).
//
// Parameters:
//     ip (net.IP): The input IP address from which the next IP address is calculated.
//
// Returns:
//     net.IP: The next sequential IP address.
//
// Usage:
//     nextIP := NextIP(net.ParseIP("192.168.1.1"))
func NextIP(ip net.IP) net.IP {
	// Convert to big.Int and increment
	ipb := big.NewInt(0).SetBytes([]byte(ip))
	ipb.Add(ipb, big.NewInt(1))

	// Add leading zeros
	b := ipb.Bytes()
	b = append(make([]byte, len(ip)-len(b)), b...)
	return net.IP(b)
}
