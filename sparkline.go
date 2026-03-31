// Package charmcharts provides terminal chart components that integrate
// with the charmbracelet ecosystem (lipgloss styling, Bubble Tea models).
package charmcharts

import (
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// Sparkline renders a compact line chart using Unicode block characters.
// It fits data into a fixed width and height, scaling values automatically.
type Sparkline struct {
	data   []float64
	width  int
	height int
	style  lipgloss.Style
	label  string
}

// blocks used for sparkline rendering, from lowest to highest.
var blocks = []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// NewSparkline creates a new sparkline chart with the given data points.
func NewSparkline(data []float64) *Sparkline {
	return &Sparkline{
		data:   data,
		width:  40,
		height: 1,
		style:  lipgloss.NewStyle(),
	}
}

// Width sets the display width in characters.
func (s *Sparkline) Width(w int) *Sparkline {
	s.width = w
	return s
}

// Height sets the display height in rows (1-8).
func (s *Sparkline) Height(h int) *Sparkline {
	if h < 1 {
		h = 1
	}
	if h > 8 {
		h = 8
	}
	s.height = h
	return s
}

// Style sets the lipgloss style for the sparkline characters.
func (s *Sparkline) Style(st lipgloss.Style) *Sparkline {
	s.style = st
	return s
}

// Label sets an optional label displayed before the sparkline.
func (s *Sparkline) Label(l string) *Sparkline {
	s.label = l
	return s
}

// Render returns the sparkline as a styled string.
func (s *Sparkline) Render() string {
	if len(s.data) == 0 {
		return ""
	}

	// Resample data to fit width
	resampled := resample(s.data, s.width)

	// Find min/max for scaling
	minVal, maxVal := resampled[0], resampled[0]
	for _, v := range resampled {
		minVal = math.Min(minVal, v)
		maxVal = math.Max(maxVal, v)
	}

	span := maxVal - minVal
	if span == 0 {
		span = 1
	}

	levels := len(blocks) * s.height

	if s.height == 1 {
		var b strings.Builder
		if s.label != "" {
			b.WriteString(s.label + " ")
		}
		for _, v := range resampled {
			normalized := (v - minVal) / span
			idx := int(normalized * float64(len(blocks)-1))
			if idx >= len(blocks) {
				idx = len(blocks) - 1
			}
			b.WriteRune(blocks[idx])
		}
		return s.style.Render(b.String())
	}

	// Multi-row rendering
	rows := make([]strings.Builder, s.height)
	for _, v := range resampled {
		normalized := (v - minVal) / span
		cellLevel := int(normalized * float64(levels-1))

		for row := 0; row < s.height; row++ {
			rowFromBottom := s.height - 1 - row
			rowBase := rowFromBottom * (len(blocks) - 1)
			localLevel := cellLevel - rowBase

			if localLevel <= 0 {
				rows[row].WriteRune(' ')
			} else if localLevel >= len(blocks)-1 {
				rows[row].WriteRune(blocks[len(blocks)-1])
			} else {
				rows[row].WriteRune(blocks[localLevel])
			}
		}
	}

	var result strings.Builder
	if s.label != "" {
		result.WriteString(s.label + "\n")
	}
	for i, row := range rows {
		result.WriteString(s.style.Render(row.String()))
		if i < len(rows)-1 {
			result.WriteRune('\n')
		}
	}
	return result.String()
}

// resample reduces or expands data to fit the target width.
func resample(data []float64, width int) []float64 {
	if len(data) == 0 {
		return nil
	}
	if len(data) <= width {
		return data
	}
	result := make([]float64, width)
	ratio := float64(len(data)) / float64(width)
	for i := 0; i < width; i++ {
		start := int(float64(i) * ratio)
		end := int(float64(i+1) * ratio)
		if end > len(data) {
			end = len(data)
		}
		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(end-start)
	}
	return result
}
