package main
import (
	"fmt"
	"strings"
)

const (
	MAX_MEMBER_ACCESS int = 50
	MAX_INDEXING int = 50
	MAX_CALL int = 50
	MAX_EXPRESSION_NESTING int = 100
	// 
	MAX_ARGS int = 255
)

type parser_t struct {
	lexer *lexer_t
	lookahead, previous *token_t
	memberAccess, indexingLevel, callLevel, expressionLevel int
}

func Parser(lexer *lexer_t) *parser_t {
	parser := new(parser_t)

	if parser == nil {
		basicError("Out of memory!!!")
	}

	parser.lexer = lexer
	parser.memberAccess = 0
	parser.indexingLevel = 0
	parser.callLevel = 0
	parser.expressionLevel = 0
	return parser
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
	return strings.Compare(p.lookahead.value, value) == 0
}

// checkS checks lookahead symbol/
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
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\".", p.lookahead.value), p.lookahead.position)
	}
}

// acceptV accepts lookahead value
// @param value string
func (p *parser_t) acceptV(value string) {
	if p.checkV(value) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, value), p.lookahead.position)
	}
}

// acceptS accepts lookahead symbol
// @param symbol string
func (p *parser_t) acceptS(symbol string) {
	if p.checkS(symbol) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, symbol), p.lookahead.position)
	}
}

// acceptK accepts lookahead keyword
// @param keyword string
func (p *parser_t) acceptK(keyword string) {
	if p.checkK(keyword) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, keyword), p.lookahead.position)
	}
}

func (p *parser_t) parseTerminal() *node_t {
	switch p.lookahead.kind {
		case TKIND_ID:
			node := TerminalNode(
				NT_ID, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_KEYWORD)
			return node
		case TKIND_INTEGER:
			node := TerminalNode(
				NT_INTEGER, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_INTEGER)
			return node
		case TKIND_BIG_INTEGER:
			node := TerminalNode(
				NT_BIG_INTEGER, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_BIG_INTEGER)
			return node
		case TKIND_OTHER_INTEGER:
			node := TerminalNode(
				NT_OTHER_INTEGER, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_OTHER_INTEGER)
			return node
		case TKIND_OTHER_BIG_INTEGER:
			node := TerminalNode(
				NT_OTHER_BIG_INTEGER, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_OTHER_BIG_INTEGER)
			return node
		case TKIND_FLOAT:
			node := TerminalNode(
				NT_FLOAT, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_FLOAT)
			return node
		case TKIND_OTHER_FLOAT:
			node := TerminalNode(
				NT_OTHER_FLOAT, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_OTHER_FLOAT)
			return node
		case TKIND_STRING:
			node := TerminalNode(
				NT_STRING, 
				p.lookahead.value, 
				p.lookahead.position,
			)
			p.acceptT(TKIND_STRING)
			return node
		case TKIND_KEYWORD: 
			if p.checkK("true") || p.checkK("false") {
				node := TerminalNode(
					NT_BOOLEAN, 
					p.lookahead.value, 
					p.lookahead.position,
				)
				p.acceptT(TKIND_KEYWORD)
				return node
			} else if p.checkK("null") {
				node := TerminalNode(
					NT_NULL, 
					p.lookahead.value, 
					p.lookahead.position,
				)
				p.acceptT(TKIND_KEYWORD)
				return node
			} else if p.checkK("self") {
				panic("Not implemented self!!!")
			} else if p.checkK("super") {
				panic("Not implemented super!!!")
			} else if p.checkK("function") {
				start := p.lookahead.position
				p.acceptK("function")
				p.acceptS("(")
				

				var parameters *[][]interface{} = new([][]interface{})

				if p.checkT(TKIND_ID) {
					param := p.lookahead.value
					p.acceptT(TKIND_ID)

					if p.checkS("...") {
						// Variadic parameter
						p.acceptS("...")
						param = fmt.Sprintf("%s...", param)
					}

					*parameters = append(*parameters, []interface{}{param, p.previous.position})

					for p.checkS(",") {
						p.acceptS(",")

						if !p.checkT(TKIND_ID) {
							raiseError(p, fmt.Sprintf("missing parameter after \",\"."), p.lookahead.position)
						}

						param := p.lookahead.value
						p.acceptT(TKIND_ID)

						if p.checkS("...") {
							// Variadic parameter
							p.acceptS("...")
							param = fmt.Sprintf("%s...", param)
						}

						*parameters = append(*parameters, []interface{}{param, p.previous.position})
					}
				}

				p.acceptS(")")
				p.acceptS("{")

				var body *[]*node_t = new([]*node_t)

				bodyN := p.parseCompoundStatement()

				for bodyN != nil {
					*body = append(*body, bodyN)
					bodyN = p.parseCompoundStatement()
				}

				p.acceptS("}")
				end := p.previous.position

				return HeadlessFunctionNode(parameters, body, start.merge(end))
			}
	}

	return nil
}

