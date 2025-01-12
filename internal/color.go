package internal

import "fmt" // an array of colors to use for the logs, we use the same color for the same task

// https://github.com/gawin/bash-colors-256
// not too dark or light
var colors = []int{
	34, 35, 36, 37, 38, 39,
	70, 71, 72, 73, 74, 75,
	106, 107, 108, 109, 110, 111,
	142, 143, 144, 145, 146, 147,
	178, 179, 180, 181, 182, 183,
	214, 215, 216, 217, 218, 219,
}

func color(x string) string {
	return fmt.Sprintf("\x1b[38;5;%dm", code(x))
}

func code(x string) int {
	code := 0
	for _, x := range x {
		code += int(x)
	}
	return colors[code%len(colors)]
}
