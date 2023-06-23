package main

import (
	"fmt"
	"net"
)

type Interface struct {
	PrivateKey string
	ListenPort uint16
	Address    []net.IPNet
	DNS        []net.IP
	MTU        uint16
}

type Peer struct {
	PublicKey           string
	AllowedIPs          []net.IPNet
	Endpoint            string
	PersistentKeepalive uint32
}

type WireguardConfig struct {
	Interface
	Peers []Peer
}

// AddPeer is a method on the WireguardConfig type that adds a new peer to the Wireguard configuration.
// It takes as input a public key (PublicKey) and a slice of allowed IP addresses (AllowedIPs), both defined in the string format.
// The method creates a new Peer struct, populating it with the provided PublicKey and AllowedIPs, and appends it to the Peers slice in the WireguardConfig.
// It returns a pointer to the newly created Peer.
// This method is typically used to add a new client to a Wireguard network.
func (wc *WireguardConfig) AddPeer(PublicKey string, AllowedIPs []net.IPNet) *Peer {

	peer := Peer{
		PublicKey:  PublicKey,
		AllowedIPs: AllowedIPs,
	}

	wc.Peers = append(wc.Peers, peer)
	return &wc.Peers[len(wc.Peers)-1]
}

// NewWireguardServerConfig is a function that creates and returns a new Wireguard server configuration.
// The function takes in a private key (PrivateKey) in string format, a slice of IPNet (Address) to specify the IP address of the server,
// and a port number (ListenPort) on which the server should listen for incoming connections.
// A new WireguardConfig struct is created and populated with these values, with the Peers field initialized as nil, indicating no connected peers.
// The populated WireguardConfig is then returned.
// This function is typically used to initialize a new Wireguard server configuration.
func NewWireguardServerConfig(PrivateKey string, Address []net.IPNet, ListenPort uint16) WireguardConfig {
	wc := WireguardConfig{
		Interface: Interface{
			PrivateKey: PrivateKey,
			Address:    Address,
			ListenPort: ListenPort,
		},
		Peers: nil,
	}

	return wc
}

// NewWireguardClientConfig is a function that creates and returns a new Wireguard client configuration.
// It takes the following parameters:
// - PrivateKey (string): the private key of the client
// - Address ([]net.IPNet): an array of IPNet objects specifying the client IP addresses
// - PublicKey (string): the public key of the server (peer) to connect to
// - AllowedIPs ([]net.IPNet): an array of IPNet objects specifying the IP addresses that are allowed for the peer
// - Endpoint (string): the endpoint of the server to connect to
// The function creates a new WireguardConfig struct, populates it with the given values and adds a peer with the provided public key, allowed IPs, and endpoint.
// The populated WireguardConfig is then returned.
// This function is typically used to create a new Wireguard client configuration.
func NewWireguardClientConfig(PrivateKey string, Address []net.IPNet, PublicKey string, AllowedIPs []net.IPNet, Endpoint string) WireguardConfig {
	wc := WireguardConfig{
		Interface: Interface{
			PrivateKey: PrivateKey,
			Address:    Address,
		},
		Peers: nil,
	}

	peer := wc.AddPeer(PublicKey, AllowedIPs)
	peer.Endpoint = Endpoint
	return wc
}

// String method on Peer struct is used to create and return a string representation of a Wireguard peer configuration.
// This method can be useful for generating peer configuration sections in Wireguard configuration files.
//
// The String method does the following:
// - It loops over the AllowedIPs slice and creates a comma-separated string representation of it.
// - It then creates a string using the PublicKey and the string representation of AllowedIPs.
// - If the Endpoint of the peer is not an empty string, it appends the Endpoint to the resulting string.
// - If the PersistentKeepalive of the peer is not 0, it appends the PersistentKeepalive to the resulting string.
//
// The resulting string is in a format that can be directly included in a Wireguard configuration file.
//
// Example:
// Given a peer with PublicKey "abcd", AllowedIPs with a single IP "10.0.0.2/32", Endpoint "10.0.0.1:51820", and PersistentKeepalive "15",
// the returned string will be:
//
// [Peer]
// PublicKey = abcd
// AllowedIPs = 10.0.0.2/32
// Endpoint = 10.0.0.1:51820
// PersistentKeepalive = 15
//
// This method does not return an error. If there are any issues with the peer configuration, those would need to be detected and handled at the point of creation of the Peer struct.
func (peer Peer) String() string {
	var allowedString string

	for i, address := range peer.AllowedIPs {
		if i != (len(peer.AllowedIPs) - 1) {
			allowedString += address.String() + ", "
		} else {
			allowedString += address.String()
		}
	}

	result := fmt.Sprintf(
		"\n[Peer]\nPublicKey = %s\nAllowedIPs = %s\n",
		peer.PublicKey, allowedString)

	if peer.Endpoint != "" {
		result += fmt.Sprintf(
			"Endpoint = %s\n",
			peer.Endpoint)
	}

	if peer.PersistentKeepalive != 0 {
		result += fmt.Sprintf("PersistentKeepalive = %d\n", peer.PersistentKeepalive)
	}

	return result
}

// String method on the WireguardConfig struct is used to create and return a string representation of a Wireguard configuration.
// This method can be useful for generating configuration files for Wireguard.
//
// The String method does the following:
// - It loops over the Address slice and creates a comma-separated string representation of it.
// - It loops over the DNS slice and creates a comma-separated string representation of it.
// - It creates a string using the PrivateKey, the string representation of Address.
// - If the ListenPort of the configuration is not 0, it appends the ListenPort to the resulting string.
// - If the DNS is not an empty string, it appends the DNS to the resulting string.
// - If the MTU of the configuration is not 0, it appends the MTU to the resulting string.
// - It then loops over the Peers slice and appends the string representation of each peer (generated by calling the String method on the Peer struct) to the resulting string.
//
// The resulting string is in a format that can be directly used as a Wireguard configuration file.
//
// Example:
// Given a configuration with PrivateKey "abcd", Address with a single IP "10.0.0.1/32", ListenPort "51820", DNS with a single IP "1.1.1.1", MTU "1500", and a single Peer,
// the returned string will be:
//
// [Interface]
// PrivateKey = abcd
// Address = 10.0.0.1/32
// ListenPort = 51820
// DNS = 1.1.1.1
// MTU = 1500
// [Peer]
// ... // (Output of peer.String() method)
//
// This method does not return an error. If there are any issues with the configuration, those would need to be detected and handled at the point of creation of the WireguardConfig struct.
func (wc WireguardConfig) String() string {
	var addressString string
	var dnsString string

	for i, address := range wc.Address {
		if i != (len(wc.Address) - 1) {
			addressString += address.String() + ", "
		} else {
			addressString += address.String()
		}
	}

	for i, address := range wc.DNS {
		if i != (len(wc.DNS) - 1) {
			dnsString += address.String() + ", "
		} else {
			dnsString += address.String()
		}
	}

	result := fmt.Sprintf(
		"[Interface]\nPrivateKey = %s\nAddress = %s\n",
		wc.PrivateKey, addressString)

	if wc.ListenPort != 0 {
		result += fmt.Sprintf("ListenPort = %d\n", wc.ListenPort)
	}

	if dnsString != "" {
		result += fmt.Sprintf("DNS = %s\n", dnsString)
	}

	if wc.MTU != 0 {
		result += fmt.Sprintf("MTU = %d\n", wc.MTU)
	}

	for _, peer := range wc.Peers {
		result += peer.String()
	}

	return result
}
