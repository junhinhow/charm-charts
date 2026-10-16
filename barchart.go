package charmcharts

import (
	"fmt"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// Limites usados para normalizar as opções do gráfico de barras.
const (
	minBarChartWidth    = 1
	maxBarChartWidth    = 500
	minBarChartBarWidth = 1
	maxBarChartBarWidth = 20
)

// verticalChartHeight é a altura (em linhas) da área de barras no modo
// vertical.
const verticalChartHeight = 10

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
// Values outside the 1..500 range are clamped to the nearest valid value.
func (b *BarChart) Width(w int) *BarChart {
	b.width = clamp(w, minBarChartWidth, maxBarChartWidth)
	return b
}

// BarWidth sets the thickness of each bar (in rows for horizontal charts).
// Values outside the 1..20 range are clamped to the nearest valid value.
func (b *BarChart) BarWidth(w int) *BarChart {
	b.barWidth = clamp(w, minBarChartBarWidth, maxBarChartBarWidth)
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

// maxValue devolve o maior valor do gráfico (nunca zero).
func (b *BarChart) maxValue() float64 {
	maxVal := 0.0
	for _, v := range b.values {
		maxVal = math.Max(maxVal, v)
	}
	if maxVal == 0 {
		maxVal = 1
	}
	return maxVal
}

// labelAt devolve o rótulo do índice informado, quando existir.
func (b *BarChart) labelAt(i int) string {
	if i < 0 || i >= len(b.labels) {
		return ""
	}
	return b.labels[i]
}

// maxLabelWidth devolve a maior largura de exibição entre os rótulos.
// Rótulos sem valor correspondente são ignorados.
func (b *BarChart) maxLabelWidth() int {
	width := 0
	for i := range b.values {
		if w := displayWidth(b.labelAt(i)); w > width {
			width = w
		}
	}
	return width
}

// maxValueWidth devolve a maior largura de exibição entre os valores
// formatados (usado para reservar a coluna de valores).
func (b *BarChart) maxValueWidth() int {
	width := 0
	for _, v := range b.values {
		if w := displayWidth(formatValue(v)); w > width {
			width = w
		}
	}
	return width
}

// formatValue formata um valor do gráfico.
func formatValue(v float64) string {
	return fmt.Sprintf("%.1f", v)
}

func (b *BarChart) renderHorizontal() string {
	width := clamp(b.width, minBarChartWidth, maxBarChartWidth)
	thickness := clamp(b.barWidth, minBarChartBarWidth, maxBarChartBarWidth)
	maxVal := b.maxValue()

	labelWidth := b.maxLabelWidth()
	valueWidth := 0
	if b.showValues {
		valueWidth = b.maxValueWidth()
	}

	// Espaço disponível para as barras: o que sobra da largura pedida
	// depois das colunas fixas de rótulo e valor.
	barArea := width - labelWidth - 1
	if b.showValues {
		barArea -= valueWidth + 1
	}
	if barArea < 1 {
		barArea = 1
	}

	lines := make([]string, 0, len(b.values)*thickness)
	for i, val := range b.values {
		barLen := int(math.Round(val / maxVal * float64(barArea)))
		if barLen < 0 {
			barLen = 0
		}
		if barLen > barArea {
			barLen = barArea
		}
		bar := strings.Repeat("█", barLen)

		// A célula da barra tem largura fixa: as barras menores são
		// completadas com espaços para que os valores fiquem alinhados.
		barCell := b.barStyle.Render(bar) + strings.Repeat(" ", barArea-barLen)

		// A espessura da barra é dada pelo número de linhas do grupo:
		// o rótulo e o valor ficam na primeira linha e as demais repete o
		// preenchimento da barra.
		for row := 0; row < thickness; row++ {
			labelPart := padLeft(b.labelAt(i), labelWidth)
			if row > 0 {
				labelPart = strings.Repeat(" ", labelWidth)
			}

			line := b.labelStyle.Render(labelPart) + " " + barCell
			if b.showValues {
				valuePart := strings.Repeat(" ", valueWidth)
				if row == 0 {
					valuePart = padLeft(formatValue(val), valueWidth)
				}
				line += " " + valuePart
			}
			lines = append(lines, line)
		}
	}

	return b.style.Render(strings.Join(lines, "\n"))
}

func (b *BarChart) renderVertical() string {
	width := clamp(b.width, minBarChartWidth, maxBarChartWidth)
	barWidth := clamp(b.barWidth, minBarChartBarWidth, maxBarChartBarWidth)
	maxVal := b.maxValue()
	count := len(b.values)

	// Largura natural de cada coluna: o maior entre a barra, o valor e o
	// rótulo, para que nada seja cortado sem necessidade.
	colWidth := barWidth
	if b.showValues {
		if w := b.maxValueWidth(); w > colWidth {
			colWidth = w
		}
	}
	if w := b.maxLabelWidth(); w > colWidth {
		colWidth = w
	}

	// A largura pedida é o limite: as colunas encolhem para caber.
	available := width - (count - 1)
	if available < count {
		available = count
	}
	if limit := available / count; colWidth > limit {
		colWidth = limit
	}
	if colWidth < 1 {
		colWidth = 1
	}

	lines := make([]string, 0, verticalChartHeight+2)

	// Linha de valores (acima do gráfico).
	if b.showValues {
		var row strings.Builder
		for i := range b.values {
			if i > 0 {
				row.WriteString(" ")
			}
			value := fitLabel(formatValue(b.values[i]), colWidth)
			row.WriteString(center(value, colWidth))
		}
		lines = append(lines, row.String())
	}

	// Área de barras.
	for r := 0; r < verticalChartHeight; r++ {
		threshold := float64(verticalChartHeight-r) / float64(verticalChartHeight)
		var row strings.Builder
		for i, val := range b.values {
			if i > 0 {
				row.WriteString(" ")
			}
			bar := strings.Repeat(" ", barWidth)
			if val/maxVal >= threshold {
				bar = b.barStyle.Render(strings.Repeat("█", barWidth))
			}
			// Centraliza a barra dentro da coluna.
			pad := colWidth - barWidth
			left := pad / 2
			row.WriteString(strings.Repeat(" ", left) + bar + strings.Repeat(" ", pad-left))
		}
		lines = append(lines, row.String())
	}

	// Linha de rótulos (abaixo do gráfico), centralizados em cada coluna.
	var labelRow strings.Builder
	for i := range b.values {
		if i > 0 {
			labelRow.WriteString(" ")
		}
		label := fitLabel(b.labelAt(i), colWidth)
		labelRow.WriteString(b.labelStyle.Render(center(label, colWidth)))
	}
	lines = append(lines, labelRow.String())

	return b.style.Render(strings.Join(lines, "\n"))
}
