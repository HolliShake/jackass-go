package main

var KEYWORDS []string = []string{
	"import",
	"from",
	"class",
	"enum",
	"function",
	"continue",
	"break",
	"return",
	"if",
	"else",
	"switch",
	"case",
	"default",
	"for",
	"do",
	"while",
	"goto",
	"label",
	"var",
	"let",
	"const",
	"true",
	"false",
	"null",
	"self",
}

func isKeyword(value string) bool {
	for _, keyword := range KEYWORDS {
		if keyword == value {
			return true
		}
	}

	return false
}
