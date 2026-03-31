package charmcharts

import (
	"fmt"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// BarChart renders horizontal or vertical bar charts in the terminal.
type BarChart struct {
	labels     []string
	values     []float64
	width      int
	barWidth   int
	horizontal bool
	style      lipgloss.Style
	barStyle   lipgloss.Style
	labelStyle lipgloss.Style
	showValues bool
}

// NewBarChart creates a new bar chart with labels and values.
func NewBarChart(labels []string, values []float64) *BarChart {
	return &BarChart{
		labels:     labels,
		values:     values,
		width:      50,
		barWidth:   1,
		horizontal: true,
		style:      lipgloss.NewStyle(),
		barStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		labelStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		showValues: true,
	}
}

// Width sets the maximum width of the chart.
func (b *BarChart) Width(w int) *BarChart {
	b.width = w
	return b
}

// BarWidth sets the thickness of each bar (in rows for horizontal charts).
func (b *BarChart) BarWidth(w int) *BarChart {
	b.barWidth = w
	return b
}

// Horizontal sets horizontal rendering mode (default).
func (b *BarChart) Horizontal() *BarChart {
	b.horizontal = true
	return b
}

// Vertical sets vertical rendering mode.
func (b *BarChart) Vertical() *BarChart {
	b.horizontal = false
	return b
}

// Style sets the overall chart style.
func (b *BarChart) Style(s lipgloss.Style) *BarChart {
	b.style = s
	return b
}

// BarStyle sets the style for the bar fills.
func (b *BarChart) BarStyle(s lipgloss.Style) *BarChart {
	b.barStyle = s
	return b
}

// LabelStyle sets the style for labels.
func (b *BarChart) LabelStyle(s lipgloss.Style) *BarChart {
	b.labelStyle = s
	return b
}

// ShowValues enables/disables value display next to bars.
func (b *BarChart) ShowValues(show bool) *BarChart {
	b.showValues = show
	return b
}

// Render returns the bar chart as a styled string.
func (b *BarChart) Render() string {
	if len(b.values) == 0 {
		return ""
	}

	if b.horizontal {
		return b.renderHorizontal()
	}
	return b.renderVertical()
}

func (b *BarChart) renderHorizontal() string {
	maxVal := 0.0
	for _, v := range b.values {
		maxVal = math.Max(maxVal, v)
	}
	if maxVal == 0 {
		maxVal = 1
	}

	// Find longest label for alignment
	maxLabelLen := 0
	for _, l := range b.labels {
		if len(l) > maxLabelLen {
			maxLabelLen = len(l)
		}
	}

	barMaxWidth := b.width - maxLabelLen - 2
	if b.showValues {
		barMaxWidth -= 10
	}
	if barMaxWidth < 5 {
		barMaxWidth = 5
	}

	var result strings.Builder
	for i, val := range b.values {
		label := ""
		if i < len(b.labels) {
			label = b.labels[i]
		}

		paddedLabel := fmt.Sprintf("%*s", maxLabelLen, label)
		barLen := int(math.Round(val / maxVal * float64(barMaxWidth)))
		if barLen < 0 {
			barLen = 0
		}

		bar := strings.Repeat("█", barLen)

		line := b.labelStyle.Render(paddedLabel) + " " + b.barStyle.Render(bar)
		if b.showValues {
			line += fmt.Sprintf(" %.1f", val)
		}

		for row := 0; row < b.barWidth; row++ {
			result.WriteString(line)
			if i < len(b.values)-1 || row < b.barWidth-1 {
				result.WriteRune('\n')
			}
		}
	}

	return b.style.Render(result.String())
}

func (b *BarChart) renderVertical() string {
	maxVal := 0.0
	for _, v := range b.values {
		maxVal = math.Max(maxVal, v)
	}
	if maxVal == 0 {
		maxVal = 1
	}

	chartHeight := 10
	colWidth := b.barWidth + 1

	// Build columns
	rows := make([]strings.Builder, chartHeight+1) // +1 for labels

	for row := 0; row < chartHeight; row++ {
		threshold := float64(chartHeight-row) / float64(chartHeight)
		for _, val := range b.values {
			normalized := val / maxVal
			if normalized >= threshold {
				rows[row].WriteString(b.barStyle.Render(strings.Repeat("█", b.barWidth)))
			} else {
				rows[row].WriteString(strings.Repeat(" ", b.barWidth))
			}
			rows[row].WriteRune(' ')
		}
	}

	// Label row
	for i := range b.values {
		label := ""
		if i < len(b.labels) {
			label = b.labels[i]
			if len(label) > colWidth {
				label = label[:colWidth-1] + "."
			}
		}
		rows[chartHeight].WriteString(b.labelStyle.Render(fmt.Sprintf("%-*s", colWidth, label)))
	}

	var result strings.Builder
	for i, row := range rows {
		result.WriteString(row.String())
		if i < len(rows)-1 {
			result.WriteRune('\n')
		}
	}

	return b.style.Render(result.String())
}
