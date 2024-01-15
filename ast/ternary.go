package ast

import "jackass/shared"

type ternary_expression_node_t struct {
	Condition, Trueval, Falseval *Node_t
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
