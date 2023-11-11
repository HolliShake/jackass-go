package main

type NodeType = int

const (
	NT_ID NodeType = 0
	NT_INTEGER NodeType = 1
	NT_BIG_INTEGER NodeType = 2
	NT_OTHER_INTEGER NodeType = 3
	NT_OTHER_BIG_INTEGER NodeType = 4
	NT_FLOAT NodeType = 5
	NT_OTHER_FLOAT NodeType = 6
	NT_STRING NodeType = 7
	NT_BOOLEAN NodeType = 8
	NT_NULL NodeType = 9
	NT_SELF NodeType = 10
	NT_SUPER NodeType = 11
	NT_HEADLESS_FUNCTION NodeType = 12
	NT_ARRAY NodeType = 13
	NT_OBJECT NodeType = 14
	NT_MEMBER_ACCESS NodeType = 15
	NT_INDEX NodeType = 16
	NT_CALL NodeType = 17
	NT_POSTFIX_EXPRESSION NodeType = 18
	NT_TERNARY_EXPRESSION NodeType = 19
	NT_UNARY_EXPRESSION NodeType = 20
	NT_UNARY_INC_DEC NodeType = 21
	NT_BINARY_EXPRESSION NodeType = 22
	NT_LOGICAL_EXPRESSION NodeType = 23
	// 
	NT_VARIABLE_DEC = 24
	NT_LOCAL_DEC = 25
	NT_CONST_DEC = 26
	NT_EMPTY_EXPRESSION = 27
	NT_EXPRESSION_STATEMENT = 28
)

type node_t struct {
	ntype NodeType
	terminal *terminal_node_t
	headlessFunction *headless_function_node_t
	array *array_node_t
	object *object_node_t
	memberAccess *member_access_node_t
	indexAccess *index_access_node_t
	call *call_node_t
	postfixExpression *postfix_expression_node_t
	ternaryExpression *ternary_expression_node_t
	unaryExpression *unary_expression_node_t
	binaryExpression *binary_expression_node_t
	// 
	declairation *variable_declairation_node_t
	expressionStatement *expression_statement_node_t
	// 
	position *position_t
}


type terminal_node_t struct {
	value string
}

type headless_function_node_t struct {
	parameters *[][]interface{}
	body *[]*node_t
}

type array_node_t struct {
	elements *[]*node_t
	size int
}

type object_node_t struct {
	members *[][]*node_t
	size int
}

type member_access_node_t struct {
	object *node_t
	member string
}

type index_access_node_t struct {
	object *node_t
	index *node_t
}

type call_node_t struct {
	object *node_t
	args *[]*node_t
}

type postfix_expression_node_t struct {
	operator string
	operand *node_t
}

type ternary_expression_node_t struct {
	condition, trueval, falseval *node_t
}

type unary_expression_node_t struct {
	operator string
	operand *node_t
}

type binary_expression_node_t struct {
	operator string
	left, right *node_t
}

type variable_declairation_node_t struct {
	declairations *[][]interface{}
}

type expression_statement_node_t struct {
	expression *node_t
}

// 
func TerminalNode(ntype NodeType, value string, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = ntype
	node.position = position
	// 
	node.terminal = new(terminal_node_t)
	node.terminal.value = value
	return node
}

func HeadlessFunctionNode(parameters *[][]interface{}, body *[]*node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_HEADLESS_FUNCTION
	node.position = position
	// 
	node.headlessFunction = new(headless_function_node_t)
	node.headlessFunction.parameters = parameters
	node.headlessFunction.body = body
	return node
}

func ArrayNode(elements *[]*node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_ARRAY
	node.position = position
	// 
	node.array = new(array_node_t)
	node.array.elements = elements
	node.array.size = len(*elements)
	return node
}

func ObjectNode(members *[][]*node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_OBJECT
	node.position = position
	// 
	node.object = new(object_node_t)
	node.object.members = members
	node.object.size = len(*members)
	return node
}

func MemberAccess(object *node_t, member string, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_MEMBER_ACCESS
	node.position = position
	// 
	node.memberAccess = new(member_access_node_t)
	node.memberAccess.object = object
	node.memberAccess.member = member
	return node
}

func IndexAccess(object, index *node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_INDEX
	node.position = position
	// 
	node.indexAccess = new(index_access_node_t)
	node.indexAccess.object = object
	node.indexAccess.index = index
	return node
}

func Call(object *node_t, args *[]*node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_CALL
	node.position = position
	// 
	node.call = new(call_node_t)
	node.call.object = object
	node.call.args = args
	return node
}

func PostfixExpressionNode(operator string, operand *node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_POSTFIX_EXPRESSION
	node.position = operand.position
	// 
	node.postfixExpression = new(postfix_expression_node_t)
	node.postfixExpression.operator = operator
	node.postfixExpression.operand = operand
	return node
}

func TernaryExpressionNode(condition, trueval, falseval *node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_TERNARY_EXPRESSION
	node.position = position
	// 
	node.ternaryExpression = new(ternary_expression_node_t)
	node.ternaryExpression.condition = condition
	node.ternaryExpression.trueval = trueval
	node.ternaryExpression.falseval = falseval
	return node
}

func UnaryExpressionNode(operator string, operand *node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_UNARY_EXPRESSION
	node.position = position
	// 
	node.unaryExpression = new(unary_expression_node_t)
	node.unaryExpression.operator = operator
	node.unaryExpression.operand = operand
	return node
}

func BinaryExpressionNode(ntype NodeType, operator string, left, right *node_t) *node_t {
	node := new(node_t)
	node.ntype = ntype
	node.position = left.position.merge(right.position)
	// 
	node.binaryExpression = new(binary_expression_node_t)
	node.binaryExpression.operator = operator
	node.binaryExpression.left = left
	node.binaryExpression.right = right

	return node
}

func VariableDeclairationNode(ntype NodeType, declairations *[][]interface{}, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = ntype
	node.position = position
	//
	node.declairation = new(variable_declairation_node_t)
	node.declairation.declairations = declairations
	return node
}

func EmptyExpressionNode(position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_EMPTY_EXPRESSION
	node.position = position
	return node
}

func ExpressionStatementNode(expression *node_t, position *position_t) *node_t {
	node := new(node_t)
	node.ntype = NT_EXPRESSION_STATEMENT
	node.position = position
	// 
	node.expressionStatement = new(expression_statement_node_t)
	node.expressionStatement.expression = expression
	return node
}