// Parsing and representation of ping targets.
package pinger

import (
	"net"
	"strings"
)

// Target represents a ping target with a host and an optional label.
type Target struct {
	Host  string
	Label string
}

// ResolvedTarget represents a target that has been resolved to an IP address.
type ResolvedTarget struct {
	IP        net.IP
	IPVersion IPVersion
	Label     string
}

// parseTarget returns addr and label. If no '=' present, label==addr.
func ParseTarget(s string) *Target {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	// check for '=' to separate host and label
	if i := strings.IndexByte(s, '='); i >= 0 {
		host := strings.TrimSpace(s[:i])
		label := strings.TrimSpace(s[i+1:])
		if label == "" {
			label = host
		}
		return &Target{
			Host:  host,
			Label: label,
		}
	}

	// no label specified, use the host as label
	return &Target{
		Host:  s,
		Label: s,
	}
}
