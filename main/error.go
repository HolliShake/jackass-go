package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"jackass/shared"
)

func basicError(message string) {
	fmt.Fprintf(os.Stderr, "JackAss::Error: %s\n", message)
	os.Exit(0x1)
}

func raiseError[T istep](step T, message string, position *shared.Position_t) {
	var error string = ""
	var padding, start, end int = 3, 0, 0
	content := step.getFileCode()
	if runtime.GOOS == "windows" {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}

	lines := strings.Split(content, "\n")
	if ((position.Lstart - 1) - padding) >= 0 {
		start = (position.Lstart - 1) - padding
	} else {
		start = 0
	}

	if (position.Lend + padding) < len(lines) {
		end = position.Lend + padding
	} else {
		end = len(lines)
	}

	error += fmt.Sprintf("Error[%s:%d:%d]: %s\n", step.getFilePath(), position.Lstart, position.Cstart, message)

	for i := start; i < end; i++ {
		lineno := fmt.Sprintf("%d", (i + 1))
		space := ""

		for j := 0; j < len(fmt.Sprintf("%d", end))-len(lineno); j++ {
			space += " "
		}

		indicator := ""

		if (i+1) == position.Lstart && position.Lstart == position.Lend {
			// single line
			indicator = ">"
		} else if ((i+1) >= position.Lstart && (i+1) <= position.Lend) && position.Lstart != position.Lend {
			indicator = "~"
		}

		error += fmt.Sprintf("%s%s | %s %s\n", space, lineno, indicator, lines[i])
	}

	fmt.Fprintf(os.Stderr, "%s", error)
	os.Exit(0x1)
}
