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
