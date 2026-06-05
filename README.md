# charm-charts

Terminal charts for the [charmbracelet](https://github.com/charmbracelet) ecosystem.

Sparklines, bar charts, and progress bars — styled with lipgloss.

> **[Leia em Portugues](README.pt-br.md)**

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
// CPU ▃▁▆▂▇▄█ ▅▃
```

Multi-row sparklines for higher resolution:
```go
spark.Height(3)  // 1-8 rows
```

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
//     Go ████████████████████████████████████████ 85.0
//   Rust ██████████████████████████████████       72.0
// Python █████████████████████████████████████████████ 95.0
//     JS ████████████████████████████████         68.0
```

Vertical mode:
```go
chart.Vertical().BarWidth(3)
```

### Progress Bar

Horizontal progress bar with percentage and custom characters.

```go
bar := charmcharts.NewProgressBar(0.73).
    Width(30).
    Label("Upload").
    FillStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("42")))

fmt.Println(bar.Render())
// Upload ██████████████████████░░░░░░░░░░  73%
```

## Styling

All charts accept lipgloss styles, so they integrate seamlessly with your app's theme:

```go
import themes "github.com/junhinhow/charm-themes"

theme := themes.Dracula
spark := charmcharts.NewSparkline(data).
    Style(lipgloss.NewStyle().Foreground(theme.Primary))
```

## API Reference

### Sparkline
| Method | Description |
|--------|-------------|
| `NewSparkline(data)` | Create with data points |
| `.Width(n)` | Set width in characters |
| `.Height(n)` | Set height in rows (1-8) |
| `.Style(s)` | Set lipgloss style |
| `.Label(s)` | Set prefix label |
| `.Render()` | Return styled string |

### BarChart
| Method | Description |
|--------|-------------|
| `NewBarChart(labels, values)` | Create with labels and values |
| `.Width(n)` | Set max width |
| `.BarWidth(n)` | Set bar thickness |
| `.Horizontal()` / `.Vertical()` | Set orientation |
| `.BarStyle(s)` / `.LabelStyle(s)` | Set styles |
| `.ShowValues(bool)` | Toggle value display |
| `.Render()` | Return styled string |

### ProgressBar
| Method | Description |
|--------|-------------|
| `NewProgressBar(value)` | Create with value (0.0-1.0) |
| `.Width(n)` | Set bar width |
| `.FillStyle(s)` / `.EmptyStyle(s)` | Set styles |
| `.Label(s)` | Set prefix label |
| `.ShowPercentage(bool)` | Toggle percentage |
| `.FillChar(r)` / `.EmptyChar(r)` | Custom characters |
| `.Render()` | Return styled string |

## License

MIT
