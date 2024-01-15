package main

import (
	"fmt"
	"strings"

	"jackass/ast"
	"jackass/util"
)

const (
	MAX_MEMBER_ACCESS      int = 50
	MAX_INDEXING           int = 50
	MAX_CALL               int = 50
	MAX_EXPRESSION_NESTING int = 100
	//
	MAX_ARGS            int = 255
	MAX_STATIC_ELEMENTS int = 512
)

type parser_t struct {
	lexer                                                   *lexer_t
	lookahead, previous                                     *token_t
	memberAccess, indexingLevel, callLevel, expressionLevel int
}

func Parser(lexer *lexer_t) *parser_t {
	parser := new(parser_t)

	//lint:ignore SA4031 possible not nil
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
		raiseError(p, fmt.Sprintf("unexpected token \"%s\".", p.lookahead.value), p.lookahead.Position)
	}
}

// acceptV accepts lookahead value
// @param value string
//
//lint:ignore U1000 disable unused
func (p *parser_t) acceptV(value string) {
	if p.checkV(value) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, value), p.lookahead.Position)
	}
}

// acceptS accepts lookahead symbol
// @param symbol string
func (p *parser_t) acceptS(symbol string) {
	if p.checkS(symbol) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, symbol), p.lookahead.Position)
	}
}

// acceptK accepts lookahead keyword
// @param keyword string
func (p *parser_t) acceptK(keyword string) {
	if p.checkK(keyword) {
		p.previous, p.lookahead = p.lookahead, p.lexer.nextToken()
	} else {
		raiseError(p, fmt.Sprintf("unexpected token \"%s\", Did you mean \"%s\"", p.lookahead.value, keyword), p.lookahead.Position)
	}
}

func (p *parser_t) parseTerminal() *ast.Node_t {
	switch p.lookahead.kind {
	case TKIND_ID:
		node := ast.TerminalNode(
			ast.NT_ID,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_ID)
		return node
	case TKIND_INTEGER:
		node := ast.TerminalNode(
			ast.NT_INTEGER,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_INTEGER)
		return node
	case TKIND_BIG_INTEGER:
		node := ast.TerminalNode(
			ast.NT_BIG_INTEGER,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_BIG_INTEGER)
		return node
	case TKIND_OTHER_INTEGER:
		node := ast.TerminalNode(
			ast.NT_OTHER_INTEGER,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_OTHER_INTEGER)
		return node
	case TKIND_OTHER_BIG_INTEGER:
		node := ast.TerminalNode(
			ast.NT_OTHER_BIG_INTEGER,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_OTHER_BIG_INTEGER)
		return node
	case TKIND_FLOAT:
		node := ast.TerminalNode(
			ast.NT_FLOAT,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_FLOAT)
		return node
	case TKIND_OTHER_FLOAT:
		node := ast.TerminalNode(
			ast.NT_OTHER_FLOAT,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_OTHER_FLOAT)
		return node
	case TKIND_STRING:
		node := ast.TerminalNode(
			ast.NT_STRING,
			p.lookahead.value,
			p.lookahead.Position,
		)
		p.acceptT(TKIND_STRING)
		return node
	case TKIND_KEYWORD:
		if p.checkK("true") || p.checkK("false") {
			node := ast.TerminalNode(
				ast.NT_BOOLEAN,
				p.lookahead.value,
				p.lookahead.Position,
			)
			p.acceptT(TKIND_KEYWORD)
			return node
		} else if p.checkK("null") {
			node := ast.TerminalNode(
				ast.NT_NULL,
				p.lookahead.value,
				p.lookahead.Position,
			)
			p.acceptT(TKIND_KEYWORD)
			return node
		} else if p.checkK("self") {
			panic("Not implemented self!!!")
		} else if p.checkK("super") {
			panic("Not implemented super!!!")
		} else if p.checkK("function") {
			start := p.lookahead.Position
			p.acceptK("function")
			p.acceptS("(")

			var count int = 0
			var parameters *[][]interface{} = new([][]interface{})

			if p.checkT(TKIND_ID) {
				count++
				param := p.lookahead.value
				p.acceptT(TKIND_ID)

				if p.checkS("...") {
					// Variadic parameter
					p.acceptS("...")
					param = fmt.Sprintf("%s...", param)
				}

				*parameters = append(*parameters, []interface{}{param, p.previous.Position})

				for p.checkS(",") {
					p.acceptS(",")

					if !p.checkT(TKIND_ID) {
						raiseError(p, "missing parameter after \",\".", p.lookahead.Position)
					}

					param := p.lookahead.value
					p.acceptT(TKIND_ID)

					if p.checkS("...") {
						// Variadic parameter
						p.acceptS("...")
						param = fmt.Sprintf("%s...", param)
					}

					*parameters = append(*parameters, []interface{}{param, p.previous.Position})
					count++

					if count > MAX_ARGS {
						break
					}
				}
			}

			p.acceptS(")")

			if count > MAX_ARGS {
				raiseError(p, "too many parameters. Try variadict function.", p.lookahead.Position)
			}

			p.acceptS("{")

			var body *[]*ast.Node_t = new([]*ast.Node_t)

			bodyN := p.parseCompoundStatement()

			for bodyN != nil {
				*body = append(*body, bodyN)
				bodyN = p.parseCompoundStatement()
			}

			p.acceptS("}")
			end := p.previous.Position

			return ast.HeadlessFunctionNode(parameters, body, start.Merge(end))
		}
	}

	return nil
}

