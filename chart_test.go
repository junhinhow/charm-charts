package charmcharts

import (
	"strings"
	"testing"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
)

// plainBarChart devolve um gráfico de barras sem códigos ANSI, para que as
// posições dos caracteres possam ser medidas diretamente.
func plainBarChart(labels []string, values []float64) *BarChart {
	chart := NewBarChart(labels, values)
	chart.Style(lipgloss.NewStyle())
	chart.BarStyle(lipgloss.NewStyle())
	chart.LabelStyle(lipgloss.NewStyle())
	return chart
}

// TestSparklineWidthInvalidoNaoPanica garante que larguras inválidas não
// causam panic (o código anterior entrava em panic com Width(0) e
// Width(-1)).
func TestSparklineWidthInvalidoNaoPanica(t *testing.T) {
	for _, largura := range []int{0, -1, -100} {
		t.Run(larguraLabel(largura), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Width(%d) causou panic: %v", largura, r)
				}
			}()
			resultado := NewSparkline([]float64{1, 2, 3}).Width(largura).Render()
			if largura < minSparklineWidth {
				largura = minSparklineWidth
			}
			if w := lipgloss.Width(resultado); w != largura {
				t.Errorf("esperava largura %d, obteve %d", largura, w)
			}
		})
	}
}

// TestSparklineHeightInvalidoNaoPanica garante que a altura é limitada a 1..8.
func TestSparklineHeightInvalidoNaoPanica(t *testing.T) {
	casos := []struct {
		altura   int
		esperado int
	}{
		{0, 1},
		{-5, 1},
		{1, 1},
		{3, 3},
		{8, 8},
		{100, 8},
	}
	for _, c := range casos {
		t.Run(larguraLabel(c.altura), func(t *testing.T) {
			resultado := NewSparkline([]float64{1, 2, 3}).Height(c.altura).Render()
			linhas := strings.Split(resultado, "\n")
			if len(linhas) != c.esperado {
				t.Errorf("Height(%d): esperava %d linhas, obteve %d", c.altura, c.esperado, len(linhas))
			}
		})
	}
}

// TestSparklineLarguraRespeitada garante que o gráfico ocupa exatamente a
// largura pedida, mesmo quando há menos pontos que colunas.
func TestSparklineLarguraRespeitada(t *testing.T) {
	for _, largura := range []int{1, 5, 10, 30} {
		t.Run(larguraLabel(largura), func(t *testing.T) {
			resultado := NewSparkline([]float64{1, 2, 3}).Width(largura).Render()
			if w := lipgloss.Width(resultado); w != largura {
				t.Errorf("esperava largura %d, obteve %d", largura, w)
			}
		})
	}
}

func TestSparklineSemDados(t *testing.T) {
	if resultado := NewSparkline(nil).Render(); resultado != "" {
		t.Errorf("esperava string vazia, obteve %q", resultado)
	}
}

// TestSparklineRotuloConsistente verifica que o rótulo aparece sempre na
// primeira linha do gráfico (inclusive com várias linhas) e que todas as
// linhas continuam com a mesma largura.
func TestSparklineRotuloConsistente(t *testing.T) {
	for _, altura := range []int{1, 2, 5} {
		t.Run(larguraLabel(altura), func(t *testing.T) {
			resultado := NewSparkline([]float64{1, 5, 3}).Width(10).Height(altura).Label("CPU").Render()
			linhas := strings.Split(resultado, "\n")
			if len(linhas) != altura {
				t.Fatalf("esperava %d linhas, obteve %d", altura, len(linhas))
			}
			if !strings.HasPrefix(linhas[0], "CPU ") {
				t.Errorf("esperava o rótulo no início da primeira linha, obteve %q", linhas[0])
			}
			for i, linha := range linhas {
				if !strings.HasPrefix(linha, "CPU ") && linha[:4] != "    " {
					t.Errorf("linha %d: esperava o mesmo recuo do rótulo, obteve %q", i, linha)
				}
				if w := lipgloss.Width(linha); w != 14 {
					t.Errorf("linha %d: esperava largura 14, obteve %d", i, w)
				}
			}
		})
	}
}

