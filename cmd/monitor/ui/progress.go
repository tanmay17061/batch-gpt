package ui

import (
	"fmt"
	"strings"
)

func renderProgress(completed, total, width int) string {
	if total == 0 {
		return progressBarStyle.Render("[No requests]")
	}

	percentage := float64(completed) / float64(total)
	filledWidth := int(float64(width) * percentage)

	filled := progressBarFilledStyle.Render(strings.Repeat("█", filledWidth))
	empty := progressBarStyle.Render(strings.Repeat("░", width-filledWidth))

	bar := fmt.Sprintf("[%s%s] %d%% (%d/%d requests completed)",
		filled, empty, int(percentage*100), completed, total)

	return progressBarStyle.Render(bar)
}

func renderProgressWithCounts(completed, total, width int) string {
	if total == 0 {
		return progressBarStyle.Render("[No requests]")
	}

	// Calculate space needed for the count suffix
	countSuffix := fmt.Sprintf(" %d/%d", completed, total)
	// Reserve space for the count and brackets
	barWidth := width - len(countSuffix) - 2 // 2 for the brackets

	percentage := float64(completed) / float64(total)
	filledWidth := int(float64(barWidth) * percentage)

	filled := progressBarFilledStyle.Render(strings.Repeat("█", filledWidth))
	empty := progressBarStyle.Render(strings.Repeat("░", barWidth-filledWidth))

	bar := fmt.Sprintf("[%s%s]%s",
		filled, empty, countSuffix)

	return progressBarStyle.Render(bar)
}