func (p *parser_t) parseGroup() *ast.Node_t {
	if p.checkS("[") {
		// '[' zeroOrOneExpression (',' mandatoryExpression)+ ']'
		start := p.lookahead.Position
		p.acceptS("[")

		var count int = 0
		var elements *[]*ast.Node_t = new([]*ast.Node_t)

		elementN := p.parseZeroOrOneExpression()
		if elementN != nil {
			count++
			*elements = append(*elements, elementN)

			for p.checkS(",") {
				p.acceptS(",")
				elementN = p.parseMandatoryExpression()
				*elements = append(*elements, elementN)
				count++

				if count >= int(util.MAX_SAFE_INDEX) {
					break
				}
			}
		}

		p.acceptS("]")
		end := p.previous.Position

		if count > MAX_STATIC_ELEMENTS {
			raiseError(p, fmt.Sprintf("too many elements %d, max %d", count, MAX_STATIC_ELEMENTS), start.Merge(end))
		}

		return ast.ArrayNode(elements, start.Merge(end))

	} else if p.checkS("{") {
		// '{' (zeroOrOneExpression ':' mandatoryExpression)? (',' zeroOrOneExpression ':' mandatoryExpression)+ '}'
		start := p.lookahead.Position
		p.acceptS("{")

		var count int = 0
		var pairs *[][]*ast.Node_t = new([][]*ast.Node_t)

		keyN := p.parseZeroOrOneExpression()
		if keyN != nil {
			count++
			p.acceptS(":")
			valueN := p.parseMandatoryExpression()
			*pairs = append(*pairs, []*ast.Node_t{keyN, valueN})

			for p.checkS(",") {
				p.acceptS(",")
				keyN = p.parseZeroOrOneExpression()
				if keyN == nil {
					raiseError(p, "missing key after \",\".", p.lookahead.Position)
				}

				p.acceptS(":")
				valueN = p.parseMandatoryExpression()
				*pairs = append(*pairs, []*ast.Node_t{keyN, valueN})
				count++

				if count >= int(util.MAX_SAFE_INDEX) {
					break
				}
			}
		}

		p.acceptS("}")
		end := p.previous.Position

		if count > MAX_STATIC_ELEMENTS {
			raiseError(p, fmt.Sprintf("too many elements %d, max %d", count, MAX_STATIC_ELEMENTS), start.Merge(end))
		}

		return ast.ObjectNode(pairs, start.Merge(end))

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

func (p *parser_t) parseMemberOrCall() *ast.Node_t {
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
				raiseError(p, "member access nesting level too deep.", node.Position.Merge(p.previous.Position))
			}

			// (. TKIND_ID)
			p.acceptS(".")

			member := p.lookahead.value
			p.acceptT(TKIND_ID)

			node = ast.MemberAccess(node, member, node.Position.Merge(p.previous.Position))
		} else if p.checkS("[") {
			p.indexingLevel += 1
			if p.indexingLevel > MAX_INDEXING {
				raiseError(p, "subscription nesting level too deep.", node.Position.Merge(p.previous.Position))
			}

			// '[' mandatoryExpression ']'
			p.acceptS("[")
			expr := p.parseMandatoryExpression()
			p.acceptS("]")

			node = ast.IndexAccess(node, expr, node.Position.Merge(p.previous.Position))
		} else if p.checkS("(") {
			p.callLevel += 1
			if p.callLevel > MAX_CALL {
				raiseError(p, "call nesting level too deep.", node.Position.Merge(p.previous.Position))
			}

			// '(' mandatoryExpression ')'
			p.acceptS("(")

			var args *[]*ast.Node_t = new([]*ast.Node_t)

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
						raiseError(p, fmt.Sprintf("too many arguments, max 255 got %d.", argc), p.lookahead.Position)
					}
				}
			}

			p.acceptS(")")

			node = ast.Call(node, args, node.Position.Merge(p.previous.Position))
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

