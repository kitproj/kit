package internal

import "fmt" // an array of colors to use for the logs, we use the same color for the same task

// we avoid red as this is for errors, and black as it is the default color
var codes = []int{
	32, // green
	33, // yellow
	34, // blue
	35, // magenta
	36, // cyan
}

func Color(x string, bold bool) string {

	code := 0
	for _, x := range x {
		code += int(x)
	}

	code = codes[code%len(codes)]

	if bold {
		code += 1
	}

	return fmt.Sprintf("\033[0;%dm", code)
}
