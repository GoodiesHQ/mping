// Providing a mechanism to resolve hostnames to IP addresses for pinging
package pinger

import (
	"context"
	"net"

	"github.com/goodieshq/goropo"
	"github.com/rs/zerolog/log"
)

// routineResolve resolves a Target to an IP address using DNS
func routineResolve(target *Target, version IPVersion) func(context.Context) (ResolvedTarget, error) {
	return func(ctx context.Context) (ResolvedTarget, error) {
		addr, err := net.ResolveIPAddr(version.Network(), target.Host)
		if err != nil {
			// resolution failed
			log.Warn().Err(err).Msgf(
				"Failed to resolve target %s (%s)",
				target.Label, target.Host,
			)

			// return a ResolvedTarget with nil IP
			return ResolvedTarget{
				IP:        nil,
				IPVersion: IPV0,
				Label:     target.Label,
			}, err
		}

		// resolution succeeded, return the resolved IP
		return ResolvedTarget{
			IP:        addr.IP,
			IPVersion: version,
			Label:     target.Label,
		}, nil
	}
}

// routineDirect return a ResolvedTarget directly for an already known IP
func routineDirect(label string, ip net.IP, ver IPVersion) func(context.Context) (ResolvedTarget, error) {
	return func(_ context.Context) (ResolvedTarget, error) {
		return ResolvedTarget{
			IP:        ip,
			IPVersion: ver,
			Label:     label,
		}, nil
	}
}

// ResolveTargets takes a list of Targets and resolves them to IP addresses
func ResolveTargets(ctx context.Context, targets []*Target, version IPVersion) []ResolvedTarget {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Validate version
	switch version {
	case IPV0, IPV4, IPV6:
	default:
		return nil
	}

	// create a goroutine pool with one worker for each target
	pool := goropo.NewPool(len(targets), 0)
	var futs []*goropo.Future[ResolvedTarget]

	for _, target := range targets {
		// check if target is already an IP
		ip, ipver := tryIP(target.Host)
		if ip != nil {
			// already an IP, no need to resolve
			log.Debug().Msgf("Target %s is a valid IP (%s), no resolution needed", target.Label, ip.String())
			futs = append(
				futs,
				goropo.Submit(pool, ctx, routineDirect(target.Label, ip, ipver)),
			)
		} else {
			// not an IP, resolve via DNS
			futs = append(
				futs,
				goropo.Submit(pool, ctx, routineResolve(target, version)),
			)
		}
	}

	// gracefully wait and close the pool when done resolving all targets
	pool.Close()

	// collect results
	var resolved = make([]ResolvedTarget, 0, len(futs))

	for _, fut := range futs {
		select {
		case <-ctx.Done():
			log.Warn().Msg("Context cancelled, aborting target resolution")
			return nil
		case res := <-fut.Done():
			if res.Ok() {
				resolved = append(resolved, res.Value)
			} else {
				log.Warn().Msgf("Failed to resolve target: %v", res.Err)
				resolved = append(resolved, ResolvedTarget{
					IP:    nil,
					Label: res.Value.Label,
				})
			}
		}
	}

	return resolved
}

// tryIP attempts to parse a string as an IP address and returns the IP and its version
func tryIP(s string) (net.IP, IPVersion) {
	if ip := net.ParseIP(s); ip == nil {
		return nil, IPV0
	} else {
		return GetIP(ip)
	}
}
