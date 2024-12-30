package parser

import (
	"testing"

	"github.com/letsmakecakes/jsonparser/internal/lexer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
		expected Node
	}{
		{
			name:     "Empty Object",
			input:    "{}",
			hasError: false,
			expected: &ObjectValue{Pairs: map[string]Value{}},
		},
		{
			name:     "Simple Key-Value Pair",
			input:    `{"key": "value"}`,
			hasError: false,
			expected: &ObjectValue{Pairs: map[string]Value{
				"key": &StringValue{Value: "value"},
			}},
		},
		{
			name: "Nested Structure",
			input: `{
                "object": {},
                "array": [],
                "nested": {"key": ["value", 123, true, null]}
            }`,
			hasError: false,
			expected: &ObjectValue{Pairs: map[string]Value{
				"object": &ObjectValue{Pairs: map[string]Value{}},
				"array":  &ArrayValue{Elements: []Value{}},
				"nested": &ObjectValue{Pairs: map[string]Value{
					"key": &ArrayValue{Elements: []Value{
						&StringValue{Value: "value"},
						&NumberValue{Value: 123},
						&BooleanValue{Value: true},
						&NullValue{},
					}},
				}},
			}},
		},
		{
			name:     "Incomplete Object",
			input:    "{",
			hasError: true,
		},
		{
			name:     "Missing Value",
			input:    `{"key"}`,
			hasError: true,
		},
		{
			name:     "Trailing Comma",
			input:    `{"key": "value",}`,
			hasError: true,
		},
		{
			name:     "Invalid Value",
			input:    `{"key": undefined}`,
			hasError: true,
		},
		{
			name:     "Valid Key-Value Pair",
			input:    `{"key": "value"}`,
			hasError: false,
			expected: &ObjectValue{Pairs: map[string]Value{
				"key": &StringValue{Value: "value"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			node, err := p.Parse()

			if tt.hasError && err == nil {
				t.Errorf("expected error for input %q, got none", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
			}
			if !tt.hasError && err == nil {
				if !compareNodes(node, tt.expected) {
					t.Errorf("expected node %v, got %v", tt.expected, node)
				}
			}
		})
	}
}

// compareNodes is a helper function to compare two AST nodes for equality.
func compareNodes(node1, node2 Node) bool {
	switch n1 := node1.(type) {
	case *ObjectValue:
		n2, ok := node2.(*ObjectValue)
		if !ok {
			return false
		}
		if len(n1.Pairs) != len(n2.Pairs) {
			return false
		}
		for key, value1 := range n1.Pairs {
			value2, exists := n2.Pairs[key]
			if !exists || !compareNodes(value1, value2) {
				return false
			}
		}
		return true
	case *ArrayValue:
		n2, ok := node2.(*ArrayValue)
		if !ok {
			return false
		}
		if len(n1.Elements) != len(n2.Elements) {
			return false
		}
		for i, elem1 := range n1.Elements {
			if !compareNodes(elem1, n2.Elements[i]) {
				return false
			}
		}
		return true
	case *StringValue:
		n2, ok := node2.(*StringValue)
		return ok && n1.Value == n2.Value
	case *NumberValue:
		n2, ok := node2.(*NumberValue)
		return ok && n1.Value == n2.Value
	case *BooleanValue:
		n2, ok := node2.(*BooleanValue)
		return ok && n1.Value == n2.Value
	case *NullValue:
		_, ok := node2.(*NullValue)
		return ok
	default:
		return false
	}
}
