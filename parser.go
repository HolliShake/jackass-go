package main
import (
	"fmt"
)

const (
	MAX_EXPRESSION_NESTING int = 100
)

type parser_t struct {
	lexer lexer_t
	lookahead, previous *token_t
	expressionLevel int
}

func Parser(lexer lexer_t) parser_t {
	return parser_t {
		lexer: lexer,
		expressionLevel: 0,
	}
}

func (p *parser_t) getFilePath() string {
	return p.lexer.getFilePath()
}

func (p *parser_t) getFileCode() string {
	return p.lexer.getFileCode()
}

// checkT checks lookahead type
// @param tokenKind TokenKind
// @returns bool
func (p *parser_t) checkT(tokenKind TokenKind) bool {
	return p.lookahead.kind == tokenKind
}

// checkV checks lookahead value
// @param value string
// @returns bool
func (p *parser_t) checkV(value string) bool {
	return p.lookahead.value == value
}

// checkS checks lookahead symbol
// @param symbol string
// @returns bool
func (p *parser_t) checkS(symbol string) bool {
	return p.checkT(TKIND_SYMBOL) && p.checkV(symbol)
}

// checkK checks lookahead keyword
// @param keyword string
// @returns bool
func (p *parser_t) checkK(keyword string) bool {
	return p.checkT(TKIND_KEYWORD) && p.checkV(keyword)
}

// acceptT accepts lookahead type
// @param tokenKind TokenKind
func (p *parser_t) acceptT(tokenKind TokenKind) {
	if p.checkT(tokenKind) {
		p.previous = p.lookahead
		p.lookahead = p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\".", p.lookahead.value), p.lookahead.position)
	}
}

// acceptV accepts lookahead value
// @param value string
func (p *parser_t) acceptV(value string) {
	if p.checkV(value) {
		p.previous = p.lookahead
		p.lookahead = p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, value), p.lookahead.position)
	}
}

// acceptS accepts lookahead symbol
// @param symbol string
func (p *parser_t) acceptS(symbol string) {
	if p.checkS(symbol) {
		p.previous = p.lookahead
		p.lookahead = p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, symbol), p.lookahead.position)
	}
}

// acceptK accepts lookahead keyword
// @param keyword string
func (p *parser_t) acceptK(keyword string) {
	if p.checkK(keyword) {
		p.previous = p.lookahead
		p.lookahead = p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, keyword), p.lookahead.position)
	}
}

func (p *parser_t) parseTerminal() *node_t {
	return nil
}

func (p *parser_t) parseMul() *node_t {
	node := p.parseTerminal()

	if (node == nil) { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("*") || p.checkS("/") || p.checkS("%") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseTerminal()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAdd() *node_t {
	node := p.parseMul()

	if (node == nil) { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("+") || p.checkS("-") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseMul()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseShift() *node_t {
	node := p.parseAdd()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("<<") || p.checkS(">>") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseAdd()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseRel() *node_t {
	node := p.parseShift()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("<") || p.checkS("<=") || p.checkS(">") || p.checkS(">=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseShift()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseEql() *node_t {
	node := p.parseRel()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("==") || p.checkS("!=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseRel()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseBit() *node_t {
	node := p.parseEql()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("&") || p.checkS("|") || p.checkS("^") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseEql()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseLog() *node_t {
	node := p.parseBit()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("&&") || p.checkS("||") || p.checkS("??") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseBit()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAss() *node_t {
	node := p.parseLog()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("="){

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseLog()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAug() *node_t {
	node := p.parseAss()

	if node == nil { 
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("*=") || p.checkS("/=") || p.checkS("%=") || p.checkS("+=") || p.checkS("-=") || p.checkS("<<=") || p.checkS(">>=") || p.checkS("&=") || p.checkS("|=") || p.checkS("^=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, fmt.Sprintf("expression nesting too deep."), p.lookahead.position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseAss()

		if rhs == nil { 
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.position)
		}

		node = BinaryExpressionNode(NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseZeroOrOneExpression() *node_t {
	return p.parseAug()
}

func (p *parser_t) parseMandatoryExpression() *node_t {
	node := p.parseZeroOrOneExpression()
	if node != nil {
		return node
	}

	// error
	raiseError(p, fmt.Sprintf("an expression is required, got \"%s\".", p.lookahead.value), p.lookahead.position)
	return nil
}

// 

func (p *parser_t) parseFile() *node_t {
	return nil
}

func (p *parser_t) parse() *node_t {
	p.lookahead = p.lexer.nextToken()
	p.previous = p.lookahead

	return p.parseFile()
}
