package ast

type Value interface{}

type Object struct {
	Pairs map[string]Value
}

type Array struct {
	Elements []Value
}

type String struct {
	Value string
}

type Number struct {
	Value string
}

type Boolean struct {
	Value string
}

type Null struct{}
