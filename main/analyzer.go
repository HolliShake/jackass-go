package main
import (
	"fmt"
	"github.com/jackass/jutil"
)

type analyzer_t struct {
	parser *parser_t
}

func Analyzer(parser *parser_t) *analyzer_t {
	analyzer := new(analyzer_t)

	//lint:ignore SA4031 possible not nil
	if analyzer != nil {
		analyzer.parser = parser
	} else {
		basicError("Out of memory!!!")
	}

	return analyzer
}

func (a *analyzer_t) visit(node *node_t) *node_t {
	switch node.ntype {
		case NT_ID:
			return a.analyzeID(node)
		case NT_INTEGER:
			return a.analyzeInteger(node)
		case NT_OTHER_INTEGER:
			return a.analyzeOtherInteger(node)
		case NT_EXPRESSION_STATEMENT:
			return a.analyzeExpressionStatement(node)
		case NT_FILE:
			return a.analyzeFile(node)
		default:
			panic(fmt.Sprintf("node not implemented %d!!!", node.ntype))
	}
}

func (a *analyzer_t) analyzeID(node *node_t) *node_t {
	return a.visit(a.parser.parse())
}

func (a *analyzer_t) analyzeInteger(node *node_t) *node_t {
	node.terminal.value = fmt.Sprintf("%d", jutil.ParseInt(node.terminal.value))
	return node
}

func (a *analyzer_t) analyzeOtherInteger(node *node_t) *node_t {
	// change to integer node
	node.ntype = NT_INTEGER
	switch node.terminal.value[0:2] {
		case "0x", "0X":
			node.terminal.value = fmt.Sprintf("%d", jutil.Parse(node.terminal.value, 16))
		case "0o", "0O":
			node.terminal.value = fmt.Sprintf("%d", jutil.Parse(node.terminal.value,  8))
		case "0b", "0B":
			node.terminal.value = fmt.Sprintf("%d", jutil.Parse(node.terminal.value,  2))
		default:
			panic(fmt.Sprintf("invalid number format %s!!!", node.terminal.value))
	}
	return node
}

// 

func (a *analyzer_t) analyzeExpressionStatement(node *node_t) *node_t {
	node.expressionStatement.expression = a.visit(node.expressionStatement.expression)
	return node
}

func (a *analyzer_t) analyzeFile(node *node_t) *node_t {
	for i := 0; i < len(*node.file.body); i ++ {
		(*node.file.body)[i] = a.visit((*node.file.body)[i])
	}
	return node
}

func (a *analyzer_t) analyze() *node_t {
	return a.visit(a.parser.parse())
}