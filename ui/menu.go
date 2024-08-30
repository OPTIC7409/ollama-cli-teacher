package ui

import (
	"fmt"
	"strings"
)

func PrintHeader(title string) {
	border := "════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════"
	borderLength := len(border)
	titleLength := len(title)

	totalPadding := (borderLength - titleLength) / 2
	leftPadding := (totalPadding / 3)
	rightPadding := (totalPadding-leftPadding)/3 - 3

	paddedTitle := strings.Repeat(" ", leftPadding) + title + strings.Repeat(" ", rightPadding)

	fmt.Printf("\n%s\n║ %s ║\n%s\n\n", border, paddedTitle, border)
}
