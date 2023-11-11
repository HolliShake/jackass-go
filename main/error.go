package main

import (
	"fmt"
	"runtime"
	"strings"
	"os"
)

func basicError(message string) {
	fmt.Fprintf(os.Stderr, "JackAss::Error: %s\n", message)
	os.Exit(0xffffffff)
}

func raiseError[T istep](step T, message string, position *position_t) {
	var error string = ""
	var padding, start, end int = 3, 0, 0
	content := step.getFileCode()
	if runtime.GOOS == "windows" {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}
	
	lines := strings.Split(content, "\n")
	if ((position.lstart - 1) - padding) >= 0 {
		start = (position.lstart - 1) - padding
	} else {
		start = 0
	}

	if (position.lend + padding) < len(lines) {
		end = position.lend + padding
	} else {
		end = len(lines)
	}

	error += fmt.Sprintf("Error[%s:%d:%d]: %s\n", step.getFilePath(), position.lstart, position.cstart, message)
	
	for i := start; i < end; i++ {
		lineno := fmt.Sprintf("%d", (i + 1))
		space := ""

		for j := 0; j <  len(fmt.Sprintf("%d", end)) - len(lineno); j++ {
			space += " "
		}

		indicator := ""

		if (i + 1) == position.lstart && position.lstart == position.lend {
			// single line
			indicator = ">"
		} else if ((i + 1) >= position.lstart && (i + 1) <= position.lend) && position.lstart != position.lend {
			indicator = "~"
		}

		error += fmt.Sprintf("%s%s | %s %s\n",  space, lineno, indicator, lines[i])
	}

	fmt.Fprintf(os.Stderr, "%s", error)
	os.Exit(0xffffffff)
}

