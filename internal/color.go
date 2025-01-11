package internal

import "fmt" // an array of colors to use for the logs, we use the same color for the same task

func Color(x string) string {

	code := 0
	for _, x := range x {
		code += int(x)
	}

	code = 16 + code%216

	return fmt.Sprintf("\x1b[38;5;%dm", code)
}
