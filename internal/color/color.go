package color

import "fmt"

const ( // color
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
	hyperStart  = "\033]8;;" // OS commnand start + ;(seperate)
	hyperEnd    = "\033\\"   // ESC
)

func Red(msg string) string {
	return colorRed + msg + colorReset
}
func Green(msg string) string {
	return colorGreen + msg + colorReset
}
func Blue(msg string) string {
	return colorBlue + msg + colorReset
}
func Yellow(msg string) string {
	return colorYellow + msg + colorReset
}
func Hyperlink(url, msg string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\\n", url, msg)
}