func (p *parser_t) parseGroup() *node_t {
	if p.checkS("[") { 
		// '[' zeroOrOneExpression (',' mandatoryExpression)+ ']'
		start := p.lookahead.position
		p.acceptS("[")

		var elements *[]*node_t = new([]*node_t)

		elementN := p.parseZeroOrOneExpression()
		if elementN != nil {
			*elements = append(*elements, elementN)

			for p.checkS(",") {
				p.acceptS(",")
				elementN = p.parseMandatoryExpression()
				*elements = append(*elements, elementN)
			}
		}

		p.acceptS("]")
		end := p.previous.position

		return ArrayNode(elements, start.merge(end))

	} else if p.checkS("{") {
		// '{' (zeroOrOneExpression ':' mandatoryExpression)? (',' zeroOrOneExpression ':' mandatoryExpression)+ '}'
		start := p.lookahead.position
		p.acceptS("{")

		var pairs *[][]*node_t = new([][]*node_t)

		keyN := p.parseZeroOrOneExpression()
		if keyN != nil {
			p.acceptS(":")
			valueN := p.parseMandatoryExpression()
			*pairs = append(*pairs, []*node_t{keyN, valueN})

			for p.checkS(",") {
				p.acceptS(",")
				keyN = p.parseZeroOrOneExpression()
				if keyN == nil {
					raiseError(p, fmt.Sprintf("missing key after \",\"."), p.lookahead.position)
				}

				p.acceptS(":")
				valueN = p.parseMandatoryExpression()
				*pairs = append(*pairs, []*node_t{keyN, valueN})
			}
		}

		p.acceptS("}")
		end := p.previous.position

		return ObjectNode(pairs, start.merge(end))

	} else if p.checkS("(") {
		// '(' mandatoryExpression ')'
		p.acceptS("(")
		expr := p.parseMandatoryExpression()
		p.acceptS(")")

		return expr
	} else {
		return p.parseTerminal()
	}
}

func (p *parser_t) parseMemberOrCall() *node_t {
	node := p.parseGroup()

	if node == nil {
		return node
	}

	// member
	tmp0 := p.memberAccess
	// indexing
	tmp1 := p.indexingLevel
	// call level
	tmp2 := p.callLevel

	for p.checkS(".") || p.checkS("[") || p.checkS("(") {
		if p.checkS(".") {
			p.memberAccess += 1
			if p.memberAccess > MAX_MEMBER_ACCESS {
				raiseError(p, fmt.Sprintf("member access nesting level too deep."), node.position.merge(p.previous.position))
			}

			// (. TKIND_ID)
			p.acceptS(".")

			member := p.lookahead.value
			p.acceptT(TKIND_ID)

			node = MemberAccess(node, member, node.position.merge(p.previous.position))
		} else if p.checkS("[") {
			p.indexingLevel += 1
			if p.indexingLevel > MAX_INDEXING {
				raiseError(p, fmt.Sprintf("subscription nesting level too deep."), node.position.merge(p.previous.position))
			}

			// '[' mandatoryExpression ']'
			p.acceptS("[")
			expr := p.parseMandatoryExpression()
			p.acceptS("]")

			node = IndexAccess(node, expr, node.position.merge(p.previous.position))
		} else if p.checkS("(") {
			p.callLevel += 1
			if p.callLevel > MAX_CALL {
				raiseError(p, fmt.Sprintf("call nesting level too deep."), node.position.merge(p.previous.position))
			}

			// '(' mandatoryExpression ')'
			p.acceptS("(")

			var args *[]*node_t = new([]*node_t)

			argN := p.parseZeroOrOneExpression()

			if argN != nil {
				argc := 1
				*args = append(*args, argN)

				for p.checkS(",") {
					p.acceptS(",")
					argN = p.parseMandatoryExpression()
					*args = append(*args, argN)

					argc += 1
					if argc > MAX_ARGS {
						raiseError(p, fmt.Sprintf("too many arguments, max 255 got %d.", argc), p.lookahead.position)
					}
				}
			}

			p.acceptS(")")

			node = Call(node, args, node.position.merge(p.previous.position))
		}
	}

	// restore member access
	p.memberAccess = tmp0
	// restore indexing
	p.indexingLevel = tmp1
	// restore call level
	p.callLevel = tmp2

	return node
}

