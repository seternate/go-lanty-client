package network

import "net"

func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	ip = conn.LocalAddr().(*net.UDPAddr).IP
	return
}
