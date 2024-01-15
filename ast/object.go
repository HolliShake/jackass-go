package ast

import "jackass/shared"

type object_node_t struct {
	Members *[][]*Node_t
	Size    int
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
