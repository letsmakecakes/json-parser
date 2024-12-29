package parser

// Node represents a node in the AST.
type Node interface {
	TokenLiteral() string
}

// Value represents a value node in the AST.
type Value interface {
	Node
	valueNode()
}

// ObjectValue represents an object value with key-value pairs.
type ObjectValue struct {
	Pairs map[string]Value
}

func (o *ObjectValue) TokenLiteral() string {
	return "{"
}

func (o *ObjectValue) valueNode() {}

// ArrayValue represents an array value with elements.
type ArrayValue struct {
	Elements []Value
}

func (a *ArrayValue) TokenLiteral() string {
	return "["
}

func (a *ArrayValue) valueNode() {}

// StringValue represents a string value.
type StringValue struct {
	Value string
}

func (s *StringValue) TokenLiteral() string {
	return s.Value
}

func (s *StringValue) valueNode() {}

// NumberValue represents a number value.
type NumberValue struct {
	Value float64
}

func (n *NumberValue) TokenLiteral() string {
	return "number"
}

func (n *NumberValue) valueNode() {}

// BooleanValue represents a boolean value.
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) TokenLiteral() string {
	return "boolean"
}

func (b *BooleanValue) valueNode() {}

// NullValue represents a null value.
type NullValue struct{}

func (n *NullValue) TokenLiteral() string {
	return "null"
}

func (n *NullValue) valueNode() {}
