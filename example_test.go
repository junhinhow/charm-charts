package charmcharts_test

import (
	"fmt"

	charmcharts "github.com/junhinhow/charm-charts"

	"charm.land/lipgloss/v2"
)

func ExampleNewSparkline() {
	spark := charmcharts.NewSparkline([]float64{4, 2, 7, 3, 8, 5, 9, 1, 6, 4}).
		Width(12).
		Style(lipgloss.NewStyle())
	fmt.Println(spark.Render())
	// Output:
	// ▂ ▄▄▂▇▄█▃ ▄▂
}

func ExampleNewBarChart() {
	chart := charmcharts.NewBarChart([]string{"Go", "Rust"}, []float64{85, 72}).
		Width(30).
		Style(lipgloss.NewStyle()).
		BarStyle(lipgloss.NewStyle()).
		LabelStyle(lipgloss.NewStyle())
	fmt.Println(chart.Render())
	// Output:
	// Go ████████████████████ 85.0
	// Rust █████████████████    72.0
}

func ExampleNewProgressBar() {
	bar := charmcharts.NewProgressBar(0.5).
		Width(10).
		FillStyle(lipgloss.NewStyle()).
		EmptyStyle(lipgloss.NewStyle())
	fmt.Println(bar.Render())
	// Output:
	// █████░░░░░  50%
}