func (p *parser_t) parsePostfix() *node_t {
	node := p.parseMemberOrCall()

	if node == nil {
		return node
	}

	if p.checkS("++") || p.checkS("--") {
		operator := p.lookahead
		p.acceptS(operator.value)

		return PostfixExpressionNode(operator.value, node, node.position.merge(p.previous.position))
	} else if p.checkS("?") {
		p.acceptS("?")

		trueval := p.parseZeroOrOneExpression()
		if trueval == nil {
			raiseError(p, fmt.Sprintf("missing true value after \"?\"."), p.lookahead.position)
		}

		p.acceptS(":")

		falseval := p.parseZeroOrOneExpression()
		if falseval == nil {
			raiseError(p, fmt.Sprintf("missing false value after \":\"."), p.lookahead.position)
		}
		
		return TernaryExpressionNode(node, trueval, falseval, node.position.merge(p.previous.position))
	}

	return node
}

func (p *parser_t) parseUnary() *node_t {
	if p.checkS("+") || p.checkS("-") || p.checkS("!") || p.checkS("~") || p.checkS("++") || p.checkS("--") {
		operator := p.lookahead
		p.acceptS(operator.value)

		// Watch recursion
		operand := p.parseUnary()
		if operand == nil {
			raiseError(p, fmt.Sprintf("missing operand after \"%s\".", operator.value), operator.position)
		}

		return UnaryExpressionNode(operator.value, operand, operator.position.merge(p.previous.position))
	} else {
		return p.parsePostfix()
	}
}

func (p *parser_t) parseMul() *node_t {
	node := p.parseUnary()

	if node == nil { 
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

		rhs := p.parseUnary()

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

	if node == nil { 
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
func (p *parser_t) parseSimpleStatement() *node_t {
	switch p.lookahead.value {
		case "var", "let", "const":
			start := p.lookahead.position
			var ntype NodeType

			if strings.Compare(p.lookahead.value, "var") == 0 {
				ntype = NT_VARIABLE_DEC
				p.acceptK("var")
			} else if strings.Compare(p.lookahead.value, "let") == 0 {
				ntype = NT_LOCAL_DEC
				p.acceptK("let")
			} else {
				ntype = NT_CONST_DEC
				p.acceptK("const")
			}

			p.acceptK(p.lookahead.value)

			if !p.checkT(TKIND_ID) {
				raiseError(p, fmt.Sprintf("variable name is required after \"%s\", got \"%s\".", p.previous.value, p.lookahead.value), p.lookahead.position)
			}

			declairations := new([][]interface{})

			for p.checkT(TKIND_ID) {
				variable := p.lookahead.value
				position := p.lookahead.position
				p.acceptT(TKIND_ID)

				var value *node_t = nil

				if p.checkS("=") {
					value = p.parseMandatoryExpression()
				}

				*declairations = append(*declairations, []interface{}{variable, position, value})
			}

			p.acceptS(";")

			return VariableDeclairationNode(ntype, declairations, start.merge(p.previous.position))
		
		case ";":
			for p.checkS(";") {
				p.acceptS(";")
			}

			return EmptyExpressionNode(p.previous.position)

		default:
			node := p.parseZeroOrOneExpression()

			if node == nil {
				return node
			}

			p.acceptS(";")

			return ExpressionStatementNode(node, node.position.merge(p.previous.position))
	}
}

func (p *parser_t) parseCompoundStatement() *node_t {
	return p.parseSimpleStatement()
}

func (p *parser_t) parseFile() *node_t {
	body := new([]*node_t)

	stmntN := p.parseCompoundStatement()

	for (!p.checkT(TKIND_EOF)) && stmntN != nil {
		*body = append(*body, stmntN)
		stmntN = p.parseCompoundStatement()
	}

	// Eof
	p.acceptT(TKIND_EOF)

	return p.parse
}

func (p *parser_t) parse() *node_t {
	p.lookahead, p.previous = p.lexer.nextToken(), p.lookahead
	return p.parseFile()
}