func (p *parser_t) parsePostfix() *ast.Node_t {
	node := p.parseMemberOrCall()

	if node == nil {
		return node
	}

	if p.checkS("++") || p.checkS("--") {
		operator := p.lookahead
		p.acceptS(operator.value)

		return ast.PostfixExpressionNode(operator.value, node, node.Position.Merge(p.previous.Position))
	} else if p.checkS("?") {
		p.acceptS("?")

		trueval := p.parseZeroOrOneExpression()
		if trueval == nil {
			raiseError(p, "missing true value after \"?\".", p.lookahead.Position)
		}

		p.acceptS(":")

		falseval := p.parseZeroOrOneExpression()
		if falseval == nil {
			raiseError(p, "missing false value after \":\".", p.lookahead.Position)
		}

		return ast.TernaryExpressionNode(node, trueval, falseval, node.Position.Merge(p.previous.Position))
	}

	return node
}

func (p *parser_t) parseUnary() *ast.Node_t {
	if p.checkS("+") || p.checkS("-") || p.checkS("!") || p.checkS("~") || p.checkS("++") || p.checkS("--") {
		operator := p.lookahead
		p.acceptS(operator.value)

		// Watch recursion
		operand := p.parseUnary()
		if operand == nil {
			raiseError(p, fmt.Sprintf("missing operand after \"%s\".", operator.value), operator.Position)
		}

		return ast.UnaryExpressionNode(operator.value, operand, operator.Position.Merge(p.previous.Position))
	} else {
		return p.parsePostfix()
	}
}

