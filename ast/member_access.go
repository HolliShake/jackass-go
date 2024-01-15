package ast

import "jackass/shared"

type member_access_node_t struct {
	Object *Node_t
	Member string
}

func MemberAccess(object *Node_t, member string, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_MEMBER_ACCESS
	node.Position = position
	//
	node.MemberAccess = &member_access_node_t{
		Object: object,
		Member: member,
	}
	return node
}
