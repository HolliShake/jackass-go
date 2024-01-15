package ast

import "jackass/shared"

type array_node_t struct {
	Elements *[]*Node_t
	Size     int
}

func ArrayNode(elements *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_ARRAY
	node.Position = position
	//
	node.Array = &array_node_t{
		Elements: elements,
		Size:     len(*elements),
	}
	return node
}
