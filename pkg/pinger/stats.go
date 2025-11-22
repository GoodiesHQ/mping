// Keeping track of ping statistics for each target
package pinger

import "time"

// PingStats holds aggregated statistics for a ping session
type PingStats struct {
	success uint32
	failure uint32
	avgRTT  time.Duration
}

func (ps *PingStats) AvgRTT() time.Duration {
	if ps.success == 0 {
		return 0
	}

	return ps.avgRTT
}

func (ps *PingStats) Success() uint32 {
	return ps.success
}

func (ps *PingStats) Failure() uint32 {
	return ps.failure
}

func (ps *PingStats) AddSuccess(rtt time.Duration) {
	ps.success += 1
	ps.addRTT(rtt)
}

func (ps *PingStats) AddFailure() {
	ps.failure += 1
}

// PacketLoss calculates the packet loss percentage based on total packets sent
func (ps *PingStats) PacketLoss() float64 {
	total := ps.success + ps.failure

	if total == 0 {
		return 0.0
	}

	if ps.success == 0 {
		return 100.0
	}

	return float64(ps.failure) / float64(total) * 100.0
}

// AddRTT updates the average RTT with a new RTT measurement
// Using: newAvg = oldAvg + (x - oldAvg) / count
func (ps *PingStats) addRTT(rtt time.Duration) {
	if ps.success == 1 {
		// first RTT measurement
		ps.avgRTT = rtt
		return
	}

	// incremental average calculation
	old := float64(ps.avgRTT)
	r := float64(rtt)
	n := float64(ps.success)

	newAvg := old + (r-old)/n
	ps.avgRTT = time.Duration(newAvg + 0.5)
}
