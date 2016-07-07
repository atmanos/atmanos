package ip

import (
	"fmt"
)

type HardwareAddr [6]byte

func (addr HardwareAddr) String() string {
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		addr[0],
		addr[1],
		addr[2],
		addr[3],
		addr[4],
		addr[5],
	)
}

type IPAddr [4]byte

const (
	EtherTypeIPv4 = 0x0800
	EtherTypeArp  = 0x0806
	EtherTypeIPv6 = 0x86DD
)

type EtherType uint16

func (t EtherType) Name() string {
	switch t {
	case EtherTypeIPv4:
		return "IPV4"
	case EtherTypeIPv6:
		return "IPV6"
	case EtherTypeArp:
		return "ARP"
	default:
		return "unknown"
	}
}

// https://en.wikipedia.org/wiki/Ethernet_frame
// https://en.wikipedia.org/wiki/EtherType
type EthernetHeader struct {
	Destination HardwareAddr
	Source      HardwareAddr
	Type        EtherType
}
