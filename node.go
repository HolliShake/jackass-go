package main

type NodeType = int

const (
	NT_ID NodeType = 0
	NT_INTEGER NodeType = 1
	NT_BIG_INTEGER NodeType = 2
	NT_OTHER_INTEGER NodeType = 3
	NT_OTHER_BIG_INTEGER NodeType = 4
	NT_FLOAT NodeType = 5
	NT_OTHER_FLOAT NodeType = 6
	NT_STRING NodeType = 7
	NT_BOOLEAN NodeType = 8
	NT_NULL NodeType = 9
	NT_BINARY_EXPRESSION NodeType = 4
	NT_LOGICAL_EXPRESSION NodeType = 5
)

type node_t struct {
	ntype NodeType
	terminal *terminal_node_t
	binaryExpression *binary_expression_node_t
	// 
	position *position_t
}


type terminal_node_t struct {
	value string
}

type binary_expression_node_t struct {
	operator string
	left, right *node_t
}

// 
func TerminalNode(ntype NodeType, value string, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = ntype
	node.position = position
	// 
	node.terminal = new(terminal_node_t)
	node.terminal.value = value
	return node
}

func BinaryExpressionNode(ntype NodeType, operator string, left, right *node_t) *node_t {
	node := new(node_t)
	node.ntype = ntype
	node.position = left.position.merge(right.position)
	// 
	node.binaryExpression = new(binary_expression_node_t)
	node.binaryExpression.operator = operator
	node.binaryExpression.left = left
	node.binaryExpression.right = right

	return node
}


