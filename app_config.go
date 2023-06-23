package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type appConfig struct {
	Server  WireguardConfig
	Clients []WireguardConfig
}

const defaultWireguardSubnet = "10.9.0.0/24"
const defaultAllowedIps = "0.0.0.0/0"
const defaultDns = "8.8.8.8, 1.1.1.1"
const defaultMtu = 1420
const defaultPersistentKeepalive = 25
const defaultClientConfigFile = "wsclient_%d.conf"
const defaultServerConfigFile = "wiresock.conf"

// clientIpNetToPeer converts a slice of IP networks into a slice of peer IP addresses.
// This function takes each IP network in the address slice, applies a /32 subnet mask to it
// to create a peer IP address (indicating a single host), and then appends it to the new slice.
//
// Parameters:
//     address ([]net.IPNet): Slice of IP networks that are to be converted to peer IP addresses.
//
// Returns:
//     []net.IPNet: Slice of peer IP addresses.
//
// Usage:
//     peerIpAddress := clientIpNetToPeer(ipAddressSlice)
func clientIpNetToPeer(address []net.IPNet) []net.IPNet {
	peerIpAddress := make([]net.IPNet, 0, len(address))
	for _, ip := range address {
		ipNet := net.IPNet{
			IP:   ip.IP,
			Mask: net.CIDRMask(32, 32),
		}
		peerIpAddress = append(peerIpAddress, ipNet)
	}

	return peerIpAddress
}

// newConfig generates a new appConfig structure that represents the Wireguard configuration
// for a VPN setup, including the server and client configurations with keys, addresses, and
// other network parameters.
//
// The function first sets up the Wireguard endpoint by retrieving the server's external IP
// address and an available UDP port.
//
// It then asks the user to input a Wireguard IPv4 subnet, using a default subnet if the user
// does not input anything.
//
// It determines the server and client IP addresses within the Wireguard subnet and parses
// the default allowed IPs which the client can connect to when the VPN is active.
//
// It generates a pair of private keys for the server and the client using the
// newWireguardPrivateKey function.
//
// The function then uses the generated keys, IP addresses, endpoint, and other parameters to
// create the server and client configurations.
//
// For the client configuration, it sets the DNS servers, the Maximum Transmission Unit (MTU),
// and the persistent keepalive interval.
//
// It then updates the appConfig structure with the new server and client configurations.
//
// Parameters:
// - config: A pointer to the appConfig structure to be updated.
//
// Returns:
// - error: An error if something goes wrong during the configuration process. If everything
//   works correctly, it returns nil.
func newConfig(config *appConfig) error {

	endpoint, serverPort := configureWireguardEndpoint()

	subnetAddressIpv4, subnetAddressIpv4Net, err := configureWireguardSubnet()

	if err != nil {
		return err
	}

	serverAddressIpv4Net := net.IPNet{
		IP:   NextIP(subnetAddressIpv4),
		Mask: subnetAddressIpv4Net.Mask,
	}

	clientAddressIpv4Net := net.IPNet{
		IP:   NextIP(serverAddressIpv4Net.IP),
		Mask: subnetAddressIpv4Net.Mask,
	}

	_, allowedIpv4Net, _ := net.ParseCIDR(defaultAllowedIps)

	allowedIPs := make([]net.IPNet, 1, 1)
	allowedIPs[0] = *allowedIpv4Net

	clientAddress := make([]net.IPNet, 1, 1)
	clientAddress[0] = clientAddressIpv4Net

	peerIpAddress := clientIpNetToPeer(clientAddress)

	serverAddress := make([]net.IPNet, 1, 1)
	serverAddress[0] = serverAddressIpv4Net

	server, _ := newWireguardPrivateKey()
	client, _ := newWireguardPrivateKey()

	serverConfig := NewWireguardServerConfig(server.base64PrivateKey(), serverAddress, uint16(serverPort))
	serverConfig.AddPeer(client.base64PublicKey(), peerIpAddress)

	clientConfig := NewWireguardClientConfig(client.base64PrivateKey(), clientAddress,
		server.base64PublicKey(), allowedIPs, endpoint)

	dns := strings.Split(defaultDns, ",")
	clientConfig.DNS = make([]net.IP, 0, len(dns))
	for i := range dns {
		dns[i] = strings.TrimSpace(dns[i])
		ip := net.ParseIP(dns[i])
		if ip != nil && ip.To4() != nil {
			clientConfig.DNS = append(clientConfig.DNS, net.ParseIP(dns[i]))
		}
	}

	clientConfig.MTU = defaultMtu
	clientConfig.Peers[0].PersistentKeepalive = defaultPersistentKeepalive

	*config = appConfig{
		Server:  serverConfig,
		Clients: nil,
	}

	config.Clients = append(config.Clients, clientConfig)

	return nil
}

