package pinger

import (
	"fmt"
	"strconv"
	"strings"
)

const widthMin = 9

type PrintOptions struct {
	showRTT bool
	widths  []int
}

// calculateColumnWidths computes the width of each column based on the length of the target labels
func calculateColumnWidths(targets []ResolvedTarget) []int {
	widths := make([]int, len(targets))
	for i, target := range targets {
		labelLen := len(target.Label)
		if labelLen > widthMin {
			widths[i] = labelLen
		} else {
			widths[i] = widthMin
		}
	}
	return widths
}

// printStats prints the aggregated statistics in a formatted table
func printStats(opts *PrintOptions, stats []PingStats) {
	var msg string

	fmt.Println()
	fmt.Print("Total ")
	for i := range stats {
		format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
		fmt.Printf(format, strconv.FormatInt(int64(stats[i].Success()+stats[i].Failure()), 10))
	}

	fmt.Println()
	fmt.Print("Fails ")
	for i := range stats {
		format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
		fmt.Printf(format, strconv.FormatInt(int64(stats[i].Failure()), 10))
	}

	fmt.Println()
	fmt.Print("Loss  ")
	for i := range stats {
		format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
		msg = fmt.Sprintf("%.1f%%", stats[i].PacketLoss())
		fmt.Printf(format, msg)
	}

	fmt.Println()
	fmt.Print("RTT   ")
	for i := range stats {
		format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
		if ms := stats[i].AvgRTT().Milliseconds(); ms == 0 {
			msg = "-"
		} else {
			msg = fmt.Sprintf("%dms", ms)
		}
		fmt.Printf(format, msg)
	}

	fmt.Println()
}

// printResults prints the results of a ping round in a formatted table
func printResults(counter uint32, opts *PrintOptions, results []PingResult, showLabels bool) {
	// widths := calculateColumnWidths(results)
	ctr := fmt.Sprintf("%5d)", counter)
	spc := strings.Repeat(" ", len(ctr))

	if showLabels {
		fmt.Println()
		fmt.Print(spc)
		for i, result := range results {
			format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
			fmt.Printf(format, result.Target.Label)
		}
		fmt.Println()
		fmt.Print(spc)
		for i := range results {
			format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
			fmt.Printf(format, strings.Repeat("-", opts.widths[i]))
		}
		fmt.Println()
	}

	fmt.Print(ctr)
	for i, result := range results {
		format := " %-" + strconv.Itoa(opts.widths[i]) + "s"
		var msg = ""
		if opts.showRTT {
			if result.Success {
				msg = strconv.Itoa(int(result.RTT.Milliseconds())) + "ms"
			} else {
				msg = "FAIL"
			}
		}
		fmt.Printf(format, msg)
	}
	fmt.Println()
}
