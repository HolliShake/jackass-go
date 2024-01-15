package main

import (
	"fmt"
	"os"
)

type jackass_t struct {
	args []string
}

func JackAss() *jackass_t {
	jackass := new(jackass_t)
	jackass.args = os.Args[1:]
	return jackass
}

func (j *jackass_t) execute(filePath string) {
	// if len(j.args) <= 0 {
	// 	os.Exit(0)
	// }

	filePath = AbsolutePath(filePath)
	content := ""

	if !IsFileReadable(filePath) {
		basicError(fmt.Sprintf("file \"%s\" is not readable...", filePath))
	}

	content = ReadFile(filePath)

	lexer := Lexer(filePath, content)
	parser := Parser(lexer)
	analyzer := Analyzer(parser)
	fmt.Println(analyzer.analyze())
}
