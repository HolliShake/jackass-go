package main
import (
	"fmt"
	"strings"
)


const MAX_ID_LENGTH int = 255

type lexer_t struct {
	filePath, fileContent string
	lookahead rune
	index, line, column int
}

func Lexer(filePath, fileContent string) *lexer_t {
	lexer := new(lexer_t)
	lexer.filePath = filePath
	lexer.fileContent = fileContent
	lexer.lookahead = lexer.nextRune()
	lexer.index = 0
	lexer.line = 1
	lexer.column = 1
	return lexer
}

func (l *lexer_t) getFilePath() string {
	return l.filePath
}

func (l *lexer_t) getFileCode() string {
	return l.fileContent
}

func (l *lexer_t) nextRune() rune {
	size := utf_sizeOfUtf(int(l.fileContent[l.index]))

	if (l.index + size) > len(l.fileContent) {
		return rune(l.fileContent[l.index])
	}

	var ord rune = 0

	switch size {
		case 1:
			return rune(l.fileContent[l.index])
		case 2:
			ord = utf_toCodePoint(
				int(l.fileContent[l.index + 0]),
				int(l.fileContent[l.index + 1]),
				0,
				0,
			)
			break
		case 3:
			ord = utf_toCodePoint(
				int(l.fileContent[l.index + 0]),
				int(l.fileContent[l.index + 1]),
				int(l.fileContent[l.index + 2]),
				0,
			)
			break
		case 4:
			ord = utf_toCodePoint(
				int(l.fileContent[l.index + 0]),
				int(l.fileContent[l.index + 1]),
				int(l.fileContent[l.index + 2]),
				int(l.fileContent[l.index + 3]),
			)
			break
	}

	l.index += (size - 1) // -1 because array index starts at 0

	return ord
}

