package validator

import (
	"testing"

	"github.com/letsmakecakes/jsonparser/internal/parser"
)

func TestValidator(t *testing.T) {
	t.Run("Test Depth Validation", func(t *testing.T) {
		v := New(3)

		// Create an AST with nested objects within the limit
		validNode := &parser.ObjectValue{Pairs: map[string]parser.Value{
			"level1": &parser.ObjectValue{Pairs: map[string]parser.Value{
				"level2": &parser.ObjectValue{Pairs: map[string]parser.Value{
					"level3": &parser.StringValue{Value: "value"},
				}},
			}},
		}}

		if err := v.Validate(validNode); err != nil {
			t.Errorf("unexpected error for valid nested object: %v", err)
		}

		// Create an AST that exceeds the depth limit
		invalidNode := &parser.ObjectValue{Pairs: map[string]parser.Value{
			"level1": &parser.ObjectValue{Pairs: map[string]parser.Value{
				"level2": &parser.ObjectValue{Pairs: map[string]parser.Value{
					"level3": &parser.ObjectValue{Pairs: map[string]parser.Value{
						"level4": &parser.StringValue{Value: "value"},
					}},
				}},
			}},
		}}

		if err := v.Validate(invalidNode); err == nil {
			t.Error("expected error when exceeding max depth")
		}
	})

	t.Run("Test String Validation", func(t *testing.T) {
		v := New(10)

		tests := []struct {
			name     string
			input    string
			hasError bool
		}{
			{"Valid String", "valid string", false},
			{"String with Null Character", "string with \u0000", true},
			{"String with Control Character", "string with \u001F", true},
			{"Normal ASCII Characters", "normal ascii {}<>", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateString(tt.input)
				if tt.hasError && err == nil {
					t.Errorf("expected error for input %q, got none", tt.input)
				}
				if !tt.hasError && err != nil {
					t.Errorf("unexpected error for input %q: %v", tt.input, err)
				}
			})
		}
	})
}
