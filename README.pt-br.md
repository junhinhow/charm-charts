# charm-charts

Graficos de terminal para o ecossistema [charmbracelet](https://github.com/charmbracelet).

Sparklines, graficos de barra e barras de progresso — estilizados com lipgloss.

> **[Read in English](README.md)**

## Instalacao

```bash
go get github.com/junhinhow/charm-charts
```

## Graficos

### Sparkline

Grafico de linha compacto usando caracteres Unicode de bloco. Perfeito para visualizacao inline de dados.

```go
data := []float64{4, 2, 7, 3, 8, 5, 9, 1, 6, 4}

spark := charmcharts.NewSparkline(data).
    Width(30).
    Label("CPU").
    Style(lipgloss.NewStyle().Foreground(lipgloss.Color("212")))

fmt.Println(spark.Render())
// CPU ▃▁▆▂▇▄█ ▅▃
```

Sparklines multi-linha para maior resolucao:
```go
spark.Height(3)  // 1-8 linhas
```

### Grafico de Barras

Graficos de barra horizontais ou verticais com labels e valores.

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

### Barra de Progresso

Barra de progresso horizontal com porcentagem e caracteres customizaveis.

```go
bar := charmcharts.NewProgressBar(0.73).
    Width(30).
    Label("Upload").
    FillStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("42")))

fmt.Println(bar.Render())
// Upload ██████████████████████░░░░░░░░░░  73%
```

## Estilizacao

Todos os graficos aceitam estilos lipgloss, integrando perfeitamente com o tema do seu app:

```go
import themes "github.com/junhinhow/charm-themes"

theme := themes.Dracula
spark := charmcharts.NewSparkline(data).
    Style(lipgloss.NewStyle().Foreground(theme.Primary))
```

## Licenca

MIT
