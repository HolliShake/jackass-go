package main

import (
	"fmt"
	"strings"

	"jackass/ast"
	"jackass/jutil"
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

func (a *analyzer_t) visit(node *ast.Node_t) *ast.Node_t {
	switch node.NodeType {
	case ast.NT_ID:
		return a.analyzeID(node)
	case ast.NT_INTEGER:
		return a.analyzeInteger(node)
	case ast.NT_OTHER_INTEGER:
		return a.analyzeOtherInteger(node)
	case ast.NT_BIG_INTEGER:
		return a.analyzeBigInteger(node)
	case ast.NT_OTHER_BIG_INTEGER:
		return a.analyzeOtherBigInteger(node)
	case ast.NT_FLOAT:
		return a.analyzeFloat(node)
	case ast.NT_OTHER_FLOAT:
		return a.analyzeOtherFloat(node)
	case ast.NT_STRING:
		return a.analyzeString(node)
	case ast.NT_BOOLEAN:
		return a.analyzeBoolean(node)
	case ast.NT_NULL:
		return a.analyzeNull(node)
	case ast.NT_ARRAY:
		return a.analyzeArray(node)
	case ast.NT_OBJECT:
		return a.analyzeObject(node)
	case ast.NT_MEMBER_ACCESS:
		return a.analyzeMemberAccess(node)
	case ast.NT_INDEX:
		return a.analyzeIndex(node)
	case ast.NT_CALL:
		return a.analyzeCall(node)
	//
	case ast.NT_VARIABLE_DEC:
		return a.analyzeVariableDeclairation(node)
	case ast.NT_EXPRESSION_STATEMENT:
		return a.analyzeExpressionStatement(node)
	case ast.NT_FILE:
		return a.analyzeFile(node)
	default:
		panic(fmt.Sprintf("node not implemented %d!!!", node.NodeType))
	}
}

func (a *analyzer_t) analyzeID(node *ast.Node_t) *ast.Node_t {
	return node
}

func (a *analyzer_t) analyzeInteger(node *ast.Node_t) *ast.Node_t {
	node.Terminal.Value = fmt.Sprintf("%d", jutil.ParseInt(node.Terminal.Value))
	return node
}

func (a *analyzer_t) analyzeOtherInteger(node *ast.Node_t) *ast.Node_t {
	// change to integer node
	node.NodeType = ast.NT_INTEGER
	switch node.Terminal.Value[0:2] {
	case "0x", "0X":
		node.Terminal.Value = fmt.Sprintf("%d", jutil.Parse(node.Terminal.Value, 16))
	case "0o", "0O":
		node.Terminal.Value = fmt.Sprintf("%d", jutil.Parse(node.Terminal.Value, 8))
	case "0b", "0B":
		node.Terminal.Value = fmt.Sprintf("%d", jutil.Parse(node.Terminal.Value, 2))
	default:
		panic(fmt.Sprintf("invalid number format %s!!!", node.Terminal.Value))
	}
	return node
}

func (a *analyzer_t) analyzeBigInteger(node *ast.Node_t) *ast.Node_t {
	return node
}

func (a *analyzer_t) analyzeOtherBigInteger(node *ast.Node_t) *ast.Node_t {
	return node
}

func (a *analyzer_t) analyzeFloat(node *ast.Node_t) *ast.Node_t {
	node.Terminal.Value = fmt.Sprintf("%f", jutil.ParseFloat(node.Terminal.Value))
	return node
}

func (a *analyzer_t) analyzeOtherFloat(node *ast.Node_t) *ast.Node_t {
	// change to float node
	node.NodeType = ast.NT_FLOAT
	node.Terminal.Value = fmt.Sprintf("%f", jutil.ParseFloat(node.Terminal.Value))
	return node
}

func (a *analyzer_t) analyzeString(node *ast.Node_t) *ast.Node_t {
	if jutil.Utf_codePointLength(node.Terminal.Value) > int(jutil.MAX_SAFE_INDEX) {
		raiseError(a, "string length is too large to represent", node.Position)
	}
	return node
}

func (a *analyzer_t) analyzeBoolean(node *ast.Node_t) *ast.Node_t {
	if !((strings.Compare(node.Terminal.Value, "true") == 0) ||
		(strings.Compare(node.Terminal.Value, "false") == 0)) {
		panic(fmt.Sprintf("invalid boolean value \"%s\"!!!", node.Terminal.Value))
	}
	return node
}

func (a *analyzer_t) analyzeNull(node *ast.Node_t) *ast.Node_t {
	if strings.Compare(node.Terminal.Value, "null") != 0 {
		panic(fmt.Sprintf("invalid null value \"%s\"!!!", node.Terminal.Value))
	}
	return node
}

func (a *analyzer_t) analyzeArray(node *ast.Node_t) *ast.Node_t {
	for i := 0; i < len(*node.Array.Elements); i++ {
		(*node.Array.Elements)[i] = a.visit((*node.Array.Elements)[i])
	}
	return node
}

func (a *analyzer_t) analyzeObject(node *ast.Node_t) *ast.Node_t {
	for i := 0; i < len(*node.Object.Members); i++ {
		(*node.Object.Members)[i][0] = a.visit((*node.Object.Members)[i][0])
		(*node.Object.Members)[i][1] = a.visit((*node.Object.Members)[i][1])
	}
	return node
}

func (a *analyzer_t) analyzeMemberAccess(node *ast.Node_t) *ast.Node_t {
	node.MemberAccess.Object = a.visit(node.MemberAccess.Object)
	return node
}

func (a *analyzer_t) analyzeIndex(node *ast.Node_t) *ast.Node_t {
	node.IndexAccess.Object = a.visit(node.IndexAccess.Object)
	node.IndexAccess.Index = a.visit(node.IndexAccess.Index)
	return node
}

func (a *analyzer_t) analyzeCall(node *ast.Node_t) *ast.Node_t {
	node.Call.Object = a.visit(node.Call.Object)

	//
	for i := 0; i < len(*node.Call.Args); i++ {
		(*node.Call.Args)[i] = a.visit((*node.Call.Args)[i])
	}

	return node
}

//

func (a *analyzer_t) analyzeVariableDeclairation(node *ast.Node_t) *ast.Node_t {

	for i := 0; i < len(*node.Declairation.Declairations); i++ {
		// var_name := (*node.Declairation.Declairations)[i][0]
		// var_post := (*node.Declairation.Declairations)[i][1]

		var_valu := ((*node.Declairation.Declairations)[i][2]).(*ast.Node_t)

		if var_valu != nil {
			(*node.Declairation.Declairations)[i][2] = a.visit(var_valu)
		}
	}

	return node
}

func (a *analyzer_t) analyzeExpressionStatement(node *ast.Node_t) *ast.Node_t {
	node.ExpressionStatement.Expression = a.visit(node.ExpressionStatement.Expression)
	return node
}

func (a *analyzer_t) analyzeFile(node *ast.Node_t) *ast.Node_t {
	for i := 0; i < len(*node.File.Body); i++ {
		(*node.File.Body)[i] = a.visit((*node.File.Body)[i])
	}
	return node
}

func (a *analyzer_t) analyze() *ast.Node_t {
	return a.visit(a.parser.parse())
}
