package validator

import "testing"

func TestValidator(t *testing.T) {
	t.Run("Test Depth Validation", func(t *testing.T) {
		v := New(3)

		// Test within limit
		for i := 0; i < 3; i++ {
			if err := v.ValidateDepth(); err != nil {
				t.Errorf("unexpected error at depth %d: %v", i, err)
			}
		}

		// Test exceeding limit
		if err := v.ValidateDepth(); err == nil {
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
