# wg-quick-config

wg-quick-config is a simple configuration tool designed for the Wiresock VPN Gateway. This tool greatly simplifies the processes of creating, managing, and implementing Wireguard configurations. You can learn more about the functionalities of Wiresock VPN Gateway [here](https://www.wiresock.net/wiresock-vpn-gateway/).

## Key Features
- **Automated Configuration:** Streamline the creation of Wireguard server and client configurations with our automated tools.
- **Client Setup Simplified:** Our system generates QR codes to facilitate a hassle-free client setup process on mobile devices.

## Usage

Kickstart your WireGuard server endpoint setup with the following command. Remember to jot down the UDP port number displayed, as it will be useful later:

```bash
wg-quick-config -add -start
```

### Other Useful Commands

- **Add New Peer & Restart WireGuard Tunnel:** 
```bash
wg-quick-config -add -restart
```
- **Stop WireGuard Tunnel:** 
```bash
wg-quick-config -stop
```
- **Start WireGuard Tunnel:** 
```bash
wg-quick-config -start
```
- **Display QR Code for First Client:** 
```bash
wg-quick-config -qrcode 1
```

## Contributing

We greatly value your contributions! If you want to contribute to this project, please feel free to open issues or create pull requests.

Whether you're looking to report a bug, suggest an improvement, or propose a new feature, we're always excited to hear from you!
