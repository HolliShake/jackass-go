package ast

import (
	"jackass/shared"
)

type NodeType = int

const (
	NT_ID                 NodeType = 0
	NT_INTEGER            NodeType = 1
	NT_BIG_INTEGER        NodeType = 2
	NT_OTHER_INTEGER      NodeType = 3
	NT_OTHER_BIG_INTEGER  NodeType = 4
	NT_FLOAT              NodeType = 5
	NT_OTHER_FLOAT        NodeType = 6
	NT_STRING             NodeType = 7
	NT_BOOLEAN            NodeType = 8
	NT_NULL               NodeType = 9
	NT_SELF               NodeType = 10
	NT_SUPER              NodeType = 11
	NT_HEADLESS_FUNCTION  NodeType = 12
	NT_ARRAY              NodeType = 13
	NT_OBJECT             NodeType = 14
	NT_MEMBER_ACCESS      NodeType = 15
	NT_INDEX              NodeType = 16
	NT_CALL               NodeType = 17
	NT_POSTFIX_EXPRESSION NodeType = 18
	NT_TERNARY_EXPRESSION NodeType = 19
	NT_UNARY_EXPRESSION   NodeType = 20
	NT_UNARY_INC_DEC      NodeType = 21
	NT_BINARY_EXPRESSION  NodeType = 22
	NT_LOGICAL_EXPRESSION NodeType = 23
	//
	NT_VARIABLE_DEC         = 24
	NT_LOCAL_DEC            = 25
	NT_CONST_DEC            = 26
	NT_EMPTY_EXPRESSION     = 27
	NT_EXPRESSION_STATEMENT = 28
	NT_FILE                 = 29
)

type Node_t struct {
	NodeType          NodeType
	Terminal          *terminal_node_t
	HeadlessFunction  *headless_function_node_t
	Array             *array_node_t
	Object            *object_node_t
	MemberAccess      *member_access_node_t
	IndexAccess       *index_access_node_t
	Call              *call_node_t
	PostfixExpression *postfix_expression_node_t
	TernaryExpression *ternary_expression_node_t
	UnaryExpression   *unary_expression_node_t
	BinaryExpression  *binary_expression_node_t
	//
	Declairation        *variable_declairation_node_t
	ExpressionStatement *expression_statement_node_t
	File                *file_node_t
	//
	Position *shared.Position_t
}

type terminal_node_t struct {
	Value string
}

type headless_function_node_t struct {
	Parameters *[][]interface{}
	Body       *[]*Node_t
}

type array_node_t struct {
	Elements *[]*Node_t
	Size     int
}

type object_node_t struct {
	Members *[][]*Node_t
	Size    int
}

type member_access_node_t struct {
	Object *Node_t
	Member string
}

type index_access_node_t struct {
	Object *Node_t
	Index  *Node_t
}

type call_node_t struct {
	Object *Node_t
	Args   *[]*Node_t
}

type postfix_expression_node_t struct {
	Operator string
	Operand  *Node_t
}

type ternary_expression_node_t struct {
	Condition, Trueval, Falseval *Node_t
}

type unary_expression_node_t struct {
	Operator string
	Operand  *Node_t
}

type binary_expression_node_t struct {
	Operator    string
	Left, Right *Node_t
}

type variable_declairation_node_t struct {
	Declairations *[][]interface{}
}

type expression_statement_node_t struct {
	Expression *Node_t
}

type file_node_t struct {
	Body *[]*Node_t
}

func TerminalNode(nType NodeType, value string, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = nType
	node.Position = position
	//
	node.Terminal = &terminal_node_t{Value: value}
	return node
}

func HeadlessFunctionNode(parameters *[][]interface{}, body *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_HEADLESS_FUNCTION
	node.Position = position
	//
	node.HeadlessFunction = &headless_function_node_t{
		Parameters: parameters,
		Body:       body,
	}
	return node
}

func ArrayNode(elements *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_ARRAY
	node.Position = position
	//
	node.Array = &array_node_t{
		Elements: elements,
		Size:     len(*elements),
	}
	return node
}

func ObjectNode(members *[][]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_OBJECT
	node.Position = position
	//
	node.Object = &object_node_t{
		Members: members,
		Size:    len(*members),
	}
	return node
}

func MemberAccess(object *Node_t, member string, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_MEMBER_ACCESS
	node.Position = position
	//
	node.MemberAccess = &member_access_node_t{
		Object: object,
		Member: member,
	}
	return node
}

func IndexAccess(object, index *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_INDEX
	node.Position = position
	//
	node.IndexAccess = &index_access_node_t{
		Object: object,
		Index:  index,
	}
	return node
}

func Call(object *Node_t, args *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_CALL
	node.Position = position
	//
	node.Call = &call_node_t{
		Object: object,
		Args:   args,
	}
	return node
}

func PostfixExpressionNode(operator string, operand *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_POSTFIX_EXPRESSION
	node.Position = operand.Position
	//
	node.PostfixExpression = &postfix_expression_node_t{
		Operator: operator,
		Operand:  operand,
	}
	return node
}

func TernaryExpressionNode(condition, trueval, falseval *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_TERNARY_EXPRESSION
	node.Position = position
	//
	node.TernaryExpression = &ternary_expression_node_t{
		Condition: condition,
		Trueval:   trueval,
		Falseval:  falseval,
	}
	return node
}

func UnaryExpressionNode(operator string, operand *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_UNARY_EXPRESSION
	node.Position = position
	//
	node.UnaryExpression = &unary_expression_node_t{
		Operator: operator,
		Operand:  operand,
	}
	return node
}

func BinaryExpressionNode(NodeType NodeType, operator string, left, right *Node_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NodeType
	node.Position = left.Position.Merge(right.Position)
	//
	node.BinaryExpression = &binary_expression_node_t{
		Operator: operator,
		Left:     left,
		Right:    right,
	}

	return node
}

func VariableDeclairationNode(NodeType NodeType, declairations *[][]interface{}, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NodeType
	node.Position = position
	//
	node.Declairation = &variable_declairation_node_t{
		Declairations: declairations,
	}
	return node
}

func EmptyExpressionNode(position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_EMPTY_EXPRESSION
	node.Position = position
	return node
}

func ExpressionStatementNode(expression *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_EXPRESSION_STATEMENT
	node.Position = position
	//
	node.ExpressionStatement = new(expression_statement_node_t)
	node.ExpressionStatement = &expression_statement_node_t{
		Expression: expression,
	}
	return node
}

func FileNode(body *[]*Node_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_FILE
	node.Position = nil
	//
	node.File = new(file_node_t)
	node.File.Body = body

	return node
}