func (l *lexer_t) forward() {
	if (l.lookahead == '\n') {
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.index++
	if l.index < len(l.fileContent) {
		l.lookahead = l.nextRune()
	} else {
		l.lookahead = rune('\x00')
	}
}

func (l *lexer_t) isWhiteSpace() bool {
	return utf_isWhiteSpace(l.lookahead)
}

func (l *lexer_t) isIdentifierStart() bool {
	return utf_isLetter(l.lookahead)
}

func (l *lexer_t) isIdentifierPart() bool {
	return utf_isLetterOrDigit(l.lookahead)
}

func (l *lexer_t) isDigit() bool {
	return utf_isDigit(l.lookahead)
}

func (l *lexer_t) isHex() bool {
	return (
		utf_isDigit(l.lookahead) ||
		(l.lookahead >= 'a' && l.lookahead <= 'f') ||
		(l.lookahead >= 'A' && l.lookahead <= 'F'))
}

func (l *lexer_t) isBin() bool {
	return (l.lookahead == '0' || l.lookahead == '1')
}

func (l *lexer_t) isOct() bool {
	return (l.lookahead >= '0' && l.lookahead <= '7')
}

func (l *lexer_t) isString() bool {
	return (l.lookahead == '"' || l.lookahead == '\'')
}

func (l *lexer_t) isEof() bool {
	return l.index >= len(l.fileContent)
}

func (l *lexer_t) skipWhiteSpace() {
	for l.isWhiteSpace() {
		l.forward()
	}
}

func (l *lexer_t) nextIdentifier() *token_t {
	var value string = ""
	var pos *position_t = Position(l.line, l.column)

	for l.isIdentifierStart() {
		value += string(l.lookahead)
		l.forward()
	}

	for l.isIdentifierPart() {
		value += string(l.lookahead)
		l.forward()
	}

	if len(value) >= MAX_ID_LENGTH {
		raiseError(l, fmt.Sprintf("identifier \"%s\"(+%dmore)... is too long.", value[0: 30], len(value) - 30), pos)
	}

	ttype := TKIND_ID

	if isKeyword(value) {
		ttype = TKIND_KEYWORD
	}

	return Token(ttype, value, pos)
}

func (l *lexer_t) nextNumber() *token_t {
	var value string = ""
	var pos *position_t = Position(l.line, l.column)

	for l.isDigit() { 
		value += string(l.lookahead)
		l.forward()
	}

	if strings.Compare(value, "0") == 0 { 

		switch l.lookahead {
			case 'x', 'X':
				value += string(l.lookahead)
				l.forward()

				if !l.isHex() {
					// Error
					raiseError(l, fmt.Sprintf("incomplete hexadecimal number \"%s\".", value), pos)
				}

				for l.isHex() {
					value += string(l.lookahead)
					l.forward()
				}

				if l.lookahead == 'n' || l.lookahead == 'N' {
					l.forward()
					return Token(TKIND_OTHER_BIG_INTEGER, value, pos)
				} else {
					return Token(TKIND_OTHER_INTEGER, value, pos)
				}
			case 'b', 'B':
				value += string(l.lookahead)
				l.forward()

				if !l.isBin() {
					// Error
					raiseError(l, fmt.Sprintf("incomplete binary number \"%s\".", value), pos)
				}

				for l.isBin() {
					value += string(l.lookahead)
					l.forward()
				}

				if l.lookahead == 'n' || l.lookahead == 'N' {
					l.forward()
					return Token(TKIND_OTHER_BIG_INTEGER, value, pos)
				} else {
					return Token(TKIND_OTHER_INTEGER, value, pos)
				}
			case 'o', 'O':
				value += string(l.lookahead)
				l.forward()

				if !l.isOct() {
					// Error
					raiseError(l, fmt.Sprintf("incomplete octal number \"%s\".", value), pos)
				}

				for l.isOct() {
					value += string(l.lookahead)
					l.forward()
				}

				if l.lookahead == 'n' || l.lookahead == 'N' {
					l.forward()
					return Token(TKIND_OTHER_BIG_INTEGER, value, pos)
				} else {
					return Token(TKIND_OTHER_INTEGER, value, pos)
				}
		}
	}

	var ttype TokenKind = TKIND_INTEGER

	if l.lookahead == 'n' || l.lookahead == 'N' {
		l.forward()
		return Token(TKIND_BIG_INTEGER, value, pos)
	} 

	if l.lookahead == '.' { 
		value += string(l.lookahead)
		l.forward()

		if !l.isDigit() {
			// Error
			raiseError(l, fmt.Sprintf("invalid number format \"%s\".", value), pos)
		}

		for l.isDigit() {
			value += string(l.lookahead)
			l.forward()
		}

		ttype = TKIND_FLOAT
	}

	if l.lookahead == 'e' || l.lookahead == 'E' { 
		value += "e"
		l.forward()

		if l.lookahead == '+' || l.lookahead == '-' {
			value += string(l.lookahead)
			l.forward()
		}

		if !l.isDigit() {
			// Error
			raiseError(l, fmt.Sprintf("invalid number format \"%s\".", value), pos)
		}

		for l.isDigit() {
			value += string(l.lookahead)
			l.forward()
		}

		ttype = TKIND_OTHER_FLOAT
	}

	return Token(ttype, value, pos)
}

func (l * lexer_t) nextString() *token_t {
	var value string = ""
	var pos *position_t = Position(l.line, l.column)
	var isopen, isclose bool = l.isString(), false
	openner := l.lookahead

	l.forward()
	isclose = l.isString() && l.lookahead == openner

	for !l.isEof() && (isopen && !isclose) { 
		if l.lookahead == '\n' {
			break
		}

		if l.lookahead == '\\' {
			l.forward()

			switch l.lookahead { 
				case 'b':
					value += "\b"
					break
				case 'n':
					value += "\n"
					break
				case 't':
					value += "\t"
					break
				case 'r':
					value += "\r"
					break
				case 'f':
					value += "\f"
					break

				case '"':
					value += "\""
					break
				case '\'':
					value += "'"
					break
				default:
					value += string(l.lookahead)
			}
		} else {
			value += string(l.lookahead)
		}

		l.forward()
		if l.lookahead == openner {
			isclose = l.isString()
		}
	}

	if !(isopen && isclose) {
		// Error
		raiseError(l, "string literal was not closed or terminated properly.", pos)
	}

	l.forward()

	return Token(TKIND_STRING, value, pos)
}

func (l *lexer_t) nextSymbol() *token_t {
	var value string = ""
	var pos *position_t = Position(l.line, l.column)

	switch l.lookahead {
		case '(':
		case ')':
		case '[':
		case ']':
		case '{':
		case '}':
		case '~':
		case ':':
		case ';':
		case ',':
			value += string(l.lookahead)
			l.forward()
			break
		case '?':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '?' { 
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '.':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '.' { 
				value += string(l.lookahead)
				l.forward()

				if l.lookahead == '.' { 
					value += string(l.lookahead)
					l.forward()
				} else {
					raiseError(l, fmt.Sprintf("invalid symbol \"%s\".", value), pos)
				}
			}
			break
		case '*':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '/':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '%':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '+':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '+' {
				value += string(l.lookahead)
				l.forward()
			} else if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '-':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '-' {
				value += string(l.lookahead)
				l.forward()
			} else if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '<':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '<' {
				value += string(l.lookahead)
				l.forward()
			} 
			
			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '>':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '>' {
				value += string(l.lookahead)
				l.forward()
			} 
			
			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '&':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '&' {
				value += string(l.lookahead)
				l.forward()
			} else if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '|':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '|' {
				value += string(l.lookahead)
				l.forward()
			} else if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '^':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '=':
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		case '!': 
			value += string(l.lookahead)
			l.forward()

			if l.lookahead == '=' {
				value += string(l.lookahead)
				l.forward()
			}
			break
		default:
			for !l.isEof() && !l.isWhiteSpace() {
				value += string(l.lookahead)
				l.forward()
			}
			raiseError(l, fmt.Sprintf("invalid symbol \"%s\".", value), pos)
	}

	return Token(TKIND_SYMBOL, value, pos)
}

func (l *lexer_t) nextEof() *token_t {
	return Token(TKIND_EOF, "", Position(l.line, l.column))
}

func (l *lexer_t) nextToken() *token_t {
	for !l.isEof() {
		if l.isWhiteSpace() {
			l.skipWhiteSpace()
		} else if l.isIdentifierStart() {
			return l.nextIdentifier()
		} else if l.isDigit() {
			return l.nextNumber()
		} else if l.isString() {
			return l.nextString()
		} else {
			return l.nextSymbol()
		}
	}
	return l.nextEof()
}