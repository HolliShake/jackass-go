package ast

import "jackass/shared"

type call_node_t struct {
	Object *Node_t
	Args   *[]*Node_t
}

func Call(object *Node_t, args *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_CALL
	node.Position = position
	//
	node.Call = &call_node_t{
		Object: object,
		Args:   args,
	}
	return node
}
