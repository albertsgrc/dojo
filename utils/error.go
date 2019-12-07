package utils

import (
	"fmt"

	"github.com/jedib0t/go-pretty/text"
)

// Error ...
func Error(message string) {
	fmt.Println(text.Colors{text.BgBlack, text.FgRed, text.Bold}.Sprint(
		"💥 ", message))

}
