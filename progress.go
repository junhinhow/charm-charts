package charmcharts

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// ProgressBar renders a horizontal progress bar with optional label and
// percentage display.
type ProgressBar struct {
	value      float64 // 0.0 to 1.0
	width      int
	fillChar   rune
	emptyChar  rune
	fillStyle  lipgloss.Style
	emptyStyle lipgloss.Style
	label      string
	showPct    bool
}

// NewProgressBar creates a new progress bar with the given value (0.0 to 1.0).
func NewProgressBar(value float64) *ProgressBar {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return &ProgressBar{
		value:      value,
		width:      40,
		fillChar:   '█',
		emptyChar:  '░',
		fillStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		emptyStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		showPct:    true,
	}
}

// Width sets the bar width in characters.
func (p *ProgressBar) Width(w int) *ProgressBar {
	p.width = w
	return p
}

// FillStyle sets the style for the filled portion.
func (p *ProgressBar) FillStyle(s lipgloss.Style) *ProgressBar {
	p.fillStyle = s
	return p
}

// EmptyStyle sets the style for the empty portion.
func (p *ProgressBar) EmptyStyle(s lipgloss.Style) *ProgressBar {
	p.emptyStyle = s
	return p
}

// Label sets a label displayed before the bar.
func (p *ProgressBar) Label(l string) *ProgressBar {
	p.label = l
	return p
}

// ShowPercentage enables/disables percentage display.
func (p *ProgressBar) ShowPercentage(show bool) *ProgressBar {
	p.showPct = show
	return p
}

// FillChar sets the character used for the filled portion.
func (p *ProgressBar) FillChar(c rune) *ProgressBar {
	p.fillChar = c
	return p
}

// EmptyChar sets the character used for the empty portion.
func (p *ProgressBar) EmptyChar(c rune) *ProgressBar {
	p.emptyChar = c
	return p
}

// Render returns the progress bar as a styled string.
func (p *ProgressBar) Render() string {
	filled := int(p.value * float64(p.width))
	if filled > p.width {
		filled = p.width
	}
	empty := p.width - filled

	bar := p.fillStyle.Render(strings.Repeat(string(p.fillChar), filled)) +
		p.emptyStyle.Render(strings.Repeat(string(p.emptyChar), empty))

	var result strings.Builder
	if p.label != "" {
		result.WriteString(p.label + " ")
	}
	result.WriteString(bar)
	if p.showPct {
		result.WriteString(fmt.Sprintf(" %3.0f%%", p.value*100))
	}
	return result.String()
}
