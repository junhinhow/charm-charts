package charmcharts

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// clamp restringe v ao intervalo [min, max].
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// displayWidth devolve a largura de exibição de uma string (em colunas),
// ignorando códigos ANSI. Nunca use len() para isso: len() conta bytes.
func displayWidth(s string) int {
	return lipgloss.Width(s)
}

// padRight preenche s com espaços à direita até atingir a largura desejada.
func padRight(s string, width int) string {
	current := displayWidth(s)
	if current >= width {
		return s
	}
	return s + strings.Repeat(" ", width-current)
}

// padLeft preenche s com espaços à esquerda até atingir a largura desejada.
func padLeft(s string, width int) string {
	current := displayWidth(s)
	if current >= width {
		return s
	}
	return strings.Repeat(" ", width-current) + s
}

// center centraliza s dentro da largura desejada.
func center(s string, width int) string {
	current := displayWidth(s)
	if current >= width {
		return s
	}
	left := (width - current) / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", width-current-left)
}

// fitLabel corta o texto para caber na largura desejada. O corte respeita os
// caracteres UTF-8 (nunca corta no meio de um rune) e usa reticências
// quando há espaço para elas.
func fitLabel(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if displayWidth(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return truncateRunes(s, width-1) + "…"
}

// truncateRunes remove caracteres do final até o texto caber na largura
// desejada, preservando a integridade do UTF-8.
func truncateRunes(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	for len(runes) > 0 && displayWidth(string(runes)) > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes)
}
