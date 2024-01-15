package ast

import "jackass/shared"

type expression_statement_node_t struct {
	Expression *Node_t
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
