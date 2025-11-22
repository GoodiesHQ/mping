package pinger

import (
	"context"
	"net"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingResult struct {
	Target  ResolvedTarget
	Success bool
	RTT     time.Duration
	Error   error
}

// Perform a single ping to the target with the specified timeout.
func pingOnce(ctx context.Context, target ResolvedTarget, timeout time.Duration) PingResult {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ip := &net.IPAddr{IP: target.IP}

	// Set up the probe
	probe := probing.New("")
	probe.SetIPAddr(ip)
	probe.Count = 1
	probe.Timeout = timeout
	probe.Interval = 1
	probe.RecordRtts = true

	// Determine network type based on IP version
	switch target.IPVersion {
	case IPV4:
		probe.SetNetwork("ip4")
	case IPV6:
		probe.SetNetwork("ip6")
	default:
		if target.IP.To4() != nil {
			probe.SetNetwork("ip4")
		} else {
			probe.SetNetwork("ip6")
		}
	}

	// use privileged ICMP
	probe.SetPrivileged(true)

	// Variable to hold RTT for the ping
	var rtt time.Duration

	probe.OnRecv = func(pkt *probing.Packet) {
		rtt = pkt.Rtt
	}

	// Run the probe with context
	err := probe.RunWithContext(ctx)

	if err == nil {
		stats := probe.Statistics()
		if stats.PacketsRecv == 0 {
			err = context.DeadlineExceeded
		} else if rtt == 0 && len(stats.Rtts) > 0 {
			rtt = stats.Rtts[0]
		}
	}

	return PingResult{
		Target:  target,
		Success: err == nil,
		RTT:     rtt,
		Error:   err,
	}
}
