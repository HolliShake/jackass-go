package ast

import "jackass/shared"

type headless_function_node_t struct {
	Parameters *[][]interface{}
	Body       *[]*Node_t
}

func HeadlessFunctionNode(parameters *[][]interface{}, body *[]*Node_t, position *shared.Position_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_HEADLESS_FUNCTION
	node.Position = position
	//
	node.HeadlessFunction = &headless_function_node_t{
		Parameters: parameters,
		Body:       body,
	}
	return node
}
