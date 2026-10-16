# charm-charts

Gráficos de terminal para o ecossistema [charmbracelet](https://github.com/charmbracelet).

Sparklines, gráficos de barra e barras de progresso — estilizados com lipgloss.

> **[Read in English](README.md)**

## Instalação

```bash
go get github.com/junhinhow/charm-charts
```

## Gráficos

### Sparkline

Gráfico de linha compacto usando caracteres Unicode de bloco. Perfeito para visualização inline de dados.

```go
data := []float64{4, 2, 7, 3, 8, 5, 9, 1, 6, 4}

spark := charmcharts.NewSparkline(data).
    Width(30).
    Label("CPU").
    Style(lipgloss.NewStyle().Foreground(lipgloss.Color("212")))

fmt.Println(spark.Render())
// CPU ▂▁▁ ▁▃▅▅▃▂▂▃▅▇▆▅▃▅▆█▆▃  ▁▃▄▄▃▂
```

O gráfico ocupa sempre exatamente a largura pedida: quando há menos pontos que
colunas a série é expandida por interpolação linear e quando há mais pontos eles
são reduzidos pela média de cada intervalo.

Sparklines multi-linha para maior resolução:
```go
spark.Height(3)  // 1-8 linhas
```

O rótulo é desenhado sempre na primeira linha do gráfico; as demais linhas
recebem o mesmo espaçamento para que as colunas permaneçam alinhadas.

### Gráfico de Barras

Gráficos de barra horizontais ou verticais com rótulos e valores.

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

`.BarWidth(n)` define a espessura de cada barra: no modo horizontal cada barra
ocupa `n` linhas (o rótulo e o valor aparecem na primeira linha do grupo) e no
modo vertical define a largura de cada coluna.

Modo vertical:
```go
chart.Vertical().BarWidth(3)
```

Os gráficos verticais respeitam `.Width(n)` (as colunas encolhem para caber) e
`.ShowValues(bool)` (os valores são impressos acima das colunas).

### Barra de Progresso

Barra de progresso horizontal com porcentagem e caracteres customizáveis.

```go
bar := charmcharts.NewProgressBar(0.73).
    Width(30).
    Label("Upload").
    FillStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("42")))

fmt.Println(bar.Render())
// Upload █████████████████████░░░░░░░░░  73%
```

## Estilização

Todos os gráficos aceitam estilos lipgloss, integrando perfeitamente com o tema do seu app:

```go
import (
    themes "github.com/junhinhow/charm-themes"

    "charm.land/lipgloss/v2"
)

theme := themes.Dracula
spark := charmcharts.NewSparkline(data).
    Style(lipgloss.NewStyle().Foreground(theme.Primary))
```

## Referência da API

### Sparkline
| Método | Descrição |
|--------|-----------|
| `NewSparkline(data)` | Cria com os pontos de dados |
| `.Width(n)` | Define a largura em caracteres (limitada a 1-500) |
| `.Height(n)` | Define a altura em linhas (limitada a 1-8) |
| `.Style(s)` | Define o estilo lipgloss |
| `.Label(s)` | Define o rótulo do prefixo |
| `.Render()` | Devolve a string estilizada |

### BarChart
| Método | Descrição |
|--------|-----------|
| `NewBarChart(labels, values)` | Cria com rótulos e valores |
| `.Width(n)` | Define a largura máxima (limitada a 1-500) |
| `.BarWidth(n)` | Define a espessura da barra (limitada a 1-20) |
| `.Horizontal()` / `.Vertical()` | Define a orientação |
| `.BarStyle(s)` / `.LabelStyle(s)` | Define os estilos |
| `.ShowValues(bool)` | Alterna a exibição dos valores |
| `.Render()` | Devolve a string estilizada |

### ProgressBar
| Método | Descrição |
|--------|-----------|
| `NewProgressBar(value)` | Cria com o valor (0.0-1.0) |
| `.Width(n)` | Define a largura da barra (limitada a 0-500) |
| `.FillStyle(s)` / `.EmptyStyle(s)` | Define os estilos |
| `.Label(s)` | Define o rótulo do prefixo |
| `.ShowPercentage(bool)` | Alterna a porcentagem |
| `.FillChar(r)` / `.EmptyChar(r)` | Caracteres customizados |
| `.Render()` | Devolve a string estilizada |

## Observações

- Larguras e alturas inválidas são ajustadas para o limite mais próximo, em vez
  de causar panic.
- Os rótulos são medidos pela largura de exibição (nunca por `len()`), de modo
  que letras acentuadas e caracteres CJK permanecem alinhados.
- Rótulos que não cabem na coluna são truncados sem quebrar o UTF-8.

## Licença

MIT