func (p *parser_t) parseMul() *ast.Node_t {
	node := p.parseUnary()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("*") || p.checkS("/") || p.checkS("%") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseUnary()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAdd() *ast.Node_t {
	node := p.parseMul()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("+") || p.checkS("-") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseMul()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseShift() *ast.Node_t {
	node := p.parseAdd()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("<<") || p.checkS(">>") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseAdd()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseRel() *ast.Node_t {
	node := p.parseShift()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("<") || p.checkS("<=") || p.checkS(">") || p.checkS(">=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseShift()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseEql() *ast.Node_t {
	node := p.parseRel()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("==") || p.checkS("!=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseRel()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseBit() *ast.Node_t {
	node := p.parseEql()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("&") || p.checkS("|") || p.checkS("^") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseEql()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_BINARY_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseLog() *ast.Node_t {
	node := p.parseBit()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("&&") || p.checkS("||") || p.checkS("??") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseBit()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAss() *ast.Node_t {
	node := p.parseLog()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseLog()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseAug() *ast.Node_t {
	node := p.parseAss()

	if node == nil {
		return node
	}

	tmp := p.expressionLevel
	for p.checkS("*=") || p.checkS("/=") || p.checkS("%=") || p.checkS("+=") || p.checkS("-=") || p.checkS("<<=") || p.checkS(">>=") || p.checkS("&=") || p.checkS("|=") || p.checkS("^=") {

		p.expressionLevel += 1
		if p.expressionLevel > MAX_EXPRESSION_NESTING {
			raiseError(p, "expression nesting too deep.", p.lookahead.Position)
		}

		operator := p.lookahead
		p.acceptS(operator.value)

		rhs := p.parseAss()

		if rhs == nil {
			raiseError(p, fmt.Sprintf("missing right-hand expression after \"%s\".", p.previous.value), operator.Position)
		}

		node = ast.BinaryExpressionNode(ast.NT_LOGICAL_EXPRESSION, operator.value, node, rhs)
	}

	p.expressionLevel = tmp
	return node
}

func (p *parser_t) parseZeroOrOneExpression() *ast.Node_t {
	return p.parseAug()
}

func (p *parser_t) parseMandatoryExpression() *ast.Node_t {
	node := p.parseZeroOrOneExpression()
	if node != nil {
		return node
	}

	// error
	raiseError(p, fmt.Sprintf("an expression is required, got \"%s\".", p.lookahead.value), p.lookahead.Position)
	return nil
}

func (p *parser_t) parseSimpleStatement() *ast.Node_t {
	switch p.lookahead.value {
	case "var", "let", "const":
		start := p.lookahead.Position
		var ntype ast.NodeType

		if strings.Compare(p.lookahead.value, "var") == 0 {
			ntype = ast.NT_VARIABLE_DEC
			p.acceptK("var")
		} else if strings.Compare(p.lookahead.value, "let") == 0 {
			ntype = ast.NT_LOCAL_DEC
			p.acceptK("let")
		} else {
			ntype = ast.NT_CONST_DEC
			p.acceptK("const")
		}

		if !p.checkT(TKIND_ID) {
			raiseError(p, fmt.Sprintf("variable name is required after \"%s\", got \"%s\".", p.previous.value, p.lookahead.value), p.lookahead.Position)
		}

		declairations := new([][]interface{})

		for p.checkT(TKIND_ID) {
			variable := p.lookahead.value
			position := p.lookahead.Position
			p.acceptT(TKIND_ID)

			var value *ast.Node_t = nil

			if p.checkS("=") {
				p.acceptS("=")
				value = p.parseMandatoryExpression()
			}

			*declairations = append(*declairations, []interface{}{variable, position, value})

			if !p.checkS(",") {
				break
			} else {
				p.acceptS(",")
			}

			if !p.checkT(TKIND_ID) {
				raiseError(p, fmt.Sprintf("variable name is required after \"%s\", got \"%s\".", p.previous.value, p.lookahead.value), p.previous.Position)
			}
		}

		p.acceptS(";")

		return ast.VariableDeclairationNode(ntype, declairations, start.Merge(p.previous.Position))

	case ";":
		for p.checkS(";") {
			p.acceptS(";")
		}

		return ast.EmptyExpressionNode(p.previous.Position)

	default:
		node := p.parseZeroOrOneExpression()

		if node == nil {
			return node
		}

		p.acceptS(";")

		return ast.ExpressionStatementNode(node, node.Position.Merge(p.previous.Position))
	}
}

func (p *parser_t) parseCompoundStatement() *ast.Node_t {
	return p.parseSimpleStatement()
}

func (p *parser_t) parseFile() *ast.Node_t {
	body := new([]*ast.Node_t)

	stmntN := p.parseCompoundStatement()

	for stmntN != nil {
		*body = append(*body, stmntN)
		stmntN = p.parseCompoundStatement()
	}

	// Eof
	p.acceptT(TKIND_EOF)

	return ast.FileNode(body)
}

func (p *parser_t) parse() *ast.Node_t {
	p.lookahead, p.previous = p.lexer.nextToken(), p.lookahead
	return p.parseFile()
}
