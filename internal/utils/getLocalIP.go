package utils

import (
	"net"
)

/*
Function to retrieve an IP address that can be used to access the Go server on the network
*/
func getLocalIp() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addresses {
		ipNet, isValid := address.(*net.IPNet)
		if !isValid {
			continue
		}
		ip := ipNet.IP
		if ip.IsLoopback() {
			continue
		}
		if ip.To4() == nil {
			continue
		}
		return ip.String()
	}
	return ""
}

