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
	}{
		{
			name:     "Empty Object",
			input:    "{}",
			hasError: false,
		},
		{
			name:     "Simple Key-Value Pair",
			input:    `{"key": "value"}`,
			hasError: false,
		},
		{
			name: "Nested Structure",
			input: `{
				"object": {},
				"array": [],
				"nested": {"key": ["value", 123, true, null]}
			}`,
			hasError: false,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			err := p.Parse()

			if tt.hasError && err == nil {
				t.Errorf("expected error for input %q, got none", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
			}
		})
	}
}
