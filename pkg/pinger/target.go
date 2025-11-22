// Parsing and representation of ping targets.
package pinger

import (
	"net"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
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

func parseHostFromURL(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to parse URL: %s", uri)
		return uri
	}
	return u.Hostname()
}

// parseTarget returns addr and label. If no '=' present, label==addr.
func ParseTarget(s string) *Target {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	var host, label string
	if strings.Contains(s, "://") {
		host = parseHostFromURL(s)
		log.Info().Msg("Parsed host from URL: " + host)
		return &Target{
			Host:  host,
			Label: host,
		}
	}

	// check for '=' to separate host and label
	if i := strings.LastIndexByte(s, '='); i >= 0 {
		host = strings.TrimSpace(s[:i])
		label = strings.TrimSpace(s[i+1:])
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
