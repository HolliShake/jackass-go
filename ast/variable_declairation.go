package ast

import "jackass/shared"

type variable_declairation_node_t struct {
	Declairations *[][]interface{}
}

func VariableDeclairationNode(NodeType NodeType, declairations *[][]interface{}, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NodeType
	node.Position = position
	//
	node.Declairation = &variable_declairation_node_t{
		Declairations: declairations,
	}
	return node
}
