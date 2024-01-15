package ast

import "jackass/shared"

type postfix_expression_node_t struct {
	Operator string
	Operand  *Node_t
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
