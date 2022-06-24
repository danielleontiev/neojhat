package printing

type Color string

const (
	cReset  Color = "\033[0m"
	cBold   Color = "\033[1m"
	cRed    Color = "\033[31m"
	cGreen  Color = "\033[32m"
	cYellow Color = "\033[33m"
	cBlue   Color = "\033[34m"
	cPurple Color = "\033[35m"
	cCyan   Color = "\033[36m"
	cGray   Color = "\033[37m"
	cWhite  Color = "\033[97m"
)

// addColoe adds ANSI codes of given color
// to the string
func addColor(str string, col Color) string {
	return string(col) + str + string(cReset)
}

func Bold(str string) string {
	return addColor(str, cBold)
}

func Cyan(str string) string {
	return addColor(str, cCyan)
}

func Red(str string) string {
	return addColor(str, cRed)
}

func Yellow(str string) string {
	return addColor(str, cYellow)
}

func Blue(str string) string {
	return addColor(str, cBlue)
}
