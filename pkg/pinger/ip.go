package pinger

import (
	"net"

	"github.com/rs/zerolog/log"
)

type IPVersion int

const (
	IPV0 IPVersion = 0 // unknown or unspecified
	IPV4 IPVersion = 4 // IPv4
	IPV6 IPVersion = 6 // IPv6
)

// Network returns the network string for net.Dial and friends based on the IP version
func (v IPVersion) Network() string {
	switch v {
	case IPV0:
		return "ip"
	case IPV4:
		return "ip4"
	case IPV6:
		return "ip6"
	default:
		log.Warn().Msgf("Unknown IPVersion %d, defaulting to empty network string", v)
		return ""
	}
}

// GetIP returns the appropriately sized IP address and its version
func GetIP(ip net.IP) (net.IP, IPVersion) {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4, IPV4
	}
	return ip.To16(), IPV6
}
