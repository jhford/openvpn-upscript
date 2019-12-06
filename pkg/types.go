package pkg

import "net"

// NameServer is a nameserver
type NameServer net.IP

func (n NameServer) String() string {
	return net.IP(n).String()
}
