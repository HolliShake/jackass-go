package ast

type binary_expression_node_t struct {
	Operator    string
	Left, Right *Node_t
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
