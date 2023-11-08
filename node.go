package main

type NodeType = int

const (
	NT_INTEGER NodeType = 0
	NT_BINARY_EXPRESSION NodeType = 1
	NT_LOGICAL_EXPRESSION NodeType = 2
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


