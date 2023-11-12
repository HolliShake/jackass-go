package main
import (
	"fmt"
	"strings"
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

func (a *analyzer_t) getFilePath() string {
	return a.parser.getFilePath()
}

func (a *analyzer_t) getFileCode() string {
	return a.parser.getFileCode()
}

func (a *analyzer_t) visit(node *node_t) *node_t {
	switch node.ntype {
		case NT_ID:
			return a.analyzeID(node)
		case NT_INTEGER:
			return a.analyzeInteger(node)
		case NT_OTHER_INTEGER:
			return a.analyzeOtherInteger(node)
		case NT_FLOAT:
			return a.analyzeFloat(node)
		case NT_OTHER_FLOAT:
			return a.analyzeOtherFloat(node)
		case NT_STRING:
			return a.analyzeString(node)
		case NT_BOOLEAN:
			return a.analyzeBoolean(node)
		case NT_NULL:
			return a.analyzeNull(node)
		case NT_ARRAY:
			return a.analyzeArray(node)
		case NT_OBJECT:
			return a.analyzeObject(node)
		// 
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

func (a *analyzer_t) analyzeFloat(node *node_t) *node_t {
	node.terminal.value = fmt.Sprintf("%f", jutil.ParseFloat(node.terminal.value))
	return node
}

func (a *analyzer_t) analyzeOtherFloat(node *node_t) *node_t {
	// change to float node
	node.ntype = NT_FLOAT
	node.terminal.value = fmt.Sprintf("%f", jutil.ParseFloat(node.terminal.value))
	return node
}

func (a *analyzer_t) analyzeString(node *node_t) *node_t {
	if jutil.Utf_codePointLength(node.terminal.value) > int(jutil.MAX_SAFE_INDEX) {
		raiseError(a, "string length is too large to represent", node.position)
	}
	return node
}

func (a *analyzer_t) analyzeBoolean(node *node_t) *node_t {
	if !(
		(strings.Compare(node.terminal.value,  "true") == 0) ||
		(strings.Compare(node.terminal.value, "false") == 0)) {
		panic(fmt.Sprintf("invalid boolean value \"%s\"!!!", node.terminal.value))
	}
	return node
}

func (a *analyzer_t) analyzeNull(node *node_t) *node_t {
	if (strings.Compare(node.terminal.value, "null") != 0) {
		panic(fmt.Sprintf("invalid null value \"%s\"!!!", node.terminal.value))
	}
	return node
}

func (a *analyzer_t) analyzeArray(node *node_t) *node_t {
	for i := 0; i < len(*node.array.elements); i++ {
		(*node.array.elements)[i] = a.visit((*node.array.elements)[i])
	}
	return node
}

func (a *analyzer_t) analyzeObject(node *node_t) *node_t {
	for i := 0; i < len(*node.object.members); i++ {
		(*node.object.members)[i][0] = a.visit((*node.object.members)[i][0])
		(*node.object.members)[i][1] = a.visit((*node.object.members)[i][1])
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