// TestResampleReduzEExpande cobre a redução (médias) e a expansão
// (interpolação linear) dos dados.
func TestResampleReduzEExpande(t *testing.T) {
	dados := []float64{1, 2, 3, 4}

	reduzido := resample(dados, 2)
	if len(reduzido) != 2 {
		t.Fatalf("esperava 2 pontos, obteve %d", len(reduzido))
	}
	if reduzido[0] != 1.5 || reduzido[1] != 3.5 {
		t.Errorf("esperava [1.5 3.5], obteve %v", reduzido)
	}

	expandido := resample(dados, 7)
	if len(expandido) != 7 {
		t.Fatalf("esperava 7 pontos, obteve %d", len(expandido))
	}
	if expandido[0] != 1 || expandido[6] != 4 {
		t.Errorf("esperava que a expansão preservasse os extremos, obteve %v", expandido)
	}

	if r := resample(nil, 10); r != nil {
		t.Errorf("esperava nil para dados vazios, obteve %v", r)
	}
	if r := resample(dados, 0); r != nil {
		t.Errorf("esperava nil para largura 0, obteve %v", r)
	}
	if r := resample([]float64{7}, 4); len(r) != 4 || r[0] != 7 || r[3] != 7 {
		t.Errorf("esperava 4 pontos com o valor 7, obteve %v", r)
	}
}

// TestProgressBarWidthInvalidoNaoPanica garante que a barra não entra em
// panic com largura negativa.
func TestProgressBarWidthInvalidoNaoPanica(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Width(-1) causou panic: %v", r)
		}
	}()
	resultado := NewProgressBar(0.5).Width(-1).Render()
	if resultado == "" {
		t.Error("esperava uma saída não vazia")
	}
}

// TestProgressBarLarguraRespeitada garante que a barra tem a largura pedida.
func TestProgressBarLarguraRespeitada(t *testing.T) {
	barra := NewProgressBar(0.5)
	barra.FillStyle(lipgloss.NewStyle())
	barra.EmptyStyle(lipgloss.NewStyle())
	resultado := barra.Width(10).Render()
	if !strings.Contains(resultado, strings.Repeat("█", 5)) {
		t.Errorf("esperava 5 caracteres preenchidos, obteve %q", resultado)
	}
	if !strings.Contains(resultado, strings.Repeat("░", 5)) {
		t.Errorf("esperava 5 caracteres vazios, obteve %q", resultado)
	}
}

func TestProgressBarValorClampado(t *testing.T) {
	if got := NewProgressBar(2).Render(); !strings.Contains(got, "100%") {
		t.Errorf("esperava 100%%, obteve %q", got)
	}
	if got := NewProgressBar(-1).Render(); !strings.Contains(got, "  0%") {
		t.Errorf("esperava 0%%, obteve %q", got)
	}
}

// TestBarChartBarWidthInvalidoNaoPanica garante que BarWidth não causa panic.
func TestBarChartBarWidthInvalidoNaoPanica(t *testing.T) {
	for _, largura := range []int{-1, -50, 0} {
		t.Run(larguraLabel(largura), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("BarWidth(%d) causou panic: %v", largura, r)
				}
			}()
			gráfico := NewBarChart([]string{"A"}, []float64{1}).BarWidth(largura)
			if resultado := gráfico.Render(); resultado == "" {
				t.Error("esperava uma saída não vazia")
			}
			gráfico.Vertical()
			if resultado := gráfico.Render(); resultado == "" {
				t.Error("esperava uma saída não vazia no modo vertical")
			}
		})
	}
}

// TestBarChartWidthInvalidoNaoPanica garante que Width não causa panic.
func TestBarChartWidthInvalidoNaoPanica(t *testing.T) {
	for _, largura := range []int{-1, 0} {
		t.Run(larguraLabel(largura), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Width(%d) causou panic: %v", largura, r)
				}
			}()
			gráfico := NewBarChart([]string{"A"}, []float64{1}).Width(largura)
			if resultado := gráfico.Render(); resultado == "" {
				t.Error("esperava uma saída não vazia")
			}
		})
	}
}

