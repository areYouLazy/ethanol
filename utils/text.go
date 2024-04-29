package utils

import (
	"fmt"
)

// Bold returns a *Bold* version of a given string
func Bold(s string) string {
	bs := fmt.Sprintf("\033[1m%s\033[0m", s)

	return bs
}

// Dim returns a *Dim* version of a given string
func Dim(s string) string {
	bs := fmt.Sprintf("\033[2m%s\033[0m", s)

	return bs
}

// Italic returns a *Italic* version of a given string
func Italic(s string) string {
	bs := fmt.Sprintf("\033[3m%s\033[0m", s)

	return bs
}

// Underlined returns a *Underlined* version of a given string
func Underlined(s string) string {
	bs := fmt.Sprintf("\033[4m%s\033[0m", s)

	return bs
}

// Blink returns a *Blink* version of a given string
func Blink(s string) string {
	bs := fmt.Sprintf("\033[5m%s\033[0m", s)

	return bs
}

// Reverse returns a *Reverse* version of a given string
func Reverse(s string) string {
	bs := fmt.Sprintf("\033[7m%s\033[0m", s)

	return bs
}

// Invisible returns a *Bold* version of a given string
func Invisible(s string) string {
	bs := fmt.Sprintf("\033[8m%s\033[0m", s)

	return bs
}
