package main

type TokenKind = int

const (
	TKIND_ID TokenKind = 0
	TKIND_KEYWORD TokenKind = 1
	TKIND_INTEGER TokenKind = 2
	TKIND_OTHER_INTEGER TokenKind = 3
	TKIND_BIG_INTEGER TokenKind = 4
	TKIND_OTHER_BIG_INTEGER TokenKind = 5
	TKIND_FLOAT TokenKind = 6
	TKIND_OTHER_FLOAT TokenKind = 7
	TKIND_STRING TokenKind = 8
	TKIND_SYMBOL TokenKind = 9
	TKIND_EOF TokenKind = 10
)

type token_t struct {
	kind TokenKind
	value string
	position *position_t
}

func Token(kind TokenKind, value string, position *position_t) *token_t {
	tok := new(token_t)

	//lint:ignore SA4031 possible not nil
	if tok == nil {
		basicError("Out of memory!!!")
	}

	tok.kind = kind
	tok.value = value
	tok.position = position
	return tok
}