// TestBarChartHorizontalEspessuraDaBarra garante que BarWidth aumenta a
// espessura da barra (e não duplica a linha inteira).
func TestBarChartHorizontalEspessuraDaBarra(t *testing.T) {
	rótulos := []string{"Go", "Rust"}
	gráfico := plainBarChart(rótulos, []float64{100, 50}).Width(40).BarWidth(3)
	linhas := strings.Split(gráfico.Render(), "\n")

	if len(linhas) != len(rótulos)*3 {
		t.Fatalf("esperava %d linhas, obteve %d", len(rótulos)*3, len(linhas))
	}

	// A primeira linha do grupo tem o rótulo; as demais não.
	if !strings.Contains(linhas[0], "Go") {
		t.Errorf("esperava o rótulo na primeira linha, obteve %q", linhas[0])
	}
	if strings.Contains(linhas[1], "Go") {
		t.Errorf("esperava a segunda linha sem rótulo, obteve %q", linhas[1])
	}
	if !strings.Contains(linhas[1], "█") {
		t.Errorf("esperava a barra repetida na segunda linha, obteve %q", linhas[1])
	}
	if got := strings.Count(linhas[0], "█"); got != strings.Count(linhas[1], "█") {
		t.Errorf("a barra deveria ter a mesma espessura nas linhas do grupo: %d e %d",
			got, strings.Count(linhas[1], "█"))
	}
}

// TestBarChartHorizontalLinhasMesmaLargura garante que todas as linhas do
// gráfico horizontal têm a mesma largura de exibição.
func TestBarChartHorizontalLinhasMesmaLargura(t *testing.T) {
	gráfico := plainBarChart([]string{"Go", "Rust", "Python", "JS"}, []float64{85, 72, 95, 68}).Width(50)
	linhas := strings.Split(gráfico.Render(), "\n")
	esperado := lipgloss.Width(linhas[0])
	for i, linha := range linhas {
		if w := lipgloss.Width(linha); w != esperado {
			t.Errorf("linha %d: esperava largura %d, obteve %d", i, esperado, w)
		}
	}
}

// TestBarChartRotulosAlinhados verifica que os rótulos terminam na mesma
// coluna (alinhados à direita) no modo horizontal.
func TestBarChartRotulosAlinhados(t *testing.T) {
	gráfico := plainBarChart([]string{"Go", "Rust", "Python", "JS"}, []float64{85, 72, 95, 68}).Width(50)
	linhas := strings.Split(gráfico.Render(), "\n")
	for i, linha := range linhas {
		rótulo := []string{"Go", "Rust", "Python", "JS"}[i]
		índice := strings.Index(linha, rótulo)
		if índice < 0 {
			t.Fatalf("linha %d: rótulo %q não encontrado em %q", i, rótulo, linha)
		}
		fim := índice + len([]rune(rótulo))
		if fim != 6 {
			t.Errorf("linha %d: esperava o rótulo terminando na coluna 6, obteve %d", i, fim)
		}
	}
}

// TestBarChartVerticalRespeitaWidth garante que o modo vertical respeita a
// largura pedida.
func TestBarChartVerticalRespeitaWidth(t *testing.T) {
	for _, largura := range []int{20, 30, 50, 120} {
		t.Run(larguraLabel(largura), func(t *testing.T) {
			gráfico := plainBarChart([]string{"Go", "Rust", "Python", "JS"}, []float64{85, 72, 95, 68}).
				Width(largura).Vertical().BarWidth(3)
			linhas := strings.Split(gráfico.Render(), "\n")
			for i, linha := range linhas {
				if w := lipgloss.Width(linha); w > largura {
					t.Errorf("linha %d: largura %d excede o limite pedido (%d)", i, w, largura)
				}
			}
		})
	}
}

// TestBarChartVerticalLarguraEscalaComWidth garante que a largura pedida
// realmente influence o gráfico vertical (o código anterior ignorava o
// Width() e devolvia sempre a mesma largura).
func TestBarChartVerticalLarguraEscalaComWidth(t *testing.T) {
	rótulos := []string{"Go", "Rust", "Python", "JS"}
	valores := []float64{85, 72, 95, 68}

	estreito := plainBarChart(rótulos, valores).Width(20).Vertical().BarWidth(3).Render()
	largo := plainBarChart(rótulos, valores).Width(60).Vertical().BarWidth(3).Render()

	if lipgloss.Width(largo) <= lipgloss.Width(estreito) {
		t.Errorf("o gráfico vertical deveria acompanhar a largura pedida: estreito=%d largo=%d",
			lipgloss.Width(estreito), lipgloss.Width(largo))
	}
}

// TestBarChartHorizontalRespeitaWidth garante que o gráfico horizontal não
// ultrapassa a largura pedida.
func TestBarChartHorizontalRespeitaWidth(t *testing.T) {
	gráfico := plainBarChart([]string{"Programação", "Go"}, []float64{85, 72}).Width(20)
	for i, linha := range strings.Split(gráfico.Render(), "\n") {
		if w := lipgloss.Width(linha); w > 20 {
			t.Errorf("linha %d: largura %d excede o limite pedido (20)", i, w)
		}
	}
}

