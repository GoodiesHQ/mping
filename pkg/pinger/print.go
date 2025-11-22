package pinger

import (
	"fmt"
	"strconv"
	"strings"
)

const widthMin = 9

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
func printStats(widths []int, stats []PingStats) {
	var msg string

	fmt.Println()
	fmt.Print("Fails ")
	for i := range stats {
		format := " %-" + strconv.Itoa(widths[i]) + "s"
		msg = fmt.Sprintf("%d", stats[i].Failure())
		fmt.Printf(format, msg)
	}

	fmt.Println()
	fmt.Print("Loss  ")
	for i := range stats {
		format := " %-" + strconv.Itoa(widths[i]) + "s"
		msg = fmt.Sprintf("%.1f%%", stats[i].PacketLoss())
		fmt.Printf(format, msg)
	}

	fmt.Println()
	fmt.Print("RTT   ")
	for i := range stats {
		format := " %-" + strconv.Itoa(widths[i]) + "s"
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
func printResults(counter uint32, widths []int, results []PingResult, showLabels bool) {
	// widths := calculateColumnWidths(results)
	ctr := fmt.Sprintf("%5d)", counter)
	spc := strings.Repeat(" ", len(ctr))

	if showLabels {
		fmt.Println()
		fmt.Print(spc)
		for i, result := range results {
			format := " %-" + strconv.Itoa(widths[i]) + "s"
			fmt.Printf(format, result.Target.Label)
		}
		fmt.Println()
		fmt.Print(spc)
		for i := range results {
			format := " %-" + strconv.Itoa(widths[i]) + "s"
			fmt.Printf(format, strings.Repeat("-", widths[i]))
		}
		fmt.Println()
	}

	fmt.Print(ctr)
	for i, result := range results {
		format := " %-" + strconv.Itoa(widths[i]) + "s"
		if result.Success {
			rttStr := strconv.Itoa(int(result.RTT.Milliseconds())) + "ms"
			fmt.Printf(format, rttStr)
		} else {
			fmt.Printf(format, "FAIL")
		}
	}
	fmt.Println()
}
