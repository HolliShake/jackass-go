package ast

import "jackass/shared"

func EmptyExpressionNode(position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_EMPTY_EXPRESSION
	node.Position = position
	return node
}