// addClient is a method on the appConfig struct that adds a new client to the Wireguard VPN setup.
// It first retrieves the configuration of the last client in the list to use as a base for the new client configuration.
// A new private key is generated for the new client using the newWireguardPrivateKey function.
// The IP address for the new client is calculated based on the IP of the last client, ensuring that it remains within the allowed subnet.
// If the new IP address falls outside the subnet's capacity, the program is terminated with an error message.
// Once the IP address is successfully allocated, a new client configuration is created. This configuration includes the new IP address and subnet mask,
// and the private key generated earlier. The new client is then added as a peer to the server configuration.
// Finally, the newly created client configuration is added to the list of clients in the appConfig.
func (config *appConfig) addClient() {
	// Get the configuration of the last client
	clientConfig := config.Clients[len(config.Clients)-1]

	// Generate a new private key for the new client
	client, _ := newWireguardPrivateKey()

	// Compute the IP address of the new client
	clientIpNet := net.IPNet{
		IP:   nil,
		Mask: nil,
	}
	clientIpNet.IP = NextIP(clientConfig.Address[0].IP)
	clientIpNet.Mask = clientConfig.Address[0].Mask

	// Check if the new IP address is within the allowed subnet
	if !clientConfig.Address[0].Contains(clientIpNet.IP) ||
		!clientConfig.Address[0].Contains(NextIP(clientIpNet.IP)) {
		log.Fatalf("Cant't allocate IP address. Subnet capacity has been reached!")
	}

	// Create the client configuration for the new client
	clientConfig.Address = make([]net.IPNet, 0, 1)
	clientConfig.Address = append(clientConfig.Address, clientIpNet)

	peerIpAddress := clientIpNetToPeer(clientConfig.Address)

	// Add the new client as a peer to the server
	config.Server.AddPeer(client.base64PublicKey(), peerIpAddress)

	// Set the private key of the new client
	clientConfig.PrivateKey = client.base64PrivateKey()

	// Add the new client to the Clients list
	config.Clients = append(config.Clients, clientConfig)
}

// updateWireguardConfigFiles is a method on the appConfig struct that updates the Wireguard VPN configuration files.
// It accepts a string argument, configPath, which represents the path where the configuration files should be stored.
// The method starts by formatting the default client configuration filename with the number of clients.
// It then attempts to write the latest client configuration to a file at the specified path.
// If an error occurs during this operation, the program is terminated with a relevant error message.
// If the operation is successful, a confirmation message is printed to the console.
// The same process is then repeated for the server configuration.
// As a result, both the client and server configuration files in the specified path are updated with the latest information.
func (config *appConfig) updateWireguardConfigFiles(configPath string) {
	clientFileName := fmt.Sprintf(defaultClientConfigFile, len(config.Clients))

	err := ioutil.WriteFile(configPath+clientFileName, []byte(config.Clients[len(config.Clients)-1].String()), 0666)

	if err != nil {
		log.Fatalf("\nCant't write client config into %s!", configPath+clientFileName)
	} else {
		fmt.Println("\nSuccessfully saved client configuration:", configPath+clientFileName)
	}

	err = ioutil.WriteFile(configPath+defaultServerConfigFile, []byte(config.Server.String()), 0666)

	if err != nil {
		log.Fatalf("\nCant't update server config in %s!", configPath+defaultServerConfigFile)
	} else {
		fmt.Println("\nSuccessfully saved server configuration:", configPath+defaultServerConfigFile)
	}
}

// showClientQrCode is a method on the appConfig struct that generates and displays a QR code from a client's configuration.
// It takes an integer parameter, index, which corresponds to the index of the client in the Clients slice of the appConfig instance.
// It starts by encoding the client configuration into a QR code string using the QREncodeToSmallString function.
// If there is no error in the encoding process, it prints the generated QR code to the console.
// If there is an error, it prints an error message indicating that the QR code could not be generated.
func (config *appConfig) showClientQrCode(index int) {
	qr, err := QREncodeToSmallString(config.Clients[index].String(), false, false)

	fmt.Println("\nClient configuration QR code to scan on mobile device:")

	if err == nil {
		fmt.Print(qr)
	} else {
		fmt.Println("Failed to generate the QR code from the client configuration!")
	}
}
