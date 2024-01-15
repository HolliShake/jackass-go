package ast

import "jackass/shared"

type terminal_node_t struct {
	Value string
}

func TerminalNode(nType NodeType, value string, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = nType
	node.Position = position
	//
	node.Terminal = &terminal_node_t{Value: value}
	return node
}
