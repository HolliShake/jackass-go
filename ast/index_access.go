package ast

import "jackass/shared"

type index_access_node_t struct {
	Object *Node_t
	Index  *Node_t
}

func IndexAccess(object, index *Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_INDEX
	node.Position = position
	//
	node.IndexAccess = &index_access_node_t{
		Object: object,
		Index:  index,
	}
	return node
}
