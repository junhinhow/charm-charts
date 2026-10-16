# charm-charts

Terminal charts for the [charmbracelet](https://github.com/charmbracelet) ecosystem.

Sparklines, bar charts, and progress bars — styled with lipgloss.

> **[Leia em Português](README.pt-br.md)**

## Installation

```bash
go get github.com/junhinhow/charm-charts
```

## Charts

### Sparkline

Compact line chart using Unicode block characters. Perfect for inline data visualization.

```go
data := []float64{4, 2, 7, 3, 8, 5, 9, 1, 6, 4}

spark := charmcharts.NewSparkline(data).
    Width(30).
    Label("CPU").
    Style(lipgloss.NewStyle().Foreground(lipgloss.Color("212")))

fmt.Println(spark.Render())
// CPU ▂▁▁ ▁▃▅▅▃▂▂▃▅▇▆▅▃▅▆█▆▃  ▁▃▄▄▃▂
```

The chart always occupies exactly the requested width: when there are fewer
data points than columns the series is expanded by linear interpolation, and
when there are more points they are reduced by averaging each bucket.

Multi-row sparklines for higher resolution:
```go
spark.Height(3)  // 1-8 rows
```

The label is always rendered on the first line of the chart, followed by the
remaining lines padded so the columns stay aligned.

### Bar Chart

Horizontal or vertical bar charts with labels and values.

```go
labels := []string{"Go", "Rust", "Python", "JS"}
values := []float64{85, 72, 95, 68}

chart := charmcharts.NewBarChart(labels, values).
    Width(50).
    Horizontal().
    BarStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99")))

fmt.Println(chart.Render())
//     Go ██████████████████████████████████     85.0
//   Rust █████████████████████████████          72.0
// Python ██████████████████████████████████████ 95.0
//     JS ███████████████████████████            68.0
```

`.BarWidth(n)` sets the thickness of each bar: in horizontal mode every bar
occupies `n` rows (the label and the value are shown on the first row of the
group), in vertical mode it sets the width of each column.

Vertical mode:
```go
chart.Vertical().BarWidth(3)
```

Vertical charts honor `.Width(n)` (the columns shrink to fit) and
`.ShowValues(bool)` (the values are printed above the columns).

### Progress Bar

Horizontal progress bar with percentage and custom characters.

```go
bar := charmcharts.NewProgressBar(0.73).
    Width(30).
    Label("Upload").
    FillStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("42")))

fmt.Println(bar.Render())
// Upload █████████████████████░░░░░░░░░  73%
```

## Styling

All charts accept lipgloss styles, so they integrate seamlessly with your app's theme:

```go
import (
    themes "github.com/junhinhow/charm-themes"

    "charm.land/lipgloss/v2"
)

theme := themes.Dracula
spark := charmcharts.NewSparkline(data).
    Style(lipgloss.NewStyle().Foreground(theme.Primary))
```

## API Reference

### Sparkline
| Method | Description |
|--------|-------------|
| `NewSparkline(data)` | Create with data points |
| `.Width(n)` | Set width in characters (clamped to 1-500) |
| `.Height(n)` | Set height in rows (clamped to 1-8) |
| `.Style(s)` | Set lipgloss style |
| `.Label(s)` | Set prefix label |
| `.Render()` | Return styled string |

### BarChart
| Method | Description |
|--------|-------------|
| `NewBarChart(labels, values)` | Create with labels and values |
| `.Width(n)` | Set max width (clamped to 1-500) |
| `.BarWidth(n)` | Set bar thickness (clamped to 1-20) |
| `.Horizontal()` / `.Vertical()` | Set orientation |
| `.BarStyle(s)` / `.LabelStyle(s)` | Set styles |
| `.ShowValues(bool)` | Toggle value display |
| `.Render()` | Return styled string |

### ProgressBar
| Method | Description |
|--------|-------------|
| `NewProgressBar(value)` | Create with value (0.0-1.0) |
| `.Width(n)` | Set bar width (clamped to 0-500) |
| `.FillStyle(s)` / `.EmptyStyle(s)` | Set styles |
| `.Label(s)` | Set prefix label |
| `.ShowPercentage(bool)` | Toggle percentage |
| `.FillChar(r)` / `.EmptyChar(r)` | Custom characters |
| `.Render()` | Return styled string |

## Notes

- Invalid widths and heights are clamped instead of causing a panic.
- Labels are measured with their display width (never with `len()`), so
  accented letters and CJK characters stay aligned.
- Labels that do not fit their column are truncated without breaking UTF-8.

## License

MIT
