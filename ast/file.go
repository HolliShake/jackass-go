package ast

type file_node_t struct {
	Body *[]*Node_t
}

func FileNode(body *[]*Node_t) *Node_t {
	node := new(Node_t)
	node.NodeType = NT_FILE
	node.Position = nil
	//
	node.File = new(file_node_t)
	node.File.Body = body

	return node
}
