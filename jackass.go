package main
import "os"

type jackass_t struct {
	args []string
}


func JackAss() {
	jackass := new(jackass_t)
	jackass.args = os.Args[1:]
}

func (j *jackass_t) execute(filePath string) {
	content := ""
	

	lexer := Lexer(filePath, content)
	parser := Parser(lexer)

	parser.parse()
}
