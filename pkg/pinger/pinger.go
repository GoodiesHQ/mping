package pinger

import (
	"context"
	"fmt"
	"time"

	"github.com/goodieshq/goropo"
)

type PingerOptions struct {
	ResolvedTargets []ResolvedTarget
	Interval        time.Duration
	Timeout         time.Duration
	Count           uint32
}

type Pinger struct {
	pool     *goropo.Pool
	targets  []ResolvedTarget
	count    uint32
	interval time.Duration
	timeout  time.Duration
}

func NewPinger(opts PingerOptions) *Pinger {
	return &Pinger{
		pool:     goropo.NewPool(len(opts.ResolvedTargets), len(opts.ResolvedTargets)),
		targets:  opts.ResolvedTargets,
		count:    opts.Count,
		interval: opts.Interval,
		timeout:  opts.Timeout,
	}
}

// pingAllOnce pings all targets once concurrently and aggregates their results
func (p *Pinger) pingAllOnce(ctx context.Context) ([]PingResult, error) {
	futs := make([]*goropo.Future[PingResult], 0, len(p.targets))
	results := make([]PingResult, 0, len(p.targets))

	for _, target := range p.targets {
		futs = append(futs, goropo.Submit(p.pool, ctx, func(ctx context.Context) (PingResult, error) {
			return pingOnce(ctx, target, p.timeout), nil
		}))
	}

	// wait for all pings to complete without closing the pool
	p.pool.WaitIdle()

	// collect results
	for _, fut := range futs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-fut.Done():
			results = append(results, result.Value)
		}
	}

	return results, nil
}

// Start the pinger, pinging all targets at the specified interval and count
func (p *Pinger) Start(ctx context.Context) error {
	// counter variable to track number of ping rounds
	counter := uint32(1)

	// column widths based on target labels
	widths := calculateColumnWidths(p.targets)

	// stats for each target
	stats := make([]PingStats, len(p.targets))

	// print final stats on exit
	defer func() {
		if ctx.Err() == context.Canceled {
			// clear the current line to ignore the "^C" output
			fmt.Printf("\r\033[2K")
		}
		printStats(widths, stats)
	}()

	// ticker for interval timing
	clock := time.NewTicker(p.interval)
	defer clock.Stop()

	for {
		results, err := p.pingAllOnce(ctx)
		if err != nil {
			return err
		}

		// update stats for each target
		for i, res := range results {
			if !res.Success {
				stats[i].AddFailure()
			} else {
				stats[i].AddSuccess(res.RTT)
			}
		}

		// print results for this round, include labels every 10 rounds
		printResults(counter, widths, results, counter%10 == 1)

		// check if we've reached the specified count
		if p.count > 0 && counter >= p.count {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-clock.C:
			counter++
		}
	}

	// close the pool and wait for all workers to finish
	p.pool.Close()
	return nil
}
