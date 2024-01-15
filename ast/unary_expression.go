package ast

import "jackass/shared"

type unary_expression_node_t struct {
	Operator string
	Operand  *Node_t
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
