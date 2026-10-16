// Package charmcharts provides terminal chart components that integrate
// with the charmbracelet ecosystem (lipgloss styling, Bubble Tea models).
package charmcharts

import (
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// Limites usados para normalizar as opções do sparkline.
const (
	minSparklineWidth  = 1
	maxSparklineWidth  = 500
	minSparklineHeight = 1
	maxSparklineHeight = 8
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
// Values outside the 1..500 range are clamped to the nearest valid value.
func (s *Sparkline) Width(w int) *Sparkline {
	s.width = clamp(w, minSparklineWidth, maxSparklineWidth)
	return s
}

// Height sets the display height in rows (1-8).
// Values outside the 1..8 range are clamped to the nearest valid value.
func (s *Sparkline) Height(h int) *Sparkline {
	s.height = clamp(h, minSparklineHeight, maxSparklineHeight)
	return s
}

// Style sets the lipgloss style for the sparkline characters.
func (s *Sparkline) Style(st lipgloss.Style) *Sparkline {
	s.style = st
	return s
}

// Label sets an optional label displayed before the sparkline.
// The label is always rendered on the first line of the chart, followed by
// the remaining lines padded so the columns stay aligned.
func (s *Sparkline) Label(l string) *Sparkline {
	s.label = l
	return s
}

// Render returns the sparkline as a styled string.
func (s *Sparkline) Render() string {
	if len(s.data) == 0 {
		return ""
	}

	width := clamp(s.width, minSparklineWidth, maxSparklineWidth)
	height := clamp(s.height, minSparklineHeight, maxSparklineHeight)

	// Resample data to fit width
	resampled := resample(s.data, width)

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

	// Cada linha do gráfico representa (len(blocks)-1) sub-níveis.
	levels := height * (len(blocks) - 1)

	rows := make([]string, height)
	for _, v := range resampled {
		normalized := (v - minVal) / span
		cellLevel := int(normalized * float64(levels))

		for row := 0; row < height; row++ {
			rowFromBottom := height - 1 - row
			localLevel := cellLevel - rowFromBottom*(len(blocks)-1)

			switch {
			case localLevel <= 0:
				rows[row] += " "
			case localLevel >= len(blocks)-1:
				rows[row] += string(blocks[len(blocks)-1])
			default:
				rows[row] += string(blocks[localLevel])
			}
		}
	}

	// O rótulo fica na primeira linha do gráfico; as demais linhas recebem
	// o mesmo espaçamento para manter as colunas alinhadas.
	prefix := ""
	if s.label != "" {
		prefix = s.label + " "
	}
	prefixWidth := displayWidth(prefix)

	lines := make([]string, height)
	for i, row := range rows {
		if i == 0 {
			lines[i] = prefix + row
		} else {
			lines[i] = strings.Repeat(" ", prefixWidth) + row
		}
	}

	return s.style.Render(strings.Join(lines, "\n"))
}

// resample reduces or expands data to fit the target width.
// When there are fewer points than columns the data is expanded with linear
// interpolation; when there are more points the data is reduced by averaging
// each bucket.
func resample(data []float64, width int) []float64 {
	if len(data) == 0 || width <= 0 {
		return nil
	}
	if len(data) == width {
		return append([]float64(nil), data...)
	}

	if len(data) > width {
		result := make([]float64, width)
		ratio := float64(len(data)) / float64(width)
		for i := 0; i < width; i++ {
			start := int(float64(i) * ratio)
			end := int(float64(i+1) * ratio)
			if end > len(data) {
				end = len(data)
			}
			if end <= start {
				end = start + 1
			}
			sum := 0.0
			for j := start; j < end; j++ {
				sum += data[j]
			}
			result[i] = sum / float64(end-start)
		}
		return result
	}

	// Expansão: interpolação linear entre os pontos originais.
	result := make([]float64, width)
	if len(data) == 1 {
		for i := range result {
			result[i] = data[0]
		}
		return result
	}
	if width == 1 {
		result[0] = data[0]
		return result
	}

	step := float64(len(data)-1) / float64(width-1)
	for i := 0; i < width; i++ {
		pos := float64(i) * step
		low := int(math.Floor(pos))
		high := int(math.Ceil(pos))
		if high >= len(data) {
			high = len(data) - 1
		}
		frac := pos - float64(low)
		result[i] = data[low]*(1-frac) + data[high]*frac
	}
	return result
}