// TestBarChartVerticalMostraValores verifica que ShowValues é respeitado no
// modo vertical.
func TestBarChartVerticalMostraValores(t *testing.T) {
	comValores := plainBarChart([]string{"Go", "Rust"}, []float64{85, 72}).Vertical().ShowValues(true).Render()
	linhas := strings.Split(comValores, "\n")
	if !strings.Contains(linhas[0], "85.0") || !strings.Contains(linhas[0], "72.0") {
		t.Errorf("esperava os valores na primeira linha, obteve %q", linhas[0])
	}

	semValores := plainBarChart([]string{"Go", "Rust"}, []float64{85, 72}).Vertical().ShowValues(false).Render()
	if strings.Contains(semValores, "85.0") {
		t.Errorf("esperava a saída sem valores, obteve %q", semValores)
	}
}

// TestBarChartVerticalRotulosCentralizados verifica que os rótulos do modo
// vertical não ficam colados nos rótulos vizinhos.
func TestBarChartVerticalRotulosCentralizados(t *testing.T) {
	gráfico := plainBarChart([]string{"Go", "Rust", "Python", "JS"}, []float64{85, 72, 95, 68}).
		Width(50).Vertical().BarWidth(3)
	linhas := strings.Split(gráfico.Render(), "\n")
	rótulos := linhas[len(linhas)-1]
	if !strings.Contains(rótulos, " Go ") {
		t.Errorf("esperava o rótulo centralizado na coluna, obteve %q", rótulos)
	}
	if !strings.Contains(rótulos, "Rust") {
		t.Errorf("esperava o rótulo Rust, obteve %q", rótulos)
	}
}

// TestBarChartUTF8Valido garante que a saída é sempre um UTF-8 válido,
// mesmo com rótulos multibyte truncados.
func TestBarChartUTF8Valido(t *testing.T) {
	rótulos := []string{"Ação", "日本語", "Go", "Programação"}
	gráfico := plainBarChart(rótulos, []float64{10, 20, 30, 40})
	for _, saída := range []string{gráfico.Render(), gráfico.Vertical().Render()} {
		if !utf8.ValidString(saída) {
			t.Errorf("saída inválida em UTF-8: %q", saída)
		}
	}
}

// TestBarChartRotuloLongoNaoQuebraLargura garante que rótulos longos são
// truncados sem ultrapassar a coluna reserved.
func TestBarChartRotuloLongoNaoQuebraLargura(t *testing.T) {
	gráfico := plainBarChart([]string{"Programação"}, []float64{10}).Width(20).Vertical()
	linhas := strings.Split(gráfico.Render(), "\n")
	for i, linha := range linhas {
		if w := lipgloss.Width(linha); w > 20 {
			t.Errorf("linha %d: largura %d excede o limite pedido (20)", i, w)
		}
	}
	if !utf8.ValidString(strings.Join(linhas, "\n")) {
		t.Error("saída inválida em UTF-8")
	}
}

// TestFitLabelNaoQuebraRune garante que o corte de rótulos respeita o UTF-8.
func TestFitLabelNaoQuebraRune(t *testing.T) {
	casos := []struct {
		texto   string
		largura int
	}{
		{"Programação", 5},
		{"日本語", 4},
		{"Ação", 3},
	}
	for _, c := range casos {
		resultado := fitLabel(c.texto, c.largura)
		if !utf8.ValidString(resultado) {
			t.Errorf("fitLabel(%q, %d) devolveu UTF-8 inválido: %q", c.texto, c.largura, resultado)
		}
		if w := lipgloss.Width(resultado); w > c.largura {
			t.Errorf("fitLabel(%q, %d) devolveu largura %d", c.texto, c.largura, w)
		}
	}
	if fitLabel("Go", 10) != "Go" {
		t.Errorf("fitLabel não deveria alterar textos que já cabem")
	}
}

func TestBarChartSemValores(t *testing.T) {
	if resultado := NewBarChart(nil, nil).Render(); resultado != "" {
		t.Errorf("esperava string vazia, obteve %q", resultado)
	}
}

// larguraLabel cria um nome de subteste legível.
func larguraLabel(v int) string {
	if v < 0 {
		return "negativo"
	}
	return "valor"
